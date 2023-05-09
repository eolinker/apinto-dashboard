package plugin_timer

import (
	"time"
)

func ExtenderTimer() {

	iExtender := newExtender()

	//第一次启动就去更新扩展ID插件
	go iExtender.UpdateExtender()
	secondTick := time.Tick(time.Second)
	for t := range secondTick {
		if t.Second() == 0 { //1分钟执行一次
			go iExtender.UpdateExtender()
		}
	}
}
