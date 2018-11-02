package controllers

import (
	"github.com/makdenis/Golang-Api-project/Models"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"

	//"github.com/gorilla/mux"
	"net/http"
)

func CreateForum(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println(Db)
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	//nickname := mux.Vars(request)["nickname"]
	//fmt.Println(request)
	forum := Models.Forum{}

	if err := json.NewDecoder(request.Body).Decode(&forum); err != nil {
		panic(err)
	}
	//user.Nickname = nickname

	res, forums, checkuser := AddForum(Db, &forum, respWriter)
	forum.User=checkuser.Nickname
	//fmt.Println(res)
	if res && forums==nil{
		insertUserQuery := `insert into forums (author, title, slug) values ($1, $2, $3);`

		_, _= Db.Exec(insertUserQuery, forum.User, forum.Title, forum.Slug)
		//fmt.Println(errr)
		respWriter.WriteHeader(http.StatusCreated)
		writeJSONBody(&respWriter, forum)
	}
	if res==false && forums!=nil {
		respWriter.WriteHeader(http.StatusConflict)
		writeJSONBody(&respWriter, forums[0])
	}
}


func Getforumbyname(Db *sql.DB, title string, user string) []Models.Forum {
	forums := make([]Models.Forum, 0)

	query := "SELECT author::text, title::text, slug::text FROM forums WHERE LOWER(title) = LOWER($1) or LOWER(author) = LOWER($2);"

	resultRows, _ := Db.Query(query, title, user)
	defer resultRows.Close()

	for resultRows.Next() {
		forum := new(Models.Forum)
		err := resultRows.Scan(&forum.User, &forum.Title, &forum.Slug)
		if err != nil {
			panic(err)
		}

		forums = append(forums, *forum)
	}

	return forums
}

func AddForum(Db *sql.DB, forum *Models.Forum, respWriter http.ResponseWriter) (bool, []Models.Forum, Models.User) {
	checkuser:=GetUsersByEmailOrNick(Db,"",forum.User)
	if len(checkuser)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user with nickname: "+forum.User}
		writeJSONBody(&respWriter, tmp)
		return false, nil, Models.User{}
	}
	conflictforums:= Getforumbyname(Db, forum.Title, forum.User)
	//fmt.Println(forum.User)
	if len(conflictforums) == 2 && conflictforums[0] == conflictforums[1] {
		conflictforums = conflictforums[:1]
	}

	if len(conflictforums) > 0 {
		return false, conflictforums, checkuser[0]
	}

	return true, nil, checkuser[0]
}

func GetForum(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {

	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	slug := mux.Vars(request)["slug"]
	//fmt.Println(nickname)
	//fmt.Println("get")
	forum:= make([]Models.Forum, 0)
	id, _:= strconv.Atoi(slug)

	updateUserQuery := `UPDATE forums set posts = (select count (*) from posts2 where lower (forum)=lower($1)), threads =  (select count (*) from threads where lower (forum)=lower($1))  where lower (slug)=lower ($1) or id=$2;`

	_, _ = Db.Exec(updateUserQuery,slug,id  )
	forum=GetForumBySlug(Db,slug)
	if len(forum)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user by nickname: "+slug}
		writeJSONBody(&respWriter, tmp)
	}else{
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, forum[0])
	}
}


func GetForumBySlug(Db *sql.DB, slug string) []Models.Forum {
	forums := make([]Models.Forum, 0)

	query := "SELECT author::text, title::text, slug::text, posts::integer, threads::integer FROM forums WHERE LOWER(slug) = LOWER($1)"

	resultRows, _ := Db.Query(query, slug)
	defer resultRows.Close()

	for resultRows.Next() {
		forum := new(Models.Forum)
		err := resultRows.Scan(&forum.User, &forum.Title, &forum.Slug,&forum.Posts,&forum.Threads)
		if err != nil {
			panic(err)
		}

		forums = append(forums, *forum)
	}

	return forums
}

//func GetForumDetails(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
//	forums := make([]Models.Forum, 0)
//	slug := mux.Vars(request)["slug"]
//	id, _:= strconv.Atoi(slug)
//
//	updateUserQuery := `UPDATE forums set posts2 = (select count (*) from posts2 where lower (forum)=lower($1)), threads =  (select count (*) from threads where lower (forum)=lower($1))  where lower (slug)=lower ($1) or id=$2;`
//
//	_, _ = Db.Exec(updateUserQuery,slug,id  )
//	forums=GetForumBySlug(Db, slug)
//	respWriter.WriteHeader(http.StatusOK)
//	writeJSONBody(&respWriter, forums)
//	}
//
//
//
