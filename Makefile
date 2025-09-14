# Proto source and output configuration
FILE ?= proto/memes.proto
OUT_DIR ?= memespb
SRC_DIR := $(dir $(FILE))
FILE_BASENAME := $(notdir $(FILE))

# Client configuration
ADDR ?= localhost:50051
CATEGORY ?=
TITLE ?=
URL ?=
FILE_JSON ?= json/memes.json
CMD ?= random

# Derived flag helpers (only set when variables are non-empty)
CATEGORY_FLAG :=
ifneq ($(strip $(CATEGORY)),)
  CATEGORY_FLAG := -category $(CATEGORY)
endif

TITLE_FLAG :=
ifneq ($(strip $(TITLE)),)
  TITLE_FLAG := -title $(TITLE)
endif

URL_FLAG :=
ifneq ($(strip $(URL)),)
  URL_FLAG := -url $(URL)
endif

.PHONY: compile run-server run-client client-random client-list client-stream client-upload-file client-upload-single

# Generate gRPC code from the proto definition
compile:
	mkdir -p $(OUT_DIR)
	protoc -I $(SRC_DIR) --go_out=$(OUT_DIR) --go_opt=paths=source_relative --go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative $(FILE_BASENAME)

# Run the gRPC server
run-server:
	go run ./cmd/server

# Generic client runner (use CMD to choose the subcommand)
run-client:
	go run ./cmd/client -addr $(ADDR) $(CATEGORY_FLAG) $(CMD)

# Convenience targets for each client option
client-random:
	go run ./cmd/client -addr $(ADDR) $(CATEGORY_FLAG) random

client-list:
	go run ./cmd/client -addr $(ADDR) $(CATEGORY_FLAG) list

client-stream:
	go run ./cmd/client -addr $(ADDR) $(CATEGORY_FLAG) stream

client-upload-file:
	go run ./cmd/client -addr $(ADDR) -file $(FILE_JSON) upload

client-upload-single:
	@if [ -z "$(strip $(TITLE))" ] || [ -z "$(strip $(URL))" ]; then \
		echo "ERROR: Please provide TITLE and URL. Example:"; \
		echo "  make client-upload-single TITLE='New Meme' URL='https://example.com/new.jpg' CATEGORY='misc'"; \
		exit 2; \
	fi
	go run ./cmd/client -addr $(ADDR) $(TITLE_FLAG) $(URL_FLAG) $(CATEGORY_FLAG) upload

