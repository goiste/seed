package pg

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(ctx context.Context, uri string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	config.LazyConnect = false
	config.MinConns = 1
	config.ConnConfig.PreferSimpleProtocol = true
	config.ConnConfig.RuntimeParams = map[string]string{
		"standard_conforming_strings": "on",
	}

	return pgxpool.ConnectConfig(ctx, config)
}

func ScanReturnedIds(dst *[]any, rows pgx.Rows) error {
	defer rows.Close()

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return err
		}

		if len(values) < 1 {
			return fmt.Errorf("empty returned values")
		}

		*dst = append(*dst, values[0])
	}

	return nil
}
