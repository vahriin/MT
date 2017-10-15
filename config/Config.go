package config

type AppConfig struct {
	AppDbConfig
	Server     ServerConfig
	System     SystemConfig
}

type AppDbConfig struct {
	Db DbConfig
	//Cache CacheConfig
}


type SystemConfig struct {
	Logfile string
}

type SSLParams struct {
	Sslcert string
	Sslkey string
	Sslrootcert string
}

type DbConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
	//ConnectTimeout string
	Sslmode string
	//Sslparams SSLParams
}

/*type CacheConfig struct {

}*/

type ServerConfig struct {
	Port string
}
