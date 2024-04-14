package post

import (
	"banner-service/internal/entity"
	"banner-service/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (pg *postgres) CreateBanner(ctx context.Context, feature_id int64, tag_ids []int64, is_active bool, content entity.Content) (int64, error) {
	query := `
	INSERT INTO banners (feature_id, tag_ids, is_active, title, text, url, created_at) 
	VALUES (@feature_id, @tag_ids, @is_active, @title, @text, @url, @created_at) RETURNING id`

	args := pgx.NamedArgs{
		"feature_id": feature_id,
		"tag_ids":    tag_ids,
		"is_active":  is_active,
		"title":      content.Title,
		"text":       content.Text,
		"url":        content.Url,
		"created_at": time.Now(),
	}

	result := pg.db.QueryRow(ctx, query, args)

	var id int64
	err := result.Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("unable to insert row: %w", err)
	}

	return id, nil
}

func (pg *postgres) GetBanners(ctx context.Context, feature_id, tag_id *int64, offset, limit *int) ([]entity.Banner, error) {
	query := `
	select *
	from banners `

	if feature_id != nil || tag_id != nil {
		query += ` where `
	}

	query = makeQuery(query, feature_id, tag_id, offset, limit)

	args := makeArgs(feature_id, tag_id, offset, limit)

	rows, err := pg.db.Query(ctx, query, args)

	if err != nil {
		return nil, storage.ErrInternalServer
	}

	defer rows.Close()
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Banner])

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, storage.ErrBannerNotFound
	}

	return result, err
}

func (pg *postgres) DeletBanner(ctx context.Context, id int64) error {
	query := `DELETE FROM banners WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	results, err := pg.db.Exec(ctx, query, args)

	if results.RowsAffected() == 0 {
		return storage.ErrBannerNotFound
	}

	if err != nil {
		return fmt.Errorf("unable to delete content: %w", err)
	}

	return nil
}

func makeQuery(query string, feature_id *int64, tag_id *int64, offset, limit *int) string {
	if feature_id != nil {
		query += ` feature_id = @feature_id `
	}

	if tag_id != nil {
		query += ` @tag_id = any(tag_ids) `
	}

	if offset != nil {
		query += ` offset @offset `
	}

	if limit != nil {
		query += ` limit @limit `
	}

	return query
}

func makeArgs(feature_id *int64, tag_id *int64, offset, limit *int) pgx.NamedArgs {
	if feature_id != nil {
		if tag_id != nil && offset != nil && limit != nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
				"tag_ids":    tag_id,
				"offset":     *offset,
				"limit":      *limit,
			}

			return query
		}
		if tag_id != nil && offset != nil && limit == nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
				"tag_ids":    tag_id,
				"offset":     *offset,
			}

			return query
		}
		if tag_id != nil && offset == nil && limit != nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
				"tag_ids":    tag_id,
				"limit":      *limit,
			}

			return query
		}
		if tag_id != nil && offset == nil && limit == nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
				"tag_ids":    tag_id,
			}

			return query
		}
		if tag_id == nil && offset != nil && limit != nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
				"offset":     *offset,
				"limit":      *limit,
			}

			return query
		}
		if tag_id == nil && offset != nil && limit == nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
				"offset":     *offset,
			}

			return query
		}
		if tag_id == nil && offset == nil && limit != nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
				"limit":      *limit,
			}

			return query
		}
		if tag_id == nil && offset == nil && limit == nil {
			query := pgx.NamedArgs{
				"feature_id": *feature_id,
			}

			return query
		}
	} else {
		if tag_id != nil && offset != nil && limit != nil {
			query := pgx.NamedArgs{
				"tag_ids": tag_id,
				"offset":  *offset,
				"limit":   *limit,
			}

			return query
		}
		if tag_id != nil && offset != nil && limit == nil {
			query := pgx.NamedArgs{
				"tag_ids": tag_id,
				"offset":  *offset,
			}

			return query
		}
		if tag_id != nil && offset == nil && limit != nil {
			query := pgx.NamedArgs{
				"tag_ids": tag_id,
				"limit":   *limit,
			}

			return query
		}
		if tag_id != nil && offset == nil && limit == nil {
			query := pgx.NamedArgs{
				"tag_ids": tag_id,
			}

			return query
		}
		if tag_id == nil && offset != nil && limit != nil {
			query := pgx.NamedArgs{
				"offset": *offset,
				"limit":  *limit,
			}

			return query
		}
		if tag_id == nil && offset != nil && limit == nil {
			query := pgx.NamedArgs{
				"offset": *offset,
			}

			return query
		}
		if tag_id == nil && offset == nil && limit != nil {
			query := pgx.NamedArgs{
				"limit": *limit,
			}

			return query
		}
		if tag_id == nil && offset == nil && limit == nil {
			query := pgx.NamedArgs{}

			return query
		}
	}
	return pgx.NamedArgs{}
}

/*
func (pg *postgres) PatchBanner(ctx context.Context, id int64, bannner entity.Banner) (bool, error){
	checkId, err := pg.CheckContentExists(ctx, id)

	if err != nil {
		return false, fmt.Errorf("unable to get content: %w", err)
	}

	if !checkId {
		return false, nil
	}

	args := ` UPDATE banner Set`

	if

	WHERE id = @id;

}



*/
