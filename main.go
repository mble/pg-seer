package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type unusedIndexes struct {
	TableName     string `db:"tablename"`
	IndexRelname  string `db:"indexrelname"`
	IdxScan       int    `db:"idx_scan"`
	WriteActivity int    `db:"write_activity"`
	SeqScan       int    `db:"seq_scan"`
	NLiveTup      int    `db:"n_live_tup"`
	Size          string `db:"size"`
}

func (u *unusedIndexes) String() string {
	var out bytes.Buffer
	out.WriteString(fmt.Sprintf("%q", u.IndexRelname))
	return out.String()
}

func main() {
	db, err := sqlx.Connect("postgres", "dbname=mblewitt user=mblewitt sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	unusedIdx, err := ioutil.ReadFile("queries/unused_indexes.sql")
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := db.Queryx(string(unusedIdx))
	unused := unusedIndexes{}

	for rows.Next() {
		err := rows.StructScan(&unused)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%v\n", unused)
	}
}
