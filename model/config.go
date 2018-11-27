package model

type Config struct {
	ServerSettings   ServerSettings
	DatabaseSettings DatabaseSettings
}

type ServerSettings struct {
	ServerPort        *int
	ServerReadTimeout *int
	WriteTimeout      *int
	ReadTimeout       *int
	IdleTimeout       *int
}

type DatabaseSettings struct {
	DriverName   *string
	DatabaseHost *string
	DatabasePort *string
	DBName       *string
	Password     *string
	UserName     *string
}
