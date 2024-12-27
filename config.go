package main

import (
	"encoding/json"
	"os"
	"path"
)

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type EngineConfig struct {
	Auth    AuthConfig     `json:"auth"`
	Scripts []ScriptConfig `json:"scripts"`
	Server  ServerConfig   `json:"server"`
}

type ScriptConfig struct {
	Cmd    []string `json:"cmd"`
	Topics []string `json:"topics"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func getConfig(config *EngineConfig) {
	var buffer = make([]byte, 65536)
	dataPath := os.Getenv("DATA_DIR")

	var configPath = "config.json"
	if dataPath != "" {
		configPath = path.Join(dataPath, configPath)
	}

	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	length, err := file.Read(buffer)
	if err != nil {
		panic(err)
	}

	buffer = buffer[0:length]
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		panic(err)
	}
}
