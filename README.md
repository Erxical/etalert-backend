# etalert-backend

etalert-backend is a backend service for the Etalert application, built with Go. It provides real-time alerts via WebSockets and is designed to integrate easily with various services. This repository contains the backend logic, including handlers, services, and middlewares for managing alerts.

## Features
- Real-time alerts using WebSockets
- Efficient data handling with middlewares
- Scalable and extensible architecture

## Requirements
- Go 1.18+
- Docker (optional, for containerization)

## Installation

1. Clone this repository:
    ```bash
    git clone https://github.com/Erxical/etalert-backend.git
    cd etalert-backend
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```
    
3. Create environment file:
   ```env
    # Google Map
    G_MAP_API_KEY=<PUT APIKEY HERE>

    # Gemini
    GEMINI_API_KEY=<PUT APIKEY HERE>

    # Azure Map
    AZURE_MAP_API_KEY=<PUT APIKEY HERE>

    # MongoDB
    MONGODB_URI=<PUT APIKEY HERE>

    # JWT
    JWT_SECRET=<PUT APIKEY HERE>
    ```
   
4. Run the application:
    ```bash
    go run main.go
    ```

5. To build a Docker image:
    ```bash
    docker build -t etalert-backend ./Dockerfile .
    ```

## Directory Structure
- `handler/`: Contains request handlers
- `middlewares/`: Middleware components for request processing
- `repository/`: Data storage layer and database interactions
- `service/`: Business logic of the application
- `validators/`: Input validation logic
- `websocket/`: WebSocket-related implementations

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing
Feel free to open an issue or submit a pull request. Contributions are always welcome!

