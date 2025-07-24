package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ajayykmr/edge_simulator_go/initializers"
	"github.com/ajayykmr/edge_simulator_go/machines"
	"github.com/manifoldco/promptui"
)

var sessionStart = time.Now()

var mqttServiceActive, httpServiceActive bool = false, false
var httpRunning, mqttRunning bool = false, false
var httpCancel, mqttCancel context.CancelFunc
var httpCount, mqttCount int

func main() {
	for {
		restart := runApp()
		if !restart {
			break
		}
	}
}

// returns true if the app should restart
func runApp() bool {

	clearScreen()

	fmt.Println("🚀 Welcome to Edge Simulator!")
	fmt.Println("🛠️  Simulating industrial machine data via HTTP and MQTT")
	fmt.Printf("📈 Real-time data. Graceful control. Instant feedback.\n\n")

	// Load environment variables
	initializers.LoadEnvVariables()

	//initialize MQTT client
	mqttClient, err := initializers.InitializeMQTTClient()
	if err != nil {
		mqttServiceActive = false
	} else {
		mqttServiceActive = true

		//if you comment this out, make sure to disconnect the MQTT client when done
		defer mqttClient.Disconnect()
	}

	//check if HTTP server is running
	err = machines.CheckHTTPServerStatus()
	httpServiceActive = err == nil

	if !mqttServiceActive {
		//no mqtt support
		fmt.Println("❌ MQTT service not active. Skipping MQTT data publishing.")
		fmt.Println()
	} else {
		fmt.Println("✅ MQTT service running — publishing enabled.")
		fmt.Println()
	}

	if !httpServiceActive {
		fmt.Println("❌ HTTP service not active. Skipping HTTP data transmission.")
	} else {
		fmt.Println("✅ HTTP service running — transmission enabled.")
	}

	if !mqttServiceActive && !httpServiceActive {
		fmt.Printf("\n⚠️  No active services found. Start at least one service to begin simulation.\n\n")

		// printExitMessages()
		// return false
	}

	fmt.Println()
	waitForEnter()

	for {
		clearScreen()

		menuItems := []string{}
		if httpServiceActive {
			menuItems = append(menuItems, menuItem("HTTP", httpRunning, httpCount))
		}
		if mqttServiceActive {
			menuItems = append(menuItems, menuItem("MQTT", mqttRunning, mqttCount))
		}
		menuItems = append(menuItems, "Restart", "Exit")

		prompt := promptui.Select{
			Label: "Select Action",
			Items: menuItems,
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return false
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
			if !mqttRunning {
				mqttCount = askForNumber("Enter MQTT task count")
				ctx, cancel := context.WithCancel(context.Background())
				mqttCancel = cancel

				go machines.SendMachineDataViaMQTT(ctx, mqttClient, mqttCount)
				mqttRunning = true
			} else {
				if mqttCancel != nil {
					mqttCancel() // stops the goroutines
				}
				mqttRunning = false
			}
		case result == "Restart":
			if httpCancel != nil {
				httpCancel()
			}
			if mqttCancel != nil {
				mqttCancel()
			}
			return true // tells main() to restart
		case result == "Exit":
			clearScreen()
			if httpCancel != nil {
				httpCancel() // stops the goroutines
			}
			if mqttCancel != nil {
				mqttCancel() // stops the goroutines
			}
			printExitMessages()
			return false // tells main() to exit
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
				return fmt.Errorf("please enter a valid number")
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
	// fmt.Print("\033[H\033[2J")

	time.Sleep(100 * time.Millisecond) // slight delay for better UX
	fmt.Print("\033[2J\033[H\033[3J")
	fmt.Println()
	fmt.Println("╔════════════════════════════╗")
	fmt.Println("║   🌐 Edge Simulator v1.0   ║")
	fmt.Printf("╚════════════════════════════╝\n\n")
}

func printExitMessages() {
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
}

func waitForEnter() {
	fmt.Print("👉 Press Enter to continue...")
	_, _ = bufio.NewReader(os.Stdin).ReadString('\n')
}
