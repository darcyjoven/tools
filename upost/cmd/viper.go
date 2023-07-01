package cmd

import (
	"fmt"
	"npost/structure"
	"os"

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

		// Search config in home directory with name "npost" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("npost")
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
	viper.SetDefault("logdir", "./temp")   // 日志目录
	viper.SetDefault("logname", "newpost") // 日志名称
	viper.SetDefault("loginterval", "one") // 日志名称

	// 多语言配置
	languages := structure.LanguagesConfig{}
	languages["blog"] = structure.LanguageConfig{
		Source:      "./blog",
		LanguageDir: false,
		Languages: []structure.Languages{
			{
				Name:        "en",
				Destination: "./content",
				Template:    "./archetypes/",
			},
			{
				Name:        "zh-cn",
				Destination: "./content",
				Template:    "./archetypes/",
			},
		},
	}
	languages["book"] = structure.LanguageConfig{
		Source:      "./book",
		LanguageDir: true,
		Languages: []structure.Languages{
			{
				Name:        "en",
				Destination: "./content.en",
				Template:    "./archetypes/",
			},
			{
				Name:        "zh",
				Destination: "./content.zh",
				Template:    "./archetypes/",
			},
		},
	}

	viper.SetDefault("languages", languages)
}
