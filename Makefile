BINARY := matecito-ai
INSTALL_DIR := $(HOME)/.local/bin
INSTALLED_BIN := $(INSTALL_DIR)/$(BINARY)

.PHONY: dev-install

# Build → copy → install, in one command, for iterating on the CLI without
# publishing a release. Each step must succeed before the next runs (make's
# default: abort on the first failing recipe line) so a build failure never
# reaches install with a stale binary.
dev-install:
	go build -o $(BINARY) ./cmd/matecito-ai
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALLED_BIN)
	$(INSTALLED_BIN) install -y
