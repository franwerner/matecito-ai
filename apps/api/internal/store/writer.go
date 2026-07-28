package store

import (
	"context"
	"fmt"

	"github.com/franwerner/matecito-ai/apps/api/internal/store/ent"
)

// txFunc is one atomic write command: whatever it does against tx either
// all lands or all rolls back together (R12 "escritura compuesta atómica").
type txFunc func(ctx context.Context, tx *ent.Tx) (any, error)

type writeCmd struct {
	ctx  context.Context
	fn   txFunc
	resp chan writeResult
}

type writeResult struct {
	val any
	err error
}

// writer is the single goroutine that owns every write to the persistence,
// fed by a channel (runtime/concurrency-async: "único goroutine escritor,
// alimentado por channel"). Writer methods on Engine submit a command here
// and block on its response, so there is never more than one write
// transaction in flight, and every atomic command runs inside its own
// transaction (R12).
type writer struct {
	client *ent.Client
	cmds   chan writeCmd
	done   chan struct{}
}

func newWriter(client *ent.Client) *writer {
	w := &writer{client: client, cmds: make(chan writeCmd), done: make(chan struct{})}
	go w.run()
	return w
}

func (w *writer) run() {
	defer close(w.done)
	for cmd := range w.cmds {
		cmd.resp <- w.execTx(cmd.ctx, cmd.fn)
	}
}

// execTx wraps fn in a transaction: fn's error rolls it back, its success
// commits.
func (w *writer) execTx(ctx context.Context, fn txFunc) writeResult {
	tx, err := w.client.Tx(ctx)
	if err != nil {
		return writeResult{err: fmt.Errorf("store: begin tx: %w", err)}
	}

	val, err := fn(ctx, tx)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return writeResult{err: fmt.Errorf("store: rollback after %w: %w", err, rerr)}
		}
		return writeResult{err: err}
	}

	if err := tx.Commit(); err != nil {
		return writeResult{err: fmt.Errorf("store: commit tx: %w", err)}
	}
	return writeResult{val: val}
}

// close stops accepting new commands and waits for the goroutine to exit.
func (w *writer) close() {
	close(w.cmds)
	<-w.done
}

// submit enqueues fn on the writer goroutine and blocks for its result —
// this is what makes the Writer interface's methods synchronous to their
// caller despite the single background writer. The response channel is
// buffered so the writer goroutine never blocks on a caller that gave up
// via ctx cancellation.
func submit[T any](ctx context.Context, w *writer, fn func(context.Context, *ent.Tx) (T, error)) (T, error) {
	var zero T
	resp := make(chan writeResult, 1)
	cmd := writeCmd{
		ctx: ctx,
		fn: func(ctx context.Context, tx *ent.Tx) (any, error) {
			return fn(ctx, tx)
		},
		resp: resp,
	}

	select {
	case w.cmds <- cmd:
	case <-ctx.Done():
		return zero, ctx.Err()
	}

	select {
	case r := <-resp:
		if r.err != nil {
			return zero, r.err
		}
		v, _ := r.val.(T)
		return v, nil
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}
