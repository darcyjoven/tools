package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".fast" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("sync")
	}
	// initViperDefault()
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		// 配置文件解析成功
		// logger.InitLogger()
		// 初始化日志
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}

// initViperDefault Viper 默认值
func initViperDefault() {
	name := filepath.Base(os.Args[0])
	name = strings.Replace(name, filepath.Ext(name), "", 1) // 取运行程序的名称
	viper.SetDefault("logdir", "./temp")                    // 日志目录
	viper.SetDefault("logname", name)                       // 日志名称
	viper.SetDefault("loginterval", "one")                  // 日志名称

	viper.SetDefault("localdir", "D:/app/apbak/u1/topprod") // 本地目录，未匹配到不处理
	viper.SetDefault("remotedir", "/u1/topprod")            // 远程目录
	viper.SetDefault("diffupload", "info")                  // 上传前diff询问
	viper.SetDefault("diffdownload", "warn")                // 下载前diff询问
	viper.SetDefault("format", "060102")                    // 备份文件流水号
	viper.SetDefault("author", "darcy")                     // 作者，会包括再备份文件后缀中
	viper.SetDefault("tempdir", "D:/app/apbak/u1/tmp")      //临时文件夹
	viper.SetDefault("remotestr", "tiptop@192.168.1.19")    //远程连接方式
	viper.SetDefault("gitdir", "D:/app/apbak/u1")           //git homeDir
	viper.SetDefault("gitdir", "")                          //提交comment
}
