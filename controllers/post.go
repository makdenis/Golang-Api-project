package controllers

import (
	"GODB/Models"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"time"

	//"github.com/gorilla/mux"
	"net/http"
)



func CreatePost(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println(Db)
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	var row *sql.Row
	thread := mux.Vars(request)["slug"]
	//fmt.Println(request)
	posts := make([]Models.Post, 0)
	id:=42
	//fmt.Println(post)
	if err := json.NewDecoder(request.Body).Decode(&posts); err != nil {
		panic(err)
	}
	//user.Nickname = nickname

	//res, posts, checkuser := AddForum(Db, &forum, respWriter)
	//forum.User=checkuser.Nickname
	//fmt.Println(res) if res && forums==nil
	//threads := GetThreadBySlug(Db,hread, "", "", false)
	intthread,_:=strconv.Atoi(thread)
	threads :=GetThreadById(Db,intthread,thread)
	for i, _ := range posts {
		posts[i].Thread=threads[0].ID
		//fmt.Println("aaa", strconv.Atoi(thread))
		posts[i].Created = time.Now().Format(time.RFC3339)
		if len(threads)!=0{
		posts[i].Forum = threads[0].Forum}
		insertUserQuery := `insert into posts (author, created, forum, is_edited, message, parent, thread) values ($1, $2, $3,$4,$5,$6,$7) returning id;`

			row = Db.QueryRow(insertUserQuery, posts[i].Author, posts[i].Created, posts[i].Forum, posts[i].IsEdited, posts[i].Message, posts[i].Parent, posts[i].Thread)

		row.Scan(&id)
		posts[i].ID = id
		//fmt.Println(errr)

	}

	respWriter.WriteHeader(http.StatusCreated)
	writeJSONBody(&respWriter, posts)
	//if res==false && forums!=nil {
	//	respWriter.WriteHeader(http.StatusConflict)
	//	writeJSONBody(&respWriter, forums[0])
	//}
}