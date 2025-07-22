package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ajayykmr/edge_simulator_go/initializers"
	"github.com/ajayykmr/edge_simulator_go/machines"
	"github.com/manifoldco/promptui"
)

var httpRunning, mqttRunning bool
var httpCount, mqttCount int
var httpCancel context.CancelFunc

var sessionStart = time.Now()

func waitForEnter() {
	prompt := promptui.Prompt{
		Label:     "Press Enter to continue...",
		Default:   "",    // So pressing Enter without typing anything works
		AllowEdit: false, // Disables user editing of default
	}

	_, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
	}
}

func main() {

	clearScreen()

	fmt.Println("🚀 Welcome to Edge Simulator!")
	fmt.Println("🛠️  Simulating industrial machine data via HTTP and MQTT")
	fmt.Println("📈 Real-time data. Graceful control. Instant feedback.\n")

	// waitForEnter()

	// Load environment variables
	initializers.LoadEnvVariables()

	//initialize MQTT client
	mqttClient, err := initializers.InitializeMQTTClient()
	if err != nil {
		fmt.Println("❌ Failed to initialize MQTT client: ", err.Error())
		return
	} else {
		defer mqttClient.Disconnect()
	}

	fmt.Println()
	waitForEnter()

	for {
		clearScreen()

		prompt := promptui.Select{
			Label: "Select Action",
			Items: []string{
				menuItem("HTTP", httpRunning, httpCount),
				menuItem("MQTT", mqttRunning, mqttCount),
				"Exit",
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return
		}

		switch {
		case strings.Contains(result, "HTTP"):
			if !httpRunning {
				httpCount = askForNumber("Enter HTTP task count")
				ctx, cancel := context.WithCancel(context.Background())
				httpCancel = cancel
				go machines.SendMachineDataViaHTTP(ctx, httpCount)
				httpRunning = true
			} else {
				if httpCancel != nil {
					httpCancel() // stops the goroutines
				}
				httpRunning = false
			}
		case strings.Contains(result, "MQTT"):
			mqttRunning = !mqttRunning
			if mqttRunning {
				mqttCount = askForNumber("Enter MQTT task count")
			}
		case result == "Exit":
			clearScreen()
			if httpCancel != nil {
				httpCancel() // stops the goroutines
			}
			fmt.Println("🧭 Shutting down simulator...")
			time.Sleep(400 * time.Millisecond)

			// fmt.Println("📡 Disconnecting sensors...")
			// time.Sleep(300 * time.Millisecond)

			// fmt.Println("⚙️  Stopping machines...")
			// time.Sleep(300 * time.Millisecond)

			// fmt.Println("📦 Cleaning up resources...")
			// time.Sleep(300 * time.Millisecond)

			// fmt.Println("✅ All systems idle.")
			// time.Sleep(300 * time.Millisecond)

			duration := time.Since(sessionStart)
			fmt.Printf("🕒 Session duration: %s\n\n", duration.Round(time.Second))
			time.Sleep(300 * time.Millisecond)

			// fmt.Println("👋 Thank you for using the simulator.")
			// fmt.Println("🔒 Exit complete.")
			return
		}
	}
}

func menuItem(name string, running bool, count int) string {
	if running {
		return fmt.Sprintf("Toggle %s [Running] (Current: %d tasks)", name, count)
	}
	return fmt.Sprintf("Toggle %s [Stopped]", name)
}

func askForNumber(label string) int {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			_, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("Please enter a valid number.")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Println("Invalid input.")
		return 0
	}

	num, _ := strconv.Atoi(result)
	return num
}

func clearScreen() {
	// ANSI escape code to clear the screen
	fmt.Print("\033[H\033[2J")

	fmt.Println()
	fmt.Println("╔════════════════════════════╗")
	fmt.Println("║   🌐 Edge Simulator v1.0   ║")
	fmt.Printf("╚════════════════════════════╝\n\n")
}
