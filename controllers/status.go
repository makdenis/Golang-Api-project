package controllers

import (
	"github.com/makdenis/Golang-Api-project/Models"
	"database/sql"
	"fmt"
	"net/http"
)

func Status(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println(Db)
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	var resultRows *sql.Rows
	var post int
	var forum int
	var user int
	var thread int
	//fmt.Println(request)


			query:="Select count (*) FROM posts2 "

			resultRows,err:= Db.Query(query)
			if err!=nil{
				fmt.Println(err)}
			//fmt.Println(err)
			defer resultRows.Close()

			for resultRows.Next() {

				err := resultRows.Scan(&post)
				if err != nil {		}


			}
	query="Select count (*) FROM threads "

	resultRows,err= Db.Query(query)
	if err!=nil{
		fmt.Println(err)}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {

		err := resultRows.Scan(&thread)
		if err != nil {		}


	}
	query="Select count (*) FROM users "

	resultRows,err= Db.Query(query)
	if err!=nil{
		fmt.Println(err)}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {

		err := resultRows.Scan(&user)
		if err != nil {		}


	}
	query="Select count (*) FROM forums "

	resultRows,err= Db.Query(query)
	if err!=nil{
		fmt.Println(err)}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {

		err := resultRows.Scan(&forum)
		if err != nil {		}


	}
	status:=Models.Status{forum,post,thread,user}

				respWriter.WriteHeader(http.StatusOK)

				writeJSONBody(&respWriter, status)
				return
			}


func Clear(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	insertUserQuery := "Truncate table users, forums, threads,posts2, votes;"
	_, _ = Db.Exec(insertUserQuery)
	respWriter.WriteHeader(http.StatusOK)

}