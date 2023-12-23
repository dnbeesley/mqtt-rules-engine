package main

import (
	"encoding/json"
	"os"
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
	File   string   `json:"file"`
	Topics []string `json:"topics"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

func getConfig(config *EngineConfig) {
	var buffer = make([]byte, 65536)
	file, err := os.Open("config.json")
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
