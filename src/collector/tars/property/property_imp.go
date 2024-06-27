package property

import (
	"context"
	"encoding/json"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/propertyf"
	"log"
)

type PropertyImp struct {
}

func NewPropertyImp() *PropertyImp {
	return &PropertyImp{}
}

func (imp *PropertyImp) Init() error {
	return nil
}

func (imp *PropertyImp) Destroy() {
}

func (imp *PropertyImp) ReportPropMsg(tarsCtx context.Context, statmsg map[propertyf.StatPropMsgHead]propertyf.StatPropMsgBody) (int32, error) {
	for head, body := range statmsg {
		hj, err := json.Marshal(head)
		if err != nil {
			return -1, err
		}
		bj, err := json.Marshal(body)
		if err != nil {
			return -1, err
		}
		log.Printf("[ReportPropMsg] StatPropMsgHead -> %v", string(hj))
		log.Printf("[ReportPropMsg] StatPropMsgBody -> %v", string(bj))
	}
	return 0, nil
}
