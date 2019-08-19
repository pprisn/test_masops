//http://localhost:7502/v1/

package main

import (
	"bufio"
	"flag"
	"fmt"
//	"io"
	"log"
	"net"
	"os"
	"time"
	"strings"
)

var fsrc = flag.String("fsrc", "fsrc.txt", `Файл с данными адресов для мониторинга отклика работы службы МАСОПС`)
var mode = flag.String("mode", "l", `Режим логирования отклика службы, l краткий, f полный`)

func main() {

	var err error
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

	scanner := bufio.NewScanner(f)
	var nameip string
	for scanner.Scan() {
		nameip = strings.TrimSpace(fmt.Sprintf("%s", scanner.Text()))
		fmt.Printf("%s\n", nameip)

		conn, err := net.Dial("tcp", nameip+":7502")
		if err != nil {
			// handle error
                  	log.Printf("Error %s\t%s", nameip, err)
      		        fmt.Printf("Error %s\t%s\n", nameip, err)
			continue
		}

		fmt.Fprintf(conn, "GET /v1/ HTTP/1.0\r\n\r\n")

		status := bufio.NewReader(conn) //.ReadString('\n')
		for {
			line, err := status.ReadString('\n')
			if err != nil {
				log.Printf("%s\t%s", nameip, err)
				break
			}
                
			log.Printf("%s\t%s", nameip, line)
			if *mode =="l" {
                        	break //Читаем 1ну успешную строку
			} else if *mode=="f" {
			  continue
			}


		}


	}
	if err := scanner.Err(); err != nil {
		fmt.Println(os.Stderr, "reading standard input:", err)
	}
	t1 := time.Now()
	log.Printf("Время выполнения сканирования %v сек.\n", t1.Sub(t0))

}
