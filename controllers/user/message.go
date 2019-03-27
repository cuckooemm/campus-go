package user

import (
	"campus/helper/encrypt"
	"campus/helper/oauth"
	"campus/helper/setting"
	"campus/models"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"strconv"
)

//获取评论
func CommentToMeList(c *gin.Context) {
	var perPage = setting.AppSetting.PageSize
	rsp := response.NewResponse(c)
	key, _ := c.Get("claims")
	claims, _ := key.(*oauth.Claims)
	// 获取分页加载的页数
	pg := c.DefaultQuery("page", "1")
	offset, err := strconv.Atoi(pg)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var list []models.DynamicComment
	var links response.Links
	var meta response.Meta
	//如果不是第一页  设置links 的prev
	if offset > 1 {
		links.Prev = setting.AppSetting.ServerApiHost + `v1/message?page=` + strconv.Itoa(offset-1)
	}

	// 根据页码设置偏移量
	offset = (offset - 1) * perPage
	if err := models.GetMessageToMe(claims.UID,offset, perPage, &list); err != nil {
		rsp.Error()
		return
	}
	// 从评论表获取回复
	// 根据返回结果定义data 数组的长度
	length := len(list)
	if length == perPage + 1 {
		// 说明有值 添加下一页的连接
		links.Next = setting.AppSetting.ServerApiHost + `v1/message?page=` + strconv.Itoa(offset/perPage+2)
		// 等于11时减去1 回复到需要数组的长度
		length -= 1
		// 切片数组，只去 10条数据
		list = list[0:10]
	}
	data := make([]response.MessageDataResponse, length)
	for index,value := range list{
		value.User.AvatarToB60()
		// Type = 1  消息是回复的评论  0 消息是回复的动态
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.UID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].Message = value.Content
		data[index].CreatedAt = value.CreatedAt.Unix()
		if value.ParentID != 0 {
			data[index].Content = value.Comment.Content
			data[index].Type = 1
		}else {
			data[index].Content = value.Dynamic.Content
			data[index].Type = 0
		}
	}
	// 设置meta  当前页码 和当前结果集数量
	meta.CurrentPage = offset/perPage + 1
	meta.PerPage = length
	rsp.Success("",gin.H{
		"list":  data,
		"links": links,
		"meta":  meta,
	})
}

func UserMessageList(c *gin.Context) {
	var pageCount = setting.AppSetting.PageSize
	rsp := response.NewResponse(c)
	key, _ := c.Get("claims")
	claims, _ := key.(*oauth.Claims)
	// 获取ID
	id := c.Query("id")
	var lastID int64
	if id != "" {
		var err error
		lastID, err = encrypt.Decode(id)
		if err != nil {
			rsp.ErrorMsg(e.INVALID_PARAMS)
			return
		}
	}
	var list = make([]models.DynamicComment,pageCount)
	if id != "" {
		if err := models.GetMoreUserMessage(claims.UID,lastID,pageCount, &list); err != nil {
			rsp.Error()
			return
		}
	} else {
		if err := models.GetUserMessage(claims.UID,pageCount, &list); err != nil {
			rsp.Error()
			return
		}
	}
	data := make([]response.MessageDataResponse, len(list))
	for index,value := range list{
		value.User.AvatarToB60()
		// Type = 1  回复  0 评论动态
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.UID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].Message = value.Content
		data[index].CreatedAt = value.CreatedAt.Unix()
		if value.ParentID != 0 {
			data[index].Content = value.Comment.Content
			data[index].Type = 1
		}else {
			data[index].Content = value.Dynamic.Content
			data[index].Type = 0
		}
	}
	rsp.Success("",data)
}