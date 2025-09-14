# Go gRPC Demo: MemeService

A simple gRPC demo showcasing a MemeService with examples of unary, server-streaming, and client-streaming RPCs. This repository includes:
- Server: cmd/server
- Client: cmd/client
- Generated protobuf code: memespb

## Getting Started
Quickest way to run the demo locally.

1) Generate code from the proto
```sh
make compile
```

2) Start the server (in one terminal)
```sh
make run-server
```

3) Run a client command (in another terminal)
- Random meme:
  ```sh
  make client-random
  ```
- List by category:
  ```sh
  make client-list CATEGORY=classic
  ```
- Stream by category:
  ```sh
  make client-stream CATEGORY=pokemon
  ```
- Upload from JSON (uses json/memes.json):
  ```sh
  make client-upload-file
  ```

For more options and flags, see the sections below.

## Prerequisites
- Go 1.21+ installed
- Protocol Buffers compiler (protoc)
  - macOS: `brew install protobuf`
- Protobuf Go plugins on your PATH:
  ```sh
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
  Make sure `$GOPATH/bin` (or the Go install bin directory) is on your PATH.

## Generate gRPC code from .proto
Use the provided Makefile target:
```sh
make compile
```
This generates the Go code under memespb.

## Run the server
The server listens on :50051 and has server reflection enabled.
```sh
go run ./cmd/server
```
You should see:
```
gRPC server listening on :50051
```

## Run the client
The client connects to a gRPC server and supports four subcommands: random, list, stream, upload.

General usage:
```
client -addr host:port random [-category cat]
client -addr host:port list   [-category cat]
client -addr host:port stream [-category cat]
client -addr host:port upload [-file memes.json | -title t -url u -category c]
```

Run via `go run` (examples assume the server is on localhost:50051):

- Random meme (optionally filter by category):
  ```sh
  go run ./cmd/client random
  go run ./cmd/client random -category classic
  ```

- List memes by category:
  ```sh
  go run ./cmd/client list -category classic
  ```

- Stream memes by category (server-streaming):
  ```sh
  go run ./cmd/client stream -category pokemon
  ```

- Upload memes (client-streaming) from a JSON file:
  You can use the included example at json/memes.json:
  ```sh
  go run ./cmd/client upload -file json/memes.json
  ```
  Or create your own file, e.g. memes.json:
  ```json
  [
    {"title": "Distracted Boyfriend", "url": "https://example.com/distracted.jpg", "category": "classic"},
    {"title": "Surprised Pikachu",  "url": "https://example.com/pikachu.jpg",    "category": "pokemon"}
  ]
  ```
  Then run:
  ```sh
  go run ./cmd/client upload -file memes.json
  ```

- Upload a single meme via flags:
  ```sh
  go run ./cmd/client upload -title "New Meme" -url "https://example.com/new.jpg" -category misc
  ```

## Makefile shortcuts
The Makefile includes convenient targets. The default address is ADDR=localhost:50051.

- Start the server:
  ```sh
  make run-server
  ```

- Random meme (optional category):
  ```sh
  make client-random
  make client-random CATEGORY=classic
  ```

- List memes by category:
  ```sh
  make client-list CATEGORY=classic
  ```

- Stream memes by category:
  ```sh
  make client-stream CATEGORY=pokemon
  ```

- Upload from JSON file (defaults to json/memes.json):
  ```sh
  make client-upload-file
  make client-upload-file FILE_JSON=path/to/your.json
  ```

- Upload a single meme via flags:
  ```sh
  make client-upload-single TITLE='New Meme' URL='https://example.com/new.jpg' CATEGORY=misc
  ```

- Generic client runner (specify CMD and optional CATEGORY):
  ```sh
  make run-client CMD=random
  make run-client CMD=list CATEGORY=classic
  ```

## Using grpcurl (optional)
With server reflection enabled, you can explore the API:
```sh
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext localhost:50051 describe
```

## Troubleshooting
- Port in use: ensure nothing else is listening on 50051 before starting the server.
- Client can’t connect: verify the address with `-addr` matches your server (default is localhost:50051).

## Development workflow
1. Edit the .proto definitions.
2. Regenerate code: `make compile`.
3. Run the server: `go run ./cmd/server`.
4. Run the client commands shown above to test changes.
