package main

import (
	"fmt"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/core/query"
	client2 "github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/client"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/timeseriesquery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
)

func runClientStreamQuery(client pb.TimeSeriesQueryServiceClient) {
	timeStrRanges := [][]string{
		{"2024-04-02 14:00:00", "2024-04-02 14:20:00"},
		{"2024-04-02 14:20:00", "2024-04-02 14:40:00"},
		{"2024-04-02 14:40:00", "2024-04-02 15:00:00"},
	}
	timeRanges := make([]*query.TimeRange, len(timeStrRanges))
	for i, timeStrRange := range timeStrRanges {
		tr, err := query.NewTimeRangeFromStr(timeStrRange[0], timeStrRange[1])
		if err != nil {
			log.Fatalf("[runClientStreamQuery] failed to create TimeRange -> %v", err)
		}
		timeRanges[i] = &tr
	}
	ctx, cancel := client2.NewContext()
	defer cancel()
	stream, err := client.ClientStreamQuery(ctx)
	if err != nil {
		log.Fatalf("[runClientStreamQuery] failed to get stream: %v", err)
	}
	for _, tr := range timeRanges {
		log.Printf("[runClientStreamQuery] start to query range -> (%v, %v)", tr.BeginTime, tr.EndTime)
		req := &pb.QueryRequest{
			QueryConfig: &pb.QueryConfig{
				BeginTime: tr.BeginTimeToUnix(),
				EndTime:   tr.EndTimeToUnix(),
				GroupBy:   []string{"host_id", "vpc"},
				Interval:  60,
			},
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("[runClientStreamQuery] stream.Send(%v) failed: %v", req, err)
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("[runClientStreamQuery] failed to recv: %v", err)
	}
	log.Printf("[runClientStreamQuery] Dimensions -> %v", resp.Series[rand.Intn(len(resp.Series))].Dimensions)
}

func main() {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", client2.Host, client2.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[cstream] failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTimeSeriesQueryServiceClient(conn)
	runClientStreamQuery(client)
}
