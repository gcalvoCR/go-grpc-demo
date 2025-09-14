// Command client is a CLI for interacting with the MemeService gRPC server.
// It supports four subcommands:
//   - random: fetch a single random meme (optionally filtered by category)
//   - list:   fetch a list of memes for a category
//   - stream: stream memes from the server until the server stops sending
//   - upload: upload memes either from a JSON file or a single meme via flags
//
// See usage() for detailed examples.
package main

import (
	// Standard library packages
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	// gRPC core and credentials
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Generated protobuf/gRPC client for the meme service
	memespb "github.com/gcalvocr/go-grpc-demo/memespb"
)

// Upload represents a meme entry read from a JSON file for bulk uploads.
// This mirrors the fields expected by the server's MemeUpload message.
type Upload struct {
	Title    string `json:"title"`
	Url      string `json:"url"`
	Category string `json:"category"`
}

// usage prints a short help message describing the available subcommands and flags.
func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
  client -addr host:port random [-category cat]
  client -addr host:port list   [-category cat]
  client -addr host:port stream [-category cat]
  client -addr host:port upload [-file memes.json | -title t -url u -category c]

Flags:
  -addr       server address (default "localhost:50051")
  -category   optional category filter for random/list/stream
  -file       JSON file with an array of {title,url,category} objects for upload
  -title      title for single upload
  -url        url for single upload
`)
}

func main() {
	// Define and parse command-line flags.
	addr := flag.String("addr", "localhost:50051", "gRPC server address")
	category := flag.String("category", "", "category filter")
	file := flag.String("file", "", "path to JSON file for bulk upload")
	title := flag.String("title", "", "title for single upload")
	url := flag.String("url", "", "url for single upload")
	flag.Parse()

	// Require at least one positional argument indicating the subcommand.
	if flag.NArg() < 1 {
		usage()
		os.Exit(2)
	}
	cmd := flag.Arg(0)

	// Establish a client connection to the gRPC server using insecure transport
	// credentials (appropriate only for local or test environments).
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// Create a typed client for the MemeService defined in the protobufs.
	client := memespb.NewMemeServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dispatch on the requested subcommand.
	switch cmd {
	case "random":
		// Unary RPC: request a single random meme, optionally filtered by category.
		resp, err := client.GetRandomMeme(ctx, &memespb.MemeRequest{Category: *category})
		if err != nil {
			log.Fatalf("GetRandomMeme: %v", err)
		}
		fmt.Printf("Random meme: id=%s title=%q url=%s category=%q\n", resp.GetId(), resp.GetTitle(), resp.GetUrl(), resp.GetCategory())

	case "list":
		// Unary RPC: request a list of memes for the given category.
		resp, err := client.GetMemesByCategory(ctx, &memespb.CategoryRequest{Category: *category})
		if err != nil {
			log.Fatalf("GetMemesByCategory: %v", err)
		}
		for i, m := range resp.GetMemes() {
			fmt.Printf("%d) id=%s title=%q url=%s category=%q\n", i+1, m.GetId(), m.GetTitle(), m.GetUrl(), m.GetCategory())
		}

	case "stream":
		// Server-streaming RPC: receive a stream of memes from the server.
		// The loop exits when the server closes the stream or an error occurs.
		stream, err := client.StreamMemes(ctx, &memespb.StreamRequest{Category: *category})
		if err != nil {
			log.Fatalf("StreamMemes: %v", err)
		}
		count := 0
		for {
			m, err := stream.Recv()
			if err != nil {
				// Typically you'd check for io.EOF; here any error ends the loop.
				break
			}
			count++
			fmt.Printf("%d) id=%s title=%q url=%s category=%q\n", count, m.GetId(), m.GetTitle(), m.GetUrl(), m.GetCategory())
		}
		fmt.Printf("Stream completed after %d items\n", count)

	case "upload":
		// Client-streaming RPC: send one or more memes to the server and then
		// receive a summary response after closing the stream.
		st, err := client.UploadMemes(ctx)
		if err != nil {
			log.Fatalf("UploadMemes open: %v", err)
		}

		// Build the list of memes to upload, either from a JSON file or from
		// the -title/-url flags.
		var toUpload []Upload
		if *file != "" {
			data, err := os.ReadFile(*file)
			if err != nil {
				log.Fatalf("read file: %v", err)
			}
			if err := json.Unmarshal(data, &toUpload); err != nil {
				log.Fatalf("parse json: %v", err)
			}
		} else if *title != "" && *url != "" {
			toUpload = []Upload{{Title: *title, Url: *url, Category: *category}}
		} else {
			usage()
			log.Fatalf("upload requires -file or -title and -url")
		}

		// Send each meme to the server over the client stream.
		for _, u := range toUpload {
			if err := st.Send(&memespb.MemeUpload{Title: u.Title, Url: u.Url, Category: u.Category}); err != nil {
				log.Fatalf("send: %v", err)
			}
		}

		// Close the send direction and receive the summary from the server.
		summary, err := st.CloseAndRecv()
		if err != nil {
			log.Fatalf("close: %v", err)
		}
		fmt.Printf("Uploaded %d memes. Message: %s\n", summary.GetCount(), summary.GetMessage())

	default:
		// Unknown subcommand: print usage and exit with an error code.
		usage()
		os.Exit(2)
	}
}
