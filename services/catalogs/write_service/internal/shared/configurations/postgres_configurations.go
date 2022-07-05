package configurations

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdihadeli/store-golang-microservice-sample/pkg/postgres"
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
