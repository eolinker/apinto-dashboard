package timer

import (
	"time"
)

func TaskTimer() {

	mon := newMonitorWarn()
	tick := time.Tick(time.Second)
	if time.Now().Second() > 30 {
		go mon.monitorWarn()
	}
	for t := range tick {
		if t.Second() == 30 {
			go mon.monitorWarn()
		}
	}

	//c := cron.New(cron.WithSeconds())
	//
	//_, err := c.AddFunc("30 0/1 * * * ?", mon.monitorWarn)
	//if err != nil {
	//	panic(err)
	//}
	//
	//c.Start()

}
