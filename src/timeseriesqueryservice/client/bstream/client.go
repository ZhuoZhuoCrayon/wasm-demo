package main

import (
	"fmt"
	client2 "github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/client"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/query"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/timeseriesquery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math/rand"
	"time"
)

func runBidirectionalStreamQuery(client pb.TimeSeriesQueryServiceClient) {
	tr, err := query.NewTimeRangeFromStr("2024-04-02 14:00:00", "2024-04-02 16:00:00")
	if err != nil {
		log.Fatalf("[runBidirectionalStreamQuery] failed to TimeRange -> %v", err)
	}
	ti := query.NewTimeRangeIteratorFromSegments(&tr, 10)
	ctx, cancel := client2.NewContext()
	defer cancel()
	stream, err := client.BidirectionalStreamQuery(ctx)
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
			log.Printf("[runBidirectionalStreamQuery] Dimensions -> %v", resp.Series[rand.Intn(len(resp.Series))].Dimensions)
		}
	}()
	for {
		next, end := ti.Next()
		if end {
			break
		}
		log.Printf("[runBidirectionalStreamQuery] start to query range -> (%v, %v)", next.BeginTime, next.EndTime)
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
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", client2.Host, client2.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("[cstream] failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTimeSeriesQueryServiceClient(conn)
	runBidirectionalStreamQuery(client)
}
