package migration

import (
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config"
	"github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/config/environment"
	typeMapper "github.com/mehdihadeli/go-food-delivery-microservices/internal/pkg/reflection/typemapper"

	"github.com/iancoleman/strcase"
)

type CommandType string

const (
	Up   CommandType = "up"
	Down CommandType = "down"
)

type MigrationOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	DBName        string `mapstructure:"dbName"`
	SSLMode       bool   `mapstructure:"sslMode"`
	Password      string `mapstructure:"password"`
	VersionTable  string `mapstructure:"versionTable"`
	MigrationsDir string `mapstructure:"migrationsDir"`
	SkipMigration bool   `mapstructure:"skipMigration"`
}

func ProvideConfig(environment environment.Environment) (*MigrationOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[MigrationOptions]())

	return config.BindConfigKey[*MigrationOptions](optionName, environment)
}
