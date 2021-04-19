package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jiliaevyp/web_yp_project/server"
	"github.com/jiliaevyp/web_yp_project/servfunc"
	_ "github.com/lib/pq"
)

var (
	IPaddrWeb, addrWeb, webPort string
	errserv                     int
	ErrInvalidPort              = errors.New("invalid port number")
	ErrInvalidIPaddress         = errors.New("invalid IP address")
)

const (
	answerServer = "Hello, I am a server."
	readyServer  = "I'm ready!"
)

func main() {
	var (
		err, yes int
	)
	IPaddrWeb = ""
	komand := 1
	fmt.Println("------------------------------------")
	fmt.Println("|          WEB server              |")
	fmt.Println("|    отвечаем на любые запросы!    |")
	fmt.Println("|                                  |")
	fmt.Println("|   (c) jiliaevyp@gmail.com        |")
	fmt.Println("------------------------------------")
	//// Создаем соединение с базой данных
	connStr := "user=yp password=12345 dbname=postgres sslmode=disable"
	db, err2 := sql.Open("postgres", connStr)
	if err2 != nil {
		fmt.Println("ошибка подключения к базе <geoplastdb>")
		panic(err2)
	} else {
		fmt.Println("база <geoplastdb> подключена!")
	}
	defer db.Close()
	for komand == 1 {
		err = 1
		for err == 1 {
			fmt.Println("Введите web адрес  сервера:	")
			IPaddrWeb, err = servfunc.InpIP() // добавить после 5 попыток
			if err == 1 {
				fmt.Println("Web адрес некорректен")
			}
		}
		fmt.Println("Aдрес web сервера установлен:	", IPaddrWeb)
		err = 1
		for err == 1 {
			fmt.Print("Введите порт:	")
			webPort, err = servfunc.InpPort() // добавить после 5 попыток
			if err == 1 {
				fmt.Println("Порт некорректен")
			}
		}
		addrWeb = IPaddrWeb + ":" + webPort
		fmt.Println("Сервер:  ", addrWeb, "\n")
		fmt.Println("Загрузите web страницу")
		//loadAnswer()
		fmt.Println("-------------------------------------------------")
		fmt.Println("Адрес сервера:         ", addrWeb)
		fmt.Println("-------------------------------------------------")
		fmt.Print("Запускаю сервер? (Y)   ")
		fmt.Println("Отменить?  (Enter)")
		yes = servfunc.YesNo() //yesNo()
		if yes == 1 {
			go server.Server(addrWeb, db)
			if server.Erserv != 0 {
				fmt.Print("*** Ошибка при загрузке сервера ***", "\n", "\n")
			} else {
				fmt.Println("---------------------------")
				fmt.Println(answerServer, "   ", addrWeb)
				fmt.Println(readyServer)
				fmt.Print("---------------------------", "\n")
			}
		} else {
			fmt.Print("\n", "Запуск отменен", "\n", "\n")
		}
		//<-time.After(5 * time.Second)
		//go client()
		fmt.Print("Перезапустить? (Y)   ")
		fmt.Println("Закончить?  (Enter)")
		komand = servfunc.YesNo()
	}
	fmt.Println("Рад был для Вас сделать что-то полезное !")
	fmt.Print("Обращайтесь в любое время без колебаний!", "\n", "\n")
}
