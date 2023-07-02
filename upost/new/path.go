package new

import (
	"fmt"
	"io/fs"
	"npost/structure"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func check(project string) bool {
	path := projectDir(project)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

// 获取配置中的project目录
func projectDir(project string) (path string) {
	path = viper.GetString("projects." + project + ".source")

	if path == "" {
		n := viper.GetString(
			"projects." + project + ".languages.1.name",
		)
		if n == "" {
			return ""
		} else {
			path = "./" + project
		}
	}
	return path
}

// 设计多语言配置
// dest docs/book/section
// post tiptop is fun
func setLanuage(project, dest, post string) (ok bool, err error) {

	postTemplate := make(structure.PostTemplate)
	temp := "projects." + project + "."

	// 多语言是否是独立文件夹
	languageDir := viper.GetBool(temp + "languagedir")

	// viper.get(temp + "languages")
	// 先找到多语言配置
	m := viper.Get("projects." + project + ".languages")
	ps, ok := m.([]any)
	if !ok {
		return ok, fmt.Errorf("can get right config")
	}
	for i := 0; i < len(ps); i++ {
		temp := "projects." + project + ".languages." + strconv.Itoa(i) + "."
		key := viper.GetString(temp + "name")
		if key == "" {
			// pass，不处理这个部分
			continue
		}
		// 文章文件夹
		destination := viper.GetString(temp + "destination")
		if destination == "" {
			return false, fmt.Errorf(
				"can not get destination config of the project:%s language:%s",
				project,
				key,
			)
		}
		destination = filepath.Join(destination, dest)
		filename, title, tags, err := getTitle(post, key, languageDir)
		if err != nil {
			return false, err
		}
		// dest1 := filepath.Join(project, dest)
		// 模板文件夹
		template := viper.GetString(temp + "template")
		if template == "" {
			// 未获得模板位置，设置为项目目录
			// 前面已经检查过，不需要再检查
			template = viper.GetString("projects." + project + ".source")
		}
		// dir D:/app/github.io/blog/archetypes/
		// desc docs/book1/section1
		template, err = getTemplate(template, dest)
		if err != nil {
			return false, err
		}
		order := viper.GetBool("projects." + project + ".order")
		// 变量解析
		template, err = replaceTemp(template, title, project, dest, tags, order)
		if err != nil {
			return false, err
		}
		//插入资料
		postTemplate[key] = structure.MarkDown{
			Dest:     destination,
			Title:    title,
			FileName: filename,
			Head:     template,
		}
	}
	return
}

// 取模板内容
// dir D:/app/github.io/blog/archetypes/
// desc docs/book1/section1
func getTemplate(dir, dest string) (head string, err error) {
	dest = filepath.Join(dir, dest) + "*"
	// ./archetypes/docs/book1/section1
	defalut, matched := "", ""
	// 取得模板文件名
	filepath.Walk(
		dir,
		func(path string, info fs.FileInfo, err error) error {
			if filepath.Ext(path) == ".md" {
				// default赋值
				if filepath.Base(path) == "default.md" {
					defalut = path
				}
				// 如果匹配到
				if ok, _ := filepath.Match(dest, path); ok {
					matched = path
					// 退出
					return fmt.Errorf("")
				}
			} // 文件名
			return err
		})
	if matched == "" {
		matched = defalut
	}
	h, err := os.ReadFile(matched)
	if err != nil {
		return "", err
	}
	return string(h), nil
	// fmt.Printf("defalut is %s & matched is %s", defalut, matched)
	// 替代变量
}

// 解析文章名称
// dir 是否多文件夹
func getTitle(post, key string, dir bool) (filename, title string, tags []string, err error) {

	tags = strings.Split(post, " ")
	filename = strings.Join(tags, "-")
	if dir {
		// 多文件夹
		filename = fmt.Sprintf("%s%sindex.md", filename, string(os.PathSeparator))
		// filename += string(os.PathSeparator) + "index.md"
	} else {
		// 单文件夹
		filename = fmt.Sprintf("%s%sindex.%s.md", filename, string(os.PathSeparator), key)
		// filename += string(os.PathSeparator) + "index.md"
	}
	title = capitalize(post)
	return
}

// 替换模板中的变量
// temp 模板原始值
// title 标题
// dest docs/book1/section1
// tags  {"tiptop","of","operation"}
// order 是否同文件夹排序，会将weight设置为最后一位
func replaceTemp(temp, title, project, dest string, tags []string, order bool) (template string, err error) {
	// title
	temp = strings.Replace(temp, "$title", title, 1)
	// date&lastmod
	// 2020-03-03T11:29:41+08:00
	// 2006-01-02T15:04:05Z07:00
	l, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(l).Format("2006-01-02T15:04:05+08:00")
	temp = strings.Replace(temp, "$date", t, 1)
	temp = strings.Replace(temp, "$lastmod", t, 1)
	// tags&categories
	t = ""
	for _, v := range tags {
		v = `"` + v + `"`
		t += v + ","
	}
	if len(t) > 1 {
		t = t[:len(t)-1]
	}
	temp = strings.Replace(temp, "$tags", t, 1)
	temp = strings.Replace(temp, "$categories", t, 1)

	if order {
		getOrder(project, dest)
	}

	return temp, nil
}

// 首字母大写
func capitalize(str string) string {
	if len(str) > 1 {
		str = strings.ToUpper(str[:1]) + str[1:]
	}
	if len(str) == 1 {
		str = strings.ToUpper(str)
	}
	return str
}

// 得到文件夹的最大序号
func getOrder(project, dest string) int {
	return 1
}
