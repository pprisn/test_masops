//Программа мониторинга версий МАСОПС на сети УФПС Липецкой области
//http://localhost:7502/v1/
// Читаем отклик службы и результат сохраняем в БД
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
	//         _ "github.com/mattn/go-sqlite3"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var fsrc = flag.String("fsrc", "fsrc.txt", `Файл с данными адресов для мониторинга отклика работы службы МАСОПС`)
var mode = flag.String("mode", "l", `Режим логирования отклика службы, l краткий, f полный`)

type Nsi struct {
	gorm.Model
	Name   string `gorm:"type:varchar(100);unique;not null"`
	Status string `gorm:"type:varchar(255)"`
}

func main() {

	//	var err error
        loging := os.Getenv("LOGBD")
	//db, err := gorm.Open("sqlite3", "masops.db?cache=shared&mode=rwc")
	db, err := gorm.Open("mysql",loging+"@/masops?charset=utf8&parseTime=True&loc=Local")
	//        db.SetMaxOpenConns(1)
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
		panic(err)
	}
	defer f.Close()

	err = check_nsi(db, f)

	scanner := bufio.NewScanner(f)
	var nameip string
	d := net.Dialer{Timeout: time.Second * 4}

	// Получить выборку
	rows, err := db.Raw("select name from nsis").Rows()
	defer rows.Close()

	var name string
	var version string
	var status string
	//	var ErrNew bool = false

	//	var Notfound bool
	//	var currentId uint
	// цикл по списку адресов ПК-ip

	//        var result map[string]interface{}
	var i int = 0
	for rows.Next() {
		i = i + 1
		rows.Scan(&name)
		nameip = name
		fmt.Printf("%d\t%s:7502\n", i, nameip)
		conn, err := d.Dial("tcp", nameip+":7502")
		if err != nil {
			// handle error
			log.Printf("\tError\t%s\t%s", nameip, err)
			fmt.Printf("\tError\t%s\t%s\n", nameip, err)
			continue
		}
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))

		fmt.Fprintf(conn, "GET /v1/ HTTP/1.0\r\n\r\n")
		jsonStatus := bufio.NewReader(conn) //.ReadString('\n')

		version = ""
		status = ""
		for {
			line, err := jsonStatus.ReadString('\n')
			if len(line) == 0 && err != nil {
				if err == io.EOF {
					break
				}
                                 break
				//return err
			}
			line = strings.TrimSuffix(line, "\n")

//			if err != nil {
//				if err == io.EOF {
//					break
//				} else {
//					fmt.Println(err)
//					continue
//				}
//			}

			if strings.Contains(line, `"version":`) {
				version = line[:len(line)-1]
			} else if strings.Contains(line, `"status"`) {
				status = line[:len(line)-1]
			}

			n := Nsi{}
			n.Name = nameip
			//                        db.First(&n,"name = ?",nameip)

			if *mode == "l" {
				log.Printf("%s\t%s\t%s\n", nameip, status, version)
				n.Status = fmt.Sprintf("\t%s\t%s", status, version)
				//db.Create(&n)
				//db.Save(&n)
				db.Model(&n).Where("name = ?", nameip).Update("status", n.Status)
				//				ErrNew = db.NewRecord(j)
				//				if ErrNew == true {
				//					continue
				//				} else {
				//					db.Create(&j)
				//				}
			}

			if *mode == "f" {
				///
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Println(os.Stderr, "reading standard input:", err)
		}
		t1 := time.Now()
		log.Printf("СТОП. Время выполнения %v сек.\n", t1.Sub(t0))
	}
}

//функция читает файл со списком адресов , добавляет новые в БД к имеющимся
func check_nsi(db *gorm.DB, f *os.File) error {
	scanner := bufio.NewScanner(f)
	var nameip string
	var status string = ""
	var ErrNew bool
	n := Nsi{}
	//	var ErrNew bool = false
	for scanner.Scan() {
		nameip = strings.TrimSpace(fmt.Sprintf("%s", scanner.Text()))

		// Наполним справочник данными из входящего файла
		n = Nsi{Name: nameip, Status: status}
//		db.Create(&n)
		ErrNew = db.NewRecord(n) // => returns `true` as primary key is blank
	//		//fmt.Printf("%v \n",ErrNew)
		if ErrNew == true {
			continue
		} else {
			db.Create(&n)
		}

	}
	return nil
}
