package common

import (
	"strings"
)

// IsPrefixMatched 检测字符串是否与给定前缀任以匹配
func IsPrefixMatched(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

//func FormatKv(kv map[string]string) ([]sdk.KeyVal, *sdk.Trace, error) {
//	if kv == nil || len(kv) == 0 {
//		return nil, nil, FormatKvError
//	}
//	attrs := make([]sdk.KeyVal, 0, len(kv))
//	for k, v := range kv {
//		attrs = append(attrs, sdk.KeyVal{Key: k, Val: v})
//	}
//	trace := &sdk.Trace{}
//	return attrs, trace, nil
//}
