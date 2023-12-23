package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var channel = make(chan string)
var subscribeTopics = map[string]string{}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var payload = string(msg.Payload())
	subscribeTopics[msg.Topic()] = payload
	println("Received update to topic:", msg.Topic())
	channel <- msg.Topic()
}

func main() {
	var config EngineConfig
	getConfig(&config)

	for _, script := range config.Scripts {
		for _, topic := range script.Topics {
			subscribeTopics[topic] = ""
		}
	}

	opts := mqtt.NewClientOptions()
	server := fmt.Sprintf("tcp://%s:%d", config.Server.Host, config.Server.Port)
	fmt.Println("Connecting to:", server)
	opts.AddBroker(server)
	opts.SetClientID(config.Auth.Username)
	opts.SetUsername(config.Auth.Username)
	opts.SetPassword(config.Auth.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for topic := range subscribeTopics {
		if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		fmt.Println("Subscribed to topic:", topic)
	}

	for topic := range channel {
		for _, script := range config.Scripts {
			match := false
			input := map[string]string{}
			for _, scriptTopic := range script.Topics {
				if subscribeTopics[scriptTopic] == "" {
					println("Mising data for topic:", scriptTopic)
					match = false
					break
				}

				input[scriptTopic] = subscribeTopics[scriptTopic]
				if topic == scriptTopic {
					match = true
				}
			}

			if !match {
				continue
			}

			inputBytes, err := json.Marshal(input)
			if err != nil {
				panic(err)
			}

			cmd := exec.Command("C:\\Python311\\python.exe", script.File)

			cmd.Stderr = os.Stderr
			cmd.Stdin = bytes.NewReader(inputBytes)
			var buffer bytes.Buffer
			cmd.Stdout = &buffer

			err = cmd.Run()
			if err != nil {
				panic(err)
			}

			output := map[string]string{}
			err = json.Unmarshal(buffer.Bytes(), &output)
			if err != nil {
				panic(err)
			}

			for topic, value := range output {
				client.Publish(topic, 0, true, value)
				fmt.Println("Publishing:", value, "to topic:", topic)
			}
		}
	}

	println("Finished reading channel")
}
