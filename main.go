package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	handler "github.com/ajayykmr/edge_simulator_go/http_handlers"
	"github.com/ajayykmr/edge_simulator_go/initializers"
	"github.com/ajayykmr/edge_simulator_go/simulator"
	"github.com/gin-gonic/gin"
)

func SendCNCToGateway(data simulator.CNCData, endpoint string) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Posted to %s | Status: %s", endpoint, resp.Status)
	return nil
}
func main() {

	// Load environment variables
	initializers.LoadEnvVariables()

	//initialize MQTT client
	mqttClient, err := initializers.InitializeMQTTClient()
	if err != nil {
		log.Println("Failed to initialize MQTT client: ", err.Error())
		return
	} else {
		defer mqttClient.Disconnect()
	}

	//
	endpoint := "http://localhost:8080/ingest" // Change this if hosted elsewhere
	machineID := "CNC-001"

	for {
		data := simulator.GenerateCNCData(machineID)

		err := SendToGateway(data, endpoint)
		if err != nil {
			log.Println("Failed to send data:", err)
		}

		time.Sleep(time.Millisecond * 500) // adjustable interval
	}

	//Gin router
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT environment variable not set, using default port 8080")
		port = "8080"
	}
	gin.SetMode(gin.ReleaseMode) // Set Gin to release mode (also stops debug messages)
	r := gin.Default()
	r.GET("/health", handler.HealthHandler)

	log.Println("Starting server on port: " + port)

	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}

}
