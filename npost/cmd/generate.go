/*
Copyright © 2023 darcy joven <darcy_joven@live.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 产生配置文件

// generateCmd represents the generate command
var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "产生默认配置文件",
		Long: `如果你得配置文件修改乱了或不小心被删除，通过此命令产生新的配置文件。
可以指定文件名，默认为./upost.yaml。
`,
		Run: generate,
	}
	generateFile *string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		cobra.CheckErr(err)
	}
	// generate cmd flag
	generateFile = generateCmd.Flags().StringP(
		"file", "f", filepath.Join(home, "./npost.yaml"),
		`指定文件名，默认为./npost.yaml。
允许的后缀名为`+supportExt()+`，如果无后缀名或不在允许范围，会自动增加后缀.yaml。
如果指定目录，未指定文件名，默认文件名为npost.yaml。`)
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// generate 产生配置文件
func generate(cmd *cobra.Command, args []string) {
	ext := filepath.Ext(*generateFile)
	if ext != "" {
		ext = ext[1:]
	} else {
		*generateFile = "./npost.yaml"
		ext = "yaml"
	}
	initViperDefault() // viper默认值设置
	if !stringInSlice(ext, viper.SupportedExts) {
		*generateFile = *generateFile + ".yaml"
		ext = "yaml"
	}
	err := viper.SafeWriteConfigAs(*generateFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "file write failed ", *generateFile, err.Error())
		return
	}
	fmt.Fprintln(os.Stderr, "success generate file :", *generateFile)
}

func supportExt() string {
	return strings.Join(viper.SupportedExts, ",")
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
