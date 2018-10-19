//package main
//
//import (
//	"database/sql"
//	"fmt"
//)
//const (
//	DB_USER     = "test_user"
//	DB_PASSWORD = "1"
//	DB_NAME     = "test"
//)
//func main() {
//	//dbinfo := "user=test_user password=1 dbname=test sslmode=disable"
//	db, _ := sql.Open("postgres", "postgres://test_user:1@localhost/test")
//	//checkErr(err)
//	fmt.Println("fff")
//
//	defer db.Close()
//}
//
//func checkErr(err error) {
//	if err != nil {
//		panic(err)
//	}
//}


package main

import (
	"GODB/controllers"
	"database/sql"
	"fmt"
	"github.com/FogCreek/mini"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"os/user"
)



//func format(rs []record) {
//	var max, tmp int
//	for _, v := range rs {
//		tmp = utf8.RuneCountInString(v.name)
//		if max < tmp {
//			max = tmp
//		}
//	}
//	s := "%-" + strconv.Itoa(max) + "s"
//	for _, v := range rs {
//		fmt.Printf("%3d   "+s+"   %s\n", v.id, v.name, v.phone)
//	}
//}



var Db *sql.DB

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



func main() {

	Db, err := sql.Open("postgres", params())
	chk(err)
	defer Db.Close()
	//_, _ = Db.Exec("insert into users (nickname, about, email, fullname) values ('ss', 'ss', 'ss', 'ww');")//, user.Nickname, user.About, user.Email, user.Fullname)
	insertUserQuery:="Truncate table users, forums, threads, posts;"
	_, _ = Db.Exec(insertUserQuery)
	router := mux.NewRouter()
	fmt.Println("dd")
	router.HandleFunc("/api/user/{nickname}/create", func (output http.ResponseWriter, request *http.Request) {
		controllers.CreateUser(Db, output, request)})

	router.HandleFunc("/api/user/{nickname}/profile", func (output http.ResponseWriter, request *http.Request) {
		controllers.GetUser(Db, output, request)}).Methods("GET")
	router.HandleFunc("/api/user/{nickname}/profile", func (output http.ResponseWriter, request *http.Request) {
		controllers.UpdateUser(Db, output, request)}).Methods("POST")

	router.HandleFunc("/api/forum/create", func (output http.ResponseWriter, request *http.Request) {
		controllers.CreateForum(Db, output, request)})

	router.HandleFunc("/api/forum/{slug}/details", func (output http.ResponseWriter, request *http.Request) {
		controllers.GetForum(Db, output, request)})

	router.HandleFunc("/api/forum/{slug}/create", func (output http.ResponseWriter, request *http.Request) {
		controllers.CreateThread(Db, output, request)})

	router.HandleFunc("/api/forum/{slug}/threads", func (output http.ResponseWriter, request *http.Request) {
		controllers.GetThread(Db, output, request)})

	router.HandleFunc("/api/thread/{slug}/create", func (output http.ResponseWriter, request *http.Request) {
		controllers.CreatePost(Db, output, request)})

		http.Handle("/",router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":5000", nil)
}




