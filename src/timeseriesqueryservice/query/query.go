package query

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const MagicTimeFormat string = "2006-01-02 15:04:05"

var (
	TimeParseError        = errors.New("failed to parse time")
	InvalidTimeRangeError = errors.New("invalid time range")
)

type DataPoint struct {
	Value     float64
	TimeStamp int64
}

type Series struct {
	Dimensions map[string]string
	DataPoints []*DataPoint
}

type TimeRange struct {
	BeginTime time.Time
	EndTime   time.Time
}

func NewTimeRange(beginTime, endTime time.Time) (TimeRange, error) {
	if !endTime.After(beginTime) {
		return TimeRange{}, InvalidTimeRangeError
	}
	return TimeRange{BeginTime: beginTime, EndTime: endTime}, nil
}

func NewTimeRangeFromUnix(beginTime, endTime int64) (TimeRange, error) {
	return NewTimeRange(time.Unix(beginTime, 0), time.Unix(endTime, 0))
}

func NewTimeRangeFromStr(beginTimeStr, endTimeStr string) (TimeRange, error) {
	beginTime, err := time.Parse(MagicTimeFormat, beginTimeStr)
	if err != nil {
		return TimeRange{}, TimeParseError
	}
	endTime, err := time.Parse(MagicTimeFormat, endTimeStr)
	if err != nil {
		return TimeRange{}, TimeParseError
	}
	return NewTimeRange(beginTime, endTime)
}

func (tr TimeRange) Sub() time.Duration {
	return tr.EndTime.Sub(tr.BeginTime)
}

func (tr TimeRange) BeginTimeToUnix() int64 {
	return tr.BeginTime.Unix()
}

func (tr TimeRange) EndTimeToUnix() int64 {
	return tr.EndTime.Unix()
}

type TimeRangeIterator struct {
	timeRange *TimeRange
	interval  time.Duration
	current   int
}

func NewTimeRangeIterator(timeRange *TimeRange, interval time.Duration) *TimeRangeIterator {
	return &TimeRangeIterator{
		timeRange: timeRange,
		interval:  interval,
		current:   0,
	}
}

func NewTimeRangeIteratorFromSegments(timeRange *TimeRange, segments int) *TimeRangeIterator {
	interval := timeRange.Sub() / time.Duration(segments)
	return NewTimeRangeIterator(timeRange, interval)
}

func (t *TimeRangeIterator) Next() (TimeRange, bool) {
	rangeBegin := t.timeRange.BeginTime.Add(t.interval * time.Duration(t.current))
	if rangeBegin.After(t.timeRange.EndTime) || rangeBegin.Equal(t.timeRange.EndTime) {
		return TimeRange{}, true
	}
	rangeEnd := rangeBegin.Add(t.interval)
	if rangeEnd.After(t.timeRange.EndTime) {
		rangeEnd = t.timeRange.EndTime
	}
	t.current++
	timeRange, _ := NewTimeRange(rangeBegin, rangeEnd)
	return timeRange, false
}

func (t *TimeRangeIterator) Reset() {
	t.current = 0
}

type Queryer struct {
	groupBy           []string
	timeRangeIterator *TimeRangeIterator
}

func (q *Queryer) Run() []*Series {
	cnt := rand.Intn(5) + 1
	series := make([]*Series, cnt)
	for i := 0; i < cnt; i++ {
		dimensions := make(map[string]string, len(q.groupBy))
		for _, d := range dimensions {
			dimensions[d] = fmt.Sprintf("%s-%05d", d, i)
		}
		singleSeries := &Series{
			Dimensions: map[string]string{"host_id": strconv.Itoa(10000 + i), "metric": "cpu_usage"},
			DataPoints: []*DataPoint{},
		}
		for {
			tr, end := q.timeRangeIterator.Next()
			if end {
				q.timeRangeIterator.Reset()
				break
			}
			dp := DataPoint{TimeStamp: tr.BeginTimeToUnix(), Value: rand.Float64()}
			singleSeries.DataPoints = append(singleSeries.DataPoints, &dp)
		}
		series[i] = singleSeries
	}
	// 模拟数据查询耗时
	time.Sleep(500 * time.Millisecond)
	return series
}

func NewQueryer(beginTime, endTime int64, groupBy []string, interval int) (*Queryer, error) {
	tr, err := NewTimeRangeFromUnix(beginTime, endTime)
	if err != nil {
		return nil, err
	}
	ti := NewTimeRangeIterator(&tr, time.Duration(interval)*time.Second)
	return &Queryer{groupBy: groupBy, timeRangeIterator: ti}, nil
}
