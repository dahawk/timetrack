package main

import (
	"fmt"
	"time"
)

// People assume that time is a strict progression of cause to effect,
// but actually, from a non-linear, non-subjective viewpoint,
// it's more like a big ball of wibbly-wobbly... timey-wimey... stuff
type timeywimey struct {
	extectedWorkTime string
	actualWorkTime   string
	sickdays         int
	holidays         int
}

func calculateStats(data []DisplayLog) timeywimey {
	var workTime float64
	var holidays float64
	var sickleave float64
	for _, log := range data {
		switch {
		case log.Type == "Work time":
			workTime += log.ToDate.Sub(*log.FromDate).Minutes()
		case log.Type == "Holiday":
			holidays += log.ToDate.Sub(*log.FromDate).Hours()
		case log.Type == "Sick leave":
			sickleave += log.ToDate.Sub(*log.FromDate).Hours()
		}
	}

	hours := 0
	minutes := int(workTime)
	for ; minutes >= 60; minutes -= 60 {
		hours++
	}

	return timeywimey{
		holidays:       int(holidays / 24),
		sickdays:       int(sickleave / 24),
		actualWorkTime: fmt.Sprintf("%dh %dmin", hours, minutes),
	}
}

func formatDuration(interval float64, t string) string {
	switch {
	case t == "Work time":
		hours := 0
		minutes := int(interval)
		for ; minutes >= 60; minutes -= 60 {
			hours++
		}
		return fmt.Sprintf("%dh %dmin", hours, minutes)
	case t == "Holiday" || t == "Sich leave":
		if interval == 24 {
			return fmt.Sprintf("%d day", int(interval)/24)
		}
		return fmt.Sprintf("%d days", int(interval)/24)
	}
	return ""
}

func prepareDateTime(l DisplayLog) DisplayLog {
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
		ret.DateTo = l.ToDate.Format(jsDate)
		ret.TimeTo = l.ToDate.Format(jsTime)
	}

	return ret
}

func prepareDate(l DisplayLog) DisplayLog {
	dateFrom := time.Date(l.FromDate.Year(), l.FromDate.Month(), l.FromDate.Day(), 0, 0, 0, 0, l.FromDate.Location())
	dateTo := time.Date(l.ToDate.Year(), l.ToDate.Month(), l.ToDate.Day(), 23, 59, 59, 0, l.ToDate.Location())

	return DisplayLog{
		DateFrom: dateFrom.Format(jsDate),
		DateTo:   dateTo.Format(jsDate),
		Type:     l.Type,
		EntryID:  l.EntryID,
		FromDate: l.FromDate,
		ToDate:   l.ToDate,
	}
}

func getDefaultDates() (time.Time, time.Time) {
	now := time.Now()
	defaultFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	defaultTo := defaultFrom.AddDate(0, 1, 0).Add(-1 * time.Second)

	return defaultFrom, defaultTo
}
