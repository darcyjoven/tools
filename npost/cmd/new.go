/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"npost/new"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "新建一篇文章",
		Long: `新建篇文章
	--destination 目的文件
	--post 文章名称`,
		Run: func(cmd *cobra.Command, args []string) {
			// fmt.Println("new called")
			new.Exc(*destination, *post)
		},
	}
	destination, post *string
)

func init() {
	destination = newCmd.Flags().StringP("destination", "d", "", `目的文件位置 类似‘book-docs/book1/section2’,其中book对于配置中的多语言名称，docs/book1/section2是项目下的文件夹。`)
	post = newCmd.Flags().StringP("post", "p", "demo test", "文章名称，创建文件夹时，空格将被替换为'-'。")
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
