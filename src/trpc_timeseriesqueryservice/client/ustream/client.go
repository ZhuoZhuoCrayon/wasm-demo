package main

import (
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/core/query"
	client2 "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/client"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/timeseriesquery"
	"log"
	"math/rand"
	"trpc.group/trpc-go/trpc-go"
)

func runQuery(proxy pb.TimeSeriesQueryServiceClientProxy) {
	tr, err := query.NewTimeRangeFromStr("2024-04-02 14:00:00", "2024-04-02 14:01:00")
	if err != nil {
		log.Fatalf("[runQuery] failed to TimeRange -> %v", err)
	}
	req := &pb.QueryRequest{
		QueryConfig: &pb.QueryConfig{
			BeginTime: tr.BeginTimeToUnix(),
			EndTime:   tr.EndTimeToUnix(),
			GroupBy:   []string{"host_id", "vpc"},
			Interval:  60,
		},
	}
	resp, err := proxy.Query(trpc.BackgroundContext(), req)
	if err != nil {
		log.Fatalf("[runQuery] failed to Query: %v", err)
	}
	log.Printf("[runQuery] Dimensions -> %v", resp.Series[rand.Intn(len(resp.Series))].Dimensions)
}

func main() {
	proxy := client2.NewClientProxy()
	runQuery(proxy)
}
