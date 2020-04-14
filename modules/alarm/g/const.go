package g

const (
	VERSION = "0.1.1"
)

const (
	NO_ERROR          = 0
	REQ_TIMEOUT       = 1
	INVALID_RESP_CODE = 2
	KEYWORD_UNMATCH   = 3
)

var EventStatus = map[int64]string{
	NO_ERROR:          "ok",
	REQ_TIMEOUT:       "timeout",
	INVALID_RESP_CODE: "bad resp code",
	KEYWORD_UNMATCH:   "keyword unmatch",
}

//0.0.2 fix event_id error
//0.1.0 support send sms shell
//0.1.1 use mysql store event
