package configs

type MySQLConnectionParams struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

type ServerConfig struct {
	PortToStart string                `toml:"start_port"`
	ConnParams  MySQLConnectionParams `toml:"server"`
}

func CreateConfigForServer() *ServerConfig {
	return &ServerConfig{}
}
