package structure

type PostTemplate map[string]MarkDown

type MarkDown struct {
	Dest     string // 目标文件夹
	Title    string // 标题
	FileName string // 文件名
	Head     string // 模板值
}
