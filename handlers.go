package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/views"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func getConnection() *sql.DB {
	connStr := "user=ps dbname=db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func sendMainPage(w http.ResponseWriter, r *http.Request) {
	views.ShowHome().Render(context.TODO(), w)
}
func sendProfilePage(w http.ResponseWriter, r *http.Request) {
	if !isLoged(r) {
		sendLoginError(w, r)
		return
	}
	c, _ := r.Cookie("Authorization")
	claims := decriptJWT(c.Value[7:])
	views.ShowProfile(claims.Email).Render(context.TODO(), w)

}
func sendTagPage(w http.ResponseWriter, r *http.Request) {

	db := getConnection()
	defer db.Close()

	var ids, sqlIds []int
	var sqlAuthors, sqlContent []string

	vars := mux.Vars(r)
	tag := vars["tag"]
	fmt.Println(tag)

	// get post's ids from tag

	var sqlJsonIds string
	stmt, err := db.Prepare("select posts from tagposts join tag on tag.id = tagposts.tag where name = $1;")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(tag).Scan(&sqlJsonIds)
	// TODO no rows
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(sqlJsonIds)

	//extract ids from json [1,2,...]

	err = json.Unmarshal([]byte(sqlJsonIds), &ids)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ids)

	// extract id, text, author from table with json ids
	query := "select post.id, userdata.nickname as author, post.text  from post join userdata on post.author = userdata.id WHERE post.id IN ("
	for i, id := range ids {
		if i != 0 {
			query += ", "
		}
		query += fmt.Sprintf("%d", id)
	}
	query += ");"
	fmt.Println(query)

	stmt, err = db.Prepare(query)
	if err != nil {
		fmt.Println(err)
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	var id int
	var text, author string
	for rows.Next() {
		rows.Scan(&id, &author, &text)
		sqlIds = append(sqlIds, id)
		sqlAuthors = append(sqlAuthors, author)
		sqlContent = append(sqlContent, text)
	}
	fmt.Println(sqlIds, sqlAuthors, sqlContent)

	views.ShowTag(sqlIds, sqlAuthors, sqlContent).Render(context.TODO(), w)
	// views.ShowTag([]int{1, 2}, []string{"kys", "idiot"}, []string{"adam", "rolo"}).Render(context.TODO(), w)
}
func sendInfoPage(w http.ResponseWriter, r *http.Request) {

	views.ShowInfo().Render(context.TODO(), w)
}
func sendDebug(w http.ResponseWriter, r *http.Request) {
	views.ShowDebug().Render(context.TODO(), w)
}
func sendLoginError(w http.ResponseWriter, r *http.Request) {
	views.LoginError().Render(context.TODO(), w)
}
func sendCatImg(w http.ResponseWriter, r *http.Request) {
	img, err := os.ReadFile("img.png")
	if err != nil {
		http.Error(w, "Błąd odczytu pliku obrazka", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(img)
}
func sendFullPost(w http.ResponseWriter, r *http.Request, content string, author string, comcontent []string, comauthor []string) {
	views.ShowFullPost(content, author, comcontent, comauthor).Render(context.TODO(), w)
}
func UserRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("USER")
	switch r.Method {
	case http.MethodPost:
		updateUser(w, r)
	}
}

type Changes struct {
	Credentials string `json:"credentials"`
	Target      string `json:"target"`
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UPDATE")
	body := r.Body
	defer body.Close()
	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Println(err)
	}
	var change Changes
	json.Unmarshal(data, &change)
	fmt.Println(change.Target)
	fmt.Println(change.Credentials)

	db := getConnection()
	defer db.Close()

	if strings.Compare(change.Target, "email") == 0 {
		stmt, err := db.Prepare("update userdata set email = $1 where email = $2;")
		if err != nil {
			fmt.Println(err)
		}
		defer stmt.Close()

		c, err := r.Cookie("Authorization")
		fmt.Println("cookie: ", c.Value)
		claims := decriptJWT(c.Value[7:])
		_, err = stmt.Exec(change.Credentials, claims.Email)
		if err != nil {
			fmt.Println(err)
		}
	}
	if strings.Compare(change.Target, "password") == 0 {
		stmt, err := db.Prepare("update userdata set password = $1 where email = $2;")
		if err != nil {
			fmt.Println(err)
		}
		defer stmt.Close()

		c, err := r.Cookie("Authorization")
		fmt.Println("cookie: ", c.Value)
		claims := decriptJWT(c.Value[7:])
		_, err = stmt.Exec(encryptPasswordSHA256(change.Credentials), claims.Email)
		if err != nil {
			fmt.Println(err)
		}
	}
	if strings.Compare(change.Target, "nickname") == 0 {
		stmt, err := db.Prepare("update userdata set nickname = $1 where email = $2;")
		if err != nil {
			fmt.Println(err)
		}
		defer stmt.Close()

		c, err := r.Cookie("Authorization")
		fmt.Println("cookie: ", c.Value)
		claims := decriptJWT(c.Value[7:])
		_, err = stmt.Exec(change.Credentials, claims.Email)
		if err != nil {
			fmt.Println(err)
		}
	}

}
func login(w http.ResponseWriter, r *http.Request) {
	db := getConnection()
	defer db.Close()

	fmt.Println("login:")

	r.ParseForm()
	inputEmail := r.PostForm.Get("input_email")
	inputPassword := r.PostForm.Get("input_password")

	fmt.Println(inputEmail, inputPassword)

	stmt, err := db.Prepare("Select password, nickname from userdata where email=$1")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	var password, nickname string
	row := stmt.QueryRow(inputEmail)

	row.Scan(&password, &nickname)

	fmt.Println(password, nickname)

	if strings.Compare(password, encryptPasswordSHA256(inputPassword)) == 0 {
		cookie := http.Cookie{
			Name:     "Authorization",
			Value:    "Bearer " + getJWTFrom(inputEmail, nickname).String(),
			Path:     "/",
			MaxAge:   3600,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
	} else {
		fmt.Println("Negetive")
	}

	http.Redirect(w, r, "http://localhost:1000/main-page", http.StatusSeeOther)

}
func register(w http.ResponseWriter, r *http.Request) {
	db := getConnection()
	defer db.Close()

	fmt.Println("register:")

	r.ParseForm()
	inputEmail := r.PostForm.Get("input_email")
	inputPassword := r.PostForm.Get("input_password")

	fmt.Println(inputEmail, inputPassword)

	// Przygotowanie prepared statement
	stmt, err := db.Prepare("INSERT INTO userdata (email, nickname, password) VALUES ($1, $2, $3);")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(inputEmail, "Tomek", encryptPasswordSHA256(inputPassword))
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "http://localhost:1000/main-page", http.StatusSeeOther)
}
func isLoged(r *http.Request) bool {
	c, err := r.Cookie("Authorization")
	if err != nil {
		return false
	}
	if strings.Compare(c.Value[0:7], "Bearer ") != 0 {
		return false
	}
	if _, err := verifyJWT(c.Value[7:]); err != nil {
		return false
	}
	return true
}
func post(w http.ResponseWriter, r *http.Request) {
	if !isLoged(r) {
		sendLoginError(w, r)
	}
	switch r.Method {
	case http.MethodGet:
		getPost(w, r)
	case http.MethodPost:
		createPost(w, r)
	}

}

func getPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getPost")
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println(vars, id)

	db := getConnection()
	defer db.Close()

	stmt, err := db.Prepare("select text from post where id = $1;")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	var content string
	stmt.QueryRow(id).Scan(&content)

	c, err := r.Cookie("Authorization")
	fmt.Println("cookie: ", c.Value)
	if err != nil {
		fmt.Println(err)
	}

	stmt, err = db.Prepare("select nickname from (select author from post where id=$1) p join userdata on p.author = userdata.id;")
	if err != nil {
		fmt.Println(err)
	}
	var nickname string
	stmt.QueryRow(id).Scan(&nickname)

	// komentarze

	stmt, err = db.Prepare("select comments from postcomments where postid = $1;")
	if err != nil {
		fmt.Println(err)
	}
	var sqlids string
	var comids []int
	stmt.QueryRow(id).Scan(&sqlids)

	json.Unmarshal([]byte(sqlids), &comids)

	var com string
	var comments []string
	var sqlauthor string
	var sqlauthors []string
	for _, id := range comids {
		fmt.Println("comids: ", id)
		stmt, err = db.Prepare("select userdata.nickname, comment.text  from comment join userdata on comment.author = userdata.id where comment.id = $1;")
		if err != nil {
			fmt.Println(err)
		}
		stmt.QueryRow(id).Scan(&sqlauthor, &com)
		fmt.Println("com: ", com, "sqlauthor: ", sqlauthor)
		comments = append(comments, com)
		sqlauthors = append(sqlauthors, sqlauthor)
	}

	fmt.Println("Coms: ", comments, comids, sqlids)
	sendFullPost(w, r, content, nickname, comments, sqlauthors)
}

type JSONData struct {
	Tags    string `json:"tags"`
	Content string `json:"content"`
}

func createPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST post")
	body := r.Body
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("Błąd odczytu danych z ciała żądania:", err)
		return
	}

	var jsondata JSONData

	json.Unmarshal([]byte(string(data)), &jsondata)

	// input post
	db := getConnection()
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO post (author, text) VALUES ($1, $2);")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	//get author and content
	if err != nil {
		fmt.Println(err)
	}
	result, err := stmt.Exec(getUserDBID(w, r), jsondata.Content)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("result from post insert: ", result)

	// get id of created post
	var cratedPostID int
	stmt, err = db.Prepare("SELECT id FROM post ORDER BY id DESC LIMIT 1;")
	if err != nil {
		fmt.Println(err)
	}
	err = stmt.QueryRow().Scan(&cratedPostID)
	if err != nil {
		fmt.Println(err)
	}
	// input tags
	stmt, err = db.Prepare("select name from tag;")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	var temptag string
	var sqltags []string
	for rows.Next() {
		rows.Scan(&temptag)
		sqltags = append(sqltags, temptag)
	}

	fmt.Println(sqltags)
	tags := strings.Split(jsondata.Tags, " ")
	fmt.Println(tags)
	var newTags []string
	for _, tag := range tags {
		if !contains(sqltags, tag) {
			// Jeśli tak, dodaj tag do newTags
			newTags = append(newTags, tag)
		}
	}
	fmt.Println("TUTAJ")
	fmt.Println("tagi ktorych nie ma: ", newTags)
	for _, tag := range newTags {
		stmt, err = db.Prepare("INSERT INTO tag (name) VALUES ($1);")
		if err != nil {
			fmt.Println(err)
		}
		stmt.Exec(tag)
		stmt, err = db.Prepare("select id from tag order by id desc limit 1;")
		if err != nil {
			fmt.Println(err)
		}
		var cratedTagID int
		err = stmt.QueryRow().Scan(&cratedTagID)
		if err != nil {
			fmt.Println(err)
		}
		stmt, err = db.Prepare("INSERT INTO tagposts(tag, posts) VALUES ($1, '[0]');")
		if err != nil {
			fmt.Println(err)
		}
		_, err = stmt.Exec(cratedTagID)
		fmt.Println(err)
	}
	fmt.Println("tutaj")
	// get tags' ids from
	var posts string
	for _, tag := range tags {
		tagid := getTagDBID(tag)
		stmt, err = db.Prepare("select posts from tagposts where tag = $1;")
		err = stmt.QueryRow(tagid).Scan(&posts)
		fmt.Println(err, posts, cratedPostID)
		new := posts[:len(posts)-1] + ", " + strconv.Itoa(cratedPostID) + "]"
		fmt.Println(new)
		stmt, err = db.Prepare("UPDATE tagposts set posts = $1 where tag = $2;")
		_, err = stmt.Exec(new, tagid)
	}

}
func contains(sqltags []string, tag string) bool {
	for _, t := range sqltags {
		if t == tag {
			return true
		}
	}
	return false
}
func getUserDBID(w http.ResponseWriter, r *http.Request) int {
	c, err := r.Cookie("Authorization")
	fmt.Println("cookie: ", c.Value)
	claims := decriptJWT(c.Value[7:])
	db := getConnection()
	defer db.Close()
	stmt, err := db.Prepare("select id from userdata where email=$1;")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	var id int
	stmt.QueryRow(claims.Email).Scan(&id)
	fmt.Println("getUserDBID: ", id)

	return id
}
func getTagDBID(name string) int {
	db := getConnection()
	defer db.Close()
	stmt, err := db.Prepare("select id from tag where name=$1;")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	var id int
	stmt.QueryRow(name).Scan(&id)
	fmt.Println("getTagDBID: ", id)

	return id
}
func getTags(w http.ResponseWriter, r *http.Request) {
	fmt.Println("/tags")
	db := getConnection()
	defer db.Close()

	stmt, err := db.Prepare("select name from tag")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	var tag string
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}

	search := r.FormValue("search")

	fmt.Println(search)

	for rows.Next() {
		rows.Scan(&tag)
		if strings.Contains(tag, search) {
			w.Write([]byte("<td> <a href=/tag/" + tag + ">" + tag + "</a> </td>"))
		}
	}
}
func commentRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Println("commentRouter")
	switch r.Method {
	case http.MethodPost:
		createComment(w, r)
	}
}

type JSONComment struct {
	PostId  int    `json:"postId"`
	Comment string `json:"comment"`
}

func createComment(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	// Przykładowe wykorzystanie ciała żądania (np. odczytanie danych)
	// Możesz użyć ioutil.ReadAll lub innych metod do odczytu danych z ciała
	// W tym przykładzie używamy ioutil.ReadAll
	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("Błąd odczytu danych z ciała żądania:", err)
		return
	}
	var content JSONComment
	err = json.Unmarshal([]byte(string(data)), &content)
	fmt.Println(string(data), content)
	if err != nil {
		fmt.Println(err)
	}

	db := getConnection()
	defer db.Close()

	stmt, err := db.Prepare("insert into comment(author, text) values ($1, $2);")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(getUserDBID(w, r), content.Comment)
	if err != nil {
		fmt.Println(err)
	}

	stmt, err = db.Prepare("select postid from postcomments;")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	var sqlids []int
	var sqlid int
	for rows.Next() {
		rows.Scan(&sqlid)
		sqlids = append(sqlids, sqlid)
	}
	// is id in db(postcomments)
	fmt.Println(sqlids, content.PostId)
	found := false
	for _, tag := range sqlids {
		if tag == content.PostId {
			found = true
		}
	}
	if !found {
		stmt, err = db.Prepare("INSERT INTO postcomments (postid, comments) VALUES ($1, $2);")
		if err != nil {
			fmt.Println(err)
		}
		_, err = stmt.Exec(content.PostId, "[0]")
		if err != nil {
			fmt.Println(err)
		}
	}
	// take comments
	var comments string
	stmt, err = db.Prepare("select comments from postcomments where postid=$1;")
	if err != nil {
		fmt.Println(err)
	}
	err = stmt.QueryRow(content.PostId).Scan(&comments)
	fmt.Println(err, comments)

	var cratedCommentID int
	stmt, err = db.Prepare("SELECT id FROM comment ORDER BY id DESC LIMIT 1;")
	if err != nil {
		fmt.Println(err)
	}
	err = stmt.QueryRow().Scan(&cratedCommentID)
	if err != nil {
		fmt.Println(err)
	}

	//dubug
	newcomments := comments[:len(comments)-1] + ", " + strconv.Itoa(cratedCommentID) + "]"
	fmt.Println(newcomments)

	stmt, err = db.Prepare("update postcomments set comments = $1 where postid= $2;")
	if err != nil {
		fmt.Println(err)
	}
	_, err = stmt.Exec(newcomments, content.PostId)
	if err != nil {
		fmt.Println(err)
	}
	// insert new val into comments and insert

}
func debug(w http.ResponseWriter, r *http.Request) {

	body := r.Body
	defer body.Close()

	// Przykładowe wykorzystanie ciała żądania (np. odczytanie danych)
	// Możesz użyć ioutil.ReadAll lub innych metod do odczytu danych z ciała
	// W tym przykładzie używamy ioutil.ReadAll
	data, err := io.ReadAll(body)
	if err != nil {
		fmt.Println("Błąd odczytu danych z ciała żądania:", err)
		return
	}

	fmt.Println("Dane z ciała żądania:", string(data))

}
