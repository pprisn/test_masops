package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	t := time.Now()
	_, err = database.Exec("update masops.nsis set status=?, updated_at=? where id = ?", status, t, id)

	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		name := r.FormValue("name")
		status := r.FormValue("status")
		//t := time.Now()

		_, err = database.Exec("insert into masops.nsis (name, created_at, updated_at, status) values (?, NOW(), NOW(), ?)",
			name, status)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/create.html")
	}
}

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
	db, err := sql.Open("mysql", loging+"@tcp(127.0.0.1)/masops?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	//    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/create", CreateHandler)
	router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
	router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")

	//    http.Handle("/",router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	//log.Fatal(srv.ListenAndServe())

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
