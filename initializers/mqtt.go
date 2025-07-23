package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/ajayykmr/edge_simulator_go/mqtt"
)

// initialize MQTT client and subscribe to topics
func InitializeMQTTClient() (mqtt.Client, error) {
	//initialize MQTT client
	mqttServer := os.Getenv("MQTT_BROKER")
	if mqttServer == "" {
		return mqtt.Client{}, fmt.Errorf("MQTT_BROKER environment variable is not set")
	}

	clientID := os.Getenv("MQTT_CLIENT_ID")
	if clientID == "" {
		return mqtt.Client{}, fmt.Errorf("MQTT_CLIENT_ID environment variable is not set")
	}

	mqttClient, err := mqtt.InitializeClient(mqttServer, clientID)
	if err != nil {
		return mqtt.Client{}, fmt.Errorf("failed to initialize MQTT client: %v", err)
	}

	testTopic := os.Getenv("MQTT_TEST_TOPIC")
	testTopicMessage := os.Getenv("MQTT_TEST_MESSAGE")
	if testTopic != "" && testTopicMessage != "" {
		err = mqttClient.Publish(testTopic, 0, false, testTopicMessage)
		if err != nil {
			fmt.Printf("❌ Failed to publish test message: %v", err)
		} else {
			fmt.Println("✅ Test message published to MQTT broker on topic:", testTopic)
		}
	} else {
		log.Println("MQTT_TEST_TOPIC or MQTT_TEST_MESSAGE environment variable not set, skipping test message publication")
	}

	return *mqttClient, nil
}
