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
	//"strconv"
	"sort"
	"strings"
	"sync"
	"time"

	//         _ "github.com/mattn/go-sqlite3"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//var fsrc = flag.String("fsrc", "./fsrc.txt", `Файл с данными адресов для мониторинга отклика работы службы МАСОПС`)
var ufps = flag.String("ufps", "R48", `Список ID УФПС на запуск и сканирование`)
var plog = flag.String("plog", "R48", `Признак логирования. full  полное, по умолчанию минимальное `)

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

func jsElement(port string) string {
	if "7502" == port {
		return "\"Status\": \""
	} else if "7522" == port {
		return "\"Statussdo\": \""
	} else if "7500" == port {
		return "\"Statusupd\": \""
	} else if "7501" == port {
		return "\"Statusauth\": \""
	} else if "7524" == port {
		return "\"Statustrans\": \""
	}
	return "0"
}

//структура для хранения результатов
type words struct {
	sync.Mutex //добавить в структуру мьютекс
	found      map[int]string
}

//Инициализация области памяти
func newWords() *words {
	return &words{found: map[int]string{}}
}

//Фиксируем вхождение слова
func (w *words) add(word int, WS string) {
	w.Lock()         //Заблокировать объект
	defer w.Unlock() // По завершению, разблокировать
	WorkStatus, ok := w.found[word]
	if !ok { //т.е. если ID запроса не найдено заводим новый элемент слайса
		w.found[word] = WS
		return
	}
	// слово найдено в очередной раз , увеличим счетчик у элемента слайса
	w.found[word] = WorkStatus + "," + WS
}

var (
	wg sync.WaitGroup
)

func main() {
	var wg2 sync.WaitGroup

	//Создание структуры хранения результатов
	w := newWords()
	//	var err error
	loging := os.Getenv("LOGDB")
	if loging == "" {
		loging = "root"
	}
	//! fmt.Printf("LOGDB=%s\n", loging)

	db, err := gorm.Open("mysql", loging+"@(localhost)/masops?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()

	db.AutoMigrate(&Nsi{})

	flag.Parse()
	//var floger, f *os.File
	var floger *os.File
         

	listufps := strings.Split(*ufps, ",")

	if floger, err = os.OpenFile("mas.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		panic(err)
	}
	defer floger.Close()

	log.SetOutput(floger)
	t0 := time.Now()
	log.Printf("СТАРТ %v %v \n", t0, listufps)
	///////////////////////////////////////////////////
	//	if f, err = os.Open(*fsrc); err != nil {
	//		log.Printf("Error open %s \n", *fsrc)
	//		panic(err)
	//	}
	//	defer f.Close()

	ports := [5]string{"7502", "7522", "7500", "7501", "7524"}
	var id int
	var name string
	for _, Ufps := range listufps {
		// fmt.Printf("Ufps = %s", Ufps)
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
			//		idstr = strconv.Itoa(id)
			//			fmt.Printf("Scan nsis i= %d  id = %s , name = %s \n", i, id, name)
			for _, port := range ports {
				wg2.Add(1) //!required
				go func(id int, name string, port string) {
					defer wg2.Done() //!required
					// Создание контекста с ограничением времени его жизни в 5 сек
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					wg2.Add(1) //!required
					go checkStatus(ctx, id, name, port, w, &wg2)
					time.Sleep(5050 * time.Millisecond) //!reuired more then timeout
				}(id, name, port)

			}
		}
	}
	wg2.Wait() //!required

	// To store the keys in slice in sorted order
	var keys []int
	for k, _ := range w.found {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// To perform the opertion you want
	for _, k := range keys {
		//! fmt.Println("Key:", k, "Value:", w.found[k])
		var dat map[string]interface{}
		err := json.Unmarshal([]byte("{ "+w.found[k]+" }"), &dat)
		if err != nil {
			//!fmt.Printf("ErrorUnmarshal id = %d \t%s\n", k, "{ " + w.found[k] + " }")
			continue
		} else {
					db.Exec("UPDATE nsis SET updated_at=NOW(), status=? , statussdo=? , statusupd=? , statusauth=? , statustrans=? WHERE id = ?",
						dat["Status"], dat["Statussdo"], dat["Statusupd"], dat["Statusauth"], dat["Statustrans"], k)

			//!fmt.Printf("OK   Unmarshal id = %d \t%s\n", k, "{ " + w.found[k] + " }")

		}
	}
	//	for key, value := range w.found {
	//		fmt.Println(key, "\t", value)
	//	}

	t1 := time.Now()
	log.Printf("СТОП. Время выполнения %v сек. %v \n", t1.Sub(t0), listufps)

}
func checkStatus(ctx context.Context, id int, ip string, port string, dict *words, wg2 *sync.WaitGroup) error {
	defer wg2.Done() //!required
	//Формируем структуру заголовков запроса ожидаем отклик до 4 сек
	tr := &http.Transport{}
	client := &http.Client{Transport: tr, Timeout: time.Duration(5 * time.Second)}
	//client := &http.Client{Transport: tr}
	// канал для распаковки данных anonymous struct to pack and unpack data in the channel
	c := make(chan struct {
		r   *http.Response
		err error
	}, 1)
	defer close(c)
	req, _ := http.NewRequest("GET", "http://"+ip+".main.russianpost.ru"+":"+port+"/v1", nil)
	req.WithContext(ctx)
	vStatus := ""
	wg2.Add(1) //!required
	go func() {
		defer wg2.Done() //!required
		resp, err := client.Do(req)
		//fmt.Printf("Doing http request, %s \n", id)
		//пишем в канал данные ответа сервера или ошибку
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
		//fmt.Printf("Cancel context, НЕ ДОЖДАЛИСЬ ОТВЕТА СЕРВЕРА на запрос %s\n", id)
		//Добавим результат выполнения запроса со статусом CancelContext
		//key := id + ";" + port
		vStatus = jsElement(port) + "Error cancel context" + ":" + port + "\""
		dict.add(id, vStatus)
		if *plog =="full" {
			log.Printf("%d\t%s:%s\t%s\n", id, ip, port, vStatus)
		}
		//!		fmt.Printf("Server Response %d;%s  [%s]\n", id, port, vStatus)
		return ctx.Err()
	case ok := <-c:
		err := ok.err
		resp_ := ok.r
		if err != nil {
			vStatus = jsElement(port) + "Error response" + ":" + port + "\""
		} else {
			defer resp_.Body.Close()
			body, err := ioutil.ReadAll(resp_.Body)
			if err != nil {
				vStatus = jsElement(port) + "Error Read conn" + "\""
			} else {
				var status, version string = "", ""
				var dat map[string]interface{}
				// fmt.Printf("%s", body)
				err := json.Unmarshal(body, &dat)
				if err != nil {
					vStatus = jsElement(port) + "Error Unmarshal" + "\""
				} else {
					version = fmt.Sprintf("%s", dat["version"])
					status = fmt.Sprintf("%s", dat["status"])
					vStatus = jsElement(port) + fmt.Sprintf("%s %s", strings.TrimSpace(status), strings.TrimSpace(version)) + "\""
				}
			}
		}
		//Добавим результат выполнения запроса Ответ сервера
		dict.add(id, vStatus)
		if *plog =="full" {
			log.Printf("%d\t%s:%s\t%s\n", id, ip, port, vStatus)
		}
		//!		fmt.Printf("Server Response ID=%d port=%s Status=%s\n", id, port, vStatus)
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
