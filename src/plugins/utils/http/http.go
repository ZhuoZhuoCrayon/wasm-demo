package http

import (
	"github.com/valyala/fastjson"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	JSON               string = "application/json"
	FormDataUrlEncoded string = "application/x-www-form-urlencoded"
	PlainText          string = "text/plain"
)

// HandleJSON 以 json 解析 body 并提取指定字段
func HandleJSON(body []byte, fields map[string]bool) (kv map[string]string, err error) {
	fv, err := fastjson.ParseBytes(body)
	if err != nil {
		return nil, err
	}
	kv = make(map[string]string, len(fields))
	for field := range fields {
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
		if v != "" {
			kv[field] = v
		}
	}
	return kv, nil
}

// HandleFormUrlencoded 以 form-data 解析 body 并提取指定字段
func HandleFormUrlencoded(body []byte, fields map[string]bool) (kv map[string]string, err error) {
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, err
	}
	kv = make(map[string]string, len(fields))
	for field := range fields {
		if values.Has(field) {
			kv[field] = values.Get(field)
		}
	}
	return kv, nil
}

func ExtractIncompleteJSON(body []byte, fieldPatterns map[string]string) (kv map[string]string, err error) {
	bodyStr := string(body)
	kv = make(map[string]string, len(fieldPatterns)+1)
	kv["source"] = "IncompleteJSON"

	for field := range fieldPatterns {
		reg := regexp.MustCompile(fieldPatterns[field])
		match := reg.FindStringSubmatch(bodyStr)
		if len(match) > 1 {
			kv[field] = match[1]
		}
	}

	return kv, nil
}

func ExtractBytes(body []byte, contentType string, fields map[string]bool) (kv map[string]string, err error) {
	if strings.HasPrefix(contentType, JSON) {
		if kv, err = HandleJSON(body, fields); err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(contentType, FormDataUrlEncoded) || strings.HasPrefix(contentType, PlainText) {
		if kv, err = HandleFormUrlencoded(body, fields); err != nil {
			return nil, err
		}
	}
	return kv, nil
}

// Extract HTTP body 解析
func Extract(b io.ReadCloser, contentType string, fields map[string]bool) (kv map[string]string, err error) {
	body, err := io.ReadAll(b)
	if err != nil || len(body) == 0 {
		return nil, err
	}
	return ExtractBytes(body, contentType, fields)
}
