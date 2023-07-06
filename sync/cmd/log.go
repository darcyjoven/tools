package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sync/global"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLogger() {
	rawJSON := []byte(`{
		"level": "debug",
		"encoding": "console"
	  }`)
	var cfg zap.Config

	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	file := getLogFile()
	if file == "" {
		panic("")
	}
	cfg.OutputPaths = append(cfg.OutputPaths, "stdout", file)
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("-07:00 06-01-02 15:04:05.00")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	global.L = zap.Must(cfg.Build(
		zap.Fields(zap.Int("pid", os.Getpid())),
	))
}

func getLogFile() (file string) {
	name := viper.GetString("logname")
	// 没有取到就设置为运行程序名
	if name == "" {
		name = os.Args[0]
		name = filepath.Base(name)
	}
	file = viper.GetString("logdir")
	// 没有取到设置为当前目录
	if file == "" {
		file = "./"
	}

	// 取上海时间
	local, _ := time.LoadLocation("Asia/Shanghai")
	date := time.Now().In(local)

	interval := viper.GetString("loginterval")
	// 根据时间间隔设置日志名称
	switch interval {
	case "one":
	case "every":
		name = fmt.Sprintf("%s_%s_%d", name, date.Format("2006-01-02_19.54.000"), os.Getpid())
	case "year":
		name = name + "_" + date.Format("2006")
	case "month":
		name = name + "_" + date.Format("2006-01")
	case "week":
		name = name + "_" + date.Format("2006-01-Feb")
	case "day":
		name = name + "_" + date.Format("2006-01-02")
	default:
		name = name + "_" + date.Format("2006-01-02")
	}
	name = name + ".log"
	file = filepath.Join(file, name)

	if _, err := os.Stat(file); err != nil {
		_, err = os.Create(file)
		if err != nil {
			cobra.CheckErr(err)
			return ""
		}
	}
	return file
}
