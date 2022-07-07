package infrastructure

import (
	"github.com/jackc/pgx/v4/pgxpool"
	postgres "github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres_pgx"
	"github.com/pkg/errors"
)

func (ic *infrastructureConfigurator) configPostgres() (*pgxpool.Pool, error, func()) {
	pgxConn, err := postgres.NewPgxConn(ic.cfg.Postgresql)
	if err != nil {
		return nil, errors.Wrap(err, "postgresql.NewPgxConn"), nil
	}

	ic.log.Infof("postgres connected: %v", pgxConn.Stat().TotalConns())

	return pgxConn, nil, func() {
		pgxConn.Close()
	}
}
