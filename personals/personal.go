package personals

import (
	_ "crypto/dsa"
	"database/sql"
	"fmt"
	"github.com/jiliaevyp/web_yp_project/mond"
	"github.com/jiliaevyp/web_yp_project/servfunc"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

type Tabel struct {
	PersonId   int
	MondId     int
	Forename   string
	Title      string
	Kadr       string
	Numotdel   int
	Department string
	Email      string
}

//var Newtabel []Tabel		// таблица нового табеля

type person struct { // данные по сотруднику при вводе и отображении в personal.HTML
	Id         string
	Forename   string
	Title      string
	Kadr       string
	Numotdel   string
	Department string
	Email      string
	Tarif      string // почасовая руб
	Real       string // "1" - включен в табель
	Ready      string // "1" - ввод корректен
	Errors     string // "1" - ошибка при вводе полей
	ErrEmail   string // "1"- ошибка при вводе Email
	ErrTitle   string // "1"- ошибка при вводе Title
	ErrTarif   string // "1"- ошибка при вводе тарифа
	ErrNumotd  string // "1"- ошибка при вводе номера отдела
	Empty      string // "1" - остались пустые поля
	ErrRange   string // "1" - выход за пределы диапазона
	Jetzyahre  string // для sidebar
	Jetzmonat  string // для sidebar
}

type Personfrombase struct { // строка  при чтении/записи из/в базы personaldb
	Id         int
	Forename   string
	Title      string
	Kadr       string
	Numotdel   int
	Email      string
	Tarif      int // почасовая руб
	Department string
	Real       int // 1 - включен в табель
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
	p.ErrEmail = "0"  // 1 - ошибка в Email
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
			var p Personfrombase
			err = rows.Scan( // пересылка  данных строки базы personals в "p"
				&p.Id,
				&p.Title,
				&p.Forename,
				&p.Kadr,
				&p.Numotdel,
				&p.Tarif,
				&p.Email,
				&p.Department,
				&p.Real,
			)
			if err != nil {
				fmt.Println("indexPersonals ошибка распаковки строки ")
				panic(err)
				return
			}
			var personalhtml person
			personalhtml.Id = strconv.Itoa(p.Id)
			personalhtml.Forename = p.Forename
			personalhtml.Title = p.Title
			personalhtml.Kadr = p.Kadr
			personalhtml.Tarif = strconv.Itoa(p.Tarif) // int ---> string for HTML
			personalhtml.Numotdel = strconv.Itoa(p.Numotdel)
			personalhtml.Department = p.Department
			personalhtml.Email = p.Email
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
		personalhtml.ErrEmail = "0"  // 1 - ошибка в Email
		personalhtml.ErrTarif = "0"  // 1 - ошибка в тарифе
		personalhtml.ErrNumotd = "0" // 1 - ошибка в номере отдела

		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM personals WHERE id=$1", id)

		var p Personfrombase
		err = row.Scan( // чтение строки из таблицы
			&p.Id, // int
			&p.Title,
			&p.Forename,
			&p.Kadr,
			&p.Numotdel, // int
			&p.Tarif,    // int
			&p.Email,
			&p.Department,
			&p.Real,
		)
		if err != nil {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		} // подготовка HTML
		personalhtml.Id = strconv.Itoa(p.Id)
		personalhtml.Title = p.Title
		personalhtml.Forename = p.Forename
		personalhtml.Kadr = p.Kadr
		personalhtml.Tarif = strconv.Itoa(p.Tarif)       // int ---> string
		personalhtml.Numotdel = strconv.Itoa(p.Numotdel) // int ---> string
		personalhtml.Email = p.Email
		personalhtml.Department = p.Department
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

			personalhtml.Title = req.Form["Title"][0]
			// проверка на пустые поля данных
			if personalhtml.Title == "" || personalhtml.Title == "???" {
				personalhtml.Title = "???"
				personalhtml.Empty = "1"
				personalhtml.Errors = "1"
			}
			personalhtml.Forename = req.Form["Forename"][0]
			// проверка на пустые поля данных
			if personalhtml.Forename == "" || personalhtml.Forename == "???" {
				personalhtml.Forename = "???"
				personalhtml.Empty = "1"
				personalhtml.Errors = "1"
			}
			personalhtml.Kadr = req.Form["Kadr"][0]
			if personalhtml.Kadr == "" || personalhtml.Kadr == "???" {
				personalhtml.Kadr = "???"
				personalhtml.Empty = "1"
				personalhtml.Errors = "1"
			}
			// ввод и проверка на числа
			personalhtml.Tarif = req.Form["Tarif"][0]
			if servfunc.Checknum(personalhtml.Tarif, 10, 1000) != 0 {
				personalhtml.Tarif = "???"
				personalhtml.ErrRange = "1"
				personalhtml.ErrTarif = "1"
				personalhtml.Errors = "1"
			}
			personalhtml.Numotdel = req.Form["Numotdel"][0]
			if servfunc.Checknum(personalhtml.Numotdel, 0, 4) != 0 {
				personalhtml.Numotdel = "???"
				personalhtml.ErrRange = "1"
				personalhtml.ErrNumotd = "1"
				personalhtml.Errors = "1"
			} else {
				department, _ := strconv.Atoi(personalhtml.Numotdel)
				personalhtml.Department = Department[department]
			}
			personalhtml.Email = req.Form["Email"][0]
			var errmail int
			if personalhtml.Email != "" && personalhtml.Email != "???" {
				errmail, personalhtml.Email, _ = servfunc.InpMailAddress(personalhtml.Title + "<" + personalhtml.Email + ">") // проверка Email адреса
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
			row := db.QueryRow("SELECT * FROM personals WHERE Email=$6", personalhtml.Email)
			var p Personfrombase
			err = row.Scan( // чтение строки из таблицы
				&p.Id,
				&p.Title,
				&p.Forename,
				&p.Kadr,
				&p.Numotdel, // int
				&p.Tarif,    // int
				&p.Email,
				&p.Department,
				&p.Real,
			)
			if err == nil {
				fmt.Println("Email уже был использован!")
				personalhtml.Email = "???"
				personalhtml.ErrEmail = "2"
				personalhtml.Errors = "1"
			}
			personalhtml.Real = req.Form["Real"][0]

			if personalhtml.Errors == "0" {
				personalhtml.Ready = "1"
				// подготовка записи нового персонала
				var p Personfrombase
				p.Title = personalhtml.Title //personalhtml.Title
				p.Forename = personalhtml.Forename
				p.Kadr = personalhtml.Kadr
				p.Tarif, _ = strconv.Atoi(personalhtml.Tarif)       // перевод в int для базы
				p.Numotdel, _ = strconv.Atoi(personalhtml.Numotdel) // перевод в int для базы
				p.Email = personalhtml.Email
				p.Department = Department[p.Numotdel]
				p.Real, _ = strconv.Atoi(personalhtml.Real)
				email := p.Email
				_, err := db.Exec("DELETE FROM personals WHERE Email = $1", email) // удаление  записи по считанному Email
				if err != nil {
					fmt.Println("Ошибка при удалении старой записи Email=", email)
					panic(err)
				}
				// запись нового персонала
				sqlStatement1 := `INSERT INTO personals (Title,Forename,Kadr,Numotdel,Tarif,Email,Department,Real)
 									VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
				_, err = db.Exec(sqlStatement1,
					p.Title,
					p.Forename,
					p.Kadr,
					p.Numotdel,
					p.Tarif,
					p.Email,
					p.Department,
					p.Real,
				)
				if err != nil {
					fmt.Println("Ошибка записи измененной строки в personals", "Title=", p.Title)
					//panic(err)
					personalhtml.Errors = "1"
					personalhtml.ErrEmail = "2"
					personalhtml.Ready = "0"
				}
				// запись нового табеля
				t := Tabel{}

				if p.Real == 0 {
					_, err1 := db.Exec("DELETE FROM tabels WHERE Email = $1", email) // удаляем запись табеля по Email
					if err1 != nil {
						fmt.Println("Ошибка при удалении записи Email=", email)
						panic(err1)
					} else {
						//t.PersonId = p.Id			// новая запись табеля
						t.MondId = mond.IdRealMond
						t.Forename = p.Forename
						t.Title = p.Title
						t.Kadr = p.Kadr
						t.Numotdel = p.Numotdel
						t.Department = p.Department
						t.Email = p.Email
						_, err1 := db.Exec("DELETE FROM tabels WHERE Email = $1", email) // удаление  записи табеля по считанному Email
						if err1 != nil {
							fmt.Println("Ошибка при удалении старой записи Email=", email)
							panic(err1)
						}
						// новая запись табеля
						sqlStatement := `INSERT INTO tabels (Title,Forename,Kadr,Numotdel,Department,MondId,Email) 
									VALUES ($1,$2,$3,$4,$5,$6,$7)`
						_, err = db.Exec(sqlStatement,
							t.Title,
							t.Forename,
							t.Kadr,
							t.Numotdel,
							t.Department,
							t.MondId,
							p.Email,
						)
					}

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
		var p Personfrombase
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM personals WHERE id=$1", id)
		err = row.Scan( // пересылка  данных строки базы personals в p
			&p.Id,
			&p.Title,
			&p.Forename,
			&p.Kadr,
			&p.Numotdel, // int
			&p.Tarif,    // int
			&p.Email,
			&p.Department,
			&p.Real,
		)
		if err != nil {
			fmt.Println("edit --> ошибка распаковки строки при чтении записи id=", id)
			panic(err)
		} else {
			var personalhtml person
			makeReadyHtml(&personalhtml) // подготовка флагов для HTML = 0
			personalhtml.Empty = "1"     // якобы - есть пустые поля для отображения
			personalhtml.Id = strconv.Itoa(p.Id)
			personalhtml.Title = p.Title
			personalhtml.Forename = p.Forename
			personalhtml.Kadr = p.Kadr
			personalhtml.Tarif = strconv.Itoa(p.Tarif)       // int ---> string
			personalhtml.Numotdel = strconv.Itoa(p.Numotdel) // int ---> string
			personalhtml.Department = p.Department
			personalhtml.Real = strconv.Itoa(p.Real)
			if req.Method == "POST" {
				req.ParseForm()
				makeReadyHtml(&personalhtml) // подготовка значений для web
				// ввод новых данных
				personalhtml.Errors = "0"
				personalhtml.Ready = "0"

				personalhtml.Title = req.Form["Title"][0]
				// проверка на пустые поля данных
				if personalhtml.Title == "" || personalhtml.Title == "???" {
					personalhtml.Title = "???"
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Forename = req.Form["Forename"][0]
				// проверка на пустые поля данных
				if personalhtml.Forename == "" || personalhtml.Forename == "???" {
					personalhtml.Forename = "???"
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Kadr = req.Form["Kadr"][0]
				if personalhtml.Kadr == "" || personalhtml.Kadr == "???" {
					personalhtml.Kadr = "???"
					personalhtml.Empty = "1"
					personalhtml.Errors = "1"
				}
				// ввод и проверка на числа
				personalhtml.Tarif = req.Form["Tarif"][0]
				if servfunc.Checknum(personalhtml.Tarif, 10, 1000) != 0 {
					personalhtml.Tarif = "???"
					personalhtml.ErrRange = "1"
					personalhtml.ErrTarif = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Numotdel = req.Form["Numotdel"][0]
				if servfunc.Checknum(personalhtml.Numotdel, 0, 20) != 0 {
					personalhtml.Numotdel = "???"
					personalhtml.ErrRange = "1"
					personalhtml.ErrNumotd = "1"
					personalhtml.Errors = "1"
				}
				personalhtml.Department = req.Form["Department"][0]
				personalhtml.Real = req.Form["Real"][0]
				personalhtml.Email = req.Form["Email"][0]
				var errmail int
				if personalhtml.Email != "" && personalhtml.Email != "???" {
					errmail, personalhtml.Email, _ = servfunc.InpMailAddress(personalhtml.Title + "<" + personalhtml.Email + ">") // проверка Email адреса
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

				if personalhtml.Errors == "0" {
					personalhtml.Ready = "1"
					var p Personfrombase
					p.Id, _ = strconv.Atoi(personalhtml.Id)
					p.Title = personalhtml.Title
					p.Forename = personalhtml.Forename
					p.Kadr = personalhtml.Kadr
					p.Tarif, _ = strconv.Atoi(personalhtml.Tarif)       // перевод в int для базы
					p.Numotdel, _ = strconv.Atoi(personalhtml.Numotdel) // перевод в int для базы
					p.Email = personalhtml.Email
					p.Department = personalhtml.Department
					p.Real, _ = strconv.Atoi(personalhtml.Real)
					id := p.Id
					_, err = db.Exec("DELETE FROM personals WHERE id = $1", id)
					if err != nil { // удаление старой записи
						panic(err)
					}
					sqlStatement := `INSERT INTO personals (id,Title, Forename, Kadr,Numotdel,Tarif,Email,Department,Real) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
					_, err = db.Exec(sqlStatement,
						p.Id,
						p.Title,
						p.Forename,
						p.Kadr,
						p.Numotdel,
						p.Tarif,
						p.Email,
						p.Department,
						&p.Real,
					)
					//fmt.Println(" запись edit Id=", p.id, p)
					if err != nil {
						fmt.Println("Ошибка записи измененной строки в personals", "Title=", p.Title)
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
