package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type SMS struct {
	address string
	date    int64
	body    string
	sub_id  int
}

type Database struct {
	*sql.DB
	last      int64
	db_path   string
	last_path string
}

func (d *Database) Save(l int64) error {
	if l == 0 {
		return nil
	}
	d.last = l
	b := []byte(strconv.FormatInt(d.last, 10))
	err := ioutil.WriteFile(d.last_path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Pull() {
	cmd := exec.Command("adb", "pull", "/data/data/com.android.providers.telephony/databases/mmssms.db", d.db_path)
	cmd.Run()
}

func (d *Database) Clear() {
	os.Remove(d.db_path)
}

func (d *Database) Load() error {
	b, err := ioutil.ReadFile(d.last_path)
	if err != nil {
		return err
	}
	sb := string(b)
	last, err := strconv.ParseInt(sb, 10, 64)
	if err != nil {
		return err
	}
	d.last = last
	return nil
}

func (d *Database) Init() {
	err := d.Load()
	if err != nil {
		log.Println(err)
		d.Save(time.Now().Unix()*1000)
		log.Println("Last file not found. Created and set last time to now.")
	}
}

func (d *Database) Read() ([]SMS, error) {
	d.Pull()
	defer d.Clear()
	ret := make([]SMS, 0, 10)
	db, err := sql.Open("sqlite3", d.db_path)
	if err != nil {
		return ret, err
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf("SELECT `address`, `date`, `body`, `sub_id` FROM sms WHERE `date` > %d ORDER BY `date` ASC", d.last))
	if err != nil {
		return ret, err
	}
	defer rows.Close()

	for rows.Next() {
		sms := SMS{}
		err = rows.Scan(&sms.address, &sms.date, &sms.body, &sms.sub_id)
		if err != nil {
			return ret, err
		}
		ret = append(ret, sms)
	}
	err = rows.Err()
	if err != nil {
		return ret, err
	}
	return ret, nil
}
