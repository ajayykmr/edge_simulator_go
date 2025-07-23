package machines

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/ajayykmr/edge_simulator_go/data_generator"
	"github.com/ajayykmr/edge_simulator_go/mqtt"
	"github.com/ajayykmr/edge_simulator_go/utils"
)

func startSingleMQTTMachine(ctx context.Context, mqttClient mqtt.Client, machineID string, interval time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			// log.Printf("Machine %s stopped", machineID)
			return
		default:
			data := data_generator.GenerateCNCData(machineID)

			jsonData, err := json.Marshal(data)
			if err != nil {
				// log.Printf("❌ Failed to marshal data for %s: %v", machineID, err)
				continue
			}

			//publish data to MQTT
			err = mqttClient.Publish("factory/pune/cnc/"+machineID+"/data", 0, false, jsonData)
			if err != nil {
				// log.Printf("Error publishing data for machine %s: %v", machineID, err)
				continue
			} else {
				// log.Printf("Published data for machine %s", machineID)
			}
			time.Sleep(interval)
		}
	}
}

func SendMachineDataViaMQTT(ctx context.Context, mqttClient mqtt.Client, count int) {
	wg := &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		machineID := "CNC-MQTT-" + strconv.Itoa(i+1)
		interval := time.Millisecond * time.Duration(utils.RandInt(500, 2000))

		wg.Add(1)
		go startSingleMQTTMachine(ctx, mqttClient, machineID, interval, wg)
	}

	// Run a separate goroutine to wait and log when done
	go func() {
		wg.Wait()
		// log.Println("Stopped sending data via HTTP")
	}()
}
