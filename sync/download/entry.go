package download

import (
	"sync/global"

	"go.uber.org/zap"
)

// 下载
func Exc(file *string) {
	err := Down(*file)
	if err != nil {
		global.L.Error(
			"download file failed",
			zap.String("file", *file),
			zap.Error(err),
		)
		return
	}
	global.L.Info(
		"download file success",
		zap.String("file", *file),
	)

}
