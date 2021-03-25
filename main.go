package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jiliaevyp/web_yp_project/server"
	_ "github.com/lib/pq"
	"net"
	"os"
	"strconv"
)

var (
	IPaddrWeb, addrWeb, webPort string
	errserv                     int
	ErrInvalidPort              = errors.New("invalid port number")
	ErrInvalidIPaddress         = errors.New("invalid IP address")
)

const (
	answerServer     = "Hello, I am a server."
	readyServer      = "I'm ready!"
	defaultNet       = "tcp"
	defaultIp        = "192.168.1.101"
	defaultLocalhost = "localhost"
	defaultPort      = "8181"
)

// проверка на ввод  'Y = 1
func yesNo() int {
	var yesNo string
	len := 4
	data := make([]byte, len)
	n, err := os.Stdin.Read(data)
	yesNo = string(data[0 : n-1])
	if err == nil && (yesNo == "Y" || yesNo == "y" || yesNo == "Н" || yesNo == "н") {
		return 1
	} else {
		return 0
	}
}

// ввод  IP адреса сервера
func inpIP() (string, int) {
	data := ""
	err := 1
	for err == 1 {
		fmt.Print("Локальный сервера по умолчанию:	", defaultLocalhost, "\n", "Для изменения нажмите 'Y' ")
		yes := yesNo()
		if yes != 1 {
			data = defaultLocalhost
			err = 0
		} else {
			for err == 1 {
				fmt.Print("IP адрес сервера по умолчанию:	", defaultIp, "\n", "Для изменения нажмите 'Y' ")
				yes := yesNo()
				if yes != 1 {
					data = defaultIp
					err = 0
				} else {
					fmt.Println("Введите IP адрес сервера:	")
					fmt.Scanf(
						"%s\n",
						&data,
					)
					iperr := net.ParseIP(data)
					if iperr == nil {
						fmt.Println(ErrInvalidIPaddress)
						return data, 1
					} else {
						err = 0
					}
				}
			}
		}
	}
	return data, err
}

//ввод порта сервера
func inpPort() (string, int) {
	var (
		webPort string
	)
	err := 1
	for err == 1 {
		fmt.Print("Порт по умолчанию:	", defaultPort, "\n", "Для изменения нажмите 'Y' ")
		yes := yesNo()
		if yes != 1 {
			webPort = defaultPort
			err = 0
		} else {
			fmt.Print("Введите порт:	")
			fmt.Scanf(
				"%s\n",
				&webPort,
			)
			res, err1 := strconv.ParseFloat(webPort, 16)
			res = res + 1
			err = 0
			if err1 != nil {
				fmt.Println(ErrInvalidPort)
				return ":" + webPort, 1
			}
		}
	}
	return webPort, 0
}

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
			IPaddrWeb, err = inpIP() // добавить после 5 попыток
			if err == 1 {
				fmt.Println("Web адрес некорректен")
			}
		}
		fmt.Println("Aдрес web сервера установлен:	", IPaddrWeb)
		err = 1
		for err == 1 {
			fmt.Print("Введите порт:	")
			webPort, err = inpPort() // добавить после 5 попыток
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
		yes = yesNo() //yesNo()
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
		komand = yesNo()
	}
	fmt.Println("Рад был для Вас сделать что-то полезное !")
	fmt.Print("Обращайтесь в любое время без колебаний!", "\n", "\n")
}
