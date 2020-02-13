// Программа сбора данных версий МАСОПС на сети c
// применением алгоритма конкурентного  программирования !
// Читаем отклик служб и результат сохраняем в БД
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	//	"io"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	//         _ "github.com/mattn/go-sqlite3"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var fsrc = flag.String("fsrc", "./fsrc.txt", `Файл с данными адресов для мониторинга отклика работы службы МАСОПС`)
var ufps = flag.String("ufps", "R48", `Список ID УФПС на запуск и сканирование`)

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
	Name        string `gorm:"type:varchar(100);unique;not null"`
	Status      string `gorm:"type:varchar(255)"` //:7502
	Statussdo   string `gorm:"type:varchar(255)"` //:7522
	Statusupd   string `gorm:"type:varchar(255)"` //:7500
	Statusauth  string `gorm:"type:varchar(255)"` //:7501
	Statustrans string `gorm:"type:varchar(255)"` //:7524
	Note        string `gorm:"type:varchar(255)"`
	ufpsid      string `gorm:"type:varchar(32)"`
}


//структура для хранения результатов
type words struct {
	sync.Mutex //добавить в структуру мьютекс
	found      map[string]string
}

//Инициализация области памяти
func newWords() *words {
	return &words{found: map[string]string{}}
}

//Фиксируем вхождение слова
func (w *words) add(word string, WS string) {
	w.Lock()         //Заблокировать объект
	defer w.Unlock() // По завершению, разблокировать
	WorkStatus, ok := w.found[word]
	if !ok { //т.е. если ID запроса не найдено заводим новый элемент слайса
		w.found[word] = WS
		return
	}
	// слово найдено в очередной раз , увеличим счетчик у элемента слайса
	w.found[word] = WorkStatus + " ; " + WS
}

var (
	wg sync.WaitGroup
)

func main() {

	//Создание структуры хранения результатов
	w := newWords()

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

	listufps := strings.Split(*ufps, ",")

	if floger, err = os.OpenFile("mas.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		panic(err)
	}
	defer floger.Close()

	log.SetOutput(floger)
	t0 := time.Now()
	log.Printf("СТАРТ %v %v \n", t0, listufps)

	if f, err = os.Open(*fsrc); err != nil {
		log.Printf("Error open %s \n", *fsrc)
		panic(err)
	}
	defer f.Close()

	ports := [5]string{"7502", "7522", "7500", "7501", "7524"}
	var id int
	var idstr string
	var name string

	//	d := net.Dialer{Timeout: time.Second * 4}
	for _, Ufps := range listufps {
		fmt.Printf("Ufps = %s", Ufps)
		// Получить выборку
		rows, err := db.Raw("select id, name from nsis where ufpsid = ?", Ufps).Rows()
		defer rows.Close()
		if err != nil {
			continue
		}
		i := 0
		for rows.Next() {
			i = i + 1
			rows.Scan(&id, &name)
			idstr = strconv.Itoa(id)
			fmt.Printf("Scan nsis i= %d  id = %s , name = %s \n", i, idstr, name)
			time.Sleep(6000 * time.Microsecond) //new
			for _, port := range ports {
				wg.Add(1)
                                //wg.Add(1)
                                //wg.Add(1)
				time.Sleep(37000 * time.Microsecond) //Влияет на 25000 полноту сбора данных (нивелирует действия defer
				go func(idstr string, name string, port string) {
                                        defer wg.Done() //!new
					// Создание контекста с ограничением времени его жизни в 5 сек
					ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
					defer cancel()
                                        wg.Add(1) //new !.!
						go checkStatus(ctx, idstr, name, port, w)
					wg.Wait()
				}(idstr, name, port)

				//		nameip = name
				//		vStatus7502, vStatus7522, vStatus7500, vStatus7501, vStatus7524 = "", "", "", "", ""
				//		port = "7502"
				//		vStatus7502 = checkStatus(i, id, nameip, port)
				//		port = "7522"
				//		vStatus7522 = checkStatus(i, id, nameip, port)
				//		port = "7500"
				//		vStatus7500 = checkStatus(i, id, nameip, port)
				//		port = "7501"
				//		vStatus7501 = checkStatus(i, id, nameip, port)
				//		port = "7524"
				//		vStatus7524 = checkStatus(i, id, nameip, port)
				//		//db.Model(&n).Where("name = ?", nameip).Update("status", vStatus).Error()
				//		db.Exec("UPDATE nsis SET updated_at=NOW(), status=? , statussdo=? , statusupd=? , statusauth=? , statustrans=? WHERE name = ?", vStatus7502, vStatus7522, vStatus7500, vStatus7501, vStatus7524, nameip)
       			}
		}
                wg.Wait() //!new
		time.Sleep(5500000 * time.Microsecond)
		t1 := time.Now()
		log.Printf("СТОП. Время выполнения %v сек.\n", t1.Sub(t0))

	}
}
func checkStatus(ctx context.Context, id string, ip string, port string, dict *words) error {
	defer wg.Done() //!new
	//Формируем структуру заголовков запроса ожидаем отклик до 4 сек
	tr := &http.Transport{}
	client := &http.Client{Transport: tr, Timeout: time.Duration(5 * time.Second)} 

	// канал для распаковки данных anonymous struct to pack and unpack data in the channel
	c := make(chan struct {
		r   *http.Response
		err error
	}, 1)
	req, _ := http.NewRequest("GET", "http://"+ip+".main.russianpost.ru"+":"+port+"/v1", nil)
	vStatus := ""
	wg.Add(1) //!new
	go func() {
          defer wg.Done() //!new
		resp, err := client.Do(req)
		fmt.Printf("Doing http request, %s \n", id)
		//dict.add(id, "StartWork")
		//пишем в канал данные ответа сервера или ошибку
		pack := struct {
			r   *http.Response
			err error
		}{resp, err}
		c <- pack
 	        //wg.Wait() //!
	}()
	// Кто первый того и тапки...
       //wg.Wait() //!
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for client.Do
		//fmt.Printf("Cancel context, НЕ ДОЖДАЛИСЬ ОТВЕТА СЕРВЕРА на запрос %s\n", id)
		//Добавим результат выполнения запроса со статусом CancelContext
		key := id + ";" + port
		vStatus = "Error cancel context" + ":" + port
		dict.add(key, vStatus)
		log.Printf("%s\t%s:%s\t%s\n", id, ip, port, vStatus)
		fmt.Printf("Server Response %s;%s  [%s]\n", id, port, vStatus)
		return ctx.Err()
	case ok := <-c:
		err := ok.err
		resp_ := ok.r
		if err != nil {
			vStatus = "Error response" + ":" + port
		} else {
			defer resp_.Body.Close()
			body, err := ioutil.ReadAll(resp_.Body)
			if err != nil {
				vStatus = "Error Read conn"
			} else {
				var status, version string = "", ""
				var dat map[string]interface{}
				// fmt.Printf("%s", body)
				err := json.Unmarshal(body, &dat)
				if err != nil {
					vStatus = "Error Unmarshal"
				} else {
					version = fmt.Sprintf("%s", dat["version"])
					status = fmt.Sprintf("%s", dat["status"])
					//fmt.Println(dat)
					//version = fmt.Sprintf("%s", dat["version"])
					//status = fmt.Sprintf("%s", dat["status"])
					//log.Printf("%s\t%s:%s\t%s\t%s\n", id,ip, port, status, version)
					vStatus = fmt.Sprintf("\t%s\t%s", status, version)
				}
			}
		}
		//Добавим результат выполнения запроса Ответ сервера
		key := id + ";" + port
		dict.add(key, vStatus)
		log.Printf("%s\t%s:%s\t%s\n", id, ip, port, vStatus)
		fmt.Printf("Server Response %s;%s  [%s]\n", id, port, vStatus)
	} //select

	return nil
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
		fmt.Printf("Doing http request, %s \n", id)

		//Добавим запись в результат статусов выполнения запросов
		dict.add(id, "StartWork")

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
		fmt.Printf("Cancel context, НЕ ДОЖДАЛИСЬ ОТВЕТА СЕРВЕРА на запрос %s\n", id)
		//Добавим результат выполнения запроса со статусом CancelContext
		dict.add(id, "CancelContext")
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
		fmt.Printf("Server Response %s:  [%s]\n", id, out)
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
