package plugin_timer

import "time"

func init() {
	iExtender := newExtender()
	//第一次启动就去更新扩展ID插件
	go iExtender.UpdateExtender()

	secondTick := time.Tick(time.Second)
	for t := range secondTick {
		if t.Second() == 60 {
			go iExtender.UpdateExtender()
		}
	}
}
