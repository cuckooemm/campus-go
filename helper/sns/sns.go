package sns

import (
	"campus/helper/logging"
	"campus/helper/setting"
	"github.com/xuruiray/aliyun-go/sms"
	"go.uber.org/zap"
)
func SendSns(phone, code string) error {
	logging.Info("发送手机验证码",zap.String("phone",phone),zap.String("code",code))
	messageInfo := sms.MessageBody{
		AccessKeyID:     setting.SNSSetting.AccessKeyID,
		AccessKeySecret: setting.SNSSetting.AccessKeySecret,
		PhoneNumbers:    phone,
		SignName:        setting.SNSSetting.SignName,
		TemplateCode:    setting.SNSSetting.TemplateCode,
		TemplateParam:   "{\"code\":\"" + code + "\"}",
	}
	if err := sms.SendMessage(messageInfo); err != nil {
		logging.Error("帐号短信验证码发送失败",zap.String("phone",phone),zap.String("code",code),zap.String("err",err.Error()))
		return err
	}
	return nil
}