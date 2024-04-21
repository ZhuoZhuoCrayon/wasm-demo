package client

import (
	"fmt"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/timeseriesquery"
	"log"
	"math/rand"
	"trpc.group/trpc-go/trpc-go/client"
)

func NewClientProxy() pb.TimeSeriesQueryServiceClientProxy {
	openid := fmt.Sprintf("req-010-%04d", rand.Intn(1000))
	fmt.Printf("[NewContext] openid -> %s \n", openid)
	log.Printf("[NewContext] openid -> %s", openid)
	opts := []client.Option{
		// If you want to set the client receiving window size, use the client option `WithMaxWindowSize`.
		client.WithMaxWindowSize(1 * 1024 * 1024),
		client.WithTarget("ip://127.0.0.1:9002"),
		client.WithMetaData("openid", []byte(openid)),
	}
	proxy := pb.NewTimeSeriesQueryServiceClientProxy(opts...)
	return proxy
}
