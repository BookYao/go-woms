package sqlmodule

import (
	"database/sql"
	"fmt"
)

const (
	host   = "192.168.100.186"
	port   = 5432
	user   = "womsuser"
	passwd = "womspasswd"
	dbname = "WOMSDB"
)

func ConnectDB() *sql.DB {
	pgsqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbname)

	fmt.Println("pgsqlInfo: ", pgsqlInfo)
	db, err := sql.Open("postgres", pgsqlInfo)
	if err != nil {
		fmt.Println("Open pgsql failed: ", err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Ping pgsql failed")
		return nil
	}

	fmt.Println("Connect pgsql success")

	return db
}

/*
func sqlQuery(db *sql.DB) bool {
	sqlcmd := "select ei_id, ei_creatuser, ei_time from \"T_Event_Info\";"
	row, err := db.Query(sqlcmd)
	if err != nil {
		fmt.Println("db Query failed.")
		return false
	}

	defer row.Close()

	columns, err := row.Columns()
	if err != nil {
		fmt.Println("Column is error: ", err.Error())
		return false
	}

	fmt.Println("Columes: ", columns)
	var id int
	var user, timestamp string
	for row.Next() {
		row.Scan(&id, &user, &timestamp)
		fmt.Printf("id: %d. user: %s. timestamp: %s\n", id, user, timestamp)
	}

	err = row.Err()
	if  err != nil {
		fmt.Println(err)
	}

	return true
}

func sqlUpdateFunc(db *sql.DB) bool {
	sqlcmd := fmt.Sprintf("update \"T_Event_Info\" set ei_creatuser=$1 where ei_id=$2")
	smtp, err := db.Prepare(sqlcmd)
	if err != nil {
		fmt.Println("Sqlcmd Prepare Failed.", err.Error())
		return false
	}

	result, err := smtp.Exec("15900001111", 3)
	if err != nil {
		fmt.Println("Sqlcmd Exec Failed.", err.Error())
		return false
	}

	affect, err := result.RowsAffected();
	if err != nil {
		fmt.Println("Sqlcmd Affect failed!", err.Error())
		return false
	}

	fmt.Println("Affect: ", affect)

	return true
}
*/
