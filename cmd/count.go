package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
)

var opts options

type options struct {
	DB_USER    string `short:"u" long:"user" description:"mysql user" default:"root" required:"false"`
	DB_HOST    string `short:"h" long:"host" description:"mysql host" default:"localhost" required:"true"`
	DB_PORT    string `short:"p" long:"port" description:"mysql port" default:"3306" required:"false"`
	DB_NAME    string `long:"db-name" default:"INFORMATION_SCHEMA"`
	WARN_COUNT int64  `long:"warn-count" description:"set threshold for warning" default:"0"`
	CRIT_COUNT int64  `long:"crit-count" description:"set threshold for critical" default:"0"`

	Threshold_Seconds int64 `long:"threshold" description:"threshold secons"`
}

func checkDuration(t string) (int64, error) {
	layout := "2006-01-02T15:04:05Z"
	t1, err := time.Parse(layout, t)
	if err != nil {
		return 0, err
	}
	t2 := time.Now()

	// Convert nanoseconds to seconds
	return (int64(t2.Sub(t1)) / 1000000000), nil
}

func connect() (*sql.DB, error) {
	user := opts.DB_USER
	host := opts.DB_HOST
	port := opts.DB_PORT
	pw := os.Getenv("MYSQL_PASSWORD")
	db_name := opts.DB_NAME
	var path string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", user, pw, host, port, db_name)
	var err error
	db, err := sql.Open("mysql", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Data struct {
	trxStarted string
}

func parseArgs(args []string) error {
	_, err := flags.ParseArgs(&opts, os.Args)

	if err != nil {
		return err
	}

	return nil
}

func do() error {
	chkSt := checkers.OK
	msg := "OK"

	Db, err := connect()
	if err != nil {
		chkSt = checkers.WARNING
		msg = fmt.Sprintf("connection error: %s", err.Error())
		checkers.NewChecker(chkSt, msg).Exit()
	}

	selected, err := Db.Query("SELECT trx_started FROM INNODB_TRX")
	if err != nil {
		chkSt = checkers.WARNING
		msg = fmt.Sprintf("connection error: %s", err.Error())
		checkers.NewChecker(chkSt, msg).Exit()
	}

	var overCount int64

	for selected.Next() {
		var data Data
		selected.Scan(&data.trxStarted)
		t, err := checkDuration(data.trxStarted)
		if err != nil {
			return err
		}

		if t > opts.Threshold_Seconds {
			overCount++
		}
	}

	if opts.WARN_COUNT > 0 && overCount >= opts.WARN_COUNT {
		chkSt = checkers.WARNING
		msg = fmt.Sprintf("[WARN] Number of long transactions exceeds threshold:%d", overCount)
	}
	if opts.CRIT_COUNT > 0 && overCount >= opts.CRIT_COUNT {
		chkSt = checkers.CRITICAL
		msg = fmt.Sprintf("[CRIT] Number of long transactions exceeds threshold:%d", overCount)
	}

	checkers.NewChecker(chkSt, msg).Exit()

	return nil

}

func Run() int {
	err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if do() != nil {
		return 1
	} else {
		return 0
	}
}
