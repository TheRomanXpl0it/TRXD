package db

import (
	"context"
	"database/sql"
	"fmt"
	"trxd/utils"
	"trxd/utils/consts"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

var tmp_db *sql.DB

func OpenTestDB(testDBName string) error {
	consts.Testing = true

	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}

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
	defer func() {
		_, err := tmp_db.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s;`, testDBName))
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "55006" { // database is being accessed by other users
					return
				}
			}
			fmt.Printf("failed to drop test database \"%s\": %v\n", testDBName, err)
		}
	}()

	info.PgDBName = testDBName
	err = ConnectDB(info, true)
	if err != nil {
		return err
	}

	return nil
}

func CloseTestDB() error {
	err := CloseDB()
	if err != nil {
		return err
	}

	if tmp_db == nil {
		return nil
	}

	err = tmp_db.Close()
	if err != nil {
		return err
	}

	return nil
}

func DeleteAll(ctx context.Context) error {
	err := StorageFlush(ctx)
	if err != nil {
		return err
	}

	_, err = db.Exec(`SELECT delete_all();`)
	if err != nil {
		return err
	}

	return nil
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
