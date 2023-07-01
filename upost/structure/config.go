package structure

// 多语言配置

type LanguageConfig struct {
	Source      string // 资源目录
	LanguageDir bool   // 多语言文件夹
	Languages   []Languages
}
type Languages struct {
	Name        string // 多语言名称
	Destination string //生成目的位置文件夹
	Template    string //模板位置
}

type LanguagesConfig map[string]LanguageConfig
