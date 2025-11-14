package pull_request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

func (r *Repository) FindStatus(ctx context.Context, spec FindSpecification) (*models.Status, error) {
	status := &models.Status{}

	builder := st.
		Select(spec.GetFields()...).From(r.statusTableName)

	sqlStr, params, err := spec.GetRule(builder).ToSql()

	if err != nil {
		return nil, fmt.Errorf("select status: %w", err)
	}

	err = r.db.GetContext(ctx, status, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStatusNotFound
		}
		return nil, fmt.Errorf("exec get status by id: %w", err)
	}

	return status, nil
}
