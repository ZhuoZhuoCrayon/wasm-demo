package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"log"
	"math/rand"
	"time"
)

const (
	Port = 9001
	Host = "localhost"
)

func NewContext() (context.Context, context.CancelFunc) {
	requestId := fmt.Sprintf("req-010-%04d", rand.Intn(1000))
	log.Printf("[NewContext] requestId -> %s", requestId)
	md := metadata.New(map[string]string{"request-id": requestId})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cancel
}
