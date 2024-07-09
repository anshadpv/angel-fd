package constants

const (
	PostgresDBKey = "postgres-db"
	MSSQLDBKey    = "mssql-db"

	DBDriver                         = "drivername"
	DBURL                            = "url"
	DBMaxOpenConnections             = "maxopenconnections"
	DBMaxIdleConnections             = "maxidleconnections"
	DBConnectionMaxLifetimeInSeconds = "connectionmaxlifetimeinseconds"
	DBConnectionMaxIdleTimeInSeconds = "connectionmaxidletimeinseconds"
	DBQueryTimeoutInMillisKey        = "querytimeoutinmillis"
	DBConnectionMaxTimeout           = "connectionmaxtimeout"
	DBDriverDefaultValue             = "postgres"
)

const (
	DATABASE_NAME     = "DATABASE_NAME"
	DATABASE_URL      = "DATABASE_URL"
	DATABASE_USERNAME = "DATABASE_USERNAME"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
)
