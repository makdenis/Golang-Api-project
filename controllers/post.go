package controllers

import (
	"github.com/makdenis/Golang-Api-project/Models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	//"fmt"
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
	if len(posts)>0 {
		if posts[0].Author!=""{
			user:=GetUsersByEmailOrNick(Db,"",posts[0].Author)
			//fmt.Println(user[0].Nickname)
			if len(user)==0{
				respWriter.WriteHeader(http.StatusNotFound)
				tmp:=errr{"wrong nick: "}
				writeJSONBody(&respWriter, tmp)
				return
			}
		}




		if posts[0].Parent!=0{

	query:="Select slug::text from threads where id=(SELECT thread::integer FROM posts2 WHERE id = $1)"
var chkthread string


	var resultRows *sql.Rows

	//.Println(query)
	//fmt.Println(query)
	resultRows,err:= Db.Query(query, posts[0].Parent)
	if err!=nil{
		fmt.Println(err)}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {

		err := resultRows.Scan(&chkthread)
		if err != nil {		}


	}
	if chkthread!=thread{
		respWriter.WriteHeader(http.StatusConflict)
		tmp:=errr{"wrong parent: "}
		writeJSONBody(&respWriter, tmp)
		return
	}}}
	//user.Nickname = nickname

	//res, posts, checkuser := AddForum(Db, &forum, respWriter)
	//forum.User=checkuser.Nickname
	//fmt.Println(res) if res && forums==nil
	//threads := GetThreadBySlug(Db,hread, "", "", false)
	intthread,_:=strconv.Atoi(thread)
	time:=time.Now().UTC().Format(time.RFC3339)
	threads :=GetThreadById(Db,intthread,thread)
	if len(threads)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"wrong nick: "}
		writeJSONBody(&respWriter, tmp)
		return
	}
	for i, _ := range posts {
		posts[i].Thread=threads[0].ID
		//fmt.Println("aaa", strconv.Atoi(thread))
		posts[i].Created = time
		if len(threads)!=0{
			posts[i].Forum = threads[0].Forum}

		insertUserQuery4:=`INSERT INTO posts2 (author, forum, thread, message, created, is_edited, parent, tree_path) VALUES ($1, $2, $3, $4, $5, $6, $7, ((SELECT p.tree_path FROM posts2 p WHERE p.id=$7) || (SELECT currval('posts2_id_seq')::integer))) returning id`


		row = Db.QueryRow(insertUserQuery4,	posts[i].Author, posts[i].Forum, posts[i].Thread, posts[i].Message, posts[i].Created, posts[i].IsEdited, posts[i].Parent)
		row.Scan(&id)
		posts[i].ID = id
		//tree:=strconv.Itoa(posts[i].Parent)+"."+strconv.Itoa(posts[i].ID)
		//tree2,_:=strconv.ParseFloat(strings.TrimSpace(tree), 32)


		//fmt.Println(errr)

	}




	respWriter.WriteHeader(http.StatusCreated)
	writeJSONBody(&respWriter, posts)

}

func GetPost(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	var sort string
	var since string
	var sincetree string
	var limit string
	var desc bool
var lim int
	flag:=0
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	if (request.URL.Query()["sort"]!=nil){
		sort= request.URL.Query()["sort"][0]
	}
	if (request.URL.Query()["limit"]!=nil){
		limit= request.URL.Query()["limit"][0]
	}

	if (request.URL.Query()["since"]!=nil){
		since=" and id> "+ request.URL.Query()["since"][0]
		sincetree= request.URL.Query()["since"][0]
	}
	if (request.URL.Query()["desc"]!=nil){
		if request.URL.Query()["desc"][0]=="true"{
			desc= true}else{
			desc=false
		}
	}
	path:=[]string{}
	slug := mux.Vars(request)["slug"]
	id,_:=strconv.Atoi(slug)
	//fmt.Println(nickname)
	//fmt.Println("get")
	threads:=GetThreadBySlugorID(Db, slug,id)
	posts0:= make([]Models.Post, 0)
	posts:= make([]Models.Post, 0)
	flag2:=0
	if sort=="tree"{
	fmt.Println("tree")
		query:="SELECT author::text, created::timestamp, forum::text, id::integer, is_edited::boolean, message::text, parent::integer, thread::integer FROM posts2 WHERE thread = $1 and array_length(tree_path,1)=1  order by id "
		if desc{
			query+=" desc"
		}

		var resultRows *sql.Rows

		resultRows,err:= Db.Query(query, threads[0].ID)

		if err!=nil{
			fmt.Println(err)}
		fmt.Println("tressfsfe")
		//fmt.Println(err)
		defer resultRows.Close()

		for resultRows.Next() {
			post := new(Models.Post)
			err := resultRows.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited,&post.Message, &post.Parent, &post.Thread)
			if err != nil {		}

			posts0 = append(posts0, *post)
		}
		if sincetree!=""{

			query0:="select tree_path from posts2 where id=$1"

			err:=Db.QueryRow(query0, sincetree).Scan(pq.Array(&path));
			if err!=nil{
				fmt.Println(err)}
			if desc{
				sincetree=" and tree_path<$2 "
			}else{
			sincetree=" and tree_path>$2 "}
		}
		query="SELECT author::text, created::timestamp, forum::text, id::integer, is_edited::boolean, message::text, parent::integer, thread::integer FROM posts2 WHERE tree_path && ARRAY[CAST ($1 AS INTEGER)]"+sincetree+" order by tree_path "
		//var resultRows *sql.Rows
		if desc{
			query+=" desc"
		}

		if limit!=""{
			query+=" limit "+limit
			lim,_=strconv.Atoi(limit)
			if lim<10{
				flag=1
			}}
		b:=0
		for _,i:= range posts0{


		//	fmt.Println(query)
			if sincetree!=""{
				resultRows,err= Db.Query(query, i.ID,pq.Array(path))

			}else {
				resultRows, err = Db.Query(query, i.ID)
			}
			if err!=nil{
				fmt.Println(err)}
			//fmt.Println(err)
			defer resultRows.Close()

			for resultRows.Next() {
				if b==lim{
					break;
				}
				post := new(Models.Post)
				err := resultRows.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited,&post.Message, &post.Parent, &post.Thread)
				if err != nil {		}
				posts = append(posts, *post)
				b++
			}
			if flag==1 && len(posts)>2{
				break
			}
		}

		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, posts)
		return}
	if sort=="parent_tree"{
		fmt.Println("partree")
		if sincetree!=""{
			query:="SELECT p.id, p.author, p.forum, p.thread, p.message, p.created, p.is_edited, p.parent from posts2 p  join (select id from posts2 where thread=$1 and parent=0 "


			if desc{
				query+="and array[tree_path[1]] && array(select tree_path[1] from posts2 where tree_path[1]<(select tree_path[1] from posts2 where tree_path&& ARRAY[CAST ($2 AS INTEGER)]  order by tree_path[1] desc, tree_path  limit 1 )) "


			}else{
				query+="and array[tree_path[1]] && array(select tree_path[1] from posts2 where tree_path[1]>(select tree_path[1] from posts2 where tree_path&& ARRAY[CAST ($2 AS INTEGER)]  limit 1)) "
			}
			if desc{
				query+="order by tree_path[1] desc, tree_path limit $3) as t on t.id=tree_path[1] order by tree_path[1] desc, tree_path ;"
			}else{
				query+=" limit $3) as t on t.id=tree_path[1] order by tree_path ;"

			}
			var lim int
			if limit!="" {
				lim, _ = strconv.Atoi(limit)
			}else{
				lim=100500}

			resultRows,errr:= Db.Query(query, threads[0].ID,sincetree, lim)

			//fmt.Println(query)
			if errr!=nil{
				fmt.Println(errr)}
			//fmt.Println(err)
			defer resultRows.Close()


			for resultRows.Next() {

				post := new(Models.Post)
				err := resultRows.Scan(&post.ID, &post.Author, &post.Forum, &post.Thread, &post.Message,&post.Created, &post.IsEdited, &post.Parent)
				if err != nil {		}
				//post.Created=time
				posts = append(posts, *post)}
				respWriter.WriteHeader(http.StatusOK)
				writeJSONBody(&respWriter, posts)
				return
			}
		a:=0
		if sincetree!=""{

			//query0:="select tree_path from posts2 where id=$1"
			//
			//err:=Db.QueryRow(query0, sincetree).Scan(pq.Array(&path));
			//if err!=nil{
			//	fmt.Println(err)}
			if desc{
				sincetree=" and id<$2 "
			}else{
				sincetree=" and id>$2 "
			}

		}
		query:="SELECT author::text, created::timestamp, forum::text, id::integer, is_edited::boolean, message::text, parent::integer, thread::integer FROM posts2 WHERE thread = $1 and array_length(tree_path,1)=1 order by tree_path "
		if desc{
			query+=" desc"
		}
		if limit!=""&&since==""{
			query+=" limit "+limit
		}
		var resultRows *sql.Rows
		var err error

		resultRows, err = Db.Query(query, threads[0].ID)

		if err!=nil{
			fmt.Println(err)}
		//fmt.Println(err)
		defer resultRows.Close()

		for resultRows.Next() {
			post := new(Models.Post)
			resultRows.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited,&post.Message, &post.Parent, &post.Thread)
			posts0 = append(posts0, *post)
		}

		//	time:=time.Now()
		for _,i:= range posts0{

			query:="SELECT author::text, created::timestamp, forum::text, id::integer, is_edited::boolean, message::text, parent::integer, thread::integer FROM posts2 WHERE tree_path && ARRAY[CAST ($1 AS INTEGER)] "+sincetree+"order by tree_path"
			var resultRows *sql.Rows
			//if desc&&since!=""{
			//	query+=" desc"
			//}
			if sincetree!=""{
				resultRows,err= Db.Query(query, i.ID,since[9:])

			}else {
				resultRows, err = Db.Query(query, i.ID)
			}
			//resultRows,err:= Db.Query(query, i.ID)
			//if limit!=""{
			//	query+=" limit "+limit
			//	lim,_:=strconv.Atoi(limit)
			//	if lim<10{
			//		flag=1
			//	}}
			if err!=nil{
				fmt.Println(err)}
			//fmt.Println(err)
			defer resultRows.Close()
			//if resultRows.Next(){
			//if since!=""{
			//	posts=posts[:0]
			//	posts=append(posts, i)

			flag2=0

			for resultRows.Next() {

				if since!=""&&flag2==0{
					posts=posts[:0]
					if desc==false{
						posts=append(posts, i)}
					flag2=1
				}
				if since!=""{
					a++
				}
				post := new(Models.Post)
				err := resultRows.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited,&post.Message, &post.Parent, &post.Thread)
				if err != nil {		}
				//post.Created=time
				posts = append(posts, *post)
				//if desc==true&&since!=""{
				//	break
				//}
			}
			//if flag==1 && len(posts)>2{
			//	break
			//}
		}
		if a>0&&a<len(posts0)-1 {
			posts=posts[:0]
		}
		//fmt.Println(a)
		//if a>0&&a<len(posts0)-1 && desc==true{
		//	posts=posts[:0]
		//}
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, posts)
		return
	}

	//fmt.Println(threads)
	if len(threads)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp:=errr{"Can't find user by slug: "+slug}
		writeJSONBody(&respWriter, tmp)
	}else{
		posts=GetPostByThread(Db,threads[0].ID,limit, since,desc)
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, posts)
	}
}



func GetPostByThread(Db *sql.DB, id int, limit, since string, desc bool) []Models.Post {
	posts := make([]Models.Post, 0)
	query:="SELECT author::text, created::timestamp, forum::text, id::integer, is_edited::boolean, message::text, parent::integer, thread::integer FROM posts2 WHERE thread = $1"


	if desc{
		//query+=" desc "
		if since!=""{
		out:=[]rune(since)
		out[7]='<'
		since=string(out)}
	}
	query+=since

	query+=" order by id "

	if desc{
		query+=" desc "

	}
	var resultRows *sql.Rows
	if limit!=""{
		query+=" limit "+limit
	}
//.Println(query)
	//fmt.Println(query)
	resultRows,err:= Db.Query(query, id)
	if err!=nil{
		fmt.Println(err)}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		post := new(Models.Post)
		err := resultRows.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited,&post.Message, &post.Parent, &post.Thread)
		if err != nil {		}

		posts = append(posts, *post)
	}

	return posts
}


func GetPostById(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	posts := make([]Models.Post, 0)
	id1 := mux.Vars(request)["id"]
	id, _ := strconv.Atoi(id1)
	var related string
	if (request.URL.Query()["related"]!=nil){
		related= request.URL.Query()["related"][0]
	}
	query := "SELECT author::text, created::timestamp, forum::text, id::integer, is_edited::boolean, message::text, parent::integer, thread::integer FROM posts2 WHERE id = $1"

	var resultRows *sql.Rows

	//.Println(query)
	//fmt.Println(query)
	resultRows, err := Db.Query(query, id)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		post := new(Models.Post)
		err := resultRows.Scan(&post.Author, &post.Created, &post.Forum, &post.ID, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
		if err != nil {
		}

		posts = append(posts, *post)
	}
	if len(posts) == 0 {
		respWriter.WriteHeader(http.StatusNotFound)
		tmp := errr{"Can't find user by nickname: " }
		writeJSONBody(&respWriter, tmp)
	} else {
		if related=="user"{
			users:=GetUsersByEmailOrNick(Db,"",posts[0].Author)
			respWriter.WriteHeader(http.StatusOK)
			writeJSONBody(&respWriter, Models.PostDetails2{users[0],posts[0]})

		}
		if related=="thread"{
			threads:=GetThreadById(Db, posts[0].Thread,"")
			respWriter.WriteHeader(http.StatusOK)
			writeJSONBody(&respWriter, Models.PostDetails3{threads[0],posts[0]})

		}
		if related=="forum"{
			updateUserQuery := `UPDATE forums set posts = (select count (*) from posts2 where lower (forum)=lower($1)), threads =  (select count (*) from threads where lower (forum)=lower($1))  where lower (slug)=lower ($1);`

			_, _ = Db.Exec(updateUserQuery,posts[0].Forum)

			forum:=GetForumBySlug(Db, posts[0].Forum)
			respWriter.WriteHeader(http.StatusOK)
			writeJSONBody(&respWriter, Models.PostDetails5{forum[0],posts[0]})

		}
		if related=="user,forum"{
			users:=GetUsersByEmailOrNick(Db,"",posts[0].Author)
			updateUserQuery := `UPDATE forums set posts = (select count (*) from posts2 where lower (forum)=lower($1)), threads =  (select count (*) from threads where lower (forum)=lower($1))  where lower (slug)=lower ($1);`

			_, _ = Db.Exec(updateUserQuery,posts[0].Forum)

			forum:=GetForumBySlug(Db, posts[0].Forum)
			respWriter.WriteHeader(http.StatusOK)
			writeJSONBody(&respWriter, Models.PostDetails6{users[0],forum[0],posts[0]})

		}
		if related=="thread,forum"{
			threads:=GetThreadById(Db, posts[0].Thread,"")
			updateUserQuery := `UPDATE forums set posts = (select count (*) from posts2 where lower (forum)=lower($1)), threads =  (select count (*) from threads where lower (forum)=lower($1))  where lower (slug)=lower ($1);`

			_, _ = Db.Exec(updateUserQuery,posts[0].Forum)
			forum:=GetForumBySlug(Db, posts[0].Forum)
			respWriter.WriteHeader(http.StatusOK)
			writeJSONBody(&respWriter, Models.PostDetails7{forum[0],threads[0],posts[0]})

		}
		if related=="user,thread,forum"{
			users:=GetUsersByEmailOrNick(Db,"",posts[0].Author)
			threads:=GetThreadById(Db, posts[0].Thread,"")
			updateUserQuery := `UPDATE forums set posts = (select count (*) from posts2 where lower (forum)=lower($1)), threads =  (select count (*) from threads where lower (forum)=lower($1))  where lower (slug)=lower ($1);`

			_, _ = Db.Exec(updateUserQuery,posts[0].Forum)
			forum:=GetForumBySlug(Db, posts[0].Forum)
			respWriter.WriteHeader(http.StatusOK)
			writeJSONBody(&respWriter, Models.PostDetails8{users[0],forum[0],threads[0],posts[0]})

		}
		fmt.Println(related)
		if related=="user,thread"{
			users:=GetUsersByEmailOrNick(Db,"",posts[0].Author)
			threads:=GetThreadById(Db, posts[0].Thread,"")
			respWriter.WriteHeader(http.StatusOK)
			writeJSONBody(&respWriter, Models.PostDetails4{users[0],threads[0],posts[0]})

		}
		if related==""{
		respWriter.WriteHeader(http.StatusOK)
		writeJSONBody(&respWriter, Models.PostDetails{posts[0]})}
	}
}

func UpdatePost(Db *sql.DB, respWriter http.ResponseWriter, request *http.Request) {
	//fmt.Println(Db)
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")

	id1 := mux.Vars(request)["id"]
	//fmt.Println(request)
	id, _:= strconv.Atoi(id1)
	post := Models.Post{}

	if err := json.NewDecoder(request.Body).Decode(&post); err != nil {
		panic(err)
	}
	//user.Nickname = nickname
	oldpost := []Models.Post{}
	query:="SELECT author::text, created::timestamp, forum::text, id::integer, is_edited::boolean, message::text, parent::integer, thread::integer FROM posts2 WHERE id = $1"



	var resultRows *sql.Rows
var edit bool
	//.Println(query)
	//fmt.Println(query)
	resultRows,err:= Db.Query(query, id)
	if err!=nil{
		fmt.Println(err)}
	//fmt.Println(err)
	defer resultRows.Close()

	for resultRows.Next() {
		post1 := new(Models.Post)
		err := resultRows.Scan(&post1.Author, &post1.Created, &post1.Forum, &post1.ID, &post1.IsEdited,&post1.Message, &post1.Parent, &post1.Thread)
		if err != nil {		}

		oldpost = append(oldpost, *post1)
	}
	//forum.User=checkuser.Nickname
	//fmt.Println(res)
	if len(oldpost)==0{
		respWriter.WriteHeader(http.StatusNotFound)
		tmp := errr{"Can't find user by nickname: " }
		writeJSONBody(&respWriter, tmp)
		return
	}

	if post.Forum==""&&post.Author==""&&post.Created==""&&post.Message==""&&post.Parent==0&&post.ID==0&&post.Thread==0{
		post.IsEdited=false
		edit=false
	}else{
		post.IsEdited=true
		edit=true
	}
	if post.Forum==""&&post.Author==""&&post.Created==""&&post.Message==oldpost[0].Message&&post.Parent==0&&post.ID==0&&post.Thread==0{
		post.IsEdited=false
		edit=false
	}
	if post.Author==""{
		post.Author=oldpost[0].Author
	}
	if post.ID==0{
		post.ID=oldpost[0].ID
	}
	if post.Forum==""{
		post.Forum=oldpost[0].Forum
	}
	if post.Created==""{
		post.Created=oldpost[0].Created
	}
	if post.Message==""{
		post.Message=oldpost[0].Message
	}
	if post.Parent==0{
		post.Parent=oldpost[0].Parent
	}
	if post.Thread==0{
		post.Thread=oldpost[0].Thread
	}
	//post.IsEdited=true

	updateUserQuery := `UPDATE posts2 set author =$1, forum = $2, message = $3, parent=$4, thread=$5, is_edited=$7  where  id=$6;`

	_, _ = Db.Exec(updateUserQuery,  post.Author,post.Forum,post.Message,post.Parent,post.Thread,id, edit)
	respWriter.WriteHeader(http.StatusOK)
	writeJSONBody(&respWriter, post)
}
