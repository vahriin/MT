package config

type AppConfig struct {
	DataModule DataModuleConfig
	Server     ServerConfig
	System     SystemConfig
}

type DataModuleConfig struct {
	Db DbConfig
	//Cache CacheConfig
}

type SystemConfig struct {
	Logfile string
}

type SSLParams map[string]string

type DbConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
	Sslmode string
	Sslparams SSLParams
}

/*type CacheConfig struct {

}*/

type ServerConfig struct {
	port string
}
