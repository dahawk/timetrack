//contains all functions and structs to store and retreive log entry data
package main

import (
	"fmt"
	"time"
)

//DisplayLog is the struct to hold a log entry with all info necessary to be
//displayed on the website
type DisplayLog struct {
	FromDate *time.Time `db:"begin" json:"-"`
	ToDate   *time.Time `db:"end" json:"-"`
	DateFrom string
	TimeFrom string
	DateTo   string
	TimeTo   string
	Type     string
	EntryID  string `db:"entry_id"`
	Duration string
	Active   bool `db:"active"`
}

var dbDateTime = "2006-01-02 15:04:05"
var jsFormatString = "02.01.2006 15:04"
var jsDateFormatString = "02.01.2006"

//GetEntry fetches the log entry for a given id
func GetEntry(entryID string) (DisplayLog, error) {
	Info.Printf("GetEntry(%s)\n", entryID)
	var entry DisplayLog
	err := db.Get(&entry,
		`select
			ed.begin as begin,
			ed.end as end,
			ed.type as type,
			e.entry_id,
			e.active
		from
			entry e,
			entry_data ed
		where
			e.entry_data=ed.id and
			e.entry_id=$1
		order by
			ed.begin desc`,
		entryID)
	if err != nil {
		return DisplayLog{}, err
	}
	return prepareDateTime(entry, false), nil
}

//StoreEntry stores a time entry into the database
func StoreEntry(userID string, begin, end time.Time, entryType, createType, entryID, impostorID string, active, shouldRound bool) error {
	Info.Printf("StoreEntry(%s,%v,%v,%s,%s,%s,%s,%t)\n",
		userID, begin.Format(jsDateTime), end.Format(jsDateTime), entryType, createType, entryID, impostorID, active)
	var preparedBegin, preparedEnd string
	var err error
	now := time.Now()
	if shouldRound {
		now = round(time.Now())
	}
	if active && createType == createConst {

		preparedBegin = now.Format(dbDateTime)
	} else if active && createType == updateType {
		if shouldRound {
			begin = round(begin)
		}
		preparedBegin = begin.Format(dbDateTime)
		preparedEnd = now.Format(dbDateTime)
	} else {
		if shouldRound {
			begin = round(begin)
		}
		preparedBegin = begin.Format(dbDateTime)
		if shouldRound {
			end = round(end)
		}
		preparedEnd = end.Format(dbDateTime)
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	var insertedID int64
	if preparedEnd == "" {
		err = tx.QueryRowx(`insert into entry_data
			(created,"begin",type,create_type,by) values
			($1,$2,$3,$4,$5) returning id`,
			now.Format(dbDateTime), preparedBegin, entryType, createType, impostorID,
		).Scan(&insertedID)

		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	} else {
		err = tx.QueryRowx(
			`insert into entry_data
			(created,"begin","end",type,create_type,by) values
			($1,$2,$3,$4,$5,$6) returning id`,
			now.Format(dbDateTime), preparedBegin, preparedEnd, entryType, createType, impostorID,
		).Scan(&insertedID)

		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}

	if createType == createConst {
		_, err = tx.Exec(
			`insert into entry
			(entry_data,modified,user_id,active,by) values
			($1,$2,$3,$4,$5)`,
			insertedID, now.Format(dbDateTime), userID, active, impostorID)

		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	} else if createType == updateType {
		_, err = tx.Exec(
			`update entry set
				entry_data=$1,
				active=false,
				by=$2
			where entry_id=$3`,
			insertedID, impostorID, entryID)

		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//GetLogsForUser fetches the log entries for a given user in a given timespan
func GetLogsForUser(userID string, from, to time.Time, reverse bool) ([]DisplayLog, error) {
	Info.Printf("GetLogsForUser(%s,%v,%v,%t)\n", userID, from.Format(jsDateTime), to.Format(jsDateTime), reverse)
	var logs []DisplayLog
	var err error

	query := `
		select
			ed.begin as begin,
			ed.end as end,
			ed.type as type,
			e.entry_id,
			e.active
		from
			entry e,
			entry_data ed
		where
			e.entry_data=ed.id and
			ed.begin between $1 and $2 and
			e.user_id=$3
		order by
			ed.begin desc`
	if reverse {
		query = `
			select
				ed.begin as begin,
				ed.end as end,
				ed.type as type,
				e.entry_id,
				e.active
			from
				entry e,
				entry_data ed
			where
				e.entry_data=ed.id and
				ed.begin between $1 and $2 and
				e.user_id=$3
			order by
				ed.begin asc`
	}

	rows, err := db.Queryx(query, from, to, userID)

	if err != nil {
		return make([]DisplayLog, 0), err
	}
	defer rows.Close()
	for rows.Next() {
		var l DisplayLog
		err = rows.StructScan(&l)
		if err != nil {
			continue
		}
		logs = append(logs, l)
	}
	ret := []DisplayLog{}

	for _, l := range logs {
		if l.Type == workTimeConst {
			ret = append(ret, prepareDateTime(l, true))
		} else if l.Type == holidayConst || l.Type == sickConst {
			ret = append(ret, prepareDate(l))
		}
	}
	return ret, nil
}

//GetAbsenceForUser retrievs entries for holidays and sick leave for the given period
func GetAbsenceForUser(userID string, from, to time.Time) ([]DisplayLog, error) {
	Info.Printf("GetAbsenceForUser(%s,%v,%v)\n", userID, from.Format(jsDateTime), to.Format(jsDateTime))
	var logs []DisplayLog

	rows, err := db.Queryx(
		`select
			ed.begin as begin,
			ed.end as end,
			ed.type as type,
			e.entry_id
		from
			entry e,
			entry_data ed
		where
			e.entry_data=ed.id and
			ed.begin between $1 and $2 and
			e.user_id=$3 and
			ed.type in ('Holiday','Sick leave')`,
		from, to, userID)

	if err != nil {
		return []DisplayLog{}, err
	}

	for rows.Next() {
		var row DisplayLog
		err = rows.StructScan(&row)
		if err == nil {
			logs = append(logs, row)
		}
	}

	return logs, nil
}

//DeleteEntry detaches an entry from its data.
//For the user this has the same effect as deleting the entry
func DeleteEntry(entryID, by string) error {
	Info.Printf("DeleteEntry(%s,%s)\n", entryID, by)
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	now := time.Now().Format(dbDateTime)

	_, err = tx.Exec(
		`update entry set
			entry_data=NULL,
			modified=$1,
			by=$2
		where entry_id=$3`,
		now, by, entryID)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

//ActiveEntry checks for a user if there is an active log entry
func ActiveEntry(userID string) (string, error) {
	Info.Printf("ActiveEntry(%s)", userID)
	var entryID string
	rows, err := db.Queryx(
		`select
			entry_id
		from
			entry
		where
			user_id=$1 and
			active=true and
			entry_data is not null`,
		userID)
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		return "", nil
	}
	defer rows.Close()
	err = rows.Scan(&entryID)
	if err != nil {
		return "", err
	}

	return entryID, nil
}
