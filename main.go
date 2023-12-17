package main

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var channel = make(chan bool)
var subscribeTopics = map[string][]uint8{}

func checkCondition(condition *Condition) bool {
	val, found := subscribeTopics[condition.Topic]
	if !found {
		return false
	}

	if len(val) <= condition.Index {
		return false
	}

	if condition.Type == BitMatchConditionType {
		return val[condition.Index]&condition.Value == condition.Value
	} else if condition.Type == GreaterThanConditionType {
		return val[condition.Index] > condition.Value
	}

	return false
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	var payload = msg.Payload()
	var state []uint8

	err := json.Unmarshal(payload, &state)
	if err != nil {
		fmt.Printf("Error parsing payload: %v\n", err)
		return
	}

	subscribeTopics[msg.Topic()] = state
	channel <- true
}

func main() {
	var config EngineConfig
	getConfig(&config)

	for _, output := range config.Outputs {
		for _, rule := range output.Rules {
			for _, condition := range rule.Conditions {
				subscribeTopics[condition.Topic] = make([]uint8, 0)
			}
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

	for range channel {
		for _, output := range config.Outputs {
			result := output.DefaultValue
			for _, rule := range output.Rules {
				matchedAll := true
				for _, condition := range rule.Conditions {
					matchedAll = checkCondition(&condition)
					if !matchedAll {
						break
					}
				}

				if matchedAll {
					result = rule.Result
					break
				}
			}

			if len(result) > 0 {
				client.Publish(output.Topic, 0, true, result)
				fmt.Println("Publishing:", result, "to topic:", output.Topic)
			}
		}
	}
}
