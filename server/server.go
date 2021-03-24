package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/jiliaevyp/web_yp_project/personals"
	//"github.com/jackc/pgx"
	//"github.com/jmoiron/sqlx"
)

const (
	defaultUser  = "yp"
	defaultEmail = "yp@yp.com"
	defaultPassw = "123"
)

var jetzYahre, jetzMonat string

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
var admin struct { // администратор
	User     string
	Email    string
	Passw    string
	ErrEmail string // ошибка ввода почты
	Passpass string
	Ready    string // 1 - идентификация прошла
	Errors   string // "1" - ошибка при вводе полей
	Empty    string // "1" - остались пустые поля
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

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

func server(addrWeb string, db *sql.DB) {
	http.HandleFunc("/", indexHandler)

	http.Handle("/mond_index", http.HandlerFunc(mondIndexHandler(db)))
	http.Handle("/mond_new", http.HandlerFunc(mondNewHandler(db)))
	http.Handle("/mond_show", http.HandlerFunc(mondShowHandler(db)))
	http.Handle("/mond_edit", http.HandlerFunc(mondEditHandler(db)))

	http.Handle("/personals_index", http.HandlerFunc(personalsIndexHandler(db)))
	http.Handle("/personal_new", http.HandlerFunc(personalNewhandler(db)))
	http.Handle("/personal_show", http.HandlerFunc(personalShowhandler(db)))
	http.Handle("/personal_edit", http.HandlerFunc(personalEdithandler(db)))

	http.Handle("/worktime_index", http.HandlerFunc(worktimeIndexHandler(db)))
	http.Handle("/worktime_new", http.HandlerFunc(worktimeNewHandler(db)))
	http.Handle("/worktime_show", http.HandlerFunc(worktimeShowHandler(db)))
	http.Handle("/worktime_edit", http.HandlerFunc(worktimeEditHandler(db)))

	http.Handle("/tabel_index", http.HandlerFunc(tabelIndexHandler(db)))
	http.Handle("/tabel_new", http.HandlerFunc(tabelNewHandler(db)))
	http.Handle("/tabel_show", http.HandlerFunc(tabelShowHandler(db)))
	http.Handle("/tabel_edit", http.HandlerFunc(tabelEditHandler(db)))

	http.Handle("/buchtabel_index", http.HandlerFunc(buchtabelIndexHandler(db)))
	http.Handle("/buchtabel_new", http.HandlerFunc(buchtabelNewHandler(db)))
	http.Handle("/buchtabel_show", http.HandlerFunc(buchtabelShowHandler(db)))
	http.Handle("/buchtabel_edit", http.HandlerFunc(buchtabelEditHandler(db)))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	fmt.Println("Топай на web страницу--->" + addrWeb + "!") // отладочная печать
	err := http.ListenAndServe(addrWeb, nil)
	if err != nil {
		errserv = 1
	} else {
		errserv = 0
	} // запуск сервера
	return
	//errserv := 0                                              // адрес в --> addrWeb
	//return  // форма в --> msgHandler
}

// первая страница проверка доступа
func indexHandler(w http.ResponseWriter, req *http.Request) {
	files := append(partials, "./static/index.html")
	t, err := template.ParseFiles(files...) // Parse template file.
	exit, ok := req.URL.Query()["exit"]
	if ok && len(exit[0]) > 0 {
		_exit := exit[0]
		if _exit == "exit" {
			admin.User = ""
			admin.Email = ""
			admin.Passw = ""
			admin.Ready = "0"
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
