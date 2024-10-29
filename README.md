# High-Load UDP Service for Game Clients

## Русская версия

Для русскоязычных пользователей документацию можно найти [здесь](https://github.com/ascenmmo/udp-connection/blob/stage/RU_README.md).

## Description

This project is a high-load **UDP service** designed for connecting game clients. It is developed for free deployment on servers, making it an ideal solution for game developers.

## Key Features

- **High Performance**: Optimized for handling a large number of simultaneous connections.
- **Free Deployment**: Easily deployable on any server with no additional costs.
- **Docker Support**: Available as both a binary and a Docker container, ensuring cross-platform compatibility.
- **Flexibility and Scalability**: Configurable to suit your needs and easily scalable.

## Installation

### Installation via Binary

1. **Download the latest release** from [GitHub Releases](link_to_release).
2. **Choose the appropriate binary** for your operating system:
    - For **Linux**: `udp-linux-amd64`
    - For **Windows**: `udp-windows-amd64.exe`
    - For **macOS**: `udp-darwin-amd64`
3. **Run the binary file**:
   ```bash
   ./udp-linux-amd64
   ```

### Installation via Docker

1. Make sure Docker is installed. If it’s not installed, follow the Docker installation instructions.
2. Run the command:
   ```bash
   docker compose up -d --force-recreate --build && docker image prune -f
	```

## Configuration

The project uses the env/env.go package to define the main configuration parameters. It is important that all services interacting with this UDP server use the same token to ensure a secure connection and authentication.

### Configuration Parameters

```go
package env

var (
	ServerAddress       = "127.0.0.1" // Server address
	TCPPort             = "8081"      // Port for TCP connections
	UDPPort             = "4500"      // Port for UDP connections
	TokenKey            = "_remember_token_must_be_32_bytes" // Unique token for authentication
	MaxRequestPerSecond = 200         // Maximum number of requests per second
)
```

* ServerAddress: Specifies the IP address on which the server will run.
* TCPPort: The port on which the server will listen for TCP connections.
* UDPPort: The port used for handling UDP connections.
* TokenKey: A unique token that must be the same across all services interacting with this UDP server. This ensures the security and integrity of connections.
* MaxRequestPerSecond: A limit on the maximum number of requests the server can handle per second.



##  Importance of a Single Token
### Using a single token allows:

* **Security Assurance:** All services check the token before establishing a connection, helping to prevent unauthorized access.
* **Simplified Authentication:** A single token for all services simplifies the authentication and access management process.
* **Ease of Maintenance:** If the token needs to be changed, it can be done in one place, and all services will be updated simultaneously.

Make sure that all your services are configured to use this token to ensure the correct operation and security of the system.



## Troubleshooting

If you encounter issues when starting the service, check the following:

- Ensure that ports 8081 and 4500 are not occupied by other applications.
- Verify that the configuration parameters in env/env.go are correctly specified.
- If you are using Docker, ensure that it is running and properly configured.






## Теги

`UDP`, `игровой сервер`, `высоконагруженный`, `бесплатное развертывание`, `Docker`, `кроссплатформенный`, `игровая разработка`, `сеть`, `многопользовательская игра`, `сервис для игр`, `настройка сервера`, `аутентификация`, `токены`, `Golang`, `open-source`

## Tags

`UDP`, `game server`, `high-performance`, `free deployment`, `Docker`, `cross-platform`, `game development`, `network`, `multiplayer game`, `game service`, `server setup`, `authentication`, `tokens`, `Golang`, `open-source`
