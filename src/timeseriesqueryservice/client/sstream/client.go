package main

import (
	"context"
	"fmt"
	client2 "github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/client"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/query"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/timeseriesquery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

func runServerStreamQuery(client pb.TimeSeriesQueryServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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
	stream, err := client.ServerStreamQuery(ctx, req)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[runServerStreamQuery] failed to recv: %v", err)
		}
		log.Printf("[runServerStreamQuery] resp -> %v", resp)
	}
}

func main() {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", client2.Host, client2.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[cstream] failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTimeSeriesQueryServiceClient(conn)
	runServerStreamQuery(client)
}
