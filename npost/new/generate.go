package new

import (
	"npost/global"
	"npost/structure"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// generateMarkdown 在目的位置产生文件
func generateMarkdown(genSource structure.PostTemplate) (err error) {
	for _, v := range genSource {
		path := filepath.Join(v.Dest, v.FileName)
		dir := filepath.Dir(path)
		// 目的文件夹为空时自动创建
		if _, err := os.Stat(dir); err != nil {
			err = os.MkdirAll(dir, 0750)
			if err != nil && !os.IsExist(err) {
				return err
			}
		}
		if _, err := os.Stat(dir); err != nil {
			return err
		}
		// 创建文件
		err = os.WriteFile(path, []byte(v.Head), 0660)
		if err != nil {
			return err
		}
		global.L.Info(
			"generate file success",
			zap.String("file path", path),
		)
	}
	return
}
