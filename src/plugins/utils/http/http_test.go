package http

import (
	"bufio"
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestHandleFormUrlencoded(t *testing.T) {
	tests := []struct {
		name   string
		body   []byte
		fields map[string]bool
		want   map[string]string
	}{
		{
			name:   "application/x-www-form-urlencoded; charset=utf-8",
			body:   []byte("offer_id=123&action=buyGiftPackage&openid=1234"),
			fields: map[string]bool{"openid": true},
			want:   map[string]string{"openid": "1234"},
		},
	}
	for _, tt := range tests {
		kv, _ := HandleFormUrlencoded(tt.body, tt.fields)
		if !reflect.DeepEqual(kv, tt.want) {
			t.Errorf("kv want %v but got %v", tt.want, kv)
		}
	}
}

func TestHandleJSON(t *testing.T) {
	tests := []struct {
		name   string
		body   []byte
		fields map[string]bool
		want   map[string]string
	}{
		{
			name: "application/json; charset=utf-8",
			body: []byte{
				0x7b, 0x22, 0x6d, 0x73, 0x67, 0x22, 0x3a, 0x22, 0x4f, 0x4b, 0x22, 0x2c, 0x22, 0x72, 0x65, 0x74, 0x75,
				0x72, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x3a, 0x30, 0x2c, 0x22, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e,
				0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3a, 0x30, 0x2c, 0x22, 0x73, 0x65, 0x71, 0x49, 0x64, 0x22, 0x3a,
				0x22, 0x63, 0x6f, 0x34, 0x64, 0x63, 0x62, 0x72, 0x65, 0x71, 0x35, 0x69, 0x71, 0x33, 0x39, 0x62, 0x6a,
				0x36, 0x67, 0x36, 0x67, 0x22, 0x7d,
			},
			fields: map[string]bool{"seqId": true},
			want:   map[string]string{"seqId": "co4dcbreq5iq39bj6g6g"},
		},
	}
	for _, tt := range tests {
		kv, _ := HandleJSON(tt.body, tt.fields)
		if !reflect.DeepEqual(kv, tt.want) {
			t.Errorf("kv want %v but got %v", tt.want, kv)
		}
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		fields  map[string]bool
		want    map[string]string
	}{
		{
			name: "application/json; charset=utf-8",
			payload: []byte(`POST /your-endpoint HTTP/1.1
Host: test.com
Content-Type: application/json; charset=utf-8
Content-Length: 74

{"msg":"OK","returnCode":0,"returnValue":0,"seqId":"co4dcbreq5iq39bj6g6g"}`),
			fields: map[string]bool{"seqId": true},
			want:   map[string]string{"seqId": "co4dcbreq5iq39bj6g6g"},
		},
		{
			name: "application/x-www-form-urlencoded; charset=utf-8",
			payload: []byte(`POST /your-endpoint HTTP/1.1
Host: test.com
Content-Type: application/x-www-form-urlencoded; charset=utf-8
Content-Length: 46

offer_id=123&action=buyGiftPackage&openid=1234`),
			fields: map[string]bool{"openid": true},
			want:   map[string]string{"openid": "1234"},
		},
	}

	for _, tt := range tests {
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(tt.payload)))
		if err != nil {
			t.Errorf("failed to read request: err -> %v \n", err)
		}
		kv, _ := Extract(req.Body, req.Header.Get("Content-Type"), tt.fields)
		if !reflect.DeepEqual(kv, tt.want) {
			t.Errorf("kv want %v but got %v", tt.want, kv)
		}
	}
}

func TestExtractIncompleteJSON(t *testing.T) {
	tests := []struct {
		name          string
		payload       []byte
		fieldPatterns map[string]string
		want          map[string]string
	}{
		{
			name:    "application/json; charset=utf-8",
			payload: []byte(`{"code": 500, "message": "Server Error", "result": false, "code_name":"ERROR"`),
			fieldPatterns: map[string]string{
				"code":      `"code":\s*(\d+)`,
				"result":    `"result":\s*([a-zA-Z]+)`,
				"message":   `"message":\s*"([^"]+)"`,
				"code_name": `"code_name":\s*"([^"]+)"`,
			},
			want: map[string]string{
				"code": "500", "message": "Server Error", "result": "false", "code_name": "ERROR", "source": "IncompleteJSON",
			},
		},
	}

	for _, tt := range tests {
		kv, _ := ExtractIncompleteJSON(tt.payload, tt.fieldPatterns)
		if !reflect.DeepEqual(kv, tt.want) {
			t.Errorf("kv want %v but got %v", tt.want, kv)
		}
	}
}
