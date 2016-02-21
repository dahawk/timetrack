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
}

var dbDateTime = "2006-01-02 15:04:05"
var jsFormatString = "02.01.2006 15:04"
var jsDateFormatString = "02.01.2006"

//GetEntry fetches the log entry for a given id
func GetEntry(entryID string) (DisplayLog, error) {
	var entry DisplayLog
	err := db.Get(&entry,
		`select ed.begin as begin, ed.end as end, ed.type as type, e.entry_id
		from entry e, entry_data ed
		where e.entry_data=ed.id and e.entry_id=$1
		order by ed.begin desc`, entryID)
	if err != nil {
		return DisplayLog{}, err
	}
	return prepareDateTime(entry), nil
}

//StoreEntry stores a time entry into the database
func StoreEntry(userID string, begin, end time.Time, entryType, createType, entryID string, active bool) error {
	var preparedBegin, preparedEnd string
	var err error
	if active && createType == "create" {
		now := time.Now()
		preparedBegin = now.Format(dbDateTime)
	} else if active && createType == "update" {
		preparedBegin = begin.Format(dbDateTime)

		now := time.Now()
		preparedEnd = now.Format(dbDateTime)
	} else {
		preparedBegin = begin.Format(dbDateTime)
		preparedEnd = end.Format(dbDateTime)
	}

	tx, err := db.Beginx()
	now := time.Now()
	if err != nil {
		return err
	}

	var insertedID int64
	if preparedEnd == "" {
		err = tx.QueryRowx(`insert into entry_data
			(created,"begin",type,create_type) values
			($1,$2,$3,$4) returning id`,
			now.Format(dbDateTime), preparedBegin, entryType, createType,
		).Scan(&insertedID)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	} else {
		err = tx.QueryRowx(`insert into entry_data
			(created,"begin","end",type,create_type) values
			($1,$2,$3,$4,$5) returning id`,
			now.Format(dbDateTime), preparedBegin, preparedEnd, entryType, createType,
		).Scan(&insertedID)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	}

	if createType == "create" {
		_, err := tx.Exec("insert into entry (entry_data,modified,user_id,active) values ($1,$2,$3,$4)", insertedID, now.Format(dbDateTime), userID, active)
		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			return err
		}
	} else if createType == "update" {
		_, err := tx.Exec("update entry set entry_data=$1, active=false where entry_id=$2", insertedID, entryID)
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
	var logs []DisplayLog
	var err error

	query := `select ed.begin as begin, ed.end as end, ed.type as type, e.entry_id
		from entry e, entry_data ed
		where e.entry_data=ed.id and ed.begin between $1 and $2 and e.user_id=$3
		order by ed.begin desc`
	if reverse {
		query = `select ed.begin as begin, ed.end as end, ed.type as type, e.entry_id
			from entry e, entry_data ed
			where e.entry_data=ed.id and ed.begin between $1 and $2 and e.user_id=$3
			order by ed.begin asc`
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
		if l.Type == "Work time" {
			ret = append(ret, prepareDateTime(l))
		} else if l.Type == "Holiday" || l.Type == "Sick leave" {
			ret = append(ret, prepareDate(l))
		}
	}
	return ret, nil
}

//DeleteEntry detaches an entry from its data.
//For the user this has the same effect as deleting the entry
func DeleteEntry(entryID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	now := time.Now().Format(dbDateTime)

	_, err = tx.Exec("update entry set entry_data=NULL, modified=$1 where entry_id=$2", now, entryID)
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
	var entryID string
	rows, err := db.Queryx(`select entry_id from entry
		where user_id=$1 and active=true and entry_data is not null`, userID)
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
