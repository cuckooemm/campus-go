package e

var msgFlags = map[int]string{
	SUCCESS:        "成功",
	ERROR:          "失败",
	SERVER_ERROR:   "连接拒绝",
	INVALID_PARAMS: "信息填写错误",

	CODE_SENTED:        "验证码已发送,请稍候重试",
	ERROR_CODE_UNSENT:  "请获取验证码",
	ERROR_CODE_INVALID: "验证码过期",
	SUCCESS_PRAISE_:    "点赞成功",
	ERROR_PRAISEED:     "你已经点过赞了",
	ERROR_PRAISE_FAIL:  "点赞失败",
	ERROR_DELETE_FAIL:  "删除失败",

	AUTH_INVALID:          "授权已失效",
	AUTH_LOGIN:            "请先登录",
	AUTH_UPDATE_FAIL:      "授权更新失败",
	AUTH_FORBIDDEN:        "拒绝访问",
	AUTH_FAIL:             "登录失败",
	AUTH_UPDATE_PASD_FAIL: "密码重置失败",

	ERROR_ACCOUNT_INVALID:      "无效帐号",
	ERROR_ACCOUNT_UNREGISTERED: "帐号未注册",
	ERROR_ACCOUNT_REGISTERED:   "帐号已注册",
	ERROR_NICKNAME_EXIST:       "呢称已存在",
	ERROR_REGISTER_FAIL:        "注册失败",
	ERROR_AVATAR_UP_FAIL:       "头像上传失败",
	MAX_LOGIN_FAIL:             "今日登录错误次数已达上限",
	ERROR_LOGIN_MISMATCH:       "帐号密码不匹配",


	STS_TOKEN_GET_FAIL: "授权获取失败",
}

func GetMsg(code int) string {
	msg, ok := msgFlags[code]
	if ok {
		return msg
	}
	return msgFlags[ERROR]
}
