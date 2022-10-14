package data

type MigrationParams struct {
	DbName        string `mapstructure:"dbName"`
	VersionTable  string `mapstructure:"versionTable"`
	MigrationsDir string `mapstructure:"migrationsDir"`
	TargetVersion uint   `mapstructure:"targetVersion"`
	SkipMigration bool   `mapstructure:"skipMigration"`
}
