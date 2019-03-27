package dynamic

import (
	"bytes"
	"campus/helper/encrypt"
	"campus/helper/oss"
	"campus/helper/redis/cache"
	"campus/helper/setting"
	"campus/models"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"strconv"
	"strings"
)

func List(c *gin.Context) {
	var perPage = setting.AppSetting.PageSize
	rsp := response.NewResponse(c)
	// 获取分页加载的页数
	pg := c.DefaultQuery("page", "1")
	offset, err := strconv.Atoi(pg)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var list []models.Dynamic
	var links response.Links
	var meta response.Meta
	//如果不是第一页  设置links 的prev
	if offset > 1 {
		links.Prev = setting.AppSetting.ServerApiHost + `v1/dynamic?page=` + strconv.Itoa(offset-1)
	}
	// 根据页码设置偏移量
	offset = (offset - 1) * perPage
	if err := models.GetDynamicList(offset, perPage+1, &list); err != nil {
		rsp.Error()
		return
	}

	// 根据返回结果定义data 数组的长度
	length := len(list)
	if length == perPage+1 {
		// 说明有值 添加下一页的连接
		links.Next = setting.AppSetting.ServerApiHost + `v1/dynamic?page=` + strconv.Itoa(offset/perPage+2)
		// 等于11时减去1 回复到需要数组的长度
		length -= 1
		// 切片数组，只去 10条数据
		list = list[0:10]
	}
	// 循环读出ID 向cache获取浏览量点评论等信息
	pipe := cache.Cache.Pipeline()
	for _, value := range list {
		pipe.HGetAll(cache.PREFIX_DYNAMIC_DETAILS_COUNT + strconv.FormatInt(value.ID, 10))
	}
	val, err := pipe.Exec()
	pipe.Close()
	count := make([]map[string]string, length)
	for i := 0; i < length; i++ {
		count[i] = val[i].(*redis.StringStringMapCmd).Val()
	}
	data := make([]response.DynamicData, length)
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	for index, value := range list {
		value.User.AvatarToB60()
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.User.ID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].Content = value.Content
		if value.Images != "" {
			image := strings.Split(strings.TrimRight(value.Images, ","), ",")
			var images response.Image
			if len(image) == 1 {
				buf.WriteString(setting.OSSSetting.Path)
				buf.WriteString(image[0])
				buf.WriteString(oss.IMAGE_NINE_S_600_400)
				buf.WriteString(",")
				buf2.WriteString(setting.OSSSetting.Path)
				buf2.WriteString(image[0])
				buf2.WriteString(oss.IMAGE_ORIGINAL)
				buf2.WriteString(",")
			} else {
				for _, value := range image {
					buf.WriteString(setting.OSSSetting.Path)
					buf.WriteString(value)
					buf.WriteString(oss.IMAGE_NINE_B_250)
					buf.WriteString(",")
					buf2.WriteString(setting.OSSSetting.Path)
					buf2.WriteString(value)
					buf2.WriteString(oss.IMAGE_ORIGINAL)
					buf2.WriteString(",")
				}
			}
			images.Url = strings.Split(strings.TrimRight(buf.String(), ","), ",")
			images.UrlOriginal = strings.Split(strings.TrimRight(buf2.String(), ","), ",")
			data[index].Image = &images
			buf.Reset()
			buf2.Reset()
		}
		browse, _ := strconv.Atoi(count[index][cache.FIELD_BROWSE_COUNT])
		data[index].Browse = browse
		comment, _ := strconv.Atoi(count[index][cache.FIELD_COMMENT_COUNT])
		data[index].CommentCount = comment
		praise, _ := strconv.Atoi(count[index][cache.FIELD_PRISE_COUNT])
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

func ListV2(c *gin.Context) {
	var pageCount = setting.AppSetting.PageSize
	rsp := response.NewResponse(c)
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
	var list = make([]models.Dynamic,pageCount)
	if id != "" {
		if err := models.GetMoreDynamicListOfID(pageCount, lastID, &list); err != nil {
			rsp.Error()
			return
		}
	} else {
		if err := models.GetDynamicListOfID(pageCount, &list); err != nil {
			rsp.Error()
			return
		}
	}
	length := len(list)
	// 循环读出ID 向cache获取浏览量点评论等信息
	pipe := cache.Cache.Pipeline()
	for _, value := range list {
		pipe.HGetAll(cache.PREFIX_DYNAMIC_DETAILS_COUNT + strconv.FormatInt(value.ID, 10))
	}
	val, err := pipe.Exec()
	pipe.Close()
	count := make([]map[string]string, length)
	// redis 没有错误的时候才去获取数据  有错误则忽略
	if err == nil{
		for i := 0; i < length; i++ {
			count[i] = val[i].(*redis.StringStringMapCmd).Val()
		}
	}

	data := make([]response.DynamicData, length)
	var buf bytes.Buffer
	var buf2 bytes.Buffer
	for index, value := range list {
		value.User.AvatarToB60()
		data[index].ID = encrypt.Encode(value.ID)
		data[index].UID = encrypt.Encode(value.User.ID)
		data[index].Name = value.User.Nickname
		data[index].Avatar = value.User.Avatar
		data[index].Gender = value.User.GenderToString()
		data[index].Content = value.Content
		if value.Images != "" {
			image := strings.Split(strings.TrimRight(value.Images, ","), ",")
			var images response.Image
			if len(image) == 1 {
				buf.WriteString(setting.OSSSetting.Path)
				buf.WriteString(image[0])
				buf.WriteString(oss.IMAGE_NINE_S_600_400)
				buf.WriteString(",")
				buf2.WriteString(setting.OSSSetting.Path)
				buf2.WriteString(image[0])
				buf2.WriteString(oss.IMAGE_ORIGINAL)
				buf2.WriteString(",")
			} else {
				for _, value := range image {
					buf.WriteString(setting.OSSSetting.Path)
					buf.WriteString(value)
					buf.WriteString(oss.IMAGE_NINE_B_250)
					buf.WriteString(",")
					buf2.WriteString(setting.OSSSetting.Path)
					buf2.WriteString(value)
					buf2.WriteString(oss.IMAGE_ORIGINAL)
					buf2.WriteString(",")
				}
			}
			images.Url = strings.Split(strings.TrimRight(buf.String(), ","), ",")
			images.UrlOriginal = strings.Split(strings.TrimRight(buf2.String(), ","), ",")
			data[index].Image = &images
			buf.Reset()
			buf2.Reset()
		}
		browse, _ := strconv.Atoi(count[index][cache.FIELD_BROWSE_COUNT])
		data[index].Browse = browse
		comment, _ := strconv.Atoi(count[index][cache.FIELD_COMMENT_COUNT])
		data[index].CommentCount = comment
		praise, _ := strconv.Atoi(count[index][cache.FIELD_PRISE_COUNT])
		data[index].PraiseCount = praise
		data[index].CreatedAt = value.CreatedAt.Unix()
	}
	rsp.Success("", data)
}
