package post

import (
	"banner-service/internal/entity"
	"banner-service/internal/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (pg *postgres) GetContent(ctx context.Context, feature_id int64, tag_id int64) (entity.Content, error) {
	query := `
	select title, text, url 
	from banners
	where feature_id = @feature_id and @tag_id = ANY(tag_ids)`

	args := pgx.NamedArgs{
		"feature_id": feature_id,
		"tag_id":     tag_id,
	}

	rows, err := pg.db.Query(ctx, query, args)
	if err != nil {
		return entity.Content{}, fmt.Errorf("unable to query content: %w", err)
	}

	defer rows.Close()
	d, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Content])

	if len(d) < 1 {
		return entity.Content{}, storage.ErrContentNotFound
	}

	if err != nil {
		return entity.Content{}, fmt.Errorf("unable to query content: %w", err)
	}

	return d[0], err
}
