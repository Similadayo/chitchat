# ChitChat

ChitChat is a real-time chat application built with Go. It leverages WebSockets for real-time communication and uses Gorilla Mux for routing. The application also integrates with MySQL for data storage and JWT for authentication.

## Features

- Real-time messaging
- User authentication with JWT
- Persistent storage with MySQL
- Environment configuration with `godotenv`
- Logging with Logrus

## Prerequisites

- Go 1.21.3 or higher
- MySQL database

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/similadayo/chitchat.git
    cd chitchat
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Set up environment variables:
    Create a `.env` file in the root directory and add your environment variables (e.g., database credentials).

4. Run the application:

    ```sh
    go run main.go
    ```

## Usage

1. Open your browser and navigate to `http://localhost:8080`.
2. Register or log in to start chatting in real-time.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
