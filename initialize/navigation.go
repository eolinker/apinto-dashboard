package initialize

import (
	_ "embed"
)

var (
	//go:embed navigation.yml
	navigationContent []byte
)

func InitNavigation() {
	// 初始化导航
}
