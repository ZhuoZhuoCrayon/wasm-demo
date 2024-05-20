package main

import (
	"context"
	"github.com/ZhuoZhuoCrayon/wasm-demo/src/core/query"
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/trpc_timeseriesqueryservice/timeseriesquery"
	"io"
	"math"
	"math/rand"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/errs"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/server"
)

func toPbSeries(series []*query.Series) []*pb.Series {
	pbSeries := make([]*pb.Series, len(series))
	for i, s := range series {
		pbSeries[i] = &pb.Series{Dimensions: s.Dimensions, DataPoints: make([]*pb.DataPoint, len(s.DataPoints))}
		for j, p := range s.DataPoints {
			pbSeries[i].DataPoints[j] = &pb.DataPoint{Timestamp: p.TimeStamp, Value: p.Value}
		}
	}
	return pbSeries
}

type timeSeriesQueryServiceImpl struct {
	pb.UnimplementedTimeSeriesQueryService
}

func randErr(errRate float64) error {
	if rand.Float64() < errRate {
		// 根据最佳实践，选择 > 10000 的错误码
		return errs.New(10001, "random error")
		// return errors.New("123")
	}
	return nil
}

// Query 一应一答模式
func (s *timeSeriesQueryServiceImpl) Query(ctx context.Context, req *pb.QueryRequest) (*pb.QueryResponse, error) {
	queryConfig := req.GetQueryConfig()
	queryer, err := query.NewQueryer(queryConfig.BeginTime, queryConfig.EndTime, queryConfig.GroupBy, int(queryConfig.Interval))
	log.Infof("[Query] recv range -> (%v, %v)", queryConfig.BeginTime, queryConfig.EndTime)
	if err != nil {
		return nil, err
	}
	if err = randErr(0.01); err != nil {
		return nil, err
	}
	series := queryer.Run()
	return &pb.QueryResponse{Series: toPbSeries(series)}, nil
}

// ClientStreamQuery 聚合客户端多次请求，统一返回
func (s *timeSeriesQueryServiceImpl) ClientStreamQuery(stream pb.TimeSeriesQueryService_ClientStreamQueryServer) error {
	endTime := int64(math.MinInt64)
	beginTime := int64(math.MaxInt64)
	var queryConfig *pb.QueryConfig
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Infof("[ClientStreamQuery] start to send range -> (%v, %v)", beginTime, endTime)
			queryer, err := query.NewQueryer(beginTime, endTime, queryConfig.GroupBy, int(queryConfig.Interval))
			if err != nil {
				return err
			}
			if err = randErr(0.01); err != nil {
				return err
			}
			series := queryer.Run()
			return stream.SendAndClose(&pb.QueryResponse{Series: toPbSeries(series)})
		}
		if err != nil {
			return err
		}
		queryConfig = req.QueryConfig
		log.Infof("[ClientStreamQuery] recv range -> (%v, %v)", queryConfig.BeginTime, queryConfig.EndTime)
		if beginTime > queryConfig.BeginTime {
			beginTime = queryConfig.BeginTime
		}
		if endTime < queryConfig.EndTime {
			endTime = queryConfig.EndTime
		}
	}
}

// ServerStreamQuery 服务端多次返回数据
func (s *timeSeriesQueryServiceImpl) ServerStreamQuery(req *pb.QueryRequest, stream pb.TimeSeriesQueryService_ServerStreamQueryServer) error {
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
		log.Infof("[ServerStreamQuery] start to send range -> (%v, %v)", tr.BeginTime, tr.EndTime)
		queryer, err := query.NewQueryer(tr.BeginTimeToUnix(), tr.EndTimeToUnix(), req.QueryConfig.GroupBy, int(req.QueryConfig.Interval))
		if err != nil {
			return err
		}
		if err = randErr(0.01); err != nil {
			return err
		}
		series := queryer.Run()
		if err := stream.Send(&pb.QueryResponse{Series: toPbSeries(series)}); err != nil {
			return err
		}
	}
	return nil
}

// BidirectionalStreamQuery 交替加载
func (s *timeSeriesQueryServiceImpl) BidirectionalStreamQuery(stream pb.TimeSeriesQueryService_BidirectionalStreamQueryServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		queryConfig := req.QueryConfig
		log.Infof("[BidirectionalStreamQuery] start to send range -> (%v, %v)", queryConfig.BeginTime, queryConfig.EndTime)
		queryer, err := query.NewQueryer(queryConfig.BeginTime, queryConfig.EndTime, queryConfig.GroupBy, int(req.QueryConfig.Interval))
		if err != nil {
			return err
		}
		if err = randErr(0.01); err != nil {
			return err
		}
		series := queryer.Run()
		if err := stream.Send(&pb.QueryResponse{Series: toPbSeries(series)}); err != nil {
			return err
		}
	}
}

func main() {
	s := trpc.NewServer(server.WithMaxWindowSize(1 * 1024 * 1024))
	pb.RegisterTimeSeriesQueryServiceService(s, &timeSeriesQueryServiceImpl{})
	if err := s.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
