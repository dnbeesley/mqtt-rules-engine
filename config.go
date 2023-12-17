package main

import (
	"encoding/json"
	"os"
)

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Condition struct {
	Index int           `json:"index"`
	Topic string        `json:"topic"`
	Type  ConditionType `json:"type"`
	Value uint8         `json:"value"`
}

type ConditionType string

type EngineConfig struct {
	Auth    AuthConfig   `json:"auth"`
	Outputs []Output     `json:"outputs"`
	Server  ServerConfig `json:"server"`
}

type Output struct {
	DefaultValue string `json:"defaultValue"`
	Rules        []Rule `json:"rules"`
	Topic        string `json:"topic"`
}

type Rule struct {
	Conditions []Condition `json:"conditions"`
	Result     string      `json:"result"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
}

const (
	BitMatchConditionType    ConditionType = "BitMatch"
	GreaterThanConditionType ConditionType = "GreaterThan"
)

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
