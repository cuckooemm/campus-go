package sensitiveWord

import (
	"campus/helper/logging"
	"github.com/antlinker/go-dirtyfilter"
	"github.com/antlinker/go-dirtyfilter/store"
	"os"
)

var filterManage *filter.DirtyManager

func Setup() {
	var err error
	file,err := os.Open("conf/sensitive.txt")
	if err != nil {
		logging.ErrorMsg("敏感词文件打开失败",err)
		return
	}
	defer file.Close()
	memStore,err := store.NewMemoryStore(store.MemoryConfig{
		Reader:file,
	})
	if err != nil {
		logging.ErrorMsg("敏感词模块加载失败",err)
		return
	}
	filterManage = filter.NewDirtyManager(memStore)
	logging.Info("敏感词模块加载成功")
}
func SensitiveWordReplace(word *string) error {
	result,err := filterManage.Filter().Replace(*word,'*')
	if err != nil{
		logging.ErrorMsg("铭感词模块发生异常",err)
		return err
	}
	*word = result
	return nil
}
// 检测是否包含敏感词
func SensitiveWordFilter(word string) bool {
	result, err := filterManage.Filter().Filter(word)
	if err != nil {
		logging.ErrorMsg("铭感词模块发生异常",err)
		return true
	}
	//不为空则包含敏感词
	if result != nil {
		return true
	}
	return false
}