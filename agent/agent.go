package agent

import (
	"context"
	"log"
	"morpher-agent/internal/client"
	"morpher-agent/internal/collector"
	"time"
)

type Sender interface {
	Send(ctx context.Context, payload any) error
}

type Agent struct {
	sender Sender
}

func New(sender Sender) *Agent {
	return &Agent{sender: sender}
}

func Run(ctx context.Context, serverURL string) error {
	sender := client.NewHTTPClient(serverURL, 10*time.Second)
	agent := New(sender)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			info, err := collector.CollectSystem()
			if err != nil {
				log.Printf("failed to collect system info: %v", err)
			}
			if err := agent.sender.Send(ctx, info); err != nil {
				log.Printf("failed to send system info: %v", err)
			}
		}
	}
}
