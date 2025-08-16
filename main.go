package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"morpher-agent/agent"
)

func main() {
	// cmd.Execute()

	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		cancel()
	}()

	baseURL := os.Getenv("MORPHER_AGENT_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:9000"
	}
	if err := agent.Run(ctx, baseURL); err != nil {
		log.Fatal(err)
	}
}
