package tabel

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	_ "strconv"
)

type tabel struct {
	Title     string
	Kadr      string
	Otdel     string
	Tag       string
	Hour      string
	Tagemach  string
	Hourmach  string
	Oberhour  string
	Tarif     string // почасовая руб
	Raise     string
	Urlaub    string
	Krank     string
	Ready     string // "1" - ввод корректен
	Errors    string // "1" - ошибка при вводе полей
	ErrPhone  string // "1"- ошибка при вводе телефона
	ErrEmail  string // "1"- ошибка при вводе email
	ErrTarif  string // "1"- ошибка при вводе тарифа
	ErrNumotd string // "1"- ошибка при вводе номера отдела
	Empty     string // "1" - остались пустые поля
	ErrRange  string // "1" - выход за пределы диапазона
}

type _tabelrow struct { // строка  при чтении/записи из/в базы personaldb
	Id       int
	Name     string
	Title    string
	Kadr     string
	Otdel    string
	Tag      int
	Hour     int
	Tagemach int
	Hourmach int
	Oberhour int
	Tarif    int // почасовая руб
}

type tabtab struct { // данные по сотруднику при в отображении строки в personals_index.html
	Title    string
	Kadr     string
	Otdel    string
	Tag      string
	Hour     string
	Tagemach string
	Hourmach string
	Oberhour string
	Tarif    string // почасовая руб
}

var tabeltab struct {
	Ready      string
	Buttonshow string   // просмотр сотрудника
	tabeltab   []tabtab // таблица по сотрудниам  в personals_index.html
}

func tabelIndexHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/tabel_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", tabeltab)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func tabelShowHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var tabelhtml tabel // переменная по сотруднику при вводе и отображении в personal.HTML
		//var tabrow _tabelrow

		files := append(partials, "./static/tabel_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", tabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func tabelEditHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var tabelhtml tabel // переменная по сотруднику при вводе и отображении в personal.HTML
		//var tabrow _tabelrow

		files := append(partials, "./static/tabel_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", tabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}
func tabelNewHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var tabelhtml tabel // переменная по сотруднику при вводе и отображении в personal.HTML

		files := append(partials, "./static/tabel_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		//tabelhtml.Title = personalhtml.Title
		//tabelhtml.Kadr = personalhtml.Kadr
		//tabelhtml.Otdel = personalhtml.Otdel
		//tabelhtml.Hour = "168"
		//tabelhtml.Tag = "21"
		//if req.Method == "POST" {
		//	req.ParseForm()
		//	tabelhtml.Ready = "0"    // 1 - ввод успешный
		//	tabelhtml.Errors = "0"   // 1 - ошибки при вводе
		//	tabelhtml.Empty = "0"    // 1 - есть пустые поля
		//	tabelhtml.ErrRange = "0" // 1 - выход за пределы диапазона
		//	if tabelhtml.Errors == "0" {
		//		tabelhtml.Ready = "1"
		//	}
		//}
		err = t.ExecuteTemplate(w, "base", tabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}
