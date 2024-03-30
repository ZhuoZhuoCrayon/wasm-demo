package main

import (
	"fmt"
	"github.com/valyala/fastjson"
	"strconv"
	"testing"
)

func TestStringFastJson(t *testing.T) {
	s := `{
		"openid": "0150857061015126968",
		"ret": 0,
		"err_code": "1003-131-35105"
	}`
	fv, err := fastjson.Parse(s)
	if err != nil {
		t.Errorf("failed to parse json")
	}
	for _, field := range []string{"openid", "ret", "err_code"} {
		if !fv.Exists(field) {
			continue
		}
		v := ""
		switch fv.Get(field).Type() {
		case fastjson.TypeString:
			v = string(fv.GetStringBytes(field))
		case fastjson.TypeNumber:
			v = strconv.Itoa(fv.GetInt(field))
		}
		fmt.Printf("field -> %s, get -> %s \n", field, v)
	}
}
