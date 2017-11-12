package config

import (
	"encoding/json"
	"os"
)

var ApplicationConfigDir string = "moneyteam-devel/"
var SystemConfigFile string = "system"
var DbConfigFile string = "database"
var ServerConfigFile = "server"

//var CacheConfigFile string = "cache"

func ReadConfig() *AppConfig {
	userDir := os.Getenv("HOME")
	if checkConfigDirExist(userDir) {
		return readConfig(userDir + "/.config/" + ApplicationConfigDir)
	} else {
		panic("No config directory")
	}
}

func readConfig(pathToConfigDir string) *AppConfig {
	appConfig := new(AppConfig)
	appConfig.System = readSystemConfig(pathToConfigDir + SystemConfigFile)
	appConfig.Db = readDbConfig(pathToConfigDir + DbConfigFile)
	appConfig.Server = readServerConfig(pathToConfigDir + ServerConfigFile)
	//appConfig.Cache = readCacheConfig(pathToConfigDir + CacheConfigFile)
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
