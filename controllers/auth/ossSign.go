package auth

import (
	"campus/helper/oss"
	"campus/response"
	"github.com/gin-gonic/gin"
)

func OSSSign(c *gin.Context) {
	rsp := response.NewResponse(c)
	var token oss.PolicyToken
	if err := oss.GetPolicyToken(&token);err != nil{
		rsp.Error()
		return
	}
	rsp.Success("",token)
}