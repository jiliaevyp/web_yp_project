package main

import (
	"database/sql"
	//"golang.org/x/net/dns/dnsmessage"
	"html/template"
	"log"
	"net/http"
	//"strconv"
)

type buchtabel struct {
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

type _buchtabelrow struct { // строка  при чтении/записи из/в базы personaldb
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

type buchtabeltab struct { // данные по сотруднику при в отображении строки в personals_index.html
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

var buchtab struct {
	Ready      string
	Buttonshow string         // просмотр сотрудника
	Buchtab    []buchtabeltab // таблица по сотрудниам  в personals_index.html
}

func buchtabelIndexHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/buchtabel_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", buchtab)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func buchtabelShowHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var buchtabelhtml buchtabel // переменная по сотруднику при вводе и отображении в personal.HTML

		files := append(partials, "./static/buchtabel_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", buchtabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func buchtabelEditHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var buchtabelhtml buchtabel // переменная по сотруднику при вводе и отображении в personal.HTML

		files := append(partials, "./static/buchtabel_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", buchtabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}
func buchtabelNewHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var buchtabelhtml buchtabel // переменная по сотруднику при вводе и отображении в personal.HTML

		files := append(partials, "./static/buchtabel_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", buchtabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}
