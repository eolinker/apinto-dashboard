package plugin_template_service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {

	clusterTemplatePlugin := newPluginTemplateService()

	bean.Injection(&clusterTemplatePlugin)
}
