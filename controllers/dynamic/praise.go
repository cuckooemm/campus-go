package dynamic

import (
	"campus/helper/encrypt"
	"campus/helper/oauth"
	"campus/helper/redis/cache"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"strconv"
)

func PraisePublish(c *gin.Context) {
	rsp := response.NewResponse(c)
	key, _ := c.Get("claims")
	claims, _ := key.(*oauth.Claims)
	var entry input.Praise
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	ID,err := encrypt.Decode(entry.ID)
	if err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var praise models.DynamicPraise
	praise.DynamicID = ID
	praise.UID = claims.UID
	if praise.DynamicID != 0 {
		if ok := praise.IsUserPraiseD(); !ok {
			rsp.ErrorMsg(e.ERROR_PRAISEED)
			return
		}
		if err := praise.Create(); err != nil {
			rsp.ErrorMsg(e.ERROR_PRAISE_FAIL)
			return
		}
		cache.Cache.HIncrBy(cache.PREFIX_DYNAMIC_DETAILS_COUNT+strconv.FormatInt(praise.DynamicID, 10),
			cache.FIELD_PRISE_COUNT, 1)
		rsp.Success("点赞成功", nil)
		return
	}
/*
	if praise.CommentID != 0 {
		if ok := praise.IsUserPraiseC(); !ok {
			rsp.ErrorMsg(e.ERROR_PRAISEED)
			return
		}
		if err := praise.Create(); err != nil {
			rsp.ErrorMsg(e.ERROR_PRAISE_FAIL)
			return
		}
		cache.Cache.HIncrBy(cache.PREFIX_DYNAMIC_COMMENT_DETAILS_COUNT+strconv.FormatInt(praise.CommentID, 10),
			cache.FIELD_PRISE_COUNT, 1)
		rsp.Success("点赞成功", nil)
		return
	}
*/
	rsp.ErrorMsg(e.INVALID_PARAMS)
}
