package apinto_dashboard

import "github.com/eolinker/apinto-dashboard/internal/template"

func RetTemplate(baseDir string,appends ...string)  {
	template.ResetAppendView(baseDir,appends)
}