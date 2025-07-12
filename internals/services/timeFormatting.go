package services

import "time"

func TimeFormatting() string {
	currentTime := time.Now()
	return currentTime.Format("3:04 PM")
}
