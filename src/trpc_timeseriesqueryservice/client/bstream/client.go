package main

import (
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/core/query"
	client2 "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/client"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/timeseriesquery"
	"io"
	"math/rand"
	"time"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func runBidirectionalStreamQuery(proxy pb.TimeSeriesQueryServiceClientProxy) {
	tr, err := query.NewTimeRangeFromStr("2024-04-02 14:00:00", "2024-04-02 16:00:00")
	if err != nil {
		log.Fatalf("[runBidirectionalStreamQuery] failed to TimeRange -> %v", err)
	}
	ti := query.NewTimeRangeIteratorFromSegments(&tr, 10)
	stream, err := proxy.BidirectionalStreamQuery(trpc.BackgroundContext())
	if err != nil {
		log.Fatalf("[runBidirectionalStreamQuery] failed to get stream: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("[runBidirectionalStreamQuery] failed to recv: %v", err)
			}
			log.Infof("[runBidirectionalStreamQuery] Dimensions -> %v", resp.Series[rand.Intn(len(resp.Series))].Dimensions)
		}
	}()
	for {
		next, end := ti.Next()
		if end {
			break
		}
		log.Infof("[runBidirectionalStreamQuery] start to query range -> (%v, %v)", next.BeginTime, next.EndTime)
		req := &pb.QueryRequest{
			QueryConfig: &pb.QueryConfig{
				BeginTime: next.BeginTimeToUnix(),
				EndTime:   next.EndTimeToUnix(),
				GroupBy:   []string{"host_id", "vpc"},
				Interval:  60,
			},
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("[runBidirectionalStreamQuery] stream.Send(%v) failed: %v", req, err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err := stream.CloseSend(); err != nil {
		log.Fatalf("[runBidirectionalStreamQuery] failed to CloseSend: %v", err)
	}
	<-waitc
}

func main() {
	proxy := client2.NewClientProxy()
	runBidirectionalStreamQuery(proxy)
}
