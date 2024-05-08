package main

import (
	"context"
)

// SayHelloImp servant implementation
type SayHelloImp struct {
}

// Init servant init
func (imp *SayHelloImp) Init() error {
	return nil
}

// Destroy servant destroy
func (imp *SayHelloImp) Destroy() {
}

func (imp *SayHelloImp) Add(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	*c = a + b
	return 0, nil
}
func (imp *SayHelloImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	*c = a + b
	return 0, nil
}
