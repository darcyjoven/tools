package upload

import (
	"sync/download"
	"sync/global"

	"go.uber.org/zap"
)

func Exc(file *string) {
	err := download.Up(*file)
	if err != nil {
		global.L.Error(
			"upload failed",
			zap.String("file", *file),
			zap.Error(err),
		)
	}
}
