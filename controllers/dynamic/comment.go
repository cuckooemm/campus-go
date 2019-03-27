package dynamic

import (
	"campus/helper/encrypt"
	"campus/helper/oauth"
	"campus/helper/push"
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

func CommentList(c *gin.Context) {
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
	dynamicID,err := encrypt.Decode(id)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var list []models.DynamicComment
	var links response.Links
	var meta response.Meta
	//如果不是第一页  设置links 的prev
	if offset > 1 {
		links.Prev = setting.AppSetting.ServerApiHost + `v1/dynamic/comment?page=` + strconv.Itoa(offset-1)  + `&id=` + id
	}
	// 根据页码设置偏移量
	offset = (offset - 1) * perPage
	if err := models.GetDynamicCommentList(dynamicID, offset, perPage, &list); err != nil {
		rsp.Error()
		return
	}

	// 根据返回结果定义data 数组的长度
	length := len(list)
	if length > perPage {
		// 说明有值 添加下一页的连接
		links.Next = setting.AppSetting.ServerApiHost + `v1/dynamic/comment?page=` + strconv.Itoa(offset/perPage+2)  + `&id=` + id
		// 等于11时减去1 回复到需要数组的长度
		length -= 1
		// 切片数组，只去 10条数据
		list = list[0:10]
	}
	// 循环读出ID 向cache获取浏览量点评论等信息
	pipe := cache.Cache.Pipeline()
	for _, value := range list {
		pipe.HGetAll(cache.PREFIX_DYNAMIC_COMMENT_DETAILS_COUNT + strconv.FormatInt(value.ID, 10))
	}
	val, err := pipe.Exec()
	pipe.Close()
	// 添加一条浏览
	cache.Cache.HIncrBy(cache.PREFIX_DYNAMIC_DETAILS_COUNT+strconv.FormatInt(dynamicID, 10), cache.FIELD_BROWSE_COUNT, 1)
	count := make([]map[string]string, length)
	for i := 0; i < length; i++ {
		count[i] = val[i].(*redis.StringStringMapCmd).Val()
	}
	data := make([]response.DynamicCommentData, length)
	for index, value := range list {
		value.User.AvatarToB60()
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.User.ID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].Content = value.Content
		praise, _ := strconv.Atoi(count[index][cache.FIELD_PRISE_COUNT])
		data[index].PraiseCount = praise
		reply, _ := strconv.Atoi(count[index][cache.FIELD_REPLY_COUNT])
		data[index].ReplyCount = reply
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

// 发布评论
func CommentPublish(c *gin.Context) {
	rsp := response.NewResponse(c)
	key, _ := c.Get("claims")
	claims, _ := key.(*oauth.Claims)
	var entry input.Comment
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	dynamicID,err := encrypt.Decode(entry.DynamicID)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var dynamic models.Dynamic
	// 向数据库查询 回复的动态是否存在 并返回 发布者的ID
	if err := models.GetDynamicUid(dynamicID, &dynamic); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	// 屏蔽当中的铭感字
	if err := sensitiveWord.SensitiveWordReplace(&entry.Content);err != nil {
		rsp.Error()
		return
	}
	var comment models.DynamicComment
	comment.UID = claims.UID
	comment.ReceiveID = dynamic.UID
	comment.Content = entry.Content
	comment.DynamicID = dynamicID
	if err := comment.Create(); err != nil {
		rsp.Error()
		return
	}
	go func(uid int64,content string) {
		push.PushCommentToAlias(uid,"您有一条新评论","评论",content)
	}(comment.ReceiveID,comment.Content)
	//向redis 存储当前动态的回复数量
	cache.Cache.HIncrBy(cache.PREFIX_DYNAMIC_DETAILS_COUNT+strconv.FormatInt(dynamicID, 10), cache.FIELD_COMMENT_COUNT, 1)

	rsp.Success("评论成功", nil)
}

func CommentListV2(c *gin.Context) {
	var pageCount = setting.AppSetting.PageSize
	rsp := response.NewResponse(c)
	id := c.Query("id")
	dynamicID,err := encrypt.Decode(id)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	commentID := c.Query("comment_id")
	var lastID int64
	if commentID != "" {
		var err error
		lastID, err = encrypt.Decode(commentID)
		if err != nil {
			rsp.ErrorMsg(e.INVALID_PARAMS)
			return
		}
	}
	var list = make([]models.DynamicComment,pageCount)
	if commentID != "" {
		if err := models.GetMoreDynamicCommentListOfID(dynamicID,lastID,pageCount,  &list); err != nil {
			rsp.Error()
			return
		}
	} else {
		if err := models.GetDynamicCommentListOfID(dynamicID,pageCount, &list); err != nil {
			rsp.Error()
			return
		}
	}
	length := len(list)
	// 循环读出ID 向cache获取浏览量点评论等信息
	pipe := cache.Cache.Pipeline()
	for _, value := range list {
		pipe.HGetAll(cache.PREFIX_DYNAMIC_COMMENT_DETAILS_COUNT + strconv.FormatInt(value.ID, 10))
	}
	val, err := pipe.Exec()
	pipe.Close()
	// 添加一条浏览
	cache.Cache.HIncrBy(cache.PREFIX_DYNAMIC_DETAILS_COUNT+strconv.FormatInt(dynamicID, 10), cache.FIELD_BROWSE_COUNT, 1)
	count := make([]map[string]string, length)
	if err == nil {
		for i := 0; i < length; i++ {
			count[i] = val[i].(*redis.StringStringMapCmd).Val()
		}
	}
	data := make([]response.DynamicCommentData, length)
	for index, value := range list {
		value.User.AvatarToB60()
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.User.ID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].Content = value.Content
		praise, _ := strconv.Atoi(count[index][cache.FIELD_PRISE_COUNT])
		data[index].PraiseCount = praise
		reply, _ := strconv.Atoi(count[index][cache.FIELD_REPLY_COUNT])
		data[index].ReplyCount = reply
		data[index].CreatedAt = value.CreatedAt.Unix()
	}
	rsp.Success("", data)
}