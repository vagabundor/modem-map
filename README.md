# Satellite Network Modem Map

Satellite Network Modem Map is a Golang application that displays satellite modems on a geo map. The information about the modems is obtained from the NMS database. The application supports connecting to multiple NMS databases if multiple satellite hubs are present.

![Map](./docs/graphics/map.png)

## Features

- Display satellite modems on a geo map
- Obtain information from the NMS database
- Connect to multiple NMS databases

## Roadmap

- Add SNMP support to obtain information from NMS
- Obtain information from iVantage API
- Display status and active information of modems (temperature, signal level, etc.)

## Getting Started

### Prerequisites

- Golang
- MySQL database

### Installation

1. Clone the repository: `git clone https://github.com/your_username/satellite-network-modem-map.git`
2. Navigate to the cloned directory: `cd satellite-network-modem-map`
3. Install dependencies: `go mod tidy`
4. Copy the `config.example.toml` file to `config.toml` and update the database configuration with your own settings
5. Start the application: `go run main.go`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
