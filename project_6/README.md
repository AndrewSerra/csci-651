# Reliable TCP FTP Server and Client

This is a simple implementation of a Reliable TCP FTP Server and Client in Go.

## Prerequisites

Before running or building this program, ensure you have the following installed:
- Go ([installation instructions](https://golang.org/doc/install))

## How to Build

To build the program, follow these steps:
1. Clone the repository to your local machine.
2. Navigate to the project directory.
3. Run the following command:
    ```sh
    go build
    ```

## How to Run

### Server

To run the server, execute the following command:
```sh
go run main.go -port 2001 -dst localhost:2000 -s
```

### Client

To run the client, execute the following command:
```sh
go run main.go -port 2000 -dst localhost:2001 -s
```

# Note

The implementation is complete and does not fully work. 
