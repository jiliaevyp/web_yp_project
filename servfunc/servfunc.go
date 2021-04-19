package servfunc

import (
	"errors"
	"fmt"
	"github.com/ttacon/libphonenumber"
	"net"
	"net/http"
	"net/mail"
	"os"
	"strconv"
)

const (
	answerServer     = "Hello, I am a server."
	readyServer      = "I'm ready!"
	defaultNet       = "tcp"
	defaultIp        = "192.168.1.101"
	defaultLocalhost = "localhost"
	defaultPort      = "8181"
)

var (
	IPaddrWeb, addrWeb, webPort string
	errserv                     int
	ErrInvalidPort              = errors.New("invalid port number")
	ErrInvalidIPaddress         = errors.New("invalid IP address")
)

type person struct { // данные по сотруднику при вводе и отображении в personal.HTML
	Id         string
	Forename   string
	Title      string
	Kadr       string
	Numotdel   string
	Department string
	Email      string
	Phone      string
	Address    string
	Tarif      string // почасовая руб
	Jetzyahre  string
	Jetzmonat  string
	Ready      string // "1" - ввод корректен
	Errors     string // "1" - ошибка при вводе полей
	ErrPhone   string // "1"- ошибка при вводе телефона
	ErrEmail   string // "1"- ошибка при вводе email
	//ErrTitle  string // "1"- ошибка при вводе title
	ErrTarif  string // "1"- ошибка при вводе тарифа
	ErrNumotd string // "1"- ошибка при вводе номера отдела
	Empty     string // "1" - остались пустые поля
	ErrRange  string // "1" - выход за пределы диапазона
}

// проверка корректности емайл адреса nameAddress --> "имя <email@mail.com>
func InpMailAddress(nameAddress string) (err int, email string, title string) {
	e, err1 := mail.ParseAddress(nameAddress)
	if err1 != nil {
		return 1, e.Address, e.Name //"?", "?"
	}
	return 0, e.Address, e.Name
}

// валидация  числовых вводов и диапазонов
func Checknum(checknum string, min int, max int) int {
	num, err := strconv.Atoi(checknum)
	if err != nil {
		return 1
	} else {
		if num >= min && num <= max {
			return 0
		} else {
			return 1
		}
	}
}

// проверка на ввод  'Y = 1
func YesNo() int {
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
func InpIP() (string, int) {
	data := ""
	err := 1
	for err == 1 {
		fmt.Print("Локальный сервера по умолчанию:	", defaultLocalhost, "\n", "Для изменения нажмите 'Y' ")
		yes := YesNo()
		if yes != 1 {
			data = defaultLocalhost
			err = 0
		} else {
			for err == 1 {
				fmt.Print("IP адрес сервера по умолчанию:	", defaultIp, "\n", "Для изменения нажмите 'Y' ")
				yes := YesNo()
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
func InpPort() (string, int) {
	var (
		webPort string
	)
	err := 1
	for err == 1 {
		fmt.Print("Порт по умолчанию:	", defaultPort, "\n", "Для изменения нажмите 'Y' ")
		yes := YesNo()
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

// подготовка  значений для web
func MakeReadyHtml(p *person) {
	p.Ready = "0"     // 1 - ввод успешный
	p.Errors = "0"    // 1 - ошибки при вводе
	p.Empty = "0"     // 1 - есть пустые поля
	p.ErrRange = "0"  // 1 - выход за пределы диапазона
	p.ErrPhone = "0"  // 1 - ошибка в тлф номере
	p.ErrEmail = "0"  // 1 - ошибка в email
	p.ErrTarif = "0"  // 1 - ошибка в тарифе
	p.ErrNumotd = "0" // 1 - ошибка в номере отдела
	return
}

// подготовка и ввод значений из web
func ReadFromHtml(p *person, req *http.Request) {
	p.Title = req.Form["title"][0]
	p.Forename = req.Form["forename"][0]
	p.Kadr = req.Form["kadr"][0]
	p.Tarif = req.Form["tarif"][0]
	p.Numotdel = req.Form["numotdel"][0]
	p.Email = req.Form["email"][0]
	p.Phone = req.Form["phone"][0]
	p.Address = req.Form["address"][0]
	return
}

// проверка вводимых значений
func CheckNumer(personalhtml *person) int {
	var err int
	errout := 0
	err = Checknum(personalhtml.Tarif, 10, 1000)
	if err != 0 {
		personalhtml.ErrRange = "1"
		personalhtml.ErrTarif = "1"
		errout = 1
	}
	err = Checknum(personalhtml.Numotdel, 0, 20)
	if err != 0 {
		personalhtml.ErrRange = "1"
		personalhtml.ErrNumotd = "1"
		errout = 1
	}
	err, personalhtml.Email, _ = InpMailAddress(personalhtml.Title + "<" + personalhtml.Email + ">") // проверка email адреса
	if err > 0 {
		personalhtml.ErrEmail = "1"
		errout = 1
	}
	_, err1 := libphonenumber.Parse(personalhtml.Phone, "RU")
	if err1 != nil {
		personalhtml.ErrPhone = "1"
		errout = 1
	}
	if personalhtml.Forename == "" || personalhtml.Title == "" || personalhtml.Kadr == "" || personalhtml.Address == "" {
		personalhtml.Empty = "1"
		errout = 1
	}
	return errout
}
