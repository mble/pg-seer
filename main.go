package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mble/pg-seer/version"
)

var (
	database    string
	user        string
	port        string
	versionFlag bool
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
	fields := reflect.Indirect(reflect.ValueOf(u))

	for i := 0; i < fields.NumField(); i++ {
		out.WriteString(fmt.Sprintf("%v: %v", fields.Type().Field(i).Name, fields.Field(i).Interface()))
		out.WriteString("\n")
	}
	return out.String()
}

func executeDemoQuery(database string, user string, port string) {
	connectionArgs := fmt.Sprintf("dbname=%s user=%s port=%s sslmode=disable ", database, user, port)
	db, err := sqlx.Connect("postgres", connectionArgs)
	if err != nil {
		log.Fatalln(err)
	}

	unusedIdx, err := ioutil.ReadFile("sql/unused_indexes.sql")
	if err != nil {
		log.Fatalln(err)
	}

	rows, err := db.Queryx(string(unusedIdx))
	unused := unusedIndexes{}

	fmt.Printf("Unused Indexes\n")
	fmt.Printf("==============\n\n")

	for rows.Next() {
		err := rows.StructScan(&unused)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%s\n", unused.String())
	}
}

func parseCommandLineFlags() {
	flag.BoolVar(&versionFlag, "version", false, "print version")
	flag.StringVar(&database, "database", "", "database to connect to")
	flag.StringVar(&user, "user", "", "database user to connect as")
	flag.StringVar(&port, "port", "", "port to connect on")
	flag.Parse()
	if versionFlag {
		log.Printf("Version: %s Build: %s\n", version.VERSION, version.GITCOMMIT)
		os.Exit(0)
	}

	if database == "" || user == "" {
		log.Println("must pass in both database and user flags")
		flag.Usage()
		os.Exit(1)
	} else {
		executeDemoQuery(database, user, port)
	}
}

func main() {
	parseCommandLineFlags()
}
