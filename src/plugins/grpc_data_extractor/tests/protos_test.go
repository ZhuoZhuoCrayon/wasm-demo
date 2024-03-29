package main

import (
	pb "github.com/ZhuoZhuoCrayon/wasm-demo/src/plugins/grpc_data_extractor/protos"
	"reflect"
	"strconv"
	"testing"
)

var ExpectedDataFields = map[string]bool{
	"Openid": true,
	"UserId": true,
}

func checkAndAddFieldsToMap(s interface{}, fields map[string]bool, kv map[string]string) {
	v := reflect.ValueOf(s).Elem()
	for field := range fields {
		f := v.FieldByName(field)
		if f.IsValid() {
			switch f.Kind() {
			case reflect.String:
				kv[field] = f.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				kv[field] = strconv.FormatInt(f.Int(), 10)
			}
		}
	}
}

func TestParseRequestPayload(t *testing.T) {
	payload := []byte{
		0x0a, 0x0f, 0x6f, 0x70, 0x65, 0x6e, 0x69, 0x64, 0x2d, 0x30, 0x31, 0x30, 0x2d, 0x30,
		0x30, 0x30, 0x32, 0x10, 0x02, 0x1a, 0x08, 0x4a, 0x6f, 0x68, 0x6e, 0x44, 0x6f, 0x65,
		0x22, 0x14, 0x6a, 0x6f, 0x68, 0x6e, 0x2e, 0x64, 0x6f, 0x65, 0x40, 0x65, 0x78, 0x61,
		0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x2a, 0x0c, 0x0a, 0x08, 0x35, 0x35,
		0x35, 0x2d, 0x31, 0x32, 0x33, 0x34, 0x10, 0x01, 0x2a, 0x0c, 0x0a, 0x08, 0x35, 0x35,
		0x35, 0x2d, 0x35, 0x36, 0x37, 0x38, 0x10, 0x02, 0x2a, 0x0c, 0x0a, 0x08, 0x35, 0x35,
		0x35, 0x2d, 0x39, 0x30, 0x31, 0x32, 0x10, 0x02,
	}
	resp := &pb.UserInfo{}
	if err := resp.UnmarshalVT(payload); err != nil {
		return
	}
	kv := make(map[string]string)
	checkAndAddFieldsToMap(resp, ExpectedDataFields, kv)
	want := "openid-010-0002"
	if resp.Openid != want {
		t.Errorf("openid want %s but got %s", want, resp.Openid)
	}
}
