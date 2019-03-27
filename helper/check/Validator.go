package check

import (
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

var Validate *validator.Validate

func ValidatorSetup() {
	Validate = validator.New()
}

func IsEmail(email string) error {
	err := Validate.Var(email, "email")
	if err == nil {
		return nil
	}
	return errors.New("不是有效的邮箱")
}

func IsPhone(phone string) error {
	if b, _ := regexp.MatchString(`^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\d{8}$`, phone); b {
		return nil
	} else {
		return errors.New("不是有效的手机号")
	}
}
func IsBirthday(birthday string) error {
	if m, _ := regexp.MatchString(`^(19|20)\d{2}-(1[0-2]|0?[1-9])-(0?[1-9]|[1-2][0-9]|3[0-1])$`, birthday); m {
		return nil
	} else {
		return errors.New("不是有效日期")
	}
}