package param

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/710leo/urlooker/modules/web/http/errors"
)

func String(r *http.Request, key string, defVal string) string {
	if val, ok := r.URL.Query()[key]; ok {
		if val[0] == "" {
			return defVal
		}
		return strings.TrimSpace(val[0])
	}

	if r.Form == nil {
		err := r.ParseForm()
		if err != nil {
			panic(errors.BadRequestError())
		}
	}

	val := r.Form.Get(key)
	if val == "" {
		return defVal
	}

	return strings.TrimSpace(val)
}

func MustString(r *http.Request, key string, displayName ...string) string {
	val := String(r, key, "")
	if val == "" {
		name := key
		if len(displayName) > 0 {
			name = displayName[0]
		}
		panic(errors.BadRequestError(fmt.Sprintf("%s is necessary", name)))
	}
	return val
}

func Int64(r *http.Request, key string, defVal int64) int64 {
	raw := String(r, key, "")
	if raw == "" {
		return defVal
	}

	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return defVal
	}

	return val
}

func MustInt64(r *http.Request, key string, displayName ...string) int64 {
	raw := String(r, key, "")
	if raw == "" {
		name := key
		if len(displayName) > 0 {
			name = displayName[0]
		}
		panic(errors.BadRequestError(fmt.Sprintf("%s is necessary", name)))
	}

	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		panic(errors.BadRequestError())
	}

	return val
}

func Int(r *http.Request, key string, defVal int) int {
	raw := String(r, key, "")
	if raw == "" {
		return defVal
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		return defVal
	}

	return val
}

func MustInt(r *http.Request, key string, displayName ...string) int {
	name := key
	if len(displayName) > 0 {
		name = displayName[0]
	}

	raw := String(r, key, "")
	if raw == "" {
		panic(errors.BadRequestError(fmt.Sprintf("%s is necessary", name)))
	}

	val, err := strconv.Atoi(raw)
	if err != nil {
		panic(errors.BadRequestError(fmt.Sprintf("%s should be integer", name)))
	}

	return val
}

func Float64(r *http.Request, key string, defVal float64) float64 {
	raw := String(r, key, "")
	if raw == "" {
		return defVal
	}

	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return defVal
	}

	return val
}

func MustFloat64(r *http.Request, key string, displayName ...string) float64 {
	raw := String(r, key, "")
	if raw == "" {
		name := key
		if len(displayName) > 0 {
			name = displayName[0]
		}
		panic(errors.BadRequestError(fmt.Sprintf("%s is necessary", name)))
	}

	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		panic(errors.BadRequestError())
	}

	return val
}

func Bool(r *http.Request, key string, defVal bool) bool {
	raw := String(r, key, "")
	if raw == "true" || raw == "1" || raw == "on" || raw == "checked" || raw == "yes" {
		return true
	} else if raw == "false" || raw == "0" || raw == "off" || raw == "" || raw == "no" {
		return false
	} else {
		return defVal
	}
}

func MustBool(r *http.Request, key string) bool {
	raw := String(r, key, "")
	if raw == "true" || raw == "1" || raw == "on" || raw == "checked" || raw == "yes" {
		return true
	} else if raw == "false" || raw == "0" || raw == "off" || raw == "" || raw == "no" {
		return false
	} else {
		panic(errors.BadRequestError())
	}
}
