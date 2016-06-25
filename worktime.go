package main

import (
	"database/sql"
	"fmt"
	"math"
	"time"

	"github.com/lib/pq"
)

type workedTime struct {
	from    time.Time
	to      time.Time
	minutes int
}

type workTimeRow struct {
	StartDate  time.Time   `db:"start_date"`
	EndDate    pq.NullTime `db:"end_date"`
	Mon        float64
	Tue        float64
	Wed        float64
	Thu        float64
	Fri        float64
	WeeklyTime float64
	FullTime   bool
	ID         int
}

type workTimeData struct {
	entries  []DisplayLog
	holidays []holiday
	workTime []workTimeRow
	from     time.Time
	to       time.Time
	user     string
}

var fullTimeRow = workTimeRow{
	Mon: 7.7,
	Tue: 7.7,
	Wed: 7.7,
	Thu: 7.7,
	Fri: 7.7,
}

func calculateExpectedWorkTime(user string, from, to time.Time) (workedTime, error) {
	workTimeRows, err := getWorkTimeRows(user, from, to)
	if err != nil {
		return workedTime{}, err
	}

	holidays, err := getHolidaysInRange(from, to)
	if err != nil {
		return workedTime{}, err
	}

	logs, err := GetAbsenceForUser(user, from, to)
	if err != nil {
		return workedTime{}, err
	}

	data := &workTimeData{
		from:     from,
		to:       to,
		workTime: workTimeRows,
		holidays: holidays,
		entries:  logs,
	}

	ret := workedTime{
		from: from,
		to:   to,
	}
	for cur := from; cur.Before(to) || cur.Equal(to); cur = cur.AddDate(0, 0, 1) {
		row := data.getApplicableWorkTimeRow(cur)
		if row == (workTimeRow{}) {
			row = fullTimeRow
		}

		switch {
		case cur.Weekday() == time.Monday:
			if !data.isDuringAbsence(cur) {
				ret.minutes += int((row.Mon - (data.isHoliday(cur) * row.Mon)) * 60)
			}

		case cur.Weekday() == time.Tuesday:
			if !data.isDuringAbsence(cur) {
				ret.minutes += int((row.Tue - (data.isHoliday(cur) * row.Tue)) * 60)
			}

		case cur.Weekday() == time.Wednesday:
			if !data.isDuringAbsence(cur) {
				ret.minutes += int((row.Wed - (data.isHoliday(cur) * row.Wed)) * 60)
			}

		case cur.Weekday() == time.Thursday:
			if !data.isDuringAbsence(cur) {
				ret.minutes += int((row.Thu - (data.isHoliday(cur) * row.Thu)) * 60)
			}

		case cur.Weekday() == time.Friday:
			if !data.isDuringAbsence(cur) {
				ret.minutes += int((row.Fri - (data.isHoliday(cur) * row.Fri)) * 60)
			}
		}
	}
	return ret, nil
}

func (d *workTimeData) isHoliday(date time.Time) float64 {
	if len(d.holidays) == 0 {
		return 0
	}
	for _, h := range d.holidays {
		if h.HolidayDate.Equal(date) {
			return h.Val
		}
	}
	return 0
}

func (d *workTimeData) getApplicableWorkTimeRow(date time.Time) workTimeRow {
	for _, r := range d.workTime {
		if date.After(r.StartDate) && (!r.EndDate.Valid || date.After(r.EndDate.Time)) {
			return r
		}
	}
	return workTimeRow{}
}

func (d *workTimeData) isDuringAbsence(date time.Time) bool {
	for _, a := range d.entries {
		if (a.FromDate.Before(date) && a.ToDate.After(date)) ||
			a.FromDate.Equal(date) ||
			a.ToDate.Equal(date) {
			return true
		}
	}
	return false
}

func getWorkTimeRows(user string, from, to time.Time) ([]workTimeRow, error) {
	rows, err := db.Queryx(`
    select
      start_date, end_date, mon, tue, wed, thu, fri
    from
      work_time
    where
      user_id=$1 and
      start_date between $2 and $3
    order by start_date desc`,
		user, from, to)

	if err != nil {
		return []workTimeRow{}, nil
	}
	defer rows.Close()

	var ret []workTimeRow
	for rows.Next() {
		var row workTimeRow
		err = rows.StructScan(&row)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ret = append(ret, row)
	}

	return ret, nil
}

func getCurrentWorkTime(user string) (workTimeRow, error) {
	var row workTimeRow
	err := db.Get(&row,
		`select
      id, start_date, end_date, mon, tue, wed, thu, fri
    from
      work_time
    where
      user_id=$1 and
			end_date is null
    order by start_date desc`,
		user)

	if err == sql.ErrNoRows {
		row = workTimeRow{
			FullTime:   true,
			WeeklyTime: 38.5,
		}
		return row, nil
	}
	workTime := row.Mon + row.Tue + row.Wed + row.Thu + row.Fri
	row.WeeklyTime = workTime
	workTime = math.Floor((workTime * 10) + 0.5)
	if workTime == 385 {
		row.FullTime = true
	}

	return row, nil
}

func finishCurrentWorkTime(user string) error {
	row, err := getCurrentWorkTime(user)
	if err != nil {
		return err
	}
	if row.ID == 0 {
		return nil
	}

	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(-1 * time.Second)
	if row.StartDate.After(date) {
		_, err = db.Exec("delete from work_time where id=$1", row.ID)
		return err
	}
	_, err = db.Exec(
		`update work_time set
			end_date=$1
		where
			id=$2`, date, row.ID)
	return err
}

func insertWorkTime(row workTimeRow, user string) error {
	startDate := time.Date(row.StartDate.Year(), row.StartDate.Month(), row.StartDate.Day(), 0, 0, 0, 0, row.StartDate.Location())
	_, err := db.Exec(
		`insert into work_time
			(start_date, mon, tue, wed, thu, fri, user_id)
		values
			($1,$2,$3,$4,$5,$6,$7)`,
		startDate, row.Mon, row.Tue, row.Wed, row.Thu, row.Fri, user)

	return err
}
