//Программа сбора данных версий МАСОПС на сети УФПС Липецкой области
//http://localhost:7502/v1/
// Читаем отклик службы и результат сохраняем в БД
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	//	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"sync"

	"golang.org/x/net/context"


	//         _ "github.com/mattn/go-sqlite3"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var fsrc = flag.String("fsrc", "./fsrc.txt", `Файл с данными адресов для мониторинга отклика работы службы МАСОПС`)
var mode = flag.String("mode", "l", `Режим логирования отклика службы, l краткий, f полный`)

/*
nsi
http://localhost:7502/v1/

sdo
http://localhost:7522/v1/

config
http://localhost:7500/v1/

easuser
http://localhost:7501/v1/
*/

type Nsi struct {
	gorm.Model
	Name       string `gorm:"type:varchar(100);unique;not null"`
	Status     string `gorm:"type:varchar(255)"` //:7502
	Statussdo  string `gorm:"type:varchar(255)"` //:7522
	Statusupd  string `gorm:"type:varchar(255)"` //:7500
	Statusauth string `gorm:"type:varchar(255)"` //:7501
	Statustrans string `gorm:"type:varchar(255)"` //:7524
}

func main() {

	//	var err error
	loging := os.Getenv("LOGDB")
	if loging == "" {
		loging = "root"
	}
	fmt.Printf("LOGDB=%s\n", loging)
	//	return
	//db, err := gorm.Open("sqlite3", "masops.db?cache=shared&mode=rwc")
	db, err := gorm.Open("mysql", loging+"@(localhost)/masops?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	db.AutoMigrate(&Nsi{})

	flag.Parse()
	var floger, f *os.File

	if floger, err = os.OpenFile("mas.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		panic(err)
	}
	defer floger.Close()

	log.SetOutput(floger)
	t0 := time.Now()
	log.Printf("СТАРТ %v \n", t0)

	if f, err = os.Open(*fsrc); err != nil {
		log.Printf("Error open %s \n", *fsrc)
		panic(err)
	}
	defer f.Close()

	err = check_nsi(db, f)
	//return
//!	scanner := bufio.NewScanner(f)
	var nameip string
	//	d := net.Dialer{Timeout: time.Second * 4}

	// Получить выборку
	rows, err := db.Raw("select id, name from nsis").Rows()
	defer rows.Close()

	var name string
	var id int
	//var version string
	var port string
	var vStatus7502, vStatus7522, vStatus7500, vStatus7501 , vStatus7524 string

	var i int = 0
	for rows.Next() {
		i = i + 1
		rows.Scan(&id,&name)
		wg.Add(1)
		go func(id int, name string) {
			// Создание контекста с ограничением времени его жизни в 4 сек
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()

			go work(ctx, id, w )
			wg.Wait()
		}(id, name)


		nameip = name
		vStatus7502, vStatus7522, vStatus7500, vStatus7501 , vStatus7524 = "", "", "", "",""
		port = "7502"
		vStatus7502 = checkStatus(i, id, nameip, port)
		port = "7522"
		vStatus7522 = checkStatus(i, id, nameip, port)
		port = "7500"
		vStatus7500 = checkStatus(i, id, nameip, port)
		port = "7501"
		vStatus7501 = checkStatus(i, id, nameip, port)
		port = "7524"
		vStatus7524 = checkStatus(i, id, nameip, port)




		//db.Model(&n).Where("name = ?", nameip).Update("status", vStatus).Error()
		db.Exec("UPDATE nsis SET updated_at=NOW(), status=? , statussdo=? , statusupd=? , statusauth=? , statustrans=? WHERE name = ?", vStatus7502, vStatus7522, vStatus7500, vStatus7501, vStatus7524, nameip)

	}

//!	if err := scanner.Err(); err != nil {
//!		fmt.Println(os.Stderr, "reading standard input:", err)
//!	}
	t1 := time.Now()
	log.Printf("СТОП. Время выполнения %v сек.\n", t1.Sub(t0))

}

func checkStatus(i int, id int, ip string, port string) string {
	fmt.Printf("%d\tid=%d\t%s:%s\n", i,id,ip, port)

	client := http.Client{
		Timeout: time.Duration(6 * time.Second),
	}

	resp, err := client.Get("http://" + ip + ".main.russianpost.ru" + ":" + port + "/v1")
	if err != nil {
		return "Error response" + ":" + port
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "Error Read conn"
	}
	var status, version string = "", ""

	var dat map[string]interface{}
	fmt.Printf("%s", body)
	if err := json.Unmarshal(body, &dat); err != nil {
		return "Error Unmarshal"
	} else {

		version = fmt.Sprintf("%s", dat["version"])
		status = fmt.Sprintf("%s", dat["status"])

		//fmt.Println(dat)
		version = fmt.Sprintf("%s", dat["version"])
		status = fmt.Sprintf("%s", dat["status"])

		log.Printf("%s\t%s\t%s\n", ip, status, version)
		//n.Status = fmt.Sprintf("\t%s\t%s", status, version)
		vStatus := fmt.Sprintf("\t%s\t%s", status, version)
		return vStatus
	}

}

// work() - функция выполнения запроса и получения результата.
// Результатом работы является запись в структуру значения ID-идентификатора запроса 
// и результата ответа сервера или 
// статус прерывания работы при достижении ограничения времени жизни контекста запроса
// Параметры: 
// ctx context.Context - контекст запроса
// id string идентификатор запроса
// dict *words - указатель на структуру хранения результатов выполнения запросов
func work(ctx context.Context, id string, dict *words) error {
	defer wg.Done()
	//Формируем структуру заголовков запроса
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}

	// канал для распаковки данных anonymous struct to pack and unpack data in the channel
	c := make(chan struct {
		r   *http.Response
		err error
	}, 1)

	req, _ := http.NewRequest("GET", "http://localhost:1111", nil)
	go func() {
		resp, err := client.Do(req)
		fmt.Printf("Doing http request, %s \n",id)
              
              //Добавим запись в результат статусов выполнения запросов
               dict.add(id,"StartWork")

		pack := struct {
			r   *http.Response
			err error
		}{resp, err}
		c <- pack
	}()
	
        // Кто первый того и тапки...	
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for client.Do
		fmt.Printf("Cancel context, НЕ ДОЖДАЛИСЬ ОТВЕТА СЕРВЕРА на запрос %s\n",id)
              //Добавим результат выполнения запроса со статусом CancelContext
               dict.add( id,"CancelContext")

		return ctx.Err()
	case ok := <-c:
		err := ok.err
		resp_ := ok.r
		if err != nil {
			fmt.Println("Error ", err)
			return err
		}
		defer resp_.Body.Close()
		out, _ := ioutil.ReadAll(resp_.Body)
		fmt.Printf("Server Response %s:  [%s]\n", id,out)

              //Добавим результат выполнения запроса Ответ сервера
               dict.add(id, string(out))

	}
    
	return nil
}




//функция читает файл со списком адресов , добавляет новые в БД к имеющимся
func check_nsi(db *gorm.DB, f *os.File) error {
	scanner := bufio.NewScanner(f)
	var nameip string
	var status string = "NewRecord"
	n := Nsi{}
	log.Println("func check_nsi")
	for scanner.Scan() {
		nameip = strings.TrimSpace(fmt.Sprintf("%s", scanner.Text()))
		log.Printf("scanner %s\n", nameip)
		// Наполним справочник данными из входящего файла
		n = Nsi{}
		newNsi := &Nsi{
			Name:   nameip,
			Status: status,
		}
		if err := db.Where("name = ?", nameip).First(&n).Error; err != nil {
			// error handling...
			if gorm.IsRecordNotFoundError(err) {
				db.Create(newNsi) // newNsi not nsi
				log.Printf("NewRecord %s\n", nameip)
			}
		} else {
			//db.Model(&n).Where("id = ?", nameip).Update("status", "newrecord")
			// Запись найдена в БД, переходим к следующему циклу
			continue
		}

	}
	return nil
}
