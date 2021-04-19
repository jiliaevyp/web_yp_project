package personals

import (
	_ "crypto/dsa"
	"database/sql"
	"fmt"
	"github.com/jiliaevyp/web_yp_project/mond"
	"github.com/jiliaevyp/web_yp_project/servfunc"
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

var partials = []string{
	"./static/base.html",
	"./static/mond_new.html",
	"./static/mond_show.html",
	"./static/mond_edit.html",
	"./static/monds_index.html",
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
	"./static/tabels_index.html",
	"./static/buchtabel_new.html",
	"./static/buchtabel_show.html",
	"./static/buchtabel_edit.html",
	"./static/buchtabel_index.html",
	"./static/css/footer.partial.tmpl.html",
	"./static/css/header.partial.tmpl.html",
	"./static/css/sidebar.partial.tmpl.html",
}
var Department = []string{
	"административный",
	"бухгалтерия",
	"коммерческий",
	"производственный",
	"конструкторский",
	"неизвестный",
}

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
	Ready      string // "1" - ввод корректен
	Errors     string // "1" - ошибка при вводе полей
	ErrPhone   string // "1"- ошибка при вводе телефона
	ErrEmail   string // "1"- ошибка при вводе email
	//ErrTitle  string // "1"- ошибка при вводе title
	ErrTarif  string // "1"- ошибка при вводе тарифа
	ErrNumotd string // "1"- ошибка при вводе номера отдела
	Empty     string // "1" - остались пустые поля
	ErrRange  string // "1" - выход за пределы диапазона
	Jetzyahre string // для sidebar
	Jetzmonat string // для sidebar
}

type frombase struct { // строка  при чтении/записи из/в базы personaldb
	id         int
	forename   string
	title      string
	kadr       string
	numotdel   int
	email      string
	phone      string
	address    string
	tarif      int // почасовая руб
	department string
}

var (
	personals struct {
		Ready     string
		Jetzyahre string
		Jetzmonat string
		//IdRealMond	string
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

// просмотр таблицы из personaldb
func Indexhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		//fmt.Println("idMond=", idMond, jetzMonat, jetzYahre)
		files := append(partials, "./static/personals_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		// обработка  key "id" "del"  ---> удаление записи
		del := req.URL.Query().Get("del")
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		if del == "del" {
			_, err = db.Exec("DELETE FROM personals WHERE id = $1", id)
			if err != nil { // удаление старой записи
				panic(err)
			}
		}
		// выборка таблицы
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
				&p.numotdel,
				&p.tarif,
				&p.email,
				&p.phone,
				&p.address,
				&p.department,
			)
			if err != nil {
				fmt.Println("indexPersonals ошибка распаковки строки ")
				panic(err)
				return
			}
			var personalhtml person
			personalhtml.Id = strconv.Itoa(p.id)
			personalhtml.Forename = p.forename
			personalhtml.Title = p.title
			personalhtml.Kadr = p.kadr
			personalhtml.Tarif = strconv.Itoa(p.tarif) // int ---> string for HTML
			personalhtml.Numotdel = strconv.Itoa(p.numotdel)
			personalhtml.Department = p.department
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
		// для sidebar
		personals.Jetzyahre = mond.Jetzyahre
		personals.Jetzmonat = mond.Jetzmonat
		err = t.ExecuteTemplate(w, "base", personals)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

// просмотр записи из personaldb
func Showhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
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

		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM personals WHERE id=$1", id)

		var p frombase
		err = row.Scan( // чтение строки из таблицы
			&p.id, // int
			&p.title,
			&p.forename,
			&p.kadr,
			&p.numotdel, // int
			&p.tarif,    // int
			&p.email,
			&p.phone,
			&p.address,
			&p.department,
		)
		if err != nil {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		} // подготовка HTML
		personalhtml.Id = strconv.Itoa(p.id)
		personalhtml.Title = p.title
		personalhtml.Forename = p.forename
		personalhtml.Kadr = p.kadr
		personalhtml.Tarif = strconv.Itoa(p.tarif)       // int ---> string
		personalhtml.Numotdel = strconv.Itoa(p.numotdel) // int ---> string
		personalhtml.Email = p.email
		personalhtml.Phone = p.phone
		personalhtml.Address = p.address
		personalhtml.Department = p.department
		// для sidebar
		personalhtml.Jetzyahre = mond.Jetzyahre
		personalhtml.Jetzmonat = mond.Jetzmonat
		err = t.ExecuteTemplate(w, "base", personalhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error personalShowhandler", http.StatusInternalServerError)
			return
		}
	}
}

// новая запись формы personal в базу personaldb
func Newhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
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
		//personalhtml.Errors = "1"
		//for personalhtml.Errors == "1" {
		if req.Method == "POST" {
			req.ParseForm()
			makeReadyHtml(&personalhtml) // подготовка значений для web
			//readFromHtml(&personalhtml, req)  	// ввод значений из web
			personalhtml.Errors = "0"
			personalhtml.Ready = "0"

			personalhtml.Title = req.Form["title"][0]
			// проверка на пустые поля данных
			if personalhtml.Title == "" || personalhtml.Title == "???" {
				personalhtml.Title = "???"
				personalhtml.Empty = "1"
				personalhtml.Errors = "1"
			}
			personalhtml.Forename = req.Form["forename"][0]
			// проверка на пустые поля данных
			if personalhtml.Forename == "" || personalhtml.Forename == "???" {
				personalhtml.Forename = "???"
				personalhtml.Empty = "1"
				personalhtml.Errors = "1"
			}
			personalhtml.Kadr = req.Form["kadr"][0]
			if personalhtml.Kadr == "" || personalhtml.Kadr == "???" {
				personalhtml.Kadr = "???"
				personalhtml.Empty = "1"
				personalhtml.Errors = "1"
			}
			personalhtml.Address = req.Form["address"][0]
			if personalhtml.Address == "" || personalhtml.Address == "???" {
				personalhtml.Address = "???"
				personalhtml.Empty = "1"
				personalhtml.Errors = "1"
			}
			// ввод и проверка на числа
			personalhtml.Tarif = req.Form["tarif"][0]
			if servfunc.Checknum(personalhtml.Tarif, 10, 1000) != 0 {
				personalhtml.Tarif = "???"
				personalhtml.ErrRange = "1"
				personalhtml.ErrTarif = "1"
				personalhtml.Errors = "1"
			}
			personalhtml.Numotdel = req.Form["numotdel"][0]
			if servfunc.Checknum(personalhtml.Numotdel, 0, 4) != 0 {
				personalhtml.Numotdel = "???"
				personalhtml.ErrRange = "1"
				personalhtml.ErrNumotd = "1"
				personalhtml.Errors = "1"
			} else {
				department, _ := strconv.Atoi(personalhtml.Numotdel)
				personalhtml.Department = Department[department]
			}
			personalhtml.Email = req.Form["email"][0]
			var errmail int
			if personalhtml.Email != "" && personalhtml.Email != "???" {
				errmail, personalhtml.Email, _ = servfunc.InpMailAddress(personalhtml.Title + "<" + personalhtml.Email + ">") // проверка email адреса
				if errmail > 0 {
					personalhtml.Email = "???"
					personalhtml.ErrEmail = "1"
					personalhtml.Errors = "1"
				}
			} else {
				personalhtml.Email = "???"
				personalhtml.ErrEmail = "1"
				personalhtml.Errors = "1"
			}
			row := db.QueryRow("SELECT * FROM personals WHERE email=$6", personalhtml.Email)
			var p frombase
			err = row.Scan( // чтение строки из таблицы
				&p.phone,
			)
			if err == nil {
				fmt.Println("email уже был использован!")
				personalhtml.Email = "???"
				personalhtml.ErrEmail = "2"
				personalhtml.Errors = "1"
			}
			personalhtml.Phone = req.Form["phone"][0]
			_, err1 := libphonenumber.Parse(personalhtml.Phone, "RU")
			if err1 != nil {
				personalhtml.Phone = "???"
				personalhtml.ErrPhone = "1"
				personalhtml.Errors = "1"
			}

			if personalhtml.Errors == "0" {
				personalhtml.Ready = "1"
				var p frombase
				p.title = personalhtml.Title //personalhtml.Title
				p.forename = personalhtml.Forename
				p.kadr = personalhtml.Kadr
				p.tarif, _ = strconv.Atoi(personalhtml.Tarif)       // перевод в int для базы
				p.numotdel, _ = strconv.Atoi(personalhtml.Numotdel) // перевод в int для базы
				p.email = personalhtml.Email
				p.phone = personalhtml.Phone
				p.address = personalhtml.Address
				p.department = Department[p.numotdel]
				email := p.email
				_, err := db.Exec("DELETE FROM personals WHERE email = $1", email) // удаление  записи по считанному email
				if err != nil {
					fmt.Println("Ошибка при удалении старой записи email=", email)
					panic(err)
				}
				sqlStatement := `INSERT INTO personals (title, forename, kadr,numotdel,tarif,email,phone,address,department) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
				_, err = db.Exec(sqlStatement,
					p.title,
					p.forename,
					p.kadr,
					p.numotdel,
					p.tarif,
					p.email,
					p.phone,
					p.address,
					p.department,
				)
				if err != nil {
					fmt.Println("Ошибка записи измененной строки в personals", "title=", p.title)
					//panic(err)
					personalhtml.Errors = "1"
					personalhtml.ErrEmail = "2"
					personalhtml.Ready = "0"
				}
			}
		}
		// для sidebar
		personalhtml.Jetzyahre = mond.Jetzyahre
		personalhtml.Jetzmonat = mond.Jetzmonat
		err = t.ExecuteTemplate(w, "base", personalhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

// редактирование формы personal и замена в базе personaldb
func Edithandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/personal_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		var p frombase
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM personals WHERE id=$1", id)
		err = row.Scan( // пересылка  данных строки базы personals в p
			&p.id,
			&p.title,
			&p.forename,
			&p.kadr,
			&p.numotdel, // int
			&p.tarif,    // int
			&p.email,
			&p.phone,
			&p.address,
			&p.department,
		)
		if err != nil {
			fmt.Println("edit --> ошибка распаковки строки при чтении записи id=", id)
			panic(err)
		} else {
			var personalhtml person
			makeReadyHtml(&personalhtml) // подготовка флагов для HTML = 0
			personalhtml.Empty = "1"     // якобы - есть пустые поля для отображения
			personalhtml.Id = strconv.Itoa(p.id)
			personalhtml.Title = p.title
			personalhtml.Forename = p.forename
			personalhtml.Kadr = p.kadr
			personalhtml.Tarif = strconv.Itoa(p.tarif)       // int ---> string
			personalhtml.Numotdel = strconv.Itoa(p.numotdel) // int ---> string
			personalhtml.Email = p.email
			personalhtml.Phone = p.phone
			personalhtml.Address = p.address
			personalhtml.Department = p.department
			if req.Method == "POST" {
				req.ParseForm()
				makeReadyHtml(&personalhtml) // подготовка значений для web
				// ввод новых данных
				personalhtml.Errors = "0"
				personalhtml.Ready = "0"

				personalhtml.Title = req.Form["title"][0]
				// проверка на пустые поля данных
				if personalhtml.Title == "" || personalhtml.Title == "???" {
					personalhtml.Title = "???"
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Forename = req.Form["forename"][0]
				// проверка на пустые поля данных
				if personalhtml.Forename == "" || personalhtml.Forename == "???" {
					personalhtml.Forename = "???"
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Kadr = req.Form["kadr"][0]
				if personalhtml.Kadr == "" || personalhtml.Kadr == "???" {
					personalhtml.Kadr = "???"
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Address = req.Form["address"][0]
				if personalhtml.Address == "" || personalhtml.Address == "???" {
					personalhtml.Address = "???"
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				// ввод и проверка на числа
				personalhtml.Tarif = req.Form["tarif"][0]
				if servfunc.Checknum(personalhtml.Tarif, 10, 1000) != 0 {
					personalhtml.Tarif = "???"
					personalhtml.ErrRange = "1"
					personalhtml.ErrTarif = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Numotdel = req.Form["numotdel"][0]
				if servfunc.Checknum(personalhtml.Numotdel, 0, 20) != 0 {
					personalhtml.Numotdel = "???"
					personalhtml.ErrRange = "1"
					personalhtml.ErrNumotd = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Department = req.Form["department"][0]
				personalhtml.Email = req.Form["email"][0]
				var errmail int
				if personalhtml.Email != "" && personalhtml.Email != "???" {
					errmail, personalhtml.Email, _ = servfunc.InpMailAddress(personalhtml.Title + "<" + personalhtml.Email + ">") // проверка email адреса
					if errmail > 0 {
						personalhtml.Email = "???"
						personalhtml.ErrEmail = "1"
						personalhtml.Errors = "1"
					}
				} else {
					personalhtml.Email = "???"
					personalhtml.ErrEmail = "1"
					personalhtml.Errors = "1"
				}

				personalhtml.Phone = req.Form["phone"][0]
				_, err1 := libphonenumber.Parse(personalhtml.Phone, "RU")
				if err1 != nil {
					personalhtml.Phone = "???"
					personalhtml.ErrPhone = "1"
					personalhtml.Errors = "1"
				}

				if personalhtml.Errors == "0" {
					personalhtml.Ready = "1"
					var p frombase
					p.id, _ = strconv.Atoi(personalhtml.Id)
					p.title = personalhtml.Title
					p.forename = personalhtml.Forename
					p.kadr = personalhtml.Kadr
					p.tarif, _ = strconv.Atoi(personalhtml.Tarif)       // перевод в int для базы
					p.numotdel, _ = strconv.Atoi(personalhtml.Numotdel) // перевод в int для базы
					p.email = personalhtml.Email
					p.phone = personalhtml.Phone
					p.address = personalhtml.Address
					p.department = personalhtml.Department
					id := p.id
					_, err = db.Exec("DELETE FROM personals WHERE id = $1", id)
					if err != nil { // удаление старой записи
						panic(err)
					}
					sqlStatement := `INSERT INTO personals (id,title, forename, kadr,numotdel,tarif,email,phone,address,department) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
					_, err = db.Exec(sqlStatement,
						p.id,
						p.title,
						p.forename,
						p.kadr,
						p.numotdel,
						p.tarif,
						p.email,
						p.phone,
						p.address,
						p.department,
					)
					//fmt.Println(" запись edit Id=", p.id, p)
					if err != nil {
						fmt.Println("Ошибка записи измененной строки в personals", "title=", p.title)
						panic(err)
					}
				}
			}
			// для sidebar
			personalhtml.Jetzyahre = mond.Jetzyahre
			personalhtml.Jetzmonat = mond.Jetzmonat
			err = t.ExecuteTemplate(w, "base", personalhtml)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "personal_edit Internal Server Execute Error ", http.StatusInternalServerError)
				return
			}
		}

	}
}
