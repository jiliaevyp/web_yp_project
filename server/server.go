package server

import (
	"database/sql"
	"fmt"
	"github.com/jiliaevyp/web_yp_project/mond"
	"github.com/jiliaevyp/web_yp_project/personals"
	"github.com/jiliaevyp/web_yp_project/servfunc"
	"github.com/jiliaevyp/web_yp_project/tabel"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"strconv"
)

const (
	defaultUser  = "yp"
	defaultEmail = "yp@yp.com"
	defaultPassw = "123"
)

type frombase struct { // для чтения из таблицы monds
	id       int
	yahre    int
	nummonat int
	tag      int
	hour     int
	kf       int
	blmond   int
	bltime   int
	bltabel  int
	blbuch   int
	blpers   int
	monat    string
}

var Erserv int // 1 - ошибка при запуске сервера

var yahre string // год для передачи в обработчик
var monat string // месяц для передачи в обработчик
var idMond int   // id текущей активной записи monds

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
var admin struct { // администратор
	User      string
	Email     string
	Passw     string
	ErrEmail  string // ошибка ввода почты
	Passpass  string
	Ready     string // 1 - идентификация прошла
	Errors    string // "1" - ошибка при вводе полей
	Empty     string // "1" - остались пустые поля
	IdMonat   int
	JetzYahre string
	JetzMonat string
}

// проверка корректности емайл адреса nameAddress --> "имя <email@mail.com>
func inpMailAddress(nameAddress string) (err int, email string, title string) {
	e, err1 := mail.ParseAddress(nameAddress)
	if err1 != nil {
		return 1, e.Address, e.Name //"?", "?"
	}
	return 0, e.Address, e.Name
}

func Server(addrWeb string, db *sql.DB) {
	http.Handle("/index", http.HandlerFunc(indexHandler(db)))

	http.Handle("/monds_index", http.HandlerFunc(mond.Indexhandler(db)))
	http.Handle("/mond_new", http.HandlerFunc(mond.Newhandler(db)))
	http.Handle("/mond_show", http.HandlerFunc(mond.Showhandler(db)))
	http.Handle("/mond_edit", http.HandlerFunc(mond.Edithandler(db)))
	_idMond := 20      //idMond
	_yahre := "2021"   //yahre
	_monat := "январь" //monat
	http.Handle("/personals_index", http.HandlerFunc(personals.Indexhandler(db, _yahre, _monat, _idMond)))
	http.Handle("/personal_new", http.HandlerFunc(personals.Newhandler(db, yahre, monat, idMond)))
	http.Handle("/personal_show", http.HandlerFunc(personals.Showhandler(db, yahre, monat, idMond)))
	http.Handle("/personal_edit", http.HandlerFunc(personals.Edithandler(db, yahre, monat, idMond)))

	//http.Handle("/worktime_index", http.HandlerFunc(IndexHandler(db)))
	//http.Handle("/worktime_new", http.HandlerFunc(NewHandler(db)))
	//http.Handle("/worktime_show", http.HandlerFunc(ShowHandler(db)))
	//http.Handle("/worktime_edit", http.HandlerFunc(EditHandler(db)))
	//
	http.Handle("/tabels_index", http.HandlerFunc(tabel.IndexHandler(db, yahre, monat, idMond)))
	http.Handle("/tabel_new", http.HandlerFunc(tabel.NewHandler(db, yahre, monat, idMond)))
	http.Handle("/tabel_show", http.HandlerFunc(tabel.ShowHandler(db, yahre, monat, idMond)))
	http.Handle("/tabel_edit", http.HandlerFunc(tabel.EditHandler(db, yahre, monat, idMond)))
	//
	//http.Handle("/buchtabel_index", http.HandlerFunc(IndexHandler(db)))
	//http.Handle("/buchtabel_new", http.HandlerFunc(NewHandler(db)))
	//http.Handle("/buchtabel_show", http.HandlerFunc(ShowHandler(db)))
	//http.Handle("/buchtabel_edit", http.HandlerFunc(EditHandler(db)))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Топай на web страницу--->" + addrWeb + "!") // отладочная печать
	err := http.ListenAndServe(addrWeb, nil)

	if err != nil {
		Erserv = 1
	} else {
		Erserv = 0
	} // запуск сервера
	return
	//errserv := 0                                              // адрес в --> addrWeb
	//return  // форма в --> msgHandler
}

// первая страница проверка доступа
//func indexHandler(w http.ResponseWriter, req *http.Request)

func indexHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		admin.JetzMonat = "???"
		admin.JetzYahre = "???"
		files := append(partials, "./static/index.html")
		t, err := template.ParseFiles(files...) // Parse template file.

		_id, ok := req.URL.Query()["id"]
		fmt.Println("_id =", _id)
		if ok && len(_id[0]) > 0 {
			id := _id[0]
			fmt.Println("id =", id)
			idMond, err1 := strconv.Atoi(id)
			if err1 == nil && idMond > 0 { // если id int и id действительный то можно выбрать запись из monds
				var p frombase
				row := db.QueryRow("SELECT * FROM monds WHERE id=$1", idMond)
				err2 := row.Scan( // чтение строки из таблицы
					&p.id,
					&p.yahre,
					&p.nummonat,
					&p.tag,
					&p.hour,
					&p.kf,
					&p.blmond,
					&p.bltime,
					&p.bltabel,
					&p.blbuch,
					&p.blpers,
					&p.monat,
				)
				if err2 != nil {
					fmt.Println("ошибка чтения monds id=", id)
					//panic(err2)
				}
				admin.IdMonat = p.id
				idMond = p.id
				admin.JetzMonat = p.monat
				admin.JetzYahre = strconv.Itoa(p.yahre)
				yahre = admin.JetzYahre
				monat = admin.JetzMonat
				fmt.Println(p, idMond, yahre, monat)
			}
		} else {
			admin.JetzMonat = "???"
			admin.JetzYahre = "???"
		}
		exit, ok := req.URL.Query()["exit"]
		if ok && len(exit[0]) > 0 {
			_exit := exit[0]
			if _exit == "exit" {
				admin.User = ""
				admin.Email = ""
				admin.Passw = ""
				admin.Ready = "0"
				admin.JetzMonat = "???"
				admin.JetzYahre = "???"
				admin.IdMonat = 0 // год месяц не выбран
				idMond = 0
			}
		}
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Index Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		if req.Method == "POST" {
			admin.Ready = "0"    // 1 - ввод успешный
			admin.Errors = "0"   // 1 - ошибки при вводе
			admin.ErrEmail = "0" // 1 - ошибки при вводе логина
			admin.Empty = "0"    // 1 - есть пустые поля
			admin.Passpass = "0" // 1 - пароли совпали
			req.ParseForm()
			//admin.User = req.Form["user"][0]
			//admin.Email = req.Form["email"][0]
			//admin.Passw = req.Form["passw"][0]
			admin.User = "yp"
			admin.Email = "yp@yp.com"
			admin.Passw = "123"
			if admin.User == "" || admin.Email == "" || admin.Passw == "" {
				admin.Empty = "1"
				admin.Errors = "1"
			}
			var errmail int
			errmail, admin.Email, admin.User = inpMailAddress(admin.User + " <" + admin.Email + ">") // проверка email адреса
			if errmail > 0 {
				admin.ErrEmail = "1"
				admin.Errors = "1"
			}
			if admin.User != defaultUser || admin.Email != defaultEmail || admin.Passw != defaultPassw {
				admin.Passpass = "1"
				admin.Errors = "1"
			}
			if admin.Errors == "0" {
				admin.Ready = "1"
			}
		}
		err = t.ExecuteTemplate(w, "base", admin)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}
