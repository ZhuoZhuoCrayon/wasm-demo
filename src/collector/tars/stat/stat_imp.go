package stat

import (
	"context"
	"encoding/json"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/statf"
	"log"
)

type StatImp struct {
}

func NewStatImp() *StatImp {
	return &StatImp{}
}

func (imp *StatImp) Init() error {
	return nil
}

func (imp *StatImp) Destroy() {
}

func (imp *StatImp) ReportMicMsg(tarsCtx context.Context, msg map[statf.StatMicMsgHead]statf.StatMicMsgBody, bFromClient bool) (int32, error) {
	for head, body := range msg {
		hj, err := json.Marshal(head)
		if err != nil {
			return -1, err
		}
		bj, err := json.Marshal(body)
		if err != nil {
			return -1, err
		}
		log.Printf("[ReportMicMsg] StatMicMsgHead -> %v", string(hj))
		log.Printf("[ReportMicMsg] StatMicMsgBody -> %v", string(bj))
	}
	return 0, nil
}

func (imp *StatImp) ReportSampleMsg(tarsCtx context.Context, msg []statf.StatSampleMsg) (int32, error) {
	return 0, nil
}
