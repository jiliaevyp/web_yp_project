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
	"net/mail"
	//"net"
	"net/http"
	"strconv"
)

var partials = []string{
	"./static/base.html",
	"./static/mond_new.html",
	"./static/mond_show.html",
	"./static/mond_edit.html",
	"./static/mond_index.html",
	"./static/personal_new.html",
	"./static/personal_show.html",
	"./static/personal_edit.html",
	"./static/personals_index.html",
	"./static/worktime_new.html",
	"./static/worktime_show.html",
	"./static/worktime_edit.html",
	"./static/worktime_index.html",
	"./static/tabel_new.html",
	"./static/tabel_show.html",
	"./static/tabel_edit.html",
	"./static/tabel_index.html",
	"./static/buchtabel_new.html",
	"./static/buchtabel_show.html",
	"./static/buchtabel_edit.html",
	"./static/buchtabel_index.html",
	"./static/css/footer.partial.tmpl.html",
	"./static/css/header.partial.tmpl.html",
	"./static/css/sidebar.partial.tmpl.html",
}

type person struct { // данные по сотруднику при вводе и отображении в personal.HTML
	Forename  string
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
	forename string
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

// проверка корректности емайл адреса nameAddress --> "имя <email@mail.com>
func inpMailAddress(nameAddress string) (err int, email string, title string) {
	e, err1 := mail.ParseAddress(nameAddress)
	if err1 != nil {
		return 1, e.Address, e.Name //"?", "?"
	}
	return 0, e.Address, e.Name
}

// валидация  числовых вводов и диапазонов
func checknum(checknum string, min int, max int) int {
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
	p.Forename = req.Form["forename"][0]
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
func checkNumer(personalhtml *person) int {
	var err int
	errout := 0
	err = checknum(personalhtml.Tarif, 10, 1000)
	if err != 0 {
		personalhtml.ErrRange = "1"
		personalhtml.ErrTarif = "1"
		errout = 1
	}
	err = checknum(personalhtml.Numotdel, 0, 20)
	if err != 0 {
		personalhtml.ErrRange = "1"
		personalhtml.ErrNumotd = "1"
		errout = 1
	}
	err, personalhtml.Email, _ = inpMailAddress(personalhtml.Title + "<" + personalhtml.Email + ">") // проверка email адреса
	if err > 0 {
		personalhtml.ErrEmail = "1"
		errout = 1
	}
	_, err1 := libphonenumber.Parse(personalhtml.Phone, "RU")
	if err1 != nil {
		personalhtml.ErrPhone = "1"
		errout = 1
	}
	if personalhtml.Forename == "" || personalhtml.Title == "" || personalhtml.Kadr == "" || personalhtml.Otdel == "" || personalhtml.Address == "" {
		personalhtml.Empty = "1"
		errout = 1
	}
	return errout
}

// просмотр таблицы из personaldb
func Personalshandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
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
				&p.forename,
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
			personalhtml.Forename = p.forename
			personalhtml.Title = p.title
			personalhtml.Kadr = p.kadr
			personalhtml.Tarif = strconv.Itoa(p.tarif) // int ---> string for HTML
			personalhtml.Numotdel = strconv.Itoa(p.numotdel)
			personalhtml.Otdel = p.otdel
			personalhtml.Email = p.email
			personalhtml.Phone = p.phone
			personalhtml.Address = p.address
			personalhtml.Ready = "1"
			personalhtml.Errors = "0"
			personalhtml.Empty = "0"
			// добавление строки в таблицу Personalstab для personals_index.html
			personals.Persontable = append(personals.Persontable, personalhtml)
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
func PersonalShowhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/personal_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		var personalhtml person
		personalhtml.Ready = "1"     // 1 - ввод успешный
		personalhtml.Errors = "0"    // 1 - ошибки при вводе
		personalhtml.Empty = "0"     // 1 - есть пустые поля
		personalhtml.ErrRange = "0"  // 1 - выход за пределы диапазона
		personalhtml.ErrPhone = "0"  // 1 - ошибка в тлф номере
		personalhtml.ErrEmail = "0"  // 1 - ошибка в email
		personalhtml.ErrTarif = "0"  // 1 - ошибка в тарифе
		personalhtml.ErrNumotd = "0" // 1 - ошибка в номере отдела
		title := req.URL.Query().Get("title")
		row := db.QueryRow("SELECT * FROM personals WHERE title=$1", title)

		var p frombase
		err = row.Scan( // чтение строки из таблицы
			&p.id, // int
			&p.title,
			&p.forename,
			&p.kadr,
			&p.tarif,    // int
			&p.numotdel, // int
			&p.otdel,
			&p.email,
			&p.phone,
			&p.address,
		)
		if err != nil {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		} // подготовка HTML
		personalhtml.Title = p.title
		personalhtml.Forename = p.forename
		personalhtml.Kadr = p.kadr
		personalhtml.Tarif = strconv.Itoa(p.tarif)       // int ---> string
		personalhtml.Numotdel = strconv.Itoa(p.numotdel) // int ---> string
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
func PersonalNewhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
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
			personalhtml.Forename = req.Form["forename"][0]
			personalhtml.Kadr = req.Form["kadr"][0]
			personalhtml.Tarif = req.Form["tarif"][0]
			personalhtml.Otdel = req.Form["otdel"][0]
			personalhtml.Numotdel = req.Form["numotdel"][0]
			personalhtml.Email = req.Form["email"][0]
			personalhtml.Phone = req.Form["phone"][0]
			personalhtml.Address = req.Form["address"][0]

			// проверка введенных данных
			if checkNumer(&personalhtml) == 0 {
				personalhtml.Errors = "0" // ввод корректный
				personalhtml.Ready = "1"
				//добавление записи в базу
				title := personalhtml.Title
				// удаление старой записи
				row := db.QueryRow("SELECT * FROM personals WHERE title=$1", title)
				if row != nil { // если запись есть удаляем
					_, err1 := db.Exec("DELETE FROM personals WHERE title = $1", title)
					if err1 != nil {
						fmt.Println("Ошибка при удалении старой записи в personals title = ", title)
						panic(err)
					}
				}
				var p frombase
				p.title = personalhtml.Title
				p.forename = personalhtml.Forename
				p.kadr = personalhtml.Kadr
				p.tarif, _ = strconv.Atoi(personalhtml.Tarif) // перевод в int для базы
				p.otdel = personalhtml.Otdel
				p.numotdel, _ = strconv.Atoi(personalhtml.Numotdel) // перевод в int для базы
				p.email = personalhtml.Email
				p.phone = personalhtml.Phone
				p.address = personalhtml.Address

				sqlStatement := `INSERT INTO personals (title, forename, kadr,tarif,numotdel,otdel,email,phone,address) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
				_, err2 := db.Exec(sqlStatement,
					p.title,
					p.forename,
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
					panic(err2)
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
func PersonalEdithandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/personal_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}

		var p frombase
		title := req.URL.Query().Get("title") // чтение строки из таблицы
		row := db.QueryRow("SELECT * FROM personals WHERE title=$1", title)
		err = row.Scan( // пересылка  данных строки базы personals в p
			&p.id,
			&p.title,
			&p.forename,
			&p.kadr,
			&p.tarif,    // int
			&p.numotdel, // int
			&p.otdel,
			&p.email,
			&p.phone,
			&p.address,
		)
		if err != nil {
			fmt.Println("edit --> ошибка распаковки строки при чтении записи title=", title)
			panic(err)
		} else {
			var personalhtml person
			makeReadyHtml(&personalhtml) // подготовка флагов для HTML = 0
			personalhtml.Empty = "1"     // якобы - есть пустые поля для отображения
			personalhtml.Title = p.title
			personalhtml.Forename = p.forename
			personalhtml.Kadr = p.kadr
			personalhtml.Tarif = strconv.Itoa(p.tarif)       // int ---> string
			personalhtml.Numotdel = strconv.Itoa(p.numotdel) // int ---> string
			personalhtml.Otdel = p.otdel
			personalhtml.Email = p.email
			personalhtml.Phone = p.phone
			personalhtml.Address = p.address

			if req.Method == "POST" {
				req.ParseForm()
				makeReadyHtml(&personalhtml) // подготовка значений для web
				//readFromHtml(&personalhtml, req)  	// ввод значений из web
				personalhtml.Title = req.Form["title"][0]
				personalhtml.Forename = req.Form["forename"][0]
				personalhtml.Kadr = req.Form["kadr"][0]
				personalhtml.Tarif = req.Form["tarif"][0]
				personalhtml.Numotdel = req.Form["numotdel"][0]
				personalhtml.Otdel = req.Form["otdel"][0]
				personalhtml.Email = req.Form["email"][0]
				personalhtml.Phone = req.Form["phone"][0]
				personalhtml.Address = req.Form["address"][0]

				// проверка введенных данных
				if checkNumer(&personalhtml) == 0 {
					personalhtml.Errors = "0" // ввод корректный
					personalhtml.Ready = "1"
					var p frombase
					p.title = personalhtml.Title
					p.forename = personalhtml.Forename
					p.kadr = personalhtml.Kadr
					p.tarif, _ = strconv.Atoi(personalhtml.Tarif)       // перевод в int для базы
					p.numotdel, _ = strconv.Atoi(personalhtml.Numotdel) // перевод в int для базы
					p.otdel = personalhtml.Otdel
					p.otdel = personalhtml.Otdel
					p.email = personalhtml.Email
					p.phone = personalhtml.Phone
					p.address = personalhtml.Address

					_, err := db.Exec("DELETE FROM personals WHERE title = $1", title) // удаление  записи по считанному title
					if err != nil {
						fmt.Println("Ошибка при удалении старой записи title=", title)
						panic(err)
					}
					sqlStatement := `INSERT INTO personals (title,forename,kadr,tarif,numotdel,otdel,email,phone,address) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
					_, err = db.Exec(sqlStatement,
						p.title,
						p.forename,
						p.kadr,
						p.tarif,
						p.numotdel,
						p.otdel,
						p.email,
						p.phone,
						p.address,
					)
					if err != nil {
						fmt.Println("Ошибка записи измененной строки в personals", "title=", p.title)
						panic(err)
					} else {
						fmt.Println("успешно записали edit title=", p.title)
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
