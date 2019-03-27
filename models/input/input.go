package input

//	发送验证码字段
type SendCode struct {
	Account   string `form:"account" json:"account" binding:"required,max=50"`
	Flag      string // 标识账号是邮箱还是手机号
	Operation int8 `form:"operation" json:"operation" binding:"required,max=2"`
}

//	注册字段
type Register struct {
	Account  string `form:"account" json:"account" binding:"required,max=50"`
	Nickname string `form:"nickname" json:"nickname" binding:"required,max=12,min=2"`
	Code     string `form:"code" json:"code" binding:"required,len=6"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
	Flag     string // 标识账号是邮箱还是手机号
}

//	手机邮箱登录字段
type Login struct {
	Account  string `form:"account" json:"account" binding:"required,max=50"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
	Flag     string
}

// QQ 登录字段
type QQLogin struct {
	OpenID string `form:"openid" json:"openid" binding:"required"`
	Token  string `form:"token" json:"token" binding:"required"`
}

// 微信小程序登录 字段
type WeAppLogin struct {
	Code     string `form:"wx_code" json:"wx_code" binding:"required"`
	Nickname string `form:"nickname" json:"nickname"`
	Avatar   string `form:"avatar" json:"avatar"`
	Gender   int8   `form:"gender" json:"gender"`
}

//	找回字段
type Retrieve struct {
	Account  string `form:"account" json:"account" binding:"required,max=50"`
	Code     string `form:"code" json:"code" binding:"required,len=6"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
	Flag     string
}

// 修改用户资料 字段
type UserProfile struct {
	Nickname     string `form:"nickname" json:"nickname" binding:"omitempty,max=12,min=2"`
	Avatar       string `form:"avatar" json:"avatar"`
	Gender       string `form:"gender" json:"gender"`
	Birthday     string `form:"birthday" json:"birthday"`
	Bio          string `form:"bio" json:"bio"`
	AvatarStatus bool
}

// 发布评论
type Comment struct {
	Content   string `form:"content" json:"content" binding:"required,max=512"`
	DynamicID string `form:"dynamic_id" json:"dynamic_id" binding:"required,min=16"`
}

type Reply struct {
	Content   string `form:"content" json:"content" binding:"required,max=512"`
	CommentID string `form:"comment_id" json:"comment_id" binding:"required,min=16"`
}

type Paging struct {
	ID   string `form:"id" json:"id" binding:"required,min=16"`
	Page int    `form:"page" json:"page"`
}

type Praise struct {
	ID string `form:"id" json:"id" binding:"required,min=16"`
}
