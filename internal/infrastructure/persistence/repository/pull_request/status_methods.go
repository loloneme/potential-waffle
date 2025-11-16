package pull_request

import (
	"context"
	"database/sql"
	"errors"

	"github.com/loloneme/potential-waffle/internal/infrastructure/persistence/models"
)

func (r *Repository) FindStatus(ctx context.Context, spec FindSpecification) (*models.Status, error) {
	status := &models.Status{}

	builder := st.
		Select(spec.GetFields()...).From(r.statusTableName)

	sqlStr, params, err := spec.GetRule(builder).ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.GetContext(ctx, status, sqlStr, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrStatusNotFound
		}
		return nil, err
	}

	return status, nil
}
