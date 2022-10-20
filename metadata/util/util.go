package util

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-gost/core/metadata"
)

func GetBool(md metadata.Metadata, key string) (v bool) {
	if md == nil || !md.IsExists(key) {
		return
	}
	switch vv := md.Get(key).(type) {
	case bool:
		return vv
	case int:
		return vv != 0
	case string:
		v, _ = strconv.ParseBool(vv)
		return
	}
	return
}

func GetInt(md metadata.Metadata, key string) (v int) {
	if md == nil || !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case bool:
		if vv {
			v = 1
		}
	case int:
		return vv
	case string:
		v, _ = strconv.Atoi(vv)
		return
	}
	return
}

func GetFloat(md metadata.Metadata, key string) (v float64) {
	if md == nil || !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case float64:
		return vv
	case int:
		return float64(vv)
	case string:
		v, _ = strconv.ParseFloat(vv, 64)
		return
	}
	return
}

func GetDuration(md metadata.Metadata, key string) (v time.Duration) {
	if md == nil || !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case int:
		return time.Duration(vv) * time.Second
	case string:
		v, _ = time.ParseDuration(vv)
		if v == 0 {
			n, _ := strconv.Atoi(vv)
			v = time.Duration(n) * time.Second
		}
	}
	return
}

func GetString(md metadata.Metadata, key string) (v string) {
	if md == nil || !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case string:
		v = vv
	case int:
		v = strconv.FormatInt(int64(vv), 10)
	case int64:
		v = strconv.FormatInt(vv, 10)
	case uint:
		v = strconv.FormatUint(uint64(vv), 10)
	case uint64:
		v = strconv.FormatUint(uint64(vv), 10)
	case bool:
		v = strconv.FormatBool(vv)
	case float32:
		v = strconv.FormatFloat(float64(vv), 'f', -1, 32)
	case float64:
		v = strconv.FormatFloat(float64(vv), 'f', -1, 64)
	}

	return
}

func GetStrings(md metadata.Metadata, key string) (ss []string) {
	if md == nil || !md.IsExists(key) {
		return
	}

	switch v := md.Get(key).(type) {
	case []string:
		ss = v
	case []any:
		for _, vv := range v {
			if s, ok := vv.(string); ok {
				ss = append(ss, s)
			}
		}
	}
	return
}

func GetStringMap(md metadata.Metadata, key string) (m map[string]any) {
	if md == nil || !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case map[string]any:
		return vv
	case map[any]any:
		m = make(map[string]any)
		for k, v := range vv {
			m[fmt.Sprintf("%v", k)] = v
		}
	}
	return
}

func GetStringMapString(md metadata.Metadata, key string) (m map[string]string) {
	if md == nil || !md.IsExists(key) {
		return
	}

	switch vv := md.Get(key).(type) {
	case map[string]any:
		m = make(map[string]string)
		for k, v := range vv {
			m[k] = fmt.Sprintf("%v", v)
		}
	case map[any]any:
		m = make(map[string]string)
		for k, v := range vv {
			m[fmt.Sprintf("%v", k)] = fmt.Sprintf("%v", v)
		}
	}
	return
}
