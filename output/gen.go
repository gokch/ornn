// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Gen struct {
	User User
}

func (t *Gen) Init(
	job *Job,
) {
	t.User.Init(job)
}

func (t *User) Init(
	job *Job,
) {
	t.job = job
}

type User struct {
	job *Job
}

func (t *User) Insert(
	val_id string,
	val_name string,
	val_seq int32,
) (
	lastInsertId int64,
	err error,
) {
	args := []interface{}{
		val_id,
		val_name,
		val_seq,
	}
	
	sql := fmt.Sprintf(
		"INSERT INTO user VALUES (?, ?, ?)",
	)
	
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.LastInsertId()
}

type User_select struct {
	Id   string
	Name string
	Seq  int32
}

func (t *User) Select(
	where_seq int32,
) (
	selects []*User_select,
	err error,
) {
	args := []interface{}{
		where_seq,
	}
	
	sql := fmt.Sprintf(
		"SELECT * FROM user WHERE seq = ?",
	)
	ret, err := t.job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*User_select, 0, 100)
	for ret.Next() {
		scan := &User_select{}
		err := ret.Scan(scan)
		if err != nil {
			return nil, err
		}
		selects = append(selects, scan)
	}
	
	return selects, nil
}

func (t *User) Update(
	set_id string,
	set_name string,
	set_seq int32,
	where_seq int32,
) (
	rowAffected int64,
	err error,
) {
	sql := fmt.Sprintf(
		"UPDATE user SET id = ?, name = ?, seq = ? WHERE seq = ?",
	)
	args := []interface{}{
		set_id,
		set_name,
		set_seq,
		where_seq,
	}
	
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

func (t *User) Delete(
	where_seq int32,
) (
	rowAffected int64,
	err error,
) {
	args := []interface{}{
		where_seq,
	}
	
	sql := fmt.Sprintf(
		"DELETE FROM user WHERE seq = ?",
	)
			
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

