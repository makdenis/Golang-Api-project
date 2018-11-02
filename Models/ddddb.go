package Models

import (
	"database/sql"
	"fmt"
	"github.com/FogCreek/mini"
	"os"
	"os/user"
	//"GODB/controllers"
	//"GODB/Models"
)


var Db *sql.DB
func Create() error {
	Db, err := sql.Open("postgres", params())
	chk(err)
	defer Db.Close()

	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS Models (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		"user" TEXT NOT NULL,
		slug TEXT UNIQUE,
		posts INTEGER,
		threads INTEGER
	);`)
	chk(err)
return nil}

const help = `Usage: phonebook COMMAND [ARG]...
Commands:
	add NAME PHONE - create new record;
	del ID1 ID2... - delete record;
	edit ID        - edit record;
	show           - display all records;
	show STRING    - display records which contain a given substring in the name;
	help           - display this help.`

func fatal(v interface{}) {
	fmt.Println(v)
	os.Exit(1)
}

func chk(err error) {
	if err != nil {
		fatal(err)
	}
}

func params() string {
	u, err := user.Current()
	chk(err)

	cfg, err := mini.LoadConfiguration("/home/denis/go/src/GODB/phonebookrc")
	chk(err)

	info := fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		cfg.String("host", "127.0.0.1"),
		cfg.String("port", "5432"),
		cfg.String("dbname", u.Username),
		cfg.String("sslmode", "disable"),
		cfg.String("user", u.Username),
		cfg.String("pass", ""),
	)
	return info
}


func CreateForum(body *Forum) error {
	var rows *sql.Rows
	var err error

		rows, err = Db.Query(`INSERT INTO Models(slug, title , "user", posts, threads) VALUES($1,$2,$3,$4,$5)`,
			body.Slug, body.Title, body.User, body.Posts, body.Threads)

	if err != nil {
		return  err
	}
	defer rows.Close()

	//var rs= make([]ForumStruct, 0)
	//var rec ForumStruct
	//for rows.Next() {
	//	if err = rows.Scan(&rec.Id, &rec.Name, &rec.Phone); err != nil {
	//		return nil, err
	//	}
	//	rs = append(rs, rec)
	//}
	//if err = rows.Err(); err != nil {
	//	return nil, err
	//}
	//return rs, nil
return nil}