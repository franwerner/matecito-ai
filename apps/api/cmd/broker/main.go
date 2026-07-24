package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	port := flag.String("port", "4300", "HTTP port for the broker")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "store": "n/a"})
	})

	addr := ":" + *port
	log.Info("broker listening", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Error("broker exited", "err", err)
		os.Exit(1)
	}
}
