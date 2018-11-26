package model

type Config struct {
	ServerPort        int
	DatabaseHost      string
	DatabasePort      string
	ServerReadTimeout int
	WriteTimeout      int
	ReadTimeout       int
	IdleTimeout       int
	DBName            string
	UserName          string
	Password          string
}

func NewConfig() *Config {
	return &Config{
		ServerPort:        9090,
		DatabaseHost:      "localhost",
		DatabasePort:      "5432",
		ServerReadTimeout: 300,
		WriteTimeout:      300,
		ReadTimeout:       300,
		IdleTimeout:       300,
		DBName:            "simplechat",
		UserName:          "shailendra",
		Password:          "",
	}
}
