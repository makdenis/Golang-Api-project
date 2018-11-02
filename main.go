

package main

import (
	"github.com/makdenis/Golang-Api-project/controllers"
	"database/sql"
	"fmt"
	"github.com/FogCreek/mini"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"strings"
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
	pwd, _ := os.Getwd()
	cfg, err := mini.LoadConfiguration(pwd+"/dbsettings")
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
	//Db.SetMaxOpenConns(2000000)
	//Db.SetMaxOpenConns(2000) // Sane default
	//Db.SetMaxIdleConns(200)
	chk(err)
	defer Db.Close()
	file1, err := os.Open("/home/denis/go/src/GODB/run.sql")
	file, err := ioutil.ReadAll(file1)

	if err != nil {
		// handle error
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		_,_= Db.Exec(request)
		// do whatever you need with result and error
	}
	//_, _ = Db.Exec("insert into users (nickname, about, email, fullname) values ('ss', 'ss', 'ss', 'ww');")//, user.Nickname, user.About, user.Email, user.Fullname)
	insertUserQuery:="Truncate table users, forums, threads,posts2, votes;"
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
	//router.HandleFunc("/api/forum/{slug}/details", func (output http.ResponseWriter, request *http.Request) {
	//	controllers.GetForumDetails(Db, output, request)})

	router.HandleFunc("/api/thread/{slug}/create", func (output http.ResponseWriter, request *http.Request) {
		controllers.CreatePost(Db, output, request)})

	router.HandleFunc("/api/thread/{slug}/vote", func (output http.ResponseWriter, request *http.Request) {
		controllers.Vote(Db, output, request)})

	router.HandleFunc("/api/thread/{slug}/details", func (output http.ResponseWriter, request *http.Request) {
		controllers.GetThreadDetails(Db, output, request)}).Methods("GET")
		//http.Handle("/",router)
	router.HandleFunc("/api/thread/{slug}/details", func (output http.ResponseWriter, request *http.Request) {
		controllers.UpdateThread(Db, output, request)}).Methods("POST")
	router.HandleFunc("/api/thread/{slug}/posts", func (output http.ResponseWriter, request *http.Request) {
		controllers.GetPost(Db, output, request)})
	router.HandleFunc("/api/forum/{slug}/users", func (output http.ResponseWriter, request *http.Request) {
		controllers.GetSortedUsers(Db, output, request)})
	router.HandleFunc("/api/post/{id}/details", func (output http.ResponseWriter, request *http.Request) {
		controllers.GetPostById(Db, output, request)}).Methods("GET")
	router.HandleFunc("/api/post/{id}/details", func (output http.ResponseWriter, request *http.Request) {
		controllers.UpdatePost(Db, output, request)}).Methods("POST")
	router.HandleFunc("/api/service/status", func (output http.ResponseWriter, request *http.Request) {
		controllers.Status(Db, output, request)})
	router.HandleFunc("/api/service/clear", func (output http.ResponseWriter, request *http.Request) {
		controllers.Clear(Db, output, request)})
	http.Handle("/",router)


	fmt.Println("Server is listening...")
	http.ListenAndServe(":5000", nil)
}




