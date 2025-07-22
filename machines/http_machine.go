package machines

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ajayykmr/edge_simulator_go/data_generator"
	"github.com/ajayykmr/edge_simulator_go/utils"
)

func startSingleHTTPMachine(ctx context.Context, machineID string, interval time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	endpoint := "http://localhost:8080/ingest"
	for {
		select {
		case <-ctx.Done():
			// log.Printf("Machine %s stopped", machineID)
			return
		default:
			data := data_generator.GenerateCNCData(machineID)

			payload, err := json.Marshal(data)
			if err != nil {
				// log.Printf("Error marshaling data for machine %s: %v", machineID, err)
				continue
			}

			resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				// log.Println("Error sending data via HTTP for machine", machineID, ":", err)
				continue
			}
			resp.Body.Close()

			time.Sleep(interval)
		}
	}
}

func SendMachineDataViaHTTP(ctx context.Context, count int) {
	wg := &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		machineID := "CNC-" + strconv.Itoa(i+1)
		interval := time.Millisecond * time.Duration(utils.RandInt(500, 2000))

		wg.Add(1)
		go startSingleHTTPMachine(ctx, machineID, interval, wg)
	}

	// Run a separate goroutine to wait and log when done
	go func() {
		wg.Wait()
		// log.Println("Stopped sending data via HTTP")
	}()
}
