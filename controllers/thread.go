package controllers

import (
	"github.com/makdenis/Golang-Api-project/Models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func CreateThread(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	slug := mux.Vars(request)["slug"]
	thread := Models.Thread{}
	if err := json.NewDecoder(request.Body).Decode(&thread); err != nil {
		panic(err)
	}
	res, checkThreads, forumName:= AddThread(Db, thread.Slug,thread.Author,slug, respWriter)
	if(forumName==""){
		respWriter.WriteHeader(http.StatusNotFound)
		tmp2:=errr{"Can't find user with nickname: "+slug}
		writeJSONBody(&respWriter, tmp2)
		return
	}
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
 	id:=42
	row.Scan(&id)
	thread.ID=id
		respWriter.WriteHeader(http.StatusCreated)
		writeJSONBody(&respWriter, thread)
	}else{
		respWriter.WriteHeader(http.StatusConflict)
		writeJSONBody(&respWriter, checkThreads[0])
	}}

func AddThread(Db *sql.DB, slug, author, forum string, respWriter http.ResponseWriter) (bool, []Models.Thread, string) {
	var forumName string
	tmp := []Models.Forum{}
	tmp=GetForumBySlug(Db, forum)
	if len(tmp)>0{
		forumName=tmp[0].Slug
	}
	checkuser:=GetUsersByEmailOrNick(Db,"",author)
	if len(checkuser)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user with nickname: "+author}
		writeJSONBody(&respWriter, tmp)
		return false, nil, forumName
	}
	var checkThreads []Models.Thread
	if slug!=""{
	checkThreads=GetThreadBySlugorID(Db,slug,0)}
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
	forums:=GetForumBySlug(Db, forum)
	threads:= make([]Models.Thread, 0)
	threads=GetThreadByForum(Db,forum,limit, since,desc)
	if len(forums)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user by slug: "+forum}
		writeJSONBody(&respWriter, tmp)
	}else{
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, threads)
	}
}

func GetThreadBySlugorID(Db *sql.DB, slug string,id int) []Models.Thread {
	threads := make([]Models.Thread, 0)
	query:="SELECT author::text, created::timestamp, forum::text, id::integer, message::text, slug::text,title::text, votes::integer FROM threads WHERE LOWER(slug) = LOWER($1) or id = $2 "
	var resultRows *sql.Rows
	resultRows,errr:= Db.Query(query, slug,id)
	if errr!=nil{
		fmt.Println(errr)}
	defer resultRows.Close()
	for resultRows.Next() {
		thread := new(Models.Thread)
		err := resultRows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {		}
		threads = append(threads, *thread)
	}
	resultRows.Close()
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
	resultRows,_= Db.Query(query, slug)
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
	resultRows,err:= Db.Query(query, id,slug)
	if err!=nil{
		fmt.Println(err)}
	defer resultRows.Close()
	for resultRows.Next() {
		thread := new(Models.Thread)
		err := resultRows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {		}
		threads = append(threads, *thread)
	}
	return threads
}


func UpdateThread(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	slug := mux.Vars(request)["slug"]
	id, _:= strconv.Atoi(slug)
	thread := Models.Thread{}
	if err := json.NewDecoder(request.Body).Decode(&thread); err != nil {
		panic(err)
	}
	oldthread := []Models.Thread{}
	oldthread= GetThreadBySlugorID(Db, slug,id)
	if len(oldthread)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp2:=errr{"Can't find user with nickname: "+slug}
		writeJSONBody(&respWriter, tmp2)
		return
	}
	if thread.Author==""{
		thread.Author=oldthread[0].Author
	}
	if thread.ID==0{
		thread.ID=oldthread[0].ID
	}
	if thread.Forum==""{
		thread.Forum=oldthread[0].Forum
	}
	if thread.Created==""{
		thread.Created=oldthread[0].Created
	}
	if thread.Message==""{
		thread.Message=oldthread[0].Message
	}
	if thread.Slug==""{
		thread.Slug=oldthread[0].Slug
	}
	if thread.Title==""{
		thread.Title=oldthread[0].Title
	}
	updateUserQuery := `UPDATE threads set author =$1, forum = $2, message = $3, slug=$4, title=$5  where lower (slug)=lower ($6) or id=$7;`
	_, _ = Db.Exec(updateUserQuery,  thread.Author,thread.Forum,thread.Message,thread.Slug,thread.Title,slug,id)
	respWriter.WriteHeader(http.StatusOK)
	writeJSONBody(&respWriter, thread)
	}
