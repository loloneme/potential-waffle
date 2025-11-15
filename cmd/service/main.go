package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/loloneme/potential-waffle/internal"
	generated "github.com/loloneme/potential-waffle/internal/generated/openapi"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/pull_request"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/team"
	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/repository/user"
	mw "github.com/loloneme/potential-waffle/internal/middleware"
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_create_post"
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_merge_post"
	"github.com/loloneme/potential-waffle/internal/rpc/pull_request/pr_reassign_post"
	"github.com/loloneme/potential-waffle/internal/rpc/service/adapter"
	"github.com/loloneme/potential-waffle/internal/rpc/team/team_add_post"
	"github.com/loloneme/potential-waffle/internal/rpc/team/team_get_get"
	"github.com/loloneme/potential-waffle/internal/rpc/user/users_get_review_get"
	"github.com/loloneme/potential-waffle/internal/rpc/user/users_set_is_active_post"
	"github.com/loloneme/potential-waffle/internal/usecase/create_pr"
	"github.com/loloneme/potential-waffle/internal/usecase/create_team"
	"github.com/loloneme/potential-waffle/internal/usecase/merge_pr"
	"github.com/loloneme/potential-waffle/internal/usecase/reassign_pr"
)

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := internal.NewDatabaseConnection(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error connecting to database: %v", err))
		panic(err)
	}

	defer db.Close()

	userRepo := user.NewRepository(db)
	teamRepo := team.NewRepository(db)
	prRepo := pull_request.NewRepository(db)

	createTeamService := create_team.New(teamRepo, userRepo)
	createPullRequestService := create_pr.New(userRepo, prRepo)
	mergePullRequestService := merge_pr.New(prRepo)
	reassignPullRequestService := reassign_pr.New(userRepo, prRepo)

	createTeamHandler := team_add_post.New(createTeamService)
	getTeamHandler := team_get_get.New(userRepo, teamRepo)
	createPullRequestHandler := pr_create_post.New(createPullRequestService)
	mergePullRequestHandler := pr_merge_post.New(mergePullRequestService)
	reassignPullRequestHandler := pr_reassign_post.New(reassignPullRequestService)
	getUsersReviewHandler := users_get_review_get.New(prRepo)
	setIsActiveHandler := users_set_is_active_post.New(userRepo)

	serviceAdapter := adapter.NewAdapter(
		createTeamHandler,
		getTeamHandler,
		createPullRequestHandler,
		mergePullRequestHandler,
		reassignPullRequestHandler,
		getUsersReviewHandler,
		setIsActiveHandler,
	)

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(mw.NewOpenAPIMiddleware(getOpenAPIPath()))
	e.Use(mw.SlogMiddleware(logger))

	generated.RegisterHandlers(e, serviceAdapter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	go func() {
		logger.Info(fmt.Sprintf("Starting reviewers-app HTTP server on %s\n", addr))
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("failed to start server: %v", err))
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("failed to start server: %v", err))
		panic(err)
	}

	log.Println("Server stopped")
}

func getOpenAPIPath() string {
	specFile := os.Getenv("OPENAPI_SPEC_PATH")
	if specFile == "" {
		specFile = "api/openapi.yaml"
	}

	return specFile
}
