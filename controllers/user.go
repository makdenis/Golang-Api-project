package controllers

import (
	"github.com/makdenis/Golang-Api-project/Models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type errr struct{
	Message string `json:"message"`
}
func CreateUser(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	nickname := mux.Vars(request)["nickname"]
	user := Models.User{}
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		panic(err)
	}
	user.Nickname = nickname
	res, users := AddUser(Db, &user)
	if res {
		insertUserQuery := `insert into users (nickname, about, email, fullname) values ($1, $2, $3, $4);`
		_, _ = Db.Exec(insertUserQuery, user.Nickname, user.About, user.Email, user.Fullname)
		respWriter.WriteHeader(http.StatusCreated)
		writeJSONBody(&respWriter, user)
	} else {
		respWriter.WriteHeader(http.StatusConflict)
		writeJSONBody(&respWriter, users)
	}
}

func GetUsersByEmailOrNick(Db *sql.DB, email, nickname string) []Models.User {
	users := make([]Models.User, 0)
	query := "SELECT about::text, email::text, fullname::text, nickname::text FROM users WHERE LOWER(email) = LOWER($1) OR LOWER(nickname) = LOWER($2)"
	resultRows, _ := Db.Query(query, email, nickname)
	defer resultRows.Close()
	for resultRows.Next() {
		user := new(Models.User)
		err := resultRows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		if err != nil {
			panic(err)
		}
		users = append(users, *user)
	}
	return users
}

func AddUser(Db *sql.DB, user *Models.User) (bool, []Models.User) {
	conflictUsers := GetUsersByEmailOrNick(Db, user.Email, user.Nickname)
	if len(conflictUsers) == 2 && conflictUsers[0] == conflictUsers[1] {
		conflictUsers = conflictUsers[:1]
	}
	if len(conflictUsers) > 0 {
		return false, conflictUsers
	}
	return true, nil
}

func writeJSONBody(respWriter *http.ResponseWriter, v interface{}) {
	if err := json.NewEncoder(*respWriter).Encode(v); err != nil {
		(*respWriter).WriteHeader(500)
	}
}

func GetUser(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	nickname := mux.Vars(request)["nickname"]
	user := make([]Models.User, 0)
	user=GetUsersByEmailOrNick(Db,"", nickname)
	if len(user)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user by nickname: "+nickname}
		writeJSONBody(&respWriter, tmp)
	}else{
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, user[0])
	}
}


func UpdateUser(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	nickname := mux.Vars(request)["nickname"]
	user := Models.User{}
	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		panic(err)
	}
	user.Nickname = nickname
	conflictUsers := GetUsersByEmailOrNick(Db, "", user.Nickname)
	if len(conflictUsers) == 0 {
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user by nickname: "+nickname}
		writeJSONBody(&respWriter, tmp)
		return
	}
	conflictEmails := GetUsersByEmailOrNick(Db, user.Email,"")
	if len(conflictEmails) != 0 {
		respWriter.WriteHeader(http.StatusConflict)
		tmp:=errr{"This email is already registered by user: "+conflictEmails[0].Nickname}
		writeJSONBody(&respWriter, tmp)
		return
	}
	olduser:=GetUsersByEmailOrNick(Db,"", nickname)
	if user.Fullname==""{
		user.Fullname=olduser[0].Fullname
	}
	if user.About==""{
		user.About=olduser[0].About
	}
	if user.Email==""{
		user.Email=olduser[0].Email
	}
	updateUserQuery := `UPDATE users set about =$1, email = $2, fullname = $3  where lower (nickname)=lower ($4);`
	_, _ = Db.Exec(updateUserQuery,  user.About, user.Email, user.Fullname, user.Nickname)
	respWriter.WriteHeader(http.StatusOK)
	writeJSONBody(&respWriter, user)
	}

func GetSortedUsers(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	var limit string
	var since string
	var desc bool
	forum := mux.Vars(request)["slug"]
	if (request.URL.Query()["limit"] != nil) {
		limit = request.URL.Query()["limit"][0]
	}
	if (request.URL.Query()["since"] != nil) {
		since = request.URL.Query()["since"][0]
	}
	if (request.URL.Query()["desc"] != nil) {
		if request.URL.Query()["desc"][0] == "true" {
			desc = true
		} else {
			desc = false
		}
	}
	users := make([]Models.User, 0)
	query := `SELECT about, fullname, email, nickname  FROM (
    SELECT about, fullname, email, nickname  FROM users AS u1
    JOIN threads AS t
    ON u1.nickname = t.author
    where lower(t.forum) = lower($1)
UNION SELECT about, fullname, email, nickname FROM users AS u2
    JOIN posts2 AS p
    ON u2.nickname = p.author
    where lower(p.forum) = lower($1)
    ) as res `
	if since != "" {
		if desc {
			query += ` where lower(nickname) < lower($2) COLLATE "ucs_basic" `
		} else {
			query += ` where lower(nickname) > lower($2) COLLATE "ucs_basic" `
		}
	}
	query += `ORDER BY lower(nickname) COLLATE "ucs_basic" `
	if desc {
		query += " desc"
	}
	if limit != "" {
		query += " limit " + limit
	}
	var resultRows *sql.Rows
	var err error
	if since != "" {
		resultRows, err = Db.Query(query, forum, since)
	} else {
		resultRows, err = Db.Query(query, forum)
	}
	if err != nil {
		fmt.Println(err)
	}
	defer resultRows.Close()
	for resultRows.Next() {
		user := new(Models.User)
		err := resultRows.Scan(&user.About, &user.Fullname, &user.Email, &user.Nickname)
		if err != nil {
			panic(err)
		}
		users = append(users, *user)
	}
	forums:=GetForumBySlug(Db,forum)
	if len(forums)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user by nickname: "+forum}
		writeJSONBody(&respWriter, tmp)
	}else {
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, users)
	}
}
