package migrate

type MigrationConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	DBName        string `mapstructure:"dbName"`
	SSLMode       bool   `mapstructure:"sslMode"`
	Password      string `mapstructure:"password"`
	VersionTable  string `mapstructure:"versionTable"`
	MigrationsDir string `mapstructure:"migrationsDir"`
	TargetVersion uint   `mapstructure:"targetVersion"`
	SkipMigration bool   `mapstructure:"skipMigration"`
}
