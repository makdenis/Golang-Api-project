package controllers

import (
	"GODB/Models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func Vote(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println(Db)
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	slug := mux.Vars(request)["slug"]
	//fmt.Println(request)

	vote := Models.Vote{}
	//
	if err := json.NewDecoder(request.Body).Decode(&vote); err != nil {
		panic(err)
	}
	//user.Nickname = nickname
	id,_:=strconv.Atoi(slug)
	user:=GetUsersByEmailOrNick(Db,"",vote.NickName)
	if len(user)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp2:=errr{"Can't find user with nickname: "+slug}
		writeJSONBody(&respWriter, tmp2)
		return
	}
	threads:=GetThreadBySlugorID(Db,slug,id)
		//res, checkThreads, forumName:= AddThread(Db, thread.Slug,thread.Author,slug, respWriter)
	//forum.User=checkuser.Nickname
	//fmt.Println(res)
	if len(threads)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp2:=errr{"Can't find user with nickname: "+slug}
		writeJSONBody(&respWriter, tmp2)
		return
	}else{
		var count int
		var flag bool
		insertUserQuery := `insert into votes (username, thread, voice) values ($1, $2, $3)`

		_, err:= Db.Exec(insertUserQuery, vote.NickName,threads[0].ID,vote.Voice)
		if err!=nil {
			fmt.Println(err)
			flag=true
		}
		if !flag{
		if vote.Voice>0 {
			insertUserQuery := `update threads set votes=votes+1 where LOWER(slug)=lower($1) or id = $2;`
			threads[0].Votes+=1
			_, errr:= Db.Exec(insertUserQuery, slug,id)
			if errr!=nil{
				fmt.Println(errr)}
		}


		if vote.Voice<0 {

			votes:=GetVoteByUser(Db,vote.NickName)

			if len(votes)>0{
				//if votes[0].Voice==-1{
				//	count=0
				//}else{
				count=2
			}else{
				count=1
			}
			insertUserQuery := `update threads set votes=votes-$3 where LOWER(slug)=lower($1) or id = $2;`
			threads[0].Votes-=count
			_, _ = Db.Exec(insertUserQuery, slug,id,count)

		}}

		respWriter.WriteHeader(http.StatusOK)

		writeJSONBody(&respWriter, threads[0])
	}}
func GetVoteByUser(Db *sql.DB, user string) []Models.Vote {
	votes := make([]Models.Vote, 0)
	query:="SELECT username::text, voice::integer FROM votes WHERE LOWER(username) = LOWER($1)"

	var resultRows *sql.Rows


	//fmt.Println(query)
	resultRows,err:= Db.Query(query, user)
	if err!=nil{
		fmt.Println(err)}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		vote := new(Models.Vote)
		err := resultRows.Scan(&vote.NickName, &vote.Voice)
		if err != nil {	}

		votes = append(votes, *vote)
	}

	return votes
}

func GetThreadDetails(Db *sql.DB,  respWriter http.ResponseWriter, request *http.Request)  {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	thr := make([]Models.Thread, 0)
	query:="SELECT author::text, created::timestamp, forum::text, id::integer, message::text, slug::text,title::text, votes::integer FROM threads WHERE LOWER(slug) = LOWER($1) or id=$2 "
	slug := mux.Vars(request)["slug"]
	id,_:=strconv.Atoi(slug)
	var resultRows *sql.Rows



	//fmt.Println(query)
	resultRows,_= Db.Query(query, slug, id)





	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		thread := new(Models.Thread)
		err := resultRows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.ID, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {	}

		thr = append(thr, *thread)
	}
	if len(thr)>0 {
		respWriter.WriteHeader(http.StatusOK)

		writeJSONBody(&respWriter, thr[0])
	}else{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp2:=errr{"Can't find user with nickname: "+slug}
		writeJSONBody(&respWriter, tmp2)
		return
	}
}

