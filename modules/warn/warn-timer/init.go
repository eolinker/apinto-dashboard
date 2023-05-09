package warn_timer

import "time"

func WarnTimer() {

	warn := newMonitorWarn()

	secondTick := time.Tick(time.Second)
	if time.Now().Second() > 30 {
		go warn.monitorWarn()
	}

	for t := range secondTick {
		if t.Second() == 30 {
			go warn.monitorWarn()
		}
	}
}
