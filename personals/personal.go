package personals

import (
	_ "crypto/dsa"
	"database/sql"
	"fmt"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ttacon/libphonenumber"
	_ "go/parser"
	"html/template"
	"log"
	//"net"
	"net/http"
	"strconv"
)

type person struct { // данные по сотруднику при вводе и отображении в personal.HTML
	Name      string
	Title     string
	Kadr      string
	Otdel     string
	Numotdel  string
	Email     string
	Phone     string
	Address   string
	Tarif     string // почасовая руб
	Ready     string // "1" - ввод корректен
	Errors    string // "1" - ошибка при вводе полей
	ErrPhone  string // "1"- ошибка при вводе телефона
	ErrEmail  string // "1"- ошибка при вводе email
	ErrTarif  string // "1"- ошибка при вводе тарифа
	ErrNumotd string // "1"- ошибка при вводе номера отдела
	Empty     string // "1" - остались пустые поля
	ErrRange  string // "1" - выход за пределы диапазона
}

type frombase struct { // строка  при чтении/записи из/в базы personaldb
	id       int
	name     string
	title    string
	kadr     string
	numotdel int
	otdel    string
	email    string
	phone    string
	address  string
	tarif    int // почасовая руб
}

var (
	personals struct {
		Ready       string
		Persontable []person //person // таблица по сотрудниам  в personals_index.html
	}
)

// подготовка  значений для web
func makeReadyHtml(p *person) {
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
func readFromHtml(p *person, req *http.Request) {
	p.Title = req.Form["title"][0]
	p.Name = req.Form["name"][0]
	p.Kadr = req.Form["kadr"][0]
	p.Tarif = req.Form["tarif"][0]
	p.Otdel = req.Form["otdel"][0]
	p.Numotdel = req.Form["numotdel"][0]
	p.Email = req.Form["email"][0]
	p.Phone = req.Form["phone"][0]
	p.Address = req.Form["address"][0]
	return
}

// проверка вводимых значений
func checkNumer(personalhtml *person) {
	var err int
	err = checknum(personalhtml.Tarif, 10, 1000)
	if err != 0 {
		personalhtml.ErrRange = "1"
		personalhtml.ErrTarif = "1"
		personalhtml.Errors = "1"
	}
	err = checknum(personalhtml.Numotdel, 0, 20)
	if err != 0 {
		personalhtml.ErrRange = "1"
		personalhtml.ErrNumotd = "1"
		personalhtml.Errors = "1"
	}
	err, personalhtml.Email, _ = inpMailAddress(personalhtml.Title + "<" + personalhtml.Email + ">") // проверка email адреса
	if err > 0 {
		personalhtml.Errors = "1"
		personalhtml.ErrEmail = "1"
	}
	_, err1 := libphonenumber.Parse(personalhtml.Phone, "RU")
	if err1 != nil {
		personalhtml.Errors = "1"
		personalhtml.ErrPhone = "1"
	}
	if personalhtml.Name == "" || personalhtml.Title == "" || personalhtml.Kadr == "" || personalhtml.Otdel == "" || personalhtml.Address == "" {
		personalhtml.Empty = "1"
		personalhtml.Errors = "1"
	}
}

// просмотр таблицы из personaldb
func personalsIndexHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/personals_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		del := req.URL.Query().Get("del")
		title := req.URL.Query().Get("title")
		if del == "del" {
			_, err = db.Exec("DELETE FROM personals WHERE title = $1", title)
			if err != nil { // удаление старой записи
				panic(err)
			}
		}
		personals.Persontable = nil

		rows, err1 := db.Query(`SELECT * FROM personals`)
		if err1 != nil {
			fmt.Println(" table Personals ошибка чтения ")
			panic(err1)
		}
		defer rows.Close()

		for rows.Next() {
			var p frombase
			err = rows.Scan( // пересылка  данных строки базы personals в "p"
				&p.id,
				&p.title,
				&p.name,
				&p.kadr,
				&p.tarif,
				&p.numotdel,
				&p.otdel,
				&p.email,
				&p.phone,
				&p.address,
			)
			if err != nil {
				fmt.Println("indexPersonals ошибка распаковки строки ")
				panic(err)
				return
			}
			var personalhtml person
			personalhtml.Name = p.name
			personalhtml.Title = p.title
			personalhtml.Kadr = p.kadr
			personalhtml.Tarif = strconv.Itoa(p.tarif)
			personalhtml.Numotdel = strconv.Itoa(p.numotdel)
			personalhtml.Otdel = p.otdel
			personalhtml.Email = p.email
			personalhtml.Phone = p.phone
			personalhtml.Address = p.address
			personalhtml.Ready = "1"
			personalhtml.Errors = "0"
			personalhtml.Empty = "0"
			personals.Persontable = append( // добавление строки в таблицу Personalstab для personals_index.html
				personals.Persontable,
				personalhtml,
			)
		}
		personals.Ready = "1"

		err = t.ExecuteTemplate(w, "base", personals)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

// просмотр записи из personaldb
func personalShowhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/personal_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		title := req.URL.Query().Get("title")
		row := db.QueryRow("SELECT * FROM personals WHERE title=$1", title)

		var personalhtml person
		personalhtml.Ready = "1"     // 1 - ввод успешный
		personalhtml.Errors = "0"    // 1 - ошибки при вводе
		personalhtml.Empty = "0"     // 1 - есть пустые поля
		personalhtml.ErrRange = "0"  // 1 - выход за пределы диапазона
		personalhtml.ErrPhone = "0"  // 1 - ошибка в тлф номере
		personalhtml.ErrEmail = "0"  // 1 - ошибка в email
		personalhtml.ErrTarif = "0"  // 1 - ошибка в тарифе
		personalhtml.ErrNumotd = "0" // 1 - ошибка в номере отдела

		// чтение строки из таблицы
		var p frombase
		err = row.Scan( // пересылка  данных строки базы personals в personrow
			&p.id,
			&p.title,
			&p.name,
			&p.kadr,
			&p.tarif,
			&p.numotdel,
			&p.otdel,
			&p.email,
			&p.phone,
			&p.address,
		)
		if err != nil {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		}
		personalhtml.Title = p.title
		personalhtml.Name = p.name
		personalhtml.Kadr = p.kadr
		personalhtml.Tarif = strconv.Itoa(p.tarif)
		personalhtml.Numotdel = strconv.Itoa(p.numotdel)
		personalhtml.Otdel = p.otdel
		personalhtml.Email = p.email
		personalhtml.Phone = p.phone
		personalhtml.Address = p.address

		err = t.ExecuteTemplate(w, "base", personalhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error personalShowhandler", http.StatusInternalServerError)
			return
		}
	}
}

// новая запись формы personal в базу personaldb
func personalNewhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/personal_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.

		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error personalNewhandler", http.StatusInternalServerError)
			return
		}
		var personalhtml person
		personalhtml.Ready = "0"

		if req.Method == "POST" {
			req.ParseForm()
			makeReadyHtml(&personalhtml) // подготовка значений для web
			//readFromHtml(&personalhtml, req)  	// ввод значений из web
			personalhtml.Title = req.Form["title"][0]
			personalhtml.Name = req.Form["name"][0]
			personalhtml.Kadr = req.Form["kadr"][0]
			personalhtml.Tarif = req.Form["tarif"][0]
			personalhtml.Otdel = req.Form["otdel"][0]
			personalhtml.Numotdel = req.Form["numotdel"][0]
			personalhtml.Email = req.Form["email"][0]
			personalhtml.Phone = req.Form["phone"][0]
			personalhtml.Address = req.Form["address"][0]
			checkNumer(&personalhtml) // проверка числовых значений
			if personalhtml.Errors == "0" {
				personalhtml.Ready = "1"
				//добавление записи в базу
				title := personalhtml.Title
				// удаление старой записи
				_, err1 := db.Exec("DELETE FROM personals WHERE title = $1", title)
				if err1 != nil {
					fmt.Println("Ошибка при удалении старой записи в personals title = ", title)
					panic(err)
				}
				var p frombase
				p.title = personalhtml.Title
				p.name = personalhtml.Name
				p.kadr = personalhtml.Kadr
				p.tarif, _ = strconv.Atoi(personalhtml.Tarif)
				p.otdel = personalhtml.Otdel
				p.numotdel, _ = strconv.Atoi(personalhtml.Numotdel)
				p.email = personalhtml.Email
				p.phone = personalhtml.Phone
				p.address = personalhtml.Address

				sqlStatement := `INSERT INTO personals (title,name,kadr,tarif,numotdel,otdel,email,phone,address) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
				_, err2 := db.Exec(sqlStatement,
					p.title,
					p.name,
					p.kadr,
					p.tarif,
					p.numotdel,
					p.otdel,
					p.email,
					p.phone,
					p.address,
				)
				if err2 != nil {
					fmt.Println("Ошибка записи новой строки в personalNew")
				}
			}
		}
		err = t.ExecuteTemplate(w, "base", personalhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

// редактирование формы personal и замена в базе personaldb
func personalEdithandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/personal_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		title := req.URL.Query().Get("title")
		row := db.QueryRow("SELECT * FROM personals WHERE title=$1", title)

		var personalhtml person
		makeReadyHtml(&personalhtml) // подготовка значений для web

		personalhtml.Empty = "1" // 1 - есть пустые поля

		// чтение строки из таблицы
		var p frombase
		err = row.Scan( // пересылка  данных строки базы personals в p
			&p.id,
			&p.title,
			&p.name,
			&p.kadr,
			&p.tarif,
			&p.numotdel,
			&p.otdel,
			&p.email,
			&p.phone,
			&p.address,
		)
		if err != nil {
			fmt.Println("edit --> ошибка распаковки строки при чтении записи title=", title)
			panic(err)
		} else {
			personalhtml.Title = p.title
			personalhtml.Name = p.name
			personalhtml.Kadr = p.kadr
			personalhtml.Tarif = strconv.Itoa(p.tarif)
			personalhtml.Numotdel = strconv.Itoa(p.numotdel)
			personalhtml.Otdel = p.otdel
			personalhtml.Email = p.email
			personalhtml.Phone = p.phone
			personalhtml.Address = p.address

			if req.Method == "POST" {
				req.ParseForm()
				makeReadyHtml(&personalhtml) // подготовка значений для web
				//readFromHtml(&personalhtml, req)  	// ввод значений из web
				personalhtml.Title = req.Form["title"][0]
				personalhtml.Name = req.Form["name"][0]
				personalhtml.Kadr = req.Form["kadr"][0]
				personalhtml.Tarif = req.Form["tarif"][0]
				personalhtml.Otdel = req.Form["otdel"][0]
				personalhtml.Numotdel = req.Form["numotdel"][0]
				personalhtml.Email = req.Form["email"][0]
				personalhtml.Phone = req.Form["phone"][0]
				personalhtml.Address = req.Form["address"][0]
				var p frombase
				p.tarif = 10
				p.tarif, err = strconv.Atoi(personalhtml.Tarif)
				if err != nil {
					personalhtml.ErrRange = "1"
					personalhtml.ErrTarif = "1"
					personalhtml.Errors = "1"
				}
				p.numotdel = 100
				p.numotdel, err = strconv.Atoi(personalhtml.Numotdel)
				if err != nil {
					personalhtml.ErrRange = "1"
					personalhtml.ErrNumotd = "1"
					personalhtml.Errors = "1"
				}
				errmail, _, _ := inpMailAddress(personalhtml.Title +
					"<" + personalhtml.Email + ">") // проверка email адреса
				if errmail > 0 {
					personalhtml.Errors = "1"
					personalhtml.ErrEmail = "1"
				}
				_, err := libphonenumber.Parse(personalhtml.Phone, "RU")
				if err != nil {
					personalhtml.Errors = "1"
					personalhtml.ErrPhone = "1"
				}
				if personalhtml.Name == "" || personalhtml.Title == "" || personalhtml.Kadr == "" || personalhtml.Otdel == "" {
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				if personalhtml.Errors == "0" {
					personalhtml.Ready = "1"
					p.title = personalhtml.Title
					p.name = personalhtml.Name
					p.kadr = personalhtml.Kadr
					p.otdel = personalhtml.Otdel
					p.email = personalhtml.Email
					p.phone = personalhtml.Phone
					p.address = personalhtml.Address

					_, err := db.Exec("DELETE FROM personals WHERE title = $1", title)
					if err != nil {
						fmt.Println("Ошибка при удалении старой записи title=", title)
						panic(err)
					}
					sqlStatement := `INSERT INTO personals (title,name,kadr,otdel,numotdel,email,phone,address,tarif) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
					_, err = db.Exec(sqlStatement,
						p.title,
						p.name,
						p.kadr,
						p.otdel,
						p.numotdel,
						p.email,
						p.phone,
						p.address,
						p.tarif,
					)
					if err != nil {
						fmt.Println("Ошибка записи измененной строки в personals", "title=", p.title)
						panic(err)
					}

				}
			}
			err = t.ExecuteTemplate(w, "base", personalhtml)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
				return
			}
		}

	}
}
