package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"main/views"
	"net/http"
	"os"
	"strings"

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
	c, _ := r.Cookie("Authorization")
	fmt.Println(c.Value[7:])
	views.ShowProfile("Tomek").Render(context.TODO(), w)
}
func sendTagPage(w http.ResponseWriter, r *http.Request) {
	views.ShowTag().Render(context.TODO(), w)
}
func sendInfoPage(w http.ResponseWriter, r *http.Request) {

	views.ShowInfo().Render(context.TODO(), w)
}
func sendDebug(w http.ResponseWriter, r *http.Request) {
	views.ShowDebug().Render(context.TODO(), w)
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
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}
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

	http.Redirect(w, r, "http://localhost:2000/main-page", http.StatusSeeOther)

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
	defer stmt.Close()
	if err != nil {
		fmt.Println(err)
	}

	_, err = stmt.Exec(inputEmail, "Tomek", encryptPasswordSHA256(inputPassword))
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "http://localhost:2000/main-page", http.StatusSeeOther)
}
func myFunc(w http.ResponseWriter, r *http.Request) {
	db := getConnection()
	defer db.Close()
	var nickname string
	rows, err := db.QueryContext(context.TODO(), "select nickname from userinfo;")

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&nickname); err != nil {
			log.Fatal(err)
		}
		fmt.Println(nickname)
	}

	w.Write([]byte("All check"))
}
