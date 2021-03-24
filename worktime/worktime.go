package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	//"strconv"
)

var work struct {
	Monat    string
	Yahre    string
	Title    string
	Kadr     string
	Otdel    string
	Tag      string
	Hour     string
	Koef     string
	Tagemach string
	Hourmach string
	Oberhour string
	Raise    string
	Urlaub   string
	Krank    string
	Ready    string // "1" - ввод корректен
	Errors   string // "1" - ошибка при вводе полей
	Empty    string // "1" - остались пустые поля
	Range    string // "1" - выход за пределы диапазона
}

type worktime struct {
	Monat     string
	Yahre     string
	Title     string
	Kadr      string
	Otdel     string
	Tag       string
	Hour      string
	Koef      string
	Tagemach  string
	Hourmach  string
	Oberhour  string
	Ready     string // "1" - ввод корректен
	Errors    string // "1" - ошибка при вводе полей
	ErrPhone  string // "1"- ошибка при вводе телефона
	ErrEmail  string // "1"- ошибка при вводе email
	ErrTarif  string // "1"- ошибка при вводе тарифа
	ErrNumotd string // "1"- ошибка при вводе номера отдела
	Empty     string // "1" - остались пустые поля
	ErrRange  string // "1" - выход за пределы диапазона
}

type _worktabelrow struct { // строка  при чтении/записи из/в базы personaldb
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

type worktabeltab struct { // данные по сотруднику при в отображении строки в personals_index.html
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

var worktab struct {
	Ready      string
	Buttonshow string         // просмотр сотрудника
	tabeltab   []worktabeltab // таблица по сотрудниам  в personals_index.html
}

func worktimeIndexHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/worktime_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", worktab)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func worktimeShowHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var worktimehtml worktime

		files := append(partials, "./static/worktime_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", worktimehtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func worktimeEditHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var worktimehtml worktime

		files := append(partials, "./static/worktime_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", worktimehtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func worktimeNewHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var worktimehtml worktime

		files := append(partials, "./static/worktime_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		err = t.ExecuteTemplate(w, "base", worktimehtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}
