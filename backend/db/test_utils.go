package db

import (
	"database/sql"
	"fmt"
	"trxd/utils"

	"github.com/joho/godotenv"
)

var tmp_db *sql.DB

func OpenTestDB(testDBName string) error {
	godotenv.Load(".env")

	info, err := utils.GetDBInfoFromEnv()
	if err != nil {
		return fmt.Errorf("failed to get DB info from env: %v", err)
	}

	tmp_db, err = sql.Open("postgres", info.ConnectionString())
	if err != nil {
		return err
	}

	_, err = tmp_db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, testDBName))
	if err != nil {
		return err
	}
	_, err = tmp_db.Exec(fmt.Sprintf(`CREATE DATABASE %s;`, testDBName))
	if err != nil {
		return err
	}
	defer tmp_db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, testDBName))

	info.DBName = testDBName
	err = setupTestDB(info)
	if err != nil {
		return err
	}

	return nil
}

func setupTestDB(info *utils.DBInfo) error {
	err := ConnectDB(info, true)
	if err != nil {
		return err
	}

	return nil
}

func CloseTestDB() {
	CloseDB()
	if tmp_db != nil {
		tmp_db.Close()
	}
}

func DeleteAll() error {
	_, err := db.Exec(`SELECT delete_all();`)
	return err
}
