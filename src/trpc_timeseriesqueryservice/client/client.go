package client

import (
	"fmt"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/timeseriesquery"
	"math/rand"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/client"
	"trpc.group/trpc-go/trpc-go/log"
)

func NewClientProxy() pb.TimeSeriesQueryServiceClientProxy {
	trpc.NewServer()
	openid := fmt.Sprintf("req-010-%04d", rand.Intn(1000))
	log.Infof("[NewContext] openid -> %s", openid)
	opts := []client.Option{
		// If you want to set the client receiving window size, use the client option `WithMaxWindowSize`.
		client.WithMaxWindowSize(1 * 1024 * 1024),
		client.WithMetaData("openid", []byte(openid)),
	}
	proxy := pb.NewTimeSeriesQueryServiceClientProxy(opts...)
	return proxy
}
