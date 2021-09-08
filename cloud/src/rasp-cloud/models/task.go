//Copyright 2021-2021 corecna Inc.

package models

import "rasp-cloud/tools"

func TaskCleanUpHosts() {
	// init crontab
	timer1 := &tools.CronTabTime{
		Hour: 0,
		Min:  0,
		Sec:  0,
		Nsec: 0,
	}
	tools.CronTabTimer(CleanOfflineHosts, timer1, 1)
}
