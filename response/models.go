package response

type UserInfo struct {
	UID      string `json:"uid"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	QQID     string `json:"qq_id"`
	WeappID  string `json:"weapp_id"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
	School   int16  `json:"school"`
	Birthday string `json:"birthday"`
	Bio      string `json:"bio"`
}

type Links struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
}

type DynamicData struct {
	ID           string `json:"id"`
	UID          string `json:"uid"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	Gender       string `json:"gender"`
	Content      string `json:"content"`
	Image        *Image `json:"image,omitempty"`
	Browse       int    `json:"browse"`
	CommentCount int    `json:"comment_count"`
	PraiseCount  int    `json:"praise_count"`
	CreatedAt    int64  `json:"created_at"`
}

type Image struct {
	Url         []string `json:"url"`
	UrlOriginal []string `json:"url_original"`
}
type DynamicCommentData struct {
	ID          string `json:"id"`
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Gender      string `json:"gender"`
	Content     string `json:"content"`
	PraiseCount int    `json:"praise_count"`
	ReplyCount  int    `json:"reply_count"`
	CreatedAt   int64  `json:"created_at"`
}

type ReplyDataRespond struct {
	ID          string `json:"id"`
	UID         string `json:"uid"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Gender      string `json:"gender"`
	RName       string `json:"r_name"`
	Content     string `json:"content"`
	PraiseCount int    `json:"praise_count"`
	CreatedAt   int64  `json:"created_at"`
}

type MessageDataResponse struct {
	ID        string `json:"id"`
	Type      int    `json:"type"`
	UID       string `json:"uid"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar"`
	Gender    string `json:"gender"`
	Content   string `json:"content"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
}
