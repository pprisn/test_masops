package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	u "github.com/pprisn/test_masops/server/utils"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Nsi struct {
	ID          uint       `json:"ID"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
	DeletedAt   *time.Time `json:"DeletedAt"`
	Name        string     `json:"Name"`
	Status      string     `json:"Status"`      //:7502 RussianPostEASnsi
	Statussdo   string     `json:"Statussdo"`   //:7522 RussianPostEASsdo
	Statusupd   string     `json:"Statusupd"`   //:7500 RussianPostEASConfiguration
	Statusauth  string     `json:"Statusauth"`  //:7501 RussianpostEASuser
	Statustrans string     `json:"Statustrans"` //:7524 RussianpostEAStruns
	Note        string     `json:"Note"`        //Примечание
	Ufpsid      string     `-`
}

type Ufps struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

var database *sql.DB
var loging string

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

		result, err := database.Exec("insert into masops.nsis (name, created_at, updated_at, status) values (?, NOW(), NOW(), ?)",
			name, status)
		if err != nil {
			log.Println(err)
		} else {
			ID_, err := result.LastInsertId()
			if err == nil {
				log.Printf("Insert ID=%d Name=%s, Status=%s\n", ID_, name, status)
			}
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/create.html")
	}
}

func MCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		name := r.FormValue("name")
		status := r.FormValue("status")
		statussdo := r.FormValue("statussdo")
		statusupd := r.FormValue("statusupd")
		statusauth := r.FormValue("statusauth")
		statustrans := r.FormValue("statustrans")
		note := r.FormValue("note")
		var ufpsid string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "myufpsid" {
				ufpsid = cookie.Value
			}
		}
		result, err := database.Exec("insert into masops.nsis (name, created_at, updated_at, status, statussdo,statusupd,statusauth,statustrans, note,ufpsid) values (?, NOW(), NOW(), ?, ?, ?, ?, ?, ?, ?)", name, status, statussdo, statusupd, statusauth, statustrans, note, ufpsid)
		if err != nil {
			//log.Println(err)
			log.Printf("%s\t%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.URL, err)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(err)))
		} else {
			ID_, err := result.LastInsertId()
			if err == nil {
				log.Printf("%s\t%s\t%s Insert ID=%d Name=%s, Status=%s\n",
					r.RemoteAddr, r.Method, r.URL, ID_, name, status)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint("pass")))
		}
		//		http.Redirect(w, r, "/", 301)
	}
}

// получаем измененные данные и сохраняем их в БД
func MEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		var myufps string
		for _, cookie := range r.Cookies() {
			if cookie.Name == "myufpsid" {
				myufps = cookie.Value
			}
		}

		id := r.FormValue("id")
		status := r.FormValue("status")
		statussdo := r.FormValue("statussdo")
		statusupd := r.FormValue("statusupd")
		statusauth := r.FormValue("statusauth")
		statustrans := r.FormValue("statustrans")
		note := r.FormValue("note")

		//	t := time.Now()
		_, err = database.Exec("update masops.nsis set status=?, statussdo=?, statusupd=?,statusauth=?, statustrans=?, note=?, updated_at= NOW() , ufpsid= ? where id = ?", status, statussdo, statusupd, statusauth, statustrans, note, myufps, id)

		if err != nil {
			log.Printf("%s\t%s\t%s\t%s\n", r.RemoteAddr, r.Method, r.URL, err)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint(err)))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprint("pass")))

		}
		//	http.Redirect(w, r, "/", 301)
	}
}

func DemoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "templates/demo.html")
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	nsis := []Nsi{}
	//	if r.Method == "POST" {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	log.Println(r.Cookies())
	ufpsid := r.FormValue("myufpsid")
	rows, err := database.Query("select id,created_at,updated_at,deleted_at,name,status,statussdo,statusupd,statusauth,statustrans,note from masops.nsis where ufpsid = ?", ufpsid)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	defer rows.Close()

	for rows.Next() {
		p := Nsi{}
		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.Name, &p.Status, &p.Statussdo, &p.Statusupd, &p.Statusauth, &p.Statustrans, &p.Note)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nsis = append(nsis, p)
	}
	//	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	err = tmpl.Execute(w, nsis)
	log.Println(err)

}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	nsis := []Nsi{}
	//	if r.Method == "POST" {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	log.Println(r.Cookies())
	//	log.Printf("%+v", r)
	ufpsid := r.FormValue("myufpsid")
	log.Printf("Запрос данных %+v", ufpsid)
	rows, err := database.Query("select id,created_at,updated_at,deleted_at,name,status,statussdo,statusupd,statusauth,statustrans,note from masops.nsis where ufpsid = ?", ufpsid)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	defer rows.Close()

	for rows.Next() {
		p := Nsi{}
		err := rows.Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt, &p.Name, &p.Status, &p.Statussdo, &p.Statusupd, &p.Statusauth, &p.Statustrans, &p.Note)
		if err != nil {
			fmt.Println(err)
			continue
		}
		nsis = append(nsis, p)
	}
	//	}

	jsonString, err := json.Marshal(nsis)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	//resp := u.Message(true, "Message")
	//	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	//	resp["data"] = nsis
	//	u.Respond(w, resp)
	//	log.Print(resp)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonString))
}

func UfpsHandler(w http.ResponseWriter, r *http.Request) {
	ufps := []Ufps{}
	ufpsid := ""
	for _, cookie := range r.Cookies() {
		if cookie.Name == "myufpsid" {
			ufpsid = cookie.Value
			break
		}
	}
	if ufpsid == "" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := u.Message(false, "Message")
		resp["Name"] = ""
		log.Print(resp)
		u.Respond(w, resp)
		return
	}

	rows, err := database.Query("select id, name from masops.ufps where id = ?", ufpsid)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
	}
	defer rows.Close()

	for rows.Next() {
		p := Ufps{}
		err := rows.Scan(&p.ID, &p.Name)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			break
		}
		ufps = append(ufps, p)
		break
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := u.Message(true, "Message")
	resp["Name"] = ufps[0].Name
	log.Print(resp)
	u.Respond(w, resp)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print(r.RemoteAddr, "\t", r.Method, "\t", r.URL)
		db, err := sql.Open("mysql", loging+"@tcp(127.0.0.1)/masops?charset=utf8&parseTime=True&loc=Local")
		defer db.Close()
		if err != nil {
			log.Println(err)
		}
		database = db
		h.ServeHTTP(w, r)
	})
}

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "./static/", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	loging = os.Getenv("LOGDB")
	if loging == "" {
		loging = "root"
	}

	// db, err := sql.Open("mysql", loging+"@tcp(127.0.0.1)/masops?charset=utf8&parseTime=True&loc=Local")
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	database = db
	//	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/mydata", DataHandler)
	router.HandleFunc("/myufps", UfpsHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	//	router.HandleFunc("/create", CreateHandler)
	router.HandleFunc("/mcreate", MCreateHandler)
	router.HandleFunc("/medit", MEdit)
	router.HandleFunc("/demo", DemoHandler)
	//	router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
	//	router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")

	router.Use(Middleware)

	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:3000",
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
