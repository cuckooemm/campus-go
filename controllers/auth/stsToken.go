package auth

import (
	"campus/helper/logging"
	"campus/helper/setting"
	"campus/response"
	"github.com/denverdino/aliyungo/sts"
	"github.com/gin-gonic/gin"
	"time"
)


func STSToken(c *gin.Context) {
	rsp := response.NewResponse(c)
	client := sts.NewClient(setting.STSSetting.AccesskeyID, setting.STSSetting.AccesskeySecret)
	client.SetDebug(setting.STSSetting.Debug)
	role := sts.AssumeRoleRequest{RoleArn: setting.STSSetting.RoleArn, RoleSessionName:setting.STSSetting.RoleSessionName, DurationSeconds: setting.STSSetting.DurationSeconds}
	token, err := client.AssumeRole(role)
	if err != nil {
		logging.WarnMsg("sts token 生成失败", err)
		rsp.Error()
		return
	}
	t,_ := time.Parse("206-01-02T15:04:05Z",token.Credentials.Expiration)
	rsp.Success("", gin.H{
		"AccessKeyId":     token.Credentials.AccessKeyId,
		"AccessKeySecret": token.Credentials.AccessKeySecret,
		"SecurityToken":   token.Credentials.SecurityToken,
		"Expiration":      t.Unix(),
	})

}
