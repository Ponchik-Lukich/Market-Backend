package validators

import (
	"fmt"
	"regexp"
	"sort"
)

func ValidateTime(time string) bool {
	re := regexp.MustCompile(`^[0-9][0-9]:[0-9][0-9]-[0-9][0-9]:[0-9][0-9]$`)
	if re.MatchString(time) {
		var startHour, startMinute, endHour, endMinute int
		_, err := fmt.Sscanf(time, "%d:%d-%d:%d", &startHour, &startMinute, &endHour, &endMinute)
		if err != nil {
			return false
		}
		if startHour < 0 || startHour > 23 || endHour < 0 || endHour > 23 || startMinute < 0 || startMinute > 59 || endMinute < 0 || endMinute > 59 {
			return false
		} else {
			if startHour > endHour || (startHour == endHour && (startMinute > endMinute || startMinute == endMinute)) {
				return false
			} else {
				return true
			}
		}
	} else {
		return false
	}
}

func ValidateTimeIntervals(intervals []string) bool {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i] < intervals[j]
	})
	// Check for overlaps
	var latestEnd int
	for _, interval := range intervals {
		var startHour, startMinute, endHour, endMinute int
		_, err := fmt.Sscanf(interval, "%d:%d-%d:%d", &startHour, &startMinute, &endHour, &endMinute)
		if err != nil {
			return false
		}

		start := startHour*60 + startMinute
		end := endHour*60 + endMinute
		if start < latestEnd {
			return false
		}
		if end > latestEnd {
			latestEnd = end
		}
	}
	return true
}
