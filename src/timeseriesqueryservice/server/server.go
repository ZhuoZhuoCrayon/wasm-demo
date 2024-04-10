package main

import (
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/core/query"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/timeseriesqueryservice/timeseriesquery"
	"google.golang.org/grpc"
	"io"
	"log"
	"math"
	"net"
	"time"
)

func toGrpcSeries(series []*query.Series) []*pb.Series {
	grpcSeries := make([]*pb.Series, len(series))
	for i, s := range series {
		grpcSeries[i] = &pb.Series{Dimensions: s.Dimensions, DataPoints: make([]*pb.DataPoint, len(s.DataPoints))}
		for j, p := range s.DataPoints {
			grpcSeries[i].DataPoints[j] = &pb.DataPoint{Timestamp: p.TimeStamp, Value: p.Value}
		}
	}
	return grpcSeries
}

type timeSeriesQueryServer struct {
	pb.UnimplementedTimeSeriesQueryServiceServer
}

// ClientStreamQuery 聚合客户端多次请求，统一返回
func (s *timeSeriesQueryServer) ClientStreamQuery(stream pb.TimeSeriesQueryService_ClientStreamQueryServer) error {
	endTime := int64(math.MinInt64)
	beginTime := int64(math.MaxInt64)
	var queryConfig *pb.QueryConfig
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[ClientStreamQuery] start to send range -> (%v, %v)", beginTime, endTime)
			queryer, err := query.NewQueryer(beginTime, endTime, queryConfig.GroupBy, int(queryConfig.Interval))
			if err != nil {
				return err
			}
			series := queryer.Run()
			return stream.SendAndClose(&pb.QueryResponse{Series: toGrpcSeries(series)})
		}
		if err != nil {
			return err
		}
		// 模拟预计算耗时
		time.Sleep(500 * time.Millisecond)
		queryConfig = req.QueryConfig
		log.Printf("[ClientStreamQuery] reciver range -> (%v, %v)", queryConfig.BeginTime, queryConfig.EndTime)
		if beginTime > queryConfig.BeginTime {
			beginTime = queryConfig.BeginTime
		}
		if endTime < queryConfig.EndTime {
			endTime = queryConfig.EndTime
		}
	}
}

// ServerStreamQuery 服务端多次返回数据
func (s *timeSeriesQueryServer) ServerStreamQuery(req *pb.QueryRequest, stream pb.TimeSeriesQueryService_ServerStreamQueryServer) error {
	queryTimeRange, err := query.NewTimeRangeFromUnix(req.QueryConfig.BeginTime, req.QueryConfig.EndTime)
	if err != nil {
		return err
	}
	// 对客户端请求的时间范围进行分片
	ti := query.NewTimeRangeIteratorFromSegments(&queryTimeRange, 5)
	for {
		tr, end := ti.Next()
		if end {
			break
		}
		log.Printf("[ServerStreamQuery] start to send range -> (%v, %v)", tr.BeginTime, tr.EndTime)
		queryer, err := query.NewQueryer(tr.BeginTimeToUnix(), tr.EndTimeToUnix(), req.QueryConfig.GroupBy, int(req.QueryConfig.Interval))
		if err != nil {
			return err
		}
		series := queryer.Run()
		if err := stream.Send(&pb.QueryResponse{Series: toGrpcSeries(series)}); err != nil {
			return err
		}
	}
	return nil
}

// BidirectionalStreamQuery 交替加载
func (s *timeSeriesQueryServer) BidirectionalStreamQuery(stream pb.TimeSeriesQueryService_BidirectionalStreamQueryServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		queryConfig := req.QueryConfig
		log.Printf("[BidirectionalStreamQuery] start to send range -> (%v, %v)", queryConfig.BeginTime, queryConfig.EndTime)
		queryer, err := query.NewQueryer(queryConfig.BeginTime, queryConfig.EndTime, queryConfig.GroupBy, int(req.QueryConfig.Interval))
		if err != nil {
			return err
		}
		series := queryer.Run()
		if err := stream.Send(&pb.QueryResponse{Series: toGrpcSeries(series)}); err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTimeSeriesQueryServiceServer(grpcServer, &timeSeriesQueryServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
