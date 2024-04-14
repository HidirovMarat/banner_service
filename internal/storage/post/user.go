package post

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var (
	Admin = "admin"
	User = "user"
)

func (pg *postgres) GetAccessLevel(ctx context.Context, login string, password string) (string, error) {
	query := `
	SELECT accesses
	FROM "user"
	WHERE login = @login AND password = @password`

	args := pgx.NamedArgs{
		"login":    login,
		"password": password,
	}

	var accessLevel string
	err := pg.db.QueryRow(ctx, query, args).Scan(&accessLevel)

	if err != nil {
		return "", fmt.Errorf("unable to query user: %w", err)
	}

	return accessLevel, nil
}
