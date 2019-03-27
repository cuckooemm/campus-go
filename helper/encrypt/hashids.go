package encrypt

import (
	"campus/helper/logging"
	"errors"
	"github.com/speps/go-hashids"
)

var h *hashids.HashID

func Setup()  {
	hd := hashids.NewData()
	hd.Salt = "chuangqianmingyueguangyishidishangshuang"
	hd.MinLength = 18
	var err error
	h,err = hashids.NewWithData(hd)
	if err != nil {
		logging.ErrorMsg("hashids 初始化失败",err)
		return
	}
}

func Encode(data int64) string {
	id,_ := h.EncodeInt64([]int64{data})
	return id
}

func Decode(data string) (int64,error) {
	id,err := h.DecodeInt64WithError(data)
	if err != nil || len(id) == 0{
		return 0,errors.New("id 解码出错")
	}
	return id[0], err
}