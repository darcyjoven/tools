package new

import (
	"npost/global"
	"strings"

	"go.uber.org/zap"
)

// 新建一篇文章
func Exc(destination, post string) {
	// 3. 依据destination创建文件/文件夹
	// 4. 将模板中的变量解析并重写
	// 1. 查看项目文件夹是否存在 blog/book
	project, dest := splitDest(destination)
	if project == "" || dest == "" {
		global.L.Error(
			"can not get the project&destination",
			zap.String("destination", destination),
		)
	}
	if !check(project) {
		global.L.Error(
			"no project,or Dir is not exist",
			zap.String("project", project),
			zap.String("Dir", projectDir(project)),
		)
	}
	// 2. 查看项目下的模板是否能匹配上destination
	_, err := setLanuage(project, dest, post)
	if err != nil {
		global.L.Error(
			"set language falied",
			zap.Error(err),
		)
		return
	}
}

// splitDest 将destination解析为project和dest
func splitDest(destination string) (project, dest string) {
	i := strings.Index(destination, "-")
	if i == -1 || i == len(destination)-1 {
		return "", ""
	}
	return destination[:i], destination[i+1:]
}
