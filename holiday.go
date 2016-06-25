//contains all functions and structs related to storing and retreiving holiday infos

package main

import (
	"fmt"
	"time"
)

type holiday struct {
	Name        string
	HolidayDate time.Time `db:"holiday_date"`
	Val         float64
}

func getHolidaysInRange(from, to time.Time) ([]holiday, error) {
	Info.Printf("getHolidaysInRange(%v,%v)\n", from.Format(jsDateTime), to.Format(jsDateTime))
	var ret []holiday

	rows, err := db.Queryx(`
    select
      name,
      holiday_date,
      val
    from
      holidays
    where
      holiday_date between $1 and $2`,
		from, to)
	if err != nil {
		fmt.Println(err)
		return []holiday{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var h holiday
		err = rows.StructScan(&h)

		if err == nil {
			ret = append(ret, h)
		}
	}
	return ret, nil
}
