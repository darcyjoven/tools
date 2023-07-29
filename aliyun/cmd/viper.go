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
		viper.SetConfigName(".aliyun")
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
	viper.SetDefault("loginterval", "day")                  // 日志名称

	viper.SetDefault("domain", "baidu.com")
	viper.SetDefault("rr", "nas")
	viper.SetDefault("accesskeyid", "xx")
	viper.SetDefault("accesskeysecret", "xx")
	viper.SetDefault("lastip", "192.168.x.x")
}
