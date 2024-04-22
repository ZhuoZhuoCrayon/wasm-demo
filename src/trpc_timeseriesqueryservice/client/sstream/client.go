package main

import (
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/core/query"
	client2 "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/client"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/timeseriesquery"
	"io"
	"math/rand"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
)

func runServerStreamQuery(proxy pb.TimeSeriesQueryServiceClientProxy) {
	tr, err := query.NewTimeRangeFromStr("2024-04-02 14:00:00", "2024-04-02 15:00:00")
	if err != nil {
		log.Fatalf("[runServerStreamQuery] failed to TimeRange -> %v", err)
	}
	req := &pb.QueryRequest{
		QueryConfig: &pb.QueryConfig{
			BeginTime: tr.BeginTimeToUnix(),
			EndTime:   tr.EndTimeToUnix(),
			GroupBy:   []string{"host_id", "vpc"},
			Interval:  60,
		},
	}
	stream, err := proxy.ServerStreamQuery(trpc.BackgroundContext(), req)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[runServerStreamQuery] failed to recv: %v", err)
		}
		log.Infof("[runServerStreamQuery] Dimensions -> %v", resp.Series[rand.Intn(len(resp.Series))].Dimensions)
	}
}

func main() {
	proxy := client2.NewClientProxy()
	runServerStreamQuery(proxy)
}
