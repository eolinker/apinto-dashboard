package apinto_dashboard

import (
	"github.com/eolinker/apinto-dashboard/internal/template"
	"os"
)

func RetTemplate(baseDir string, appends ...string) {
	template.ResetAppendView(baseDir, appends)
}

//判断目录是否存在
func IsDirExist(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}
