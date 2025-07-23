# Edge Simulator

The Edge Simulator is a Go-based CLI application designed to simulate industrial machine data, specifically CNC (Computer Numerical Control) machine data, and transmit it via HTTP and MQTT protocols. It provides a real-time simulation environment with graceful control and instant feedback, making it ideal for testing data pipelines, IoT platforms, and edge computing solutions.

## ✨ Features

  * **CNC Data Simulation:** Generates realistic CNC machine data.
  * **Multi-Protocol Transmission:** Transmits simulated data via HTTP and MQTT protocols.
  * **Kafka Integration:** Supports producing data to Kafka topics.
  * **Interactive CLI:** Provides a command-line interface to manage HTTP and MQTT simulations.
  * **Configurable:** Utilizes environment variables for easy setup.
  * **Service Monitoring:** Includes health checks for active services.

## 🛠️ Installation

To get started with the Edge Simulator, ensure you have Go installed (Go 1.24.5 or higher is required).

1.  **Clone the Repository:**

    ```bash
    git clone https://github.com/ajayykmr/edge_simulator_go.git
    cd edge_simulator_go
    ```

2.  **Download Dependencies:**

    ```bash
    go mod download
    ```

3.  **Build the Application (Optional):**

    ```bash
    go build -o edge-simulator
    ```

    This will create an executable named `edge-simulator` in your current directory.

## 🚀 Usage

The Edge Simulator is run as a CLI application. It requires an HTTP server (e.g., running on `localhost:8080` with an `/ingest` endpoint and `/health` endpoint) and optionally an MQTT broker.

1.  **Environment Configuration:**
    Create a `.env` file in the root directory of the project to configure your MQTT and other settings. The `.gitignore` file includes `.env`, so it won't be committed to your repository.

    Example `.env` file:

    ```
    MQTT_BROKER=tcp://localhost:1883
    MQTT_CLIENT_ID=edge-simulator-client
    MQTT_TEST_TOPIC=test/topic
    MQTT_TEST_MESSAGE=HelloFromEdgeSimulator
    ENV=development
    ```

      * `MQTT_BROKER`: The address of your MQTT broker (e.g., `tcp://localhost:1883`).
      * `MQTT_CLIENT_ID`: A unique client ID for the MQTT connection.
      * `MQTT_TEST_TOPIC`: An optional topic to publish a test message upon MQTT client initialization.
      * `MQTT_TEST_MESSAGE`: The message content for the test topic.
      * `ENV`: Set to `production` to skip loading the `.env` file.

2.  **Run the Simulator:**

    ```bash
    go run .
    ```

    If you built the executable:

    ```bash
    ./edge-simulator
    ```

3.  **Interactive Menu:**
    Upon running, the application will check the status of HTTP and MQTT services and then present an interactive menu:

    ```
    🚀 Welcome to Edge Simulator!
    🛠️  Simulating industrial machine data via HTTP and MQTT
    📈 Real-time data. Graceful control. Instant feedback.

    ✅ MQTT service running — publishing enabled.

    ✅ HTTP service running — transmission enabled.

    👉 Press Enter to continue...

    ╔════════════════════════════╗
    ║   🌐 Edge Simulator v1.0   ║
    ╚════════════════════════════╝

    Select Action:
    > Toggle HTTP [Stopped]
      Toggle MQTT [Stopped]
      Exit
    ```

      * Use the arrow keys to navigate and Enter to select an option.
      * **Toggle HTTP:** Starts or stops HTTP data transmission. If starting, you'll be prompted to "Enter HTTP task count" to specify how many concurrent HTTP machine simulations to run.
      * **Toggle MQTT:** Starts or stops MQTT data publishing. If starting, you'll be prompted to "Enter MQTT task count" to specify how many concurrent MQTT machine simulations to run.
      * **Exit:** Shuts down the simulator gracefully.

## ⚙️ Project Structure

```
.
├── .gitignore             # Specifies intentionally untracked files to ignore
├── data_generator
│   └── cnc.go             # Generates CNC machine data
├── go.mod                 # Go module file with project dependencies
├── http_handlers
│   └── health.go          # HTTP health check handler
├── initializers
│   ├── load_env_variables.go # Loads environment variables from .env
│   └── mqtt.go            # Initializes MQTT client
├── kafka
│   ├── consumer
│   │   └── consumer.go    # Kafka consumer implementation
│   └── producer
│       └── producer.go    # Kafka producer implementation
├── machines
│   ├── http_machine.go    # Simulates data transmission via HTTP
│   └── mqtt_machine.go    # Simulates data publishing via MQTT
├── main.go                # Main application entry point and CLI logic
└── utils
    └── utils.go           # Utility functions (e.g., random number generation)
```

## 🤝 Contributing

We welcome contributions to the Edge Simulator\! If you have suggestions for improvements, bug fixes, or new features, please feel free to contribute.

1.  **Reporting Issues:** If you find a bug or have a feature request, please open an issue on the GitHub repository.
2.  **Making Changes:**
      * Fork the repository.
      * Create a new branch for your changes.
      * Make your modifications and ensure they are well-tested.
      * Submit a pull request with a clear description of your changes.

We appreciate your help in making this simulator better\!

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/search?q=LICENSE) file for details.
