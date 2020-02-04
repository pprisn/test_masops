package main

//import s "strings"
import "fmt"
import "encoding/json"

func main(){

byt := []byte(`{"mid":"sdo01","version":"1.2.115","name":"Почта России. Сводная денежная отчетность.","type":"system","status":"OK","now":"2020-01-31T14:05:10.967","post":398000,"win":0,"develop":{"name":"ЦАИТС - ОСП ФГУП \"Почта России\"","phone":["8-812-327-39-70"],"email":["itsm.support-c00@russianpost.ru"],"web":["www.pochta.ru"]},"support":{"name":"ЦАИТС - ОСП ФГУП \"Почта России\"","phone":["8-812-327-39-70"],"email":["itsm.support-c00@russianpost.ru"],"web":["www.pochta.ru"]},"services":[],"configstatus":"OK","confighost":"http://R48-398000-N:7500/v1","hostname":"R48-398000-N"}`)

var dat map[string]interface{}

if err := json.Unmarshal(byt, &dat); err != nil {
    panic(err)
}
fmt.Println(dat)
fmt.Println(dat["version"])
fmt.Println(dat["status"])
}