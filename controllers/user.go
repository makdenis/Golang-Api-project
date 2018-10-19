package controllers

import (
	"GODB/Models"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type err struct{
	message string `json:"message"`
}
func CreateUser(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println(Db)
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	nickname := mux.Vars(request)["nickname"]
	//fmt.Println(request)
	user := Models.User{}

	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		panic(err)
	}
	user.Nickname = nickname

	res, users := AddUser(Db, &user)
	//fmt.Println(res)
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
	//fmt.Println(nickname)
	//fmt.Println("get")
	user := make([]Models.User, 0)

	user=GetUsersByEmailOrNick(Db,"", nickname)
	if len(user)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=err{"Can't find user by nickname: "+nickname}
		writeJSONBody(&respWriter, tmp)
	}else{
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, user[0])
	}
}




func UpdateUser(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println("ppp")
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	nickname := mux.Vars(request)["nickname"]

	//fmt.Println(request)
	user := Models.User{}

	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		panic(err)
	}









	user.Nickname = nickname
	conflictUsers := GetUsersByEmailOrNick(Db, "", user.Nickname)
	if len(conflictUsers) == 0 {
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=err{"Can't find user by nickname: "+nickname}
		writeJSONBody(&respWriter, tmp)
		return
	}
//
// fmt.Println(user.Email)
	conflictEmails := GetUsersByEmailOrNick(Db, user.Email,"")
	if len(conflictEmails) != 0 {
		respWriter.WriteHeader(http.StatusConflict)
		tmp:=err{"This email is already registered by user: "+conflictEmails[0].Nickname}
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
	//res, users := AddUser(Db, &user)
	//fmt.Println(user)

		updateUserQuery := `UPDATE users set about =$1, email = $2, fullname = $3  where lower (nickname)=lower ($4);`

		_, _ = Db.Exec(updateUserQuery,  user.About, user.Email, user.Fullname, user.Nickname)
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, user)
	}
