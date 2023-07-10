package download

import (
	"fmt"
	"os"
	"sync/global"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 上传文件
func Up(file string) (err error) {
	// 1.file 是否和localdir匹配上
	if !checkDir(file) {
		return fmt.Errorf("文件不在配置文件localdir中")
	}
	// 2.file 是否需要保存
	if err := beforeUp(file); err != nil {
		return err
	}
	// 3. 远程文件下载到tmp目录
	err = down(getRemoteDir(file), getTempDir(file))
	if err == nil {
		// 4. 本地文件是否与远程文件有差异
		ok, err := diff(file, getTempDir(file))
		if !ok {
			// 无差异，退出
			global.L.Warn(
				"has no diff with remote file",
				zap.String("file", file),
			)
			return nil
		}
		if err != nil {
			return err
		}
		// 5. 远程文件备份（将temp file 上传到remote）
		// todo: 检查是否已存在同名文件
		ok = isBakExist(file)
		if ok {
			// 存在备份文件，不需要备份
			global.L.Info(
				"bakfile is exist,no bakup again",
				zap.String("bakup file", file),
			)
		} else {
			err = up(getRemoteBak(file), getTempDir(file))
			if err != nil {
				return err
			}
		}
	} else {
		return err
	}
	defer func() {
		os.Remove(getTempDir(file))
	}()

	// 6. 本地文件上传到远程
	return up(getRemoteDir(file), file)
}

// 下载前提示
func beforeUp(file string) (err error) {
	// 是否有diff
	ok, err := checkDiff(file)
	if !ok {
		// 无差异
		return nil
	}
	if err != nil {
		return err
	}
	level := viper.GetString("diffupload")
	switch level {
	case "warn":
		var ct string
		fmt.Print("文件与git仓库有差异记录，输入n/N取消上传，输入其它键继续:")
		fmt.Scan(&ct)
		if ct == "N" || ct == "n" {
			global.L.Info(
				"用户已取消操作",
				zap.String("file", file),
			)
			os.Exit(-1)
			// return nil
		}
	case "error":
		return fmt.Errorf("文件与git仓库有差异记录，请提交后再尝试！")
	default:
		// info
		global.L.Info(
			"文件与git仓库有差异记录，但会忽略",
			zap.String("file", file),
		)
	}
	return
}

// 远程是否存在备份文件
func isBakExist(file string) (ok bool) {
	localDir := getLocalBak(file)
	err := down(getRemoteBak(file), localDir)
	// todo: err 是否是文件不存在？
	if err != nil {
		return false
	}
	if _, err = os.Stat(localDir); err != nil {
		return false
	}
	return true
}
