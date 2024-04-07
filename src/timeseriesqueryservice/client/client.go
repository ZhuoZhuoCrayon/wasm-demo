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
	openid := fmt.Sprintf("req-010-%04d", rand.Intn(1000))
	log.Printf("[NewContext] openid -> %s", openid)
	md := metadata.New(map[string]string{"openid": openid})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, cancel
}
