package utils

// ScriptConfig 脚本基础配置
type ScriptConfig struct {
	URL               string `yaml:"url"` // 对应接口
	ArticleBusinessID string `yaml:"article_business_id"`
	XCSRF             string `yaml:"xcsrf"`
	Cookie            string `yaml:"cookie"`
}

// Resp response
type Resp struct {
	Code    int64                  `json:"code"`
	Data    map[string]interface{} `json:"data"`
	Message string                 `json:"message"`
}
