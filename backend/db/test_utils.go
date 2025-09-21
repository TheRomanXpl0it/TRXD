package db

import (
	"database/sql"
	"fmt"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/joho/godotenv"
)

var tmp_db *sql.DB

func OpenTestDB(testDBName string) error {
	consts.Testing = true

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

	info.PgDBName = testDBName
	err = ConnectDB(info, true)
	if err != nil {
		return err
	}

	err = StorageFlush()
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

func InsertMockData() error {
	_, err := db.Exec(`SELECT insert_mock_data();`)
	if err != nil {
		return fmt.Errorf("failed to insert mock data: %v", err)
	}

	_, err = db.Exec(`SELECT insert_mock_submissions();`)
	if err != nil {
		return fmt.Errorf("failed to insert mock submissions: %v", err)
	}

	return nil
}
