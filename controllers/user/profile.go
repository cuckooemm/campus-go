package user

import (
	"campus/helper/check"
	"campus/helper/logging"
	"campus/helper/oauth"
	"campus/helper/oss"
	"campus/helper/sensitiveWord"
	"campus/helper/setting"
	"campus/models"
	"campus/models/input"
	"campus/response"
	"campus/response/e"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Profile(c *gin.Context) {
	rsp := response.NewResponse(c)
	key, _ := c.Get("claims")
	claims, _ := key.(*oauth.Claims)

	var entry input.UserProfile
	// 设置此字段 用于标记是否修改的是头像
	entry.AvatarStatus = false
	if err := c.ShouldBind(&entry); err != nil {
		rsp.ErrorMsg(e.INVALID_PARAMS)
		return
	}
	var user models.User
	if err := models.GetUserProfile(claims.UID,&user);err != nil {
		rsp.Error()
		logging.DatabaseError(err)
		return
	}
	userInfo := make(map[string]interface{})
	if entry.Nickname != "" {
		if sensitiveWord.SensitiveWordFilter(entry.Nickname) {
			rsp.ErrorMsg(e.ERROR_NICKNAME_EXIST)
			return
		}
		entry.Nickname = models.ReplaceName(entry.Nickname)
		// 向数据库发送查询请求
		if models.IsNameExist(entry.Nickname) {
			rsp.ErrorMsg(e.ERROR_NICKNAME_EXIST)
			return
		}
		userInfo["nickname"] = entry.Nickname
	}
	if entry.Avatar != "" {
		if !oss.Find(entry.Avatar) {
			rsp.ErrorMsg(e.ERROR_AVATAR_UP_FAIL)
			return
		}
		// 对比头像是否为默认的头像
		if user.Avatar != "avatar/default" && models.IsAvatarOss(user.Avatar) {
			// 开启一个协程 删除图片
			go func(object []string) {
				oss.Detele(object)
			}([]string{user.Avatar})
		}
		userInfo["avatar"] = entry.Avatar
		entry.AvatarStatus = true
	}
	if entry.Gender != "" {
		gender, err := strconv.Atoi(entry.Gender)
		if err != nil {
			rsp.ErrorMsg(e.INVALID_PARAMS)
			return
		}
		if gender < 0 && gender > 2 {
			rsp.ErrorMsg(e.INVALID_PARAMS)
			return
		}
		userInfo["gender"] = gender
	}
	if entry.Birthday != "" {
		if err := check.IsBirthday(entry.Birthday); err != nil {
			rsp.ErrorMsg(e.INVALID_PARAMS)
			return
		}
		userInfo["birthday"] = entry.Birthday
	}
	if entry.Bio != "" {
		if err := sensitiveWord.SensitiveWordReplace(&entry.Bio);err != nil {
			rsp.Error()
			return
		}
		userInfo["bio"] = entry.Bio
	}
	if err := models.UpdateProfile(claims.UID, userInfo); err != nil {
		logging.DatabaseError(err)
		rsp.Error()
		return
	}
	if entry.AvatarStatus {
		avatar := setting.OSSSetting.Path + entry.Avatar + oss.IMAGE_ORIGINAL
		rsp.Success("头像修改成功",avatar)
		return
	}
	rsp.Success("更新成功",nil)

}
