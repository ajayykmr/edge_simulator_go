package mqtt

import (
	"fmt"
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client MQTT.Client
}

var defaultHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Printf("TOPIC: %s, ", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func InitializeClient(server string, client_id string) (*Client, error) {

	opts := MQTT.NewClientOptions().AddBroker(server)
	opts.SetClientID(client_id)
	opts.SetDefaultPublishHandler(defaultHandler)

	mqttClient := MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{client: mqttClient}, nil
}

func (c *Client) Subscribe(topic string, qos byte, handler MQTT.MessageHandler) error {
	token := c.client.Subscribe(topic, qos, handler)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("%v", token.Error())
	}
	return nil
}

// Disconnect disconnects the client
func (c *Client) Disconnect() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
		// fmt.Println("🔴 MQTT client disconnected successfully\n")
	}
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload any) error {
	if token := c.client.Publish(topic, qos, retained, payload); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %w", token.Error())
	}
	return nil
}
