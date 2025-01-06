package cloudpg

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func CloudSQL(dsn, instanceConnectionName string) (*sql.DB, func() error, error) {
	d, err := cloudsqlconn.NewDialer(context.Background(), cloudsqlconn.WithIAMAuthN())
	if err != nil {
		return nil, nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}
	var opts []cloudsqlconn.DialOption
	opts = append(opts, cloudsqlconn.WithPrivateIP())

	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, nil, err
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	sDb, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, nil, fmt.Errorf("sql.Open: %w", err)
	}
	err = sDb.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("error pinging db: %v", err)
	}
	return sDb, func() error {
		if d != nil {
			_ = d.Close()
		}
		if sDb != nil {
			_ = sDb.Close()
		}
		return nil
	}, nil
}
