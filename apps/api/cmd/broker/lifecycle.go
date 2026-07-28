package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/franwerner/matecito-ai/apps/api/internal/config"
	"github.com/franwerner/matecito-ai/apps/api/internal/store"
	"github.com/franwerner/matecito-ai/apps/api/internal/transport"
)

const shutdownTimeout = 5 * time.Second

// run wires the broker's lifecycle in the order runtime/error-handling and
// contracts/mcp-server require: config -> logger -> store (open + migrate)
// -> MCP availability check -> transport, with clean shutdown on ctx
// cancellation. args/getenv/stdout are injected so main() owns the only
// process-global wiring (os.Args, os.Getenv, real OS signals) and this
// function stays testable.
func run(ctx context.Context, args []string, getenv func(string) string, stdout io.Writer) int {
	cfg, err := config.Load(args, getenv)
	if err != nil {
		_, _ = fmt.Fprintf(stdout, "broker: invalid configuration: %v\n", err)
		return 1
	}

	logger := newLogger(cfg, stdout)

	machineID, err := store.EnsureMachineID(cfg.PersistDir)
	if err != nil {
		logger.Error("machine id setup failed", "err", err)
		return 1
	}

	// R1/R14: the daemon must not serve with an unmigrated or partially
	// migrated schema, so opening (and migrating) the store happens before
	// anything else can bind or check readiness.
	st, err := store.Open(ctx, cfg.PersistDir, machineID)
	if err != nil {
		logger.Error("store open/migrate failed", "err", err)
		return 1
	}
	defer func() {
		if err := st.Close(); err != nil {
			logger.Error("store close error", "err", err)
		}
	}()
	logger.Info("store opened and migrated")

	// contracts/mcp-server: the daemon verifies MCP availability before it
	// ever opens the transport boundary (Stub B in Phase 1 — always ok).
	if err := transport.CheckMCPAvailability(ctx); err != nil {
		logger.Error("mcp availability check failed", "err", err)
		return 1
	}

	// Stub C: the health status is a local, on-demand check against the
	// already-open store (observability/health-checks) — no external
	// dependency, and no separate readiness state to track here.
	storeStatus := func() string { return string(st.StoreStatus(ctx)) }
	srv := transport.New(fmt.Sprintf(":%d", cfg.Port), storeStatus)
	if err := srv.Listen(); err != nil {
		logger.Error("broker failed to bind", "err", err)
		return 1
	}
	logger.Info("broker listening", "addr", srv.Addr())

	serveErrCh := make(chan error, 1)
	go func() { serveErrCh <- srv.Serve() }()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case err := <-serveErrCh:
		if err != nil {
			logger.Error("broker exited unexpectedly", "err", err)
			return 1
		}
		return 0
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", "err", err)
		return 1
	}
	<-serveErrCh // wait for Serve() to return so resources are released before exit
	logger.Info("broker stopped cleanly")
	return 0
}
