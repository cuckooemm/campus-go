package controllers

import (
	"campus/helper/logging"
	"campus/helper/push"
	"campus/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"time"
)

func Test(c *gin.Context) {
	//insertDynamic()
	//insertDynamicComment()
	//send()
}

func send()  {
	body := push.AndroidBody{Ticker:"ces ",Title: "ces title",Text: "ces message", AfterOpen:"go_app",Activity: "cn.campuswall.campus.MainActivity"}
	message := push.AndroidBody{AfterOpen:"go_custom",Custom:gin.H{
		"key1": "key1",
		"key2": 151572121,
		"key3":true,
	}}
	var id int64= 1539420473786777600
	push.PushAliasNotification(&id,&body)
	push.PushAliasMessage(&id,&message)
}

func insertDynamicComment()  {
	var comment models.DynamicComment
	comment.UID = 1544351957333196800
	comment.ReceiveID = 1544351957333196800
	comment.DynamicID = 1543393152340475904
	for i := 0; i < 50; i++ {
		comment.Content = "测试评论数据" + strconv.Itoa(i)
		comment.Create()
	}
}
func insertDynamic() {
	var dynamic models.Dynamic
	// 上次插入到 I条 数据
	dynamic.UID = 1543141806148304896
	img := [6]string {
		"dynamic/cd5b9377-d223-45f1-ae00-987273043ac7,dynamic/ac1d49f1-0818-48e1-807c-0e19dcfc57e9",
		"dynamic/6756f41d-293a-4524-a4c8-ecd3e2c36caa,dynamic/fc722992-218a-4f39-aefc-b5f00b84e4b0,dynamic/66a4a157-9dd5-4e88-aa39-53482dfb3b63,dynamic/36fd59b4-54af-40cf-adf0-56e5624a2d6b,dynamic/1a7dbfdc-87ff-4c17-a4a9-78a844d0c4dc",
		"dynamic/1e9cbee0-8547-46c4-8c31-298562a1fcbd,dynamic/43449378-b5dd-4fc8-a280-762b962beb74,dynamic/680169d3-6f5d-4ade-b6de-f6a027f48caf,dynamic/359b0a9a-9f86-443c-905b-65ffaf7ac7a8,dynamic/7c8b3a3e-83a7-4540-9d5a-f3b0a7fe9879,dynamic/40384ee3-0528-44ed-a1fc-00b1acf717ec,dynamic/6cd9ed8c-2507-4d5f-afc0-a9ca0dcce708,dynamic/dd667fdd-a4dd-498a-b9d5-4bb558bc24d0,dynamic/3bf9e712-6f48-4b42-ab7d-fd66d3903c51",
		"dynamic/e3d1a19d-34f6-4cb0-bc9e-99b3e893a4c3,dynamic/b67475cd-74ef-4de5-980c-a7d6e0efb417",
		"dynamic/376b692a-eb64-4cf7-963e-1bd732eaadc9",
		"dynamic/10fe9e36-339e-4a23-a298-dbf45f2b2dda",
	}
	for i := 0; i < 200000; i++ {
		dynamic.Content = "测试数据 " + strconv.Itoa(i)
		dynamic.Images = img[rand.Intn(5)]
		dynamic.Create()
		time.Sleep(1 * time.Millisecond)
	}
}
func selectlimitoffset() {

	var list []models.Dynamic
	up := time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		models.GetDynamicList(0, 10, &list)
	}
	logging.Warn("1000 次循环所用时间", zap.Int64("nano", time.Now().UnixNano()-up))
	println(time.Now().UnixNano() - up)

	up = time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		models.GetDynamicList(0, 100, &list)
	}
	logging.Warn("1000 limit 100 次循环所用时间", zap.Int64("nano", time.Now().UnixNano()-up))
	println(time.Now().UnixNano() - up)

	up = time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		models.GetDynamicList(1000, 10, &list)
	}
	logging.Warn("1000次循环 offset 1000 所用时间", zap.Int64("nano", time.Now().UnixNano()-up))
	println(time.Now().UnixNano() - up)

	up = time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		models.GetDynamicList(10000, 10, &list)
	}
	logging.Warn("1000次循环 offset 10000 所用时间", zap.Int64("nano", time.Now().UnixNano()-up))
	println(time.Now().UnixNano() - up)

	up = time.Now().UnixNano()
	for i := 0; i < 1000; i++ {
		models.GetDynamicList(40000, 10, &list)
	}
	logging.Warn("1000次循环 offset 40000 所用时间", zap.Int64("nano", time.Now().UnixNano()-up))
	println(time.Now().UnixNano() - up)

}
