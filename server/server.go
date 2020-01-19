package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

/*
type Nsi struct {
	ID	   uint  `sql:int(10);unsigned NOT NULL AUTO_INCREMENT`
	CreatedAt time.Time `sql: datetime DEFAULT NULL`
	UpdatedAt time.Time `sql: datetime DEFAULT NULL`
	DeletedAt *time.Time `sql:"index"`
	Name   string `sql:varchar(100);unique;not null`
	Status string `sql:varchar(255) DEFAULT NULL`
}
*/

type Nsi struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string
	Status    string
}

var database *sql.DB

// возвращаем пользователю страницу для редактирования объекта
func EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	row := database.QueryRow("select * from masops.nsis where id = ?", id)
	nsi := Nsi{}
	err := row.Scan(&nsi.ID, &nsi.CreatedAt, &nsi.UpdatedAt, &nsi.DeletedAt, &nsi.Name, &nsi.Status)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, err := template.ParseFiles("templates/edit.html")
		if err != nil {
			log.Println(err)

		}
		err = tmpl.Execute(w, nsi)
		if err != nil {
			log.Println(err)
		}

	}
}

// получаем измененные данные и сохраняем их в БД
func EditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	status := r.FormValue("status")

	_, err = database.Exec("update masops.nsis set status=? where id = ?", status, id)

	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

/*
func CreateHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {

        err := r.ParseForm()
        if err != nil {
            log.Println(err)
        }
        model := r.FormValue("model")
        company := r.FormValue("company")
        price := r.FormValue("price")

        _, err = database.Exec("insert into productdb.Products (model, company, price) values (?, ?, ?)",
          model, company, price)

        if err != nil {
            log.Println(err)
        }
        http.Redirect(w, r, "/", 301)
    }else{
        http.ServeFile(w,r, "templates/create.html")
    }
}

*/

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from masops.nsis ")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	defer rows.Close()
	nsis := []Nsi{}

	for rows.Next() {
		p := Nsi{}
		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.Name, &p.Status)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nsis = append(nsis, p)
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	err = tmpl.Execute(w, nsis)
	log.Println(err)

}

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./static/", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	loging := os.Getenv("LOGDB")
	db, err := sql.Open("mysql", loging+"@tcp(localhost)/masops?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	//    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//    router.HandleFunc("/create", CreateHandler)
	router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
	router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")

	//    http.Handle("/",router)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	//    fmt.Println("Server is listening...")
	//    http.ListenAndServe(":3000", nil)
}
