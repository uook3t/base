package response

var (
	messages = make(map[int32]string)
)

const (
	CodeSuccess       = int32(0)
	CodeFailed        = int32(1)
	CodeParamError    = int32(2)
	CodeNotExist      = int32(3)
	CodeForbidden     = int32(4)
	CodeLoginRequired = int32(5)
	CodeTimeout       = int32(6)
)

const (
	MsgSuccess       = "Success"
	MsgFailed        = "Failed"
	MsgParamError    = "Param Error"
	MsgNotExist      = "Not Exist"
	MsgForbidden     = "Forbidden"
	MsgLoginRequired = "Login Required"
	MsgTimeout       = "Timeout"
)

func init() {
	messages[CodeSuccess] = MsgSuccess
	messages[CodeFailed] = MsgFailed
	messages[CodeParamError] = MsgParamError
	messages[CodeNotExist] = MsgNotExist
	messages[CodeForbidden] = MsgForbidden
	messages[CodeLoginRequired] = MsgLoginRequired
	messages[CodeTimeout] = MsgTimeout
}

func getMessageByCode(code int32) string {
	if msg, ok := messages[code]; ok {
		return msg
	}
	return MsgFailed
}
