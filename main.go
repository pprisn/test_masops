//Программа сбора данных версий МАСОПС на сети УФПС Липецкой области
//http://localhost:7502/v1/
// Читаем отклик службы и результат сохраняем в БД
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
        "encoding/json"
	"log"
	"net"
	"os"
	"strings"
	"time"
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

update
http://localhost:7504/v1/

easuser
http://localhost:7501/v1/
*/

type Nsi struct {
	gorm.Model
	Name       string `gorm:"type:varchar(100);unique;not null"`
	Status     string `gorm:"type:varchar(255)"`
	Statussdo  string `gorm:"type:varchar(255)"`
	Statusupd  string `gorm:"type:varchar(255)"`
	Statusauth string `gorm:"type:varchar(255)"`
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
	scanner := bufio.NewScanner(f)
	var nameip string
	d := net.Dialer{Timeout: time.Second * 4}

	// Получить выборку
	rows, err := db.Raw("select name from nsis").Rows()
	defer rows.Close()

	var name string
	//var version string
	var port string
	var vStatus7502, vStatus7522, vStatus7504, vStatus7501 string

	var i int = 0
	for rows.Next() {
		i = i + 1
		rows.Scan(&name)
		nameip = name
		vStatus7502, vStatus7522, vStatus7504, vStatus7501 = "", "", "", ""
		port = "7502"
		vStatus7502 = checkStatus(d, i, nameip, port)
		port = "7522"
		vStatus7522 = checkStatus(d, i, nameip, port)
		port = "7504"
		vStatus7504 = checkStatus(d, i, nameip, port)
		port = "7501"
		vStatus7501 = checkStatus(d, i, nameip, port)

		//db.Model(&n).Where("name = ?", nameip).Update("status", vStatus).Error()
		db.Exec("UPDATE nsis SET updated_at=NOW(), status=? , statussdo=? , statusupd=? , statusauth=? WHERE name = ?", vStatus7502, vStatus7522, vStatus7504, vStatus7501, nameip)

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(os.Stderr, "reading standard input:", err)
	}
	t1 := time.Now()
	log.Printf("СТОП. Время выполнения %v сек.\n", t1.Sub(t0))

}

func checkStatus(d net.Dialer, i int, ip string, port string) string {
	fmt.Printf("%d\t%s:%s\n", i, ip, port)
	conn, err := d.Dial("tcp", ip+".main.russianpost.ru"+":"+port)
	if err != nil {
		// handle error
		log.Printf("\tError\t%s\t%s", ip, err)
		er := fmt.Sprintf("Error Dial:[%s]{%s}", ip, err)
		return er
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	fmt.Fprintf(conn, "GET /v1/ HTTP/1.0\r\n\r\n")
	jsonStatus := bufio.NewReader(conn) //.ReadString('\n')

	var status, version string = "", ""

	/*
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

			if strings.Contains(line, `"version":`) {
				version = line[:len(line)-1]

			} else if strings.Contains(line, `"status"`) {
				status = line[:len(line)-1]

			}
			if len(version) > 0 && len(status) > 0 {
				break
			}

		} //
	*/

	var dat map[string]interface{}
	if err := json.Unmarshal(jsonStatus, &dat); err != nil {
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

/*


	fmt.Printf("%d\t%s:7502\n", i, nameip)
	conn, err := d.Dial("tcp", nameip+".main.russianpost.ru"+":7502")
	if err != nil {
		// handle error
		log.Printf("\tError\t%s\t%s", nameip, err)
		fmt.Printf("\tError\t%s\t%s\n", nameip, err)
		continue
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	fmt.Fprintf(conn, "GET /v1/ HTTP/1.0\r\n\r\n")
	jsonStatus := bufio.NewReader(conn) //.ReadString('\n')


	fmt.Printf("%d\t%s:7522\n", i, nameip)
	conn, err := d.Dial("tcp", nameip+".main.russianpost.ru"+":7522")
	if err != nil {
		// handle error
		log.Printf("\tError\t%s\t%s", nameip, err)
		fmt.Printf("\tError\t%s\t%s\n", nameip, err)
		continue
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	fmt.Fprintf(conn, "GET /v1/ HTTP/1.0\r\n\r\n")
	jsonStatus7522 := bufio.NewReader(conn) //.ReadString('\n')


	fmt.Printf("%d\t%s:7504\n", i, nameip)
	conn, err := d.Dial("tcp", nameip+".main.russianpost.ru"+":7504")
	if err != nil {
		// handle error
		log.Printf("\tError\t%s\t%s", nameip, err)
		fmt.Printf("\tError\t%s\t%s\n", nameip, err)
		continue
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	fmt.Fprintf(conn, "GET /v1/ HTTP/1.0\r\n\r\n")
	jsonStatus7504 := bufio.NewReader(conn) //.ReadString('\n')


	fmt.Printf("%d\t%s:7501\n", i, nameip)
	conn, err := d.Dial("tcp", nameip+".main.russianpost.ru"+":7501")
	if err != nil {
		// handle error
		log.Printf("\tError\t%s\t%s", nameip, err)
		fmt.Printf("\tError\t%s\t%s\n", nameip, err)
		continue
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	fmt.Fprintf(conn, "GET /v1/ HTTP/1.0\r\n\r\n")
	jsonStatus7501 := bufio.NewReader(conn) //.ReadString('\n')

*/
