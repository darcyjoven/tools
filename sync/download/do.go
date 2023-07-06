package download

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync/global"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 从from下载到dest
// 如果dest的文件夹不存在会新建
func down(remote, local string) (err error) {
	// 生成dest文件夹
	if err = genLocalDir(local); err != nil {
		return err
	}
	remote = remoteDir(remote)
	_, err = exec.Command("scp", "-r", remote, local).Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s", e.Stderr)
		} else {
			return err
		}
	}
	global.L.Info(
		"download successfule",
		zap.String("local", local),
		zap.String("remote", remote),
	)
	return
}

// 检查并生成路径中不存在的文件夹
func genLocalDir(dir string) (err error) {
	_, err = os.Stat(filepath.Dir(dir))
	if err != nil {
		err = os.MkdirAll(filepath.Dir(dir), 0750)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}
	return
}

// 上传文件覆盖
func up(remote, local string) (err error) {
	if _, err = os.Stat(local); err != nil {
		return err
	}
	remote = remoteDir(remote)
	_, err = exec.Command("scp", "-r", local, remote).Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s", e.Stderr)
		} else {
			return err
		}
	}
	global.L.Info(
		"upload successfule",
		zap.String("local", local),
		zap.String("remote", remote),
	)
	return
}

// 组合远程路径
func remoteDir(dir string) string {
	return viper.GetString("remotestr") + ":" + dir
}

// 检查文件是否符和配置文件
func checkDir(dir string) (ok bool) {
	// match := viper.GetString("localdir") + string(os.PathSeparator) + "*"
	localdir := viper.GetString("localdir")

	localdir, _ = filepath.Abs(localdir)
	for {
		// file := append(file, filepath.Base(path))
		dir = filepath.Dir(dir)
		dir, _ := filepath.Abs(dir)
		if dir == localdir {
			return true
		}
		if dir == "" {
			return false
		}
	}
}

// 检查本地文件是否diff，利用git检查
func checkDiff(path string) (ok bool, err error) {
	gitDir := viper.GetString("gitdir")
	output, err := exec.Command("git", "-C", gitDir, "diff", path).Output()
	if err != nil {
		return false, err
	}
	if len(output) == 0 {
		return false, nil
	} else {
		return true, err
	}
}

// 提交git信息，只有下载后有差异才需要提交
func gitCommit(path, comment string) (err error) {
	// TODO: diff才需要git
	ok, err := checkDiff(path)
	if !ok {
		global.L.Info(
			"no git diff,do not need commit",
			zap.String("file", path),
		)
		return nil
	}
	if err != nil {
		return err
	}
	//git add
	gitDir := viper.GetString("gitdir")
	_, err = exec.Command("git", "-C", gitDir, "add", path).Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s", e.Stderr)
		} else {
			return err
		}
	}
	// git commit
	_, err = exec.Command("git", "-C", gitDir, "commit", "-m", comment).Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("%s", e.Stderr)
		} else {
			return err
		}
	}
	global.L.Info(
		"git commit",
		zap.String("comment", comment),
	)
	return
}

// 获得comment
func getComment(path string) (coment string) {
	coment = viper.GetString("gitcomment")
	coment = strings.ReplaceAll(coment, "$fullfilename", path)

	localdir := viper.GetString("localdir")
	file := []string{}
	for {
		file = append(file, filepath.Base(path))
		path = filepath.Dir(path)
		p, _ := filepath.Abs(path)
		q, _ := filepath.Abs(localdir)
		if p == q || p == "" {
			break
		}
	}
	sort.Slice(file, func(i, j int) bool {
		return i > j
	})
	coment = strings.ReplaceAll(coment, "$filename", strings.Join(file, "/"))
	return
}

// 比较文件差异
func diff(a, b string) (ok bool, err error) {
	// 创建句柄
	fia, err := os.Open(a)
	if err != nil {
		return true, err
	}
	defer fia.Close()
	fib, err := os.Open(b)
	if err != nil {
		return true, err
	}
	defer fib.Close()

	// 创建 Reader
	ra := bufio.NewReader(fia)
	rb := bufio.NewReader(fib)

	// 每次读取 1024 个字节
	bufa := make([]byte, 1024)
	bufb := make([]byte, 1024)

	for {
		//func (b *Reader) Read(p []byte) (n int, err error) {}
		n, err := ra.Read(bufa)
		if err != nil && err != io.EOF {
			return true, err
		}
		if n == 0 {
			break
		}
		_, err = rb.Read(bufb)
		if err != nil && err != io.EOF {
			return true, err
		}
		if string(bufa) != string(bufb) {
			return true, err
		}
	}
	return false, err
}

// 依据本地文件获取远程文件地址
func getRemoteDir(path string) (str string) {
	localdir := viper.GetString("localdir")
	str = viper.GetString("remotedir")

	file := []string{}
	for {
		file = append(file, filepath.Base(path))
		path = filepath.Dir(path)
		p, _ := filepath.Abs(path)
		q, _ := filepath.Abs(localdir)
		if p == q || p == "" {
			break
		}
	}

	for i := len(file) - 1; i >= 0; i-- {
		str = filepath.Join(str, file[i])
	}
	return
}

// getTempDIr 获得临时文件夹路径
func getTempDir(path string) (str string) {
	localdir := viper.GetString("localdir")
	str = viper.GetString("tempdir")

	file := []string{}
	for {
		file = append(file, filepath.Base(path))
		path = filepath.Dir(path)
		p, _ := filepath.Abs(path)
		q, _ := filepath.Abs(localdir)
		if p == q || p == "" {
			break
		}
	}

	for i := len(file) - 1; i >= 0; i-- {
		str = filepath.Join(str, file[i])
	}
	return
}

// 依据本地文件获取远程文件地址
func getRemoteBak(path string) (str string) {
	str = getRemoteDir(path)
	author := viper.GetString("author")
	format := viper.GetString("format")

	l, _ := time.LoadLocation("Asia/Shanghai")
	str = str + time.Now().In(l).Format("."+author+"."+format)
	return str
}

// 依据本地temp 文件的 备份文件地址
func getLocalBak(path string) (str string) {
	str = getTempDir(path)
	author := viper.GetString("author")
	format := viper.GetString("format")

	l, _ := time.LoadLocation("Asia/Shanghai")
	str = str + time.Now().In(l).Format("."+author+"."+format)
	return str
}
