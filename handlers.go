package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"main/views"
	"net/http"
	"os"
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
	} else {
		views.ShowProfile("Tomek").Render(context.TODO(), w)
	}
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
func sendFullPost(w http.ResponseWriter, r *http.Request, content string, author string) {
	views.ShowFullPost(content, author).Render(context.TODO(), w)
}
func UserRouter(w http.ResponseWriter, r *http.Request) {

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
	c, _ := r.Cookie("Authorization")
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
	sendFullPost(w, r, content, nickname)
}
