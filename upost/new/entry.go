package new

import (
	"npost/global"

	"go.uber.org/zap"
)

// 新建一篇文章
func Exc(destination, post string) {
	global.L.Info(
		"new a post",
		zap.String("destination", destination),
		zap.String("post", post),
	)
}
