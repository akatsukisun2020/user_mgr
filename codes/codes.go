package codes

const (
	ERROR_PARMA       = 1001 // 参数错误
	ERROR_LOGINERROR  = 1002 // 登录错误
	ERROR_QUERYREDIS  = 1003 // 查询存储失败
	ERROR_NOUSER      = 1004 // 用户不存在
	ERROR_TOKENCHECK  = 1005 // TOKEN校验失败
	ERROR_TOKENEXPIRE = 1006 // TOKEN过期
	ERROR_SETREDIS    = 1007 // 写存储失败
)
