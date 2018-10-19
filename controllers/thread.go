package controllers

import (
	"GODB/Models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func CreateThread(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println(Db)
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	slug := mux.Vars(request)["slug"]
	//fmt.Println(request)

	thread := Models.Thread{}

	if err := json.NewDecoder(request.Body).Decode(&thread); err != nil {
		panic(err)
	}
	//user.Nickname = nickname

	res, checkThreads, forumName:= AddThread(Db, thread.Slug,thread.Author,slug, respWriter)
	//forum.User=checkuser.Nickname
	//fmt.Println(res)
	if(forumName==""){
		respWriter.WriteHeader(http.StatusNotFound)
		tmp2:=err{"Can't find user with nickname: "+slug}
		writeJSONBody(&respWriter, tmp2)
		return
	}
	fmt.Println(forumName)
	thread.Forum=forumName
	if !res && checkThreads==nil {
		return}
	if res && len(checkThreads)==0{
	var row *sql.Row

	if(thread.Created==""){
		insertUserQuery := `insert into threads (author,  forum, message, title, slug) values ($1, $2, $3, $4, $5)  returning id;`

		row= Db.QueryRow(insertUserQuery, thread.Author,  thread.Forum, thread.Message, thread.Title, thread.Slug)

	}else{
		insertUserQuery := `insert into threads (author, created, forum, message, title, slug) values ($1, $2, $3, $4, $5, $6)  returning id;`

		row= Db.QueryRow(insertUserQuery, thread.Author, thread.Created, thread.Forum, thread.Message, thread.Title, thread.Slug)
	}
	//fmt.Println("ddd")
	//fmt.Println(row)
	id:=42
	//defer rows.Close()

	//for rows.Next() {
	row.Scan(&id)
	//fmt.Println(thread)
	//fmt.Println(kek)
	thread.ID=id
		respWriter.WriteHeader(http.StatusCreated)
		writeJSONBody(&respWriter, thread)
	}else{

		respWriter.WriteHeader(http.StatusConflict)
		writeJSONBody(&respWriter, checkThreads[0])
	}}

func AddThread(Db *sql.DB, slug, author, forum string, respWriter http.ResponseWriter) (bool, []Models.Thread, string) {
	//fmt.Println
	var forumName string
	tmp := []Models.Forum{}
	tmp=GetForumBySlug(Db, forum)
	//thread.Slug=slug
	if len(tmp)>0{
		forumName=tmp[0].Slug
	}
	checkuser:=GetUsersByEmailOrNick(Db,"",author)
	if len(checkuser)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=err{"Can't find user with nickname: "+author}
		writeJSONBody(&respWriter, tmp)
		return false, nil, forumName
	}
	var checkThreads []Models.Thread
	if slug!=""{
	checkThreads=GetThreadBySlug(Db,slug)}
	//if len(checkThreads)==0{
	//	respWriter.WriteHeader(http.StatusNotFound)
	//	tmp:=err{"Can't find user with nickname: "+forum.User}
	//	writeJSONBody(&respWriter, tmp)
	//	return false, nil, Models.User{}
	//}
	//fmt.Println(slug)
	//fmt.Println(checkThreads)
	//if len(conflictforums) == 2 && conflictforums[0] == conflictforums[1] {
	//	conflictforums = conflictforums[:1]
	//}

	if len(checkThreads) > 0 {
		return false, checkThreads, forumName
	}

	return true, checkThreads,forumName
}


func GetThread(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	var limit string
	var since string
	var desc bool
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	if (request.URL.Query()["limit"]!=nil){
		limit= request.URL.Query()["limit"][0]
	}
	if (request.URL.Query()["since"]!=nil){
		since= request.URL.Query()["since"][0]
	}
	if (request.URL.Query()["desc"]!=nil){
		if request.URL.Query()["desc"][0]=="true"{
		desc= true}else{
			desc=false
		}
	}
	forum := mux.Vars(request)["slug"]
	//fmt.Println(nickname)
	//fmt.Println("get")
	forums:=GetForumBySlug(Db, forum)
	threads:= make([]Models.Thread, 0)

	threads=GetThreadByForum(Db,forum,limit, since,desc)
	//fmt.Println(threads)
	if len(forums)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=err{"Can't find user by slug: "+forum}
		writeJSONBody(&respWriter, tmp)
	}else{
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, threads)
	}
}

func GetThreadBySlug(Db *sql.DB, slug string) []Models.Thread {
	threads := make([]Models.Thread, 0)
	query:="SELECT author::text, created::timestamp, forum::text, id::integer, message::text, slug::text,title::text, votes::integer FROM threads WHERE LOWER(slug) = LOWER($1) "

	var resultRows *sql.Rows



	//fmt.Println(query)
	resultRows,_= Db.Query(query, slug)

	//fmt.Println(err)
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		thread := new(Models.Thread)
		err := resultRows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {		}

		threads = append(threads, *thread)
	}

	return threads
}



func GetThreadByForum(Db *sql.DB, slug string, limit, since string, desc bool) []Models.Thread {
	threads := make([]Models.Thread, 0)
	query:="SELECT author::text, created::timestamp, forum::text, id::integer, message::text, slug::text,title::text, votes::integer FROM threads WHERE LOWER(forum) = LOWER($1) "

	var resultRows *sql.Rows

	if since!="" {
		t, err := time.Parse(time.RFC3339Nano, since)
		if err != nil {
			fmt.Println(err)
		}
		since = t.UTC().Format(time.RFC3339Nano)
		if desc{
			query+="and created <= "+ "'"+since+ "'"
		}else{
			query+="and created >= "+ "'"+since+ "'"
		}

	}
	if desc{
	query+=" order by created desc  "
	}else{
		query+=" order by created "
	}
	if limit!="" {
		query +="limit "+limit
		}


	//fmt.Println(query)
		resultRows,_= Db.Query(query, slug)

	//fmt.Println(err)
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		thread := new(Models.Thread)
		err := resultRows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {		}

		threads = append(threads, *thread)
	}

	return threads
}

func GetThreadById(Db *sql.DB, id int, slug string) []Models.Thread {
	threads := make([]Models.Thread, 0)
	query:="SELECT author::text, created::timestamp, forum::text, id::integer, message::text, slug::text,title::text, votes::integer FROM threads WHERE id = $1 or LOWER(slug) = LOWER($2)"

	var resultRows *sql.Rows


	//fmt.Println(query)
	resultRows,err:= Db.Query(query, id,slug)

	fmt.Println(err)
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		thread := new(Models.Thread)
		err := resultRows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {		}

		threads = append(threads, *thread)
	}

	return threads
}