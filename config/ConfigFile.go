package config

import (
	"encoding/json"
	"os"
)

const ApplicationConfigDir  = "moneyteam/"
const SystemConfigFile = "system"
const DbConfigFile = "database"
const ServerConfigFile = "server"
var PathToConfigDir = ""

//var CacheConfigFile string = "cache"

func ReadConfig() *AppConfig {
	userDir := os.Getenv("HOME")
	if checkConfigDirExist(userDir) {
		PathToConfigDir = userDir + "/.config/" + ApplicationConfigDir
		return readConfig()
	} else {
		panic("No config directory")
	}
}

func readConfig() *AppConfig {
	appConfig := new(AppConfig)
	appConfig.System = readSystemConfig(PathToConfigDir + SystemConfigFile)
	appConfig.Db = readDbConfig(PathToConfigDir + DbConfigFile)
	appConfig.Server = readServerConfig(PathToConfigDir + ServerConfigFile)
	//appConfig.Cache = readCacheConfig(PathToConfigDir + CacheConfigFile)
	return appConfig
}

func checkConfigDirExist(home string) bool {
	_, err := os.Stat(home + "/.config/" + ApplicationConfigDir)
	return !os.IsNotExist(err)
}

func readSystemConfig(fileName string) SystemConfig {
	configFile, err := os.Open(fileName + ".json")
	if err != nil {
		panic(err) //add create config file
	}

	var systemConfig SystemConfig
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&systemConfig)
	if err != nil {
		panic(err)
	}

	systemConfig.Logfile = PathToConfigDir + systemConfig.Logfile

	return systemConfig
}

func readDbConfig(fileName string) DbConfig {
	configFile, err := os.Open(fileName + ".json")
	if err != nil {
		panic(err)
	}

	var dbConfig DbConfig
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&dbConfig)
	if err != nil {
		panic(err)
	}

	return dbConfig
}

func readServerConfig(fileName string) ServerConfig {
	configFile, err := os.Open(fileName + ".json")
	if err != nil {
		panic(err)
	}

	var sConfig ServerConfig
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&sConfig)
	if err != nil {
		panic(err)
	}

	sConfig.KeyFile = PathToConfigDir + sConfig.KeyFile
	sConfig.CertFile = PathToConfigDir + sConfig.CertFile

	return sConfig
}

/*func readCacheConfig(fileName string) CacheConfig {
	configFile, err := os.Open(fileName + ".json")
	if err != nil {
		panic(err)
	}

	var sConfig CacheConfig
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&sConfig)
	if err != nil {
		panic(err)
	}

	return sConfig
}*/
