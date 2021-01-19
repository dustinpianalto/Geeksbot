package discord_utils

import (
	"fmt"
	"time"
)

func ParseDateString(inTime time.Time) string {
	d := time.Now().Sub(inTime)
	s := int64(d.Seconds())
	days := s / 86400
	s = s - (days * 86400)
	hours := s / 3600
	s = s - (hours * 3600)
	minutes := s / 60
	seconds := s - (minutes * 60)
	dateString := ""
	if days != 0 {
		dateString += fmt.Sprintf("%v days ", days)
	}
	if hours != 0 {
		dateString += fmt.Sprintf("%v hours ", hours)
	}
	if minutes != 0 {
		dateString += fmt.Sprintf("%v minutes ", minutes)
	}
	if seconds != 0 {
		dateString += fmt.Sprintf("%v seconds ", seconds)
	}
	if dateString != "" {
		dateString += " ago."
	} else {
		dateString = "Now"
	}
	stamp := inTime.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%v\n%v", dateString, stamp)
}

func ParseDurationString(inDur time.Duration) string {
	s := int64(inDur.Seconds())
	days := s / 86400
	s = s - (days * 86400)
	hours := s / 3600
	s = s - (hours * 3600)
	minutes := s / 60
	seconds := s - (minutes * 60)
	durString := ""
	if days != 0 {
		durString += fmt.Sprintf("%v days ", days)
	}
	if hours != 0 {
		durString += fmt.Sprintf("%v hours ", hours)
	}
	if minutes != 0 {
		durString += fmt.Sprintf("%v minutes ", minutes)
	}
	if seconds != 0 {
		durString += fmt.Sprintf("%v seconds ", seconds)
	}
	if durString == "" {
		durString = "0 seconds"
	}
	return fmt.Sprintf("%v", durString)
}
