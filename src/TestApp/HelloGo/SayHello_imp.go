package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

var (
	RandomError = errors.New("random error")
)

func randSleep() {
	r := 100 + rand.Intn(200)
	time.Sleep(time.Duration(r) * time.Millisecond)
}

func randError(errRate float64) (int32, error) {
	if rand.Float64() < errRate {
		// 选择 -1 ～ -13 的错误码
		return int32(-1 * (rand.Intn(12) + 1)), RandomError
	}
	return 0, nil
}

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
	randSleep()
	*c = a + b
	log.Printf("a(%v) + b(%v) = c(%v)", a, b, c)
	return randError(0.01)
}
func (imp *SayHelloImp) Sub(ctx context.Context, a int32, b int32, c *int32) (int32, error) {
	randSleep()
	*c = a - b
	log.Printf("a(%v) - b(%v) = c(%v)", a, b, c)
	return randError(0.01)
}
