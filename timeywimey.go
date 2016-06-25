//contains all functions and structs for necessary date/time conversions
package main

import (
	"fmt"
	"time"
)

// People assume that time is a strict progression of cause to effect,
// but actually, from a non-linear, non-subjective viewpoint,
// it's more like a big ball of wibbly-wobbly... timey-wimey... stuff
type timeywimey struct {
	ExtectedWorkTime string
	ActualWorkTime   string
	Delta            string
	Sickdays         int
	Holidays         int
}

func calculateStats(data []DisplayLog, anon anonStruct) timeywimey {
	Info.Println("calculateStats()")
	var workTime float64
	var holidays float64
	var sickleave float64
	containsSick := false
	containsHoliday := false
	for _, log := range data {
		if log.ToDate == nil {
			continue
		}
		switch {
		case log.Type == workTimeConst:
			workTime += log.ToDate.Sub(*log.FromDate).Minutes()
		case log.Type == holidayConst:
			containsHoliday = true
			holidays += log.ToDate.Sub(*log.FromDate).Hours()
		case log.Type == sickConst:
			containsSick = true
			sickleave += log.ToDate.Sub(*log.FromDate).Hours()
		}
	}
	if containsHoliday {
		holidays += 24
	}
	if containsSick {
		sickleave += 24
	}

	hours := 0
	minutes := int(workTime)
	for ; minutes >= 60; minutes -= 60 {
		hours++
	}

	expected, err := calculateExpectedWorkTime(anon.User.UserID, anon.from, anon.to)
	expectedString := "NaNh NaNmin"
	if err == nil {
		h := int(expected.minutes / 60)
		mins := expected.minutes - (h * 60)
		expectedString = fmt.Sprintf("%dh%dm", h, mins)
	}

	actString := fmt.Sprintf("%dh%dm", hours, minutes)
	actDur, _ := time.ParseDuration(actString)
	expDur, _ := time.ParseDuration(expectedString)
	var delta string

	if actDur != 0 && expDur != 0 {
		deltaMins := actDur.Minutes() - expDur.Minutes()

		h := int(deltaMins / 60)
		mins := int(deltaMins) - (h * 60)
		if mins < 0 {
			mins = mins * -1
		}
		delta = fmt.Sprintf("%dh%dm", h, mins)
	} else if actDur == 0 {
		delta = fmt.Sprintf("-%s", expectedString)
	}

	return timeywimey{
		Holidays:         int(holidays / 24),
		Sickdays:         int(sickleave / 24),
		ActualWorkTime:   actString,
		ExtectedWorkTime: expectedString,
		Delta:            delta,
	}
}

func formatDuration(from, to *time.Time, t string) string {
	now := time.Now()
	var tmpTime time.Time
	if to == nil {
		tmpTime = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.FixedZone("", 0))
	} else {
		tmpTime = *to
	}

	Info.Println("formatDuration()")
	switch {
	case t == workTimeConst:
		interval := tmpTime.Sub(*from).Minutes()
		hours := 0
		minutes := int(interval)
		for ; minutes >= 60; minutes -= 60 {
			hours++
		}
		return fmt.Sprintf("%dh %dmin", hours, minutes)
	case t == holidayConst || t == sickConst:
		interval := tmpTime.Sub(*from).Hours()
		interval += 24
		if interval == 24 {
			return fmt.Sprintf("%d day", int(interval)/24)
		}
		return fmt.Sprintf("%d days", int(interval)/24)
	}
	return ""
}

func prepareDateTime(l DisplayLog, skipTo bool) DisplayLog {
	Info.Printf("prepareDateTime(%v)\n", l)
	ret := DisplayLog{
		Type:     l.Type,
		EntryID:  l.EntryID,
		FromDate: l.FromDate,
		ToDate:   l.ToDate,
	}
	if l.FromDate != nil {
		ret.DateFrom = l.FromDate.Format(jsDate)
		ret.TimeFrom = l.FromDate.Format(jsTime)
	}
	if l.ToDate != nil {
		if !skipTo {
			ret.DateTo = l.ToDate.Format(jsDate)
		}
		ret.TimeTo = l.ToDate.Format(jsTime)
	}

	ret.Duration = formatDuration(l.FromDate, l.ToDate, l.Type)

	return ret
}

func prepareDate(l DisplayLog) DisplayLog {
	Info.Printf("prepareDate(%v)\n", l)
	dateFrom := time.Date(l.FromDate.Year(), l.FromDate.Month(), l.FromDate.Day(), 0, 0, 0, 0, l.FromDate.Location())
	dateTo := time.Date(l.ToDate.Year(), l.ToDate.Month(), l.ToDate.Day(), 23, 59, 59, 0, l.ToDate.Location())

	return DisplayLog{
		DateFrom: dateFrom.Format(jsDate),
		DateTo:   dateTo.Format(jsDate),
		Type:     l.Type,
		EntryID:  l.EntryID,
		FromDate: l.FromDate,
		ToDate:   l.ToDate,
		Duration: formatDuration(l.FromDate, l.ToDate, l.Type),
	}
}

func getDefaultDates() (time.Time, time.Time) {
	Info.Println("getDefaultDates()")
	now := time.Now()
	defaultFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	defaultTo := defaultFrom.AddDate(0, 1, 0).Add(-1 * time.Second)

	return defaultFrom, defaultTo
}

func round(in time.Time) time.Time {
	if in.Minute() > 57 || (in.Minute() == 57 && in.Second() > 30) {
		return time.Date(in.Year(), in.Month(), in.Day(), (in.Hour() + 1), 0, 0, 0, in.Location())
	}

	remainder := in.Minute()
	var correct int
	for ; remainder >= 5; remainder -= 5 {
	}
	if remainder < 2 {
		correct = -remainder
	} else if remainder > 2 {
		correct = (5 - remainder)
	} else if remainder == 2 {
		if in.Second() <= 30 {
			correct = -2
		} else {
			correct = 3
		}
	}

	return time.Date(in.Year(), in.Month(), in.Day(), in.Hour(), (in.Minute() + correct), 0, 0, in.Location())

}
