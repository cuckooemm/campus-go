package email

import (
	"campus/helper/logging"
	"campus/helper/setting"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

const emailSubject  = "验证码"

func SendEmail(email, code string) error {
	logging.Info("发送邮件验证码",zap.String("email",email),zap.String("code",code))
	m := gomail.NewMessage()
	m.SetHeader("From", setting.EmailSetting.Account)
	m.SetHeader("To", email)
	m.SetAddressHeader("Cc", setting.EmailSetting.Account, setting.EmailSetting.SendName)
	//m.SetHeader("Subject", email_account)
	m.SetHeader("Subject", emailSubject)
	m.SetBody("text/html",
		`<table class="wrapper" width="100%" cellpadding="0" cellspacing="0" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; background-color: #f5f8fa; margin: 0; padding: 0; width: 100%; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 100%;">
        <tbody><tr>
            <td align="center" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box;">
                <table class="content" width="100%" cellpadding="0" cellspacing="0" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; margin: 0; padding: 0; width: 100%; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 100%;">
                    <tbody><tr>
    <td class="header" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; padding: 25px 0; text-align: center;">
        <a href="http://campus.dev" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; color: #000000; font-size: 32px; font-weight: bold; text-decoration: none; text-shadow: 0 1px 0 white;" target="_blank">校园墙</a>
    </td>
</tr>

                    
                    <tr>
                        <td class="body" width="100%" cellpadding="0" cellspacing="0" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; background-color: #FFFFFF; border-bottom: 1px solid #EDEFF2; border-top: 1px solid #EDEFF2; margin: 0; padding: 0; width: 100%; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 100%;">
                            <table class="inner-body" align="center" width="570" cellpadding="0" cellspacing="0" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; background-color: #FFFFFF; margin: 0 auto; padding: 0; width: 570px; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 570px;">
                                
                                <tbody><tr>
                                    <td class="content-cell" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; padding: 35px;">
                                        <p style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; color: #74787E; font-size: 16px; line-height: 1.5em; margin-top: 0; text-align: center;">亲爱的用户，您正使用邮箱:<span style="color: #000000">`+email+`</span>获取验证码，请确保是您本人操作，如不是请忽略</p>
<table class="panel" width="100%" cellpadding="0" cellspacing="0" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; margin: 0 0 21px;">
    <tbody><tr>
        <td class="panel-content" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; background-color: #EDEFF2; padding: 16px;">
            <table width="100%" cellpadding="0" cellspacing="0" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box;">
                <tbody><tr>
                    <td class="panel-item" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; padding: 0;">
                        <p style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; color: #74787E; font-size: 20px; line-height: 1.5em; margin-top: 0; text-align: center; margin-bottom: 0; padding-bottom: 0;">验证码：<span style="border-bottom: 1px dashed rgb(204, 204, 204); color: #000000; z-index: 1; position: static;">` + code + `</span></p>
                    </td>
                </tr>
            </tbody></table>
        </td>
    </tr>
</tbody></table>
<pre style="font-family: Avenir, Helvetica, sans-serif; color: #74787E;text-align: right; box-sizing: border-box;">该验证码10分钟内有效</pre>

                                        
                                    </td>
                                </tr>
                            </tbody></table>
                        </td>
                    </tr>

                    <tr>
    <td style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box;">
        <table class="footer" align="center" width="570" cellpadding="0" cellspacing="0" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; margin: 0 auto; padding: 0; text-align: center; width: 570px; -premailer-cellpadding: 0; -premailer-cellspacing: 0; -premailer-width: 570px;">
            <tbody><tr>
                <td class="content-cell" align="center" style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; padding: 35px;">
                    <p style="font-family: Avenir, Helvetica, sans-serif; box-sizing: border-box; line-height: 1.5em; margin-top: 0; color: #AEAEAE; font-size: 12px; text-align: center;">© 2018 校园墙. All rights reserved.</p>
                </td>
            </tr>
        </tbody></table>
    </td>
</tr>
                </tbody></table>
            </td>
        </tr>
    </tbody>
</table>
`)

	d := gomail.NewDialer(setting.EmailSetting.Host, 465, setting.EmailSetting.Account, setting.EmailSetting.Password)
	if err := d.DialAndSend(m); err != nil {
		logging.Error("帐号邮箱验证码发送失败",zap.String("email",email),zap.String("code",code),zap.String("err",err.Error()))
		return err
	}
	return nil
}
