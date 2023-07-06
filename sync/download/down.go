package download

import (
	"fmt"
	"os"
	"sync/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 文件下载
func Down(file string) (err error) {
	// 1.file 是否和localdir匹配上
	// if !checkDir(file) {
	// 	return fmt.Errorf("文件不在配置文件localdir中")
	// }
	// 2.file 是否需要保存
	if err := beforeDown(file); err != nil {
		return err
	}
	// 3. 覆盖本地文件
	err = down(getRemoteDir(file), file)
	if err != nil {
		return err
	}
	// 4. git commit
	return gitCommit(file, getComment(file))
}

// 下载前提示
func beforeDown(file string) (err error) {
	// 是否有diff
	ok, err := checkDiff(file)
	if !ok {
		// 无差异
		return nil
	}
	if err != nil {
		return err
	}
	level := viper.GetString("diffdownload")
	switch level {
	case "info":
		global.L.Info(
			"文件与git仓库有差异记录，但会忽略",
			zap.String("file", file),
		)
	case "error":
		return fmt.Errorf("文件与git仓库有差异记录，请提交后再尝试！")
	default:
		// warn
		var ct string
		fmt.Print("文件与git仓库有差异记录，输入n/N取消下载，输入其它键继续:")
		fmt.Scan(&ct)
		if ct == "N" || ct == "n" {
			global.L.Info(
				"用户已取消操作",
				zap.String("file", file),
			)
			os.Exit(-1)
			// return nil
		}
	}
	return
}
