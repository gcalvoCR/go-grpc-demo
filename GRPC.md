# gRPC in Simple Words

gRPC is a fast and efficient way for programs (services) to talk to each other, even if they’re written in different languages. You write a contract in a .proto file (what requests and responses look like), then tools generate code so both server and client can communicate without you writing boilerplate. Under the hood, gRPC uses HTTP/2 and Protocol Buffers (compact binary format).

---

## Advantages
- High performance: compact binary messages over HTTP/2.
- Strongly typed contract: schema-first design (fewer integration surprises).
- Code generation: stubs for client/server across many languages.
- Streaming support: unary, server-streaming, client-streaming, and bidirectional.
- Cross-language and cross-platform.
- Production-ready features: deadlines, cancellation, interceptors, auth, load balancing.

---

## Implementation Order (as used in this repo)

1) Define the proto
- File: `proto/memes.proto`
- This is the single source of truth for request/response shapes and service methods.

2) Generate the code using protoc
- Generated Go files live in:
  - `memespb/memes.pb.go`
  - `memespb/memes_grpc.pb.go`
- Re-run generation whenever you change `proto/memes.proto`.

3) Implement the server
- Files:
  - `cmd/server/main.go` (server bootstrap and listener)
  - `cmd/server/meme_server.go` (service implementation)

4) Implement the client
- File: `cmd/client/main.go` (simple client calling the RPCs defined in the proto)

---

## How to Run
- Start the server: go run ./cmd/server
- Run the client: go run ./cmd/client

If you use tools like grpcurl, point them at the server on localhost:50051 and the services defined in `proto/memes.proto`.

---

## Project Pointers
- Proto: `proto/memes.proto`
- Generated: `memespb/memes.pb.go`, `memespb/memes_grpc.pb.go`
- Server: `cmd/server/main.go`, `cmd/server/meme_server.go`
- Client: `cmd/client/main.go`
- Additional docs: `readme.md`

---

## Summary
- Define the proto (contract) in `proto/memes.proto`
- Generate code (stubs/types) into `memespb/`
- Implement the server under `cmd/server/`
- Implement the client under `cmd/client/`
