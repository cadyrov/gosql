package internal

const (
	ModeMySQL = "mysql"
	ModePSQL  = "psql"
)

type Config struct {
	Mode     string `long:"mode" env:"GOSQL_MODE"  required:"true"`
	Scheme   string `long:"scheme"  required:"true"`
	Table    string `long:"table"  required:"true"`
	Template string `long:"template"  required:"true"`
	Result   string `long:"result"  required:"true"`
	SQL      SQL
	MYSQL    MYSQL
	PSQL     PSQL
}

type SQL struct {
	Host     string `long:"sql-host" env:"GOSQL_SQL_HOST"`
	Port     int    `long:"sql-port" env:"GOSQL_SQL_PORT"`
	UserName string `long:"sql-user-name" env:"GOSQL_SQL_USER_NAME"`
	DBName   string `long:"sql-db-name" env:"GOSQL_SQL_DB_NAME"`
	Password string `long:"sql-password" env:"GOSQL_SQL_PASSWORD"`
}

type MYSQL struct {
}

type PSQL struct {
	SslMode        string `long:"psql-ssl-mode" env:"GOSQL_PSQL_SSL_MODE"`
	Binary         bool   `long:"psql-binary" env:"GOSQL_PSQL_BINARY"`
	MaxConnections int    `long:"psql-max-connections" env:"GOSQL_PSQL_MAX_CONNECTIONS"`
	ConnectionIdle int    `long:"psql-connection-idle" env:"GOSQL_PSQL_CONNECTION_IDLE"`
	MaxLimit       int    `long:"psql-max-limit" env:"GOSQL_PSQL_MAX_LIMIT"`
}
