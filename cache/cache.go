package cache

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	_ "github.com/mattn/go-sqlite3"
	zaia_crypt "github.com/youyo/zaia/crypt"
)

const (
	database string = "cache.db"
	driver   string = "sqlite3"
)

func flushDatabase() {
	os.Remove(database)
}

func connectDatabase() *sql.DB {
	db, err := sql.Open(driver, database)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InitializeCacheDb() {
	flushDatabase()
	db := connectDatabase()
	defer db.Close()

	sqlStmt := `
		create table credentials (
			arn text not null primary key,
			credential_values blob not null,
			expired_at timestamp not null
		);
	`
	if _, err := db.Exec(sqlStmt); err != nil {
		log.Fatal(err)
	}
}

func ReadCredentialsFromCache(arn string) (credValues credentials.Value, err error) {
	db := connectDatabase()
	defer db.Close()

	stmt, err := db.Prepare("select credential_values, expired_at from credentials where arn = ?")
	if err != nil {
		return
	}
	defer stmt.Close()

	var encodedCredValues []byte
	var expiredAt time.Time
	err = stmt.QueryRow(arn).Scan(&encodedCredValues, &expiredAt)
	if err != nil {
		return
	}

	if isExpired(expiredAt) {
		err = errors.New("expired")
		return
	}

	credValues, err = zaia_crypt.Decode(encodedCredValues)
	return
}

func WriteCredentialsToCache(arn string, encodedCredValues []byte) error {
	db, err := sql.Open("sqlite3", "cache.db")
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("replace into credentials(arn, credential_values, expired_at) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(arn, encodedCredValues, time.Now().Add(59*time.Minute))
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func isExpired(t time.Time) bool {
	return t.Before(time.Now())
}
