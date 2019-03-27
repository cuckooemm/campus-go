package dynamic

import (
	"campus/helper/encrypt"
	"campus/helper/oauth"
	"campus/helper/redis/cache"
	"campus/helper/sensitiveWord"
	"campus/helper/setting"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strconv"
)

func ReplyList(c *gin.Context) {
	var perPage = setting.AppSetting.PageSize
	rsp := response.NewResponse(c)
	// 获取分页加载的页数
	pg := c.DefaultQuery("page", "1")
	offset, err := strconv.Atoi(pg)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	id := c.Query("id")
	commentID,err := encrypt.Decode(id)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}

	var list []models.DynamicComment
	var links response.Links
	var meta response.Meta
	//如果不是第一页  设置links 的prev
	if offset > 1 {
		links.Prev = setting.AppSetting.ServerApiHost + `v1/dynamic/reply?page=` + strconv.Itoa(offset-1) + `&id=` + id
	}
	// 根据页码设置偏移量
	offset = (offset - 1) * perPage
	if err := models.GetDynamicCommentReplyList(commentID, offset, perPage, &list); err != nil {
		rsp.Error()
		return
	}
	// 根据返回结果定义data 数组的长度
	length := len(list)
	if length > perPage {
		// 说明有值 添加下一页的连接
		links.Next = setting.AppSetting.ServerApiHost + `v1/dynamic/reply?page=` + strconv.Itoa(offset/perPage+2) + `&id=` + id
		// 等于11时减去1 回复到需要数组的长度
		length -= 1
		// 切片数组，只去 10条数据
		list = list[0:10]
	}
	// 循环读出ID 向cache获取点赞 评论 数量
	pipe := cache.Cache.Pipeline()
	for _, value := range list {
		pipe.HGetAll(cache.PREFIX_DYNAMIC_COMMENT_DETAILS_COUNT + strconv.FormatInt(value.ID, 10))
	}
	val, err := pipe.Exec()
	pipe.Close()
	count := make([]map[string]string, length)
	for i := 0; i < length; i++ {
		count[i] = val[i].(*redis.StringStringMapCmd).Val()
	}

	data := make([]response.ReplyDataRespond, length)
	for index,value := range list {
		value.User.AvatarToB60()
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.User.ID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].RName = value.Receive.Nickname
		data[index].Content = value.Content
		praise,_ := strconv.Atoi(count[index][cache.FIELD_PRISE_COUNT])
		data[index].PraiseCount = praise
		data[index].CreatedAt = value.CreatedAt.Unix()
	}
	// 设置meta  当前页码 和当前结果集数量
	meta.CurrentPage = offset/perPage + 1
	meta.PerPage = length
	rsp.Success("", gin.H{
		"list":  data,
		"links": links,
		"meta":  meta,
	})

}

func ReplyPublish(c *gin.Context) {
	rsp := response.NewResponse(c)
	key, _ := c.Get("claims")
	claims, _ := key.(*oauth.Claims)

	var entry input.Reply
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	commentID,err :=encrypt.Decode(entry.CommentID)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var comment models.DynamicComment
	var reply models.DynamicComment
	if err := models.GetDynamicCommentInfo(commentID, &comment); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	// 屏蔽当中的铭感字
	if err := sensitiveWord.SensitiveWordReplace(&entry.Content);err != nil {
		rsp.Error()
		return
	}
	reply.UID = claims.UID
	reply.DynamicID = comment.DynamicID
	reply.ReceiveID = comment.UID
	reply.Content = entry.Content
	reply.ParentID = commentID
	if comment.ParentID != 0 {
		reply.ParentID = comment.ParentID
	}
	if err := reply.Create(); err != nil {
		rsp.Error()
		return
	}
	// 动态评论数 +1  动态评论的回复 +1
	pipe := cache.Cache.Pipeline()
	pipe.HIncrBy(cache.PREFIX_DYNAMIC_DETAILS_COUNT+strconv.FormatInt(reply.DynamicID, 10),
		cache.FIELD_COMMENT_COUNT, 1)
	pipe.HIncrBy(cache.PREFIX_DYNAMIC_COMMENT_DETAILS_COUNT+strconv.FormatInt(commentID, 10),
		cache.FIELD_REPLY_COUNT, 1)
	pipe.Exec()
	pipe.Close()

	// TODO 向用户发送被回复消息
	rsp.Success("回复成功", nil)
}

func ReplyListV2(c *gin.Context) {
	var pageCount = setting.AppSetting.PageSize
	rsp := response.NewResponse(c)
	id := c.Query("id")
	commentID,err := encrypt.Decode(id)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	replyID := c.Query("reply_id")
	var lastID int64
	if replyID != "" {
		var err error
		lastID, err = encrypt.Decode(replyID)
		if err != nil {
			rsp.ErrorMsg(e.INVALID_PARAMS)
			return
		}
	}
	var list = make([]models.DynamicComment,pageCount)
	if replyID != "" {
		if err := models.GetMoreDynamicCommentReplyListOfID(commentID,lastID,pageCount,&list); err != nil {
			rsp.Error()
			return
		}
	} else {
		if err := models.GetDynamicCommentReplyListOfID(commentID,pageCount,&list); err != nil {
			rsp.Error()
			return
		}
	}
	length := len(list)
	data := make([]response.ReplyDataRespond, length)
	for index,value := range list {
		value.User.AvatarToB60()
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.User.ID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].RName = value.Receive.Nickname
		data[index].Content = value.Content
		data[index].CreatedAt = value.CreatedAt.Unix()
	}
	rsp.Success("", data)
}