GOMOD=go mod
CC=go build
# CFLAGS=
LDFLAGS=""

SRC=./cmd/recon/
TARGET_DIR=./bin
TARGET=$(TARGET_DIR)/recon

.PHONY: build clean tidy

build: $(TARGET)

$(TARGET):
	$(CC) $(CFLAGS) -ldflags $(LDFLAGS) -o $@ $(SRC)

tidy:
	$(GOMOD) tidy

clean:
	rm -rf $(TARGET_DIR)
