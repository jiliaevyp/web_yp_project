package tabel

import (
	"database/sql"
	//"debug/dwarf"
	"fmt"
	"github.com/jiliaevyp/web_yp_project/mond"

	"html/template"
	"log"
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

type Tabel struct {
	Id         string
	Yahre      string
	Monat      string
	Nummonat   string
	Title      string
	Forename   string
	Email      string
	Kadr       string
	Department string
	Tarif      string // почасовая руб
	Tagemach   string
	Hourmach   string
	Oberhour   string
	Oplata     string
	Bonusproz  string
	Bonus      string
	Total      string
	MondId     string
	PersonId   string
	Jetzyahre  string
	Jetzmonat  string
	Ready      string // "1" - ввод корректен
	Errors     string // "1" - ошибка при вводе полей
	Empty      string // "1" - остались пустые поля
	ErrRange   string // "1" - выход за пределы диапазона
}

type Tabfrombase struct { // строка  при чтении/записи из/в базы personaldb
	Id         int
	Yahre      int
	Monat      string
	Nummonat   int
	Title      string
	Forename   string
	Email      string
	Kadr       string
	Department string
	Tarif      int // почасовая руб
	Tagemach   int
	Hourmach   int
	Oberhour   int
	Oplata     int
	Bonusproz  int
	Bonus      int
	Total      int
	MondId     int
	PersonId   int
}

var Tabtable struct {
	Ready       string
	Jetzyahre   string
	Jetzmonat   string
	Tag         string
	Hour        string
	Department  string
	Tabelstable []Tabel // таблица по сотрудниам  в personals_index.html
}

func IndexHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/tabels_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "tabels Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		// удаление записи
		del := req.URL.Query().Get("del")
		idhtml := req.URL.Query().Get("Id")
		id, _ := strconv.Atoi(idhtml)
		if del == "del" {
			_, err = db.Exec("DELETE FROM tabels WHERE Id = $1", id)
			if err != nil { // удаление старой записи
				panic(err)
			}
		}
		Tabtable.Tabelstable = nil
		var Tabelhtml Tabel
		// выборка записи из monds по IdRealMond
		idrealmond := mond.IdRealMond
		if idrealmond > 0 {
			row := db.QueryRow("SELECT * FROM monds WHERE Id=$1", idrealmond)
			var p mond.Mondfrombase
			err = row.Scan( // чтение строки из таблицы
				&p.Id,
				&p.Yahre,
				&p.Nummonat,
				&p.Tag,
				&p.Hour,
				&p.Kf,
				&p.Blmond,
				&p.Bltime,
				&p.Bltabel,
				&p.Blbuch,
				&p.Blpers,
				&p.Monat,
			)
			if err != nil {
				fmt.Println("Tabel ошибка распаковки строки monds")
				panic(err)
			}
			Tabelhtml.Tag = strconv.Itoa(p.Tag)
			Tabelhtml.Hour = strconv.Itoa(p.Hour)
			Tabelhtml.Yahre = strconv.Itoa(p.Yahre)
			Tabelhtml.Monat = p.Monat
			Tabelhtml.Nummonat = strconv.Itoa(p.Nummonat)
			Tabelhtml.Kf = strconv.Itoa(p.Kf)
		}
		// выборка по текущему месяцу через mond.IdRealMond
		rows, err1 := db.Query(`SELECT * FROM tabels WHERE MondId=$1`, mond.IdRealMond)
		if err1 != nil {
			fmt.Println(" table Tabels ошибка чтения ")
			fmt.Println("mondId=", mond.IdRealMond)
			panic(err1)
		}
		defer rows.Close()
		for rows.Next() {
			var p Tabfrombase
			err = rows.Scan( // пересылка  данных строки базы tabels в "p"
				&p.Id,
				&p.Yahre,
				&p.Monat,
				&p.Nummonat,
				&p.Title,
				&p.Forename,
				&p.Email,
				&p.Kadr,
				&p.Department,
				&p.Tarif,
				&p.Tagemach,
				&p.Hourmach,
				&p.Oberhour,
				&p.Oplata,
				&p.Bonusproz,
				&p.Bonus,
				&p.Total,
				&p.MondId,
				&p.PersonId,
			)
			if err != nil {
				fmt.Println("indexTabels ошибка распаковки строки ")
				panic(err)
				return
			}
			// подготовка к отображениею на tabel_index.html

			Tabelhtml.Id = strconv.Itoa(p.Id)
			Tabelhtml.Yahre = strconv.Itoa(p.Yahre)
			Tabelhtml.Monat = p.Monat
			Tabelhtml.Nummonat = strconv.Itoa(p.Nummonat)
			Tabelhtml.Title = p.Title
			Tabelhtml.Forename = p.Forename
			Tabelhtml.Email = p.Email
			Tabelhtml.Kadr = p.Kadr
			Tabelhtml.Department = p.Department
			Tabelhtml.Tarif = strconv.Itoa(p.Tarif) // int ---> string for HTML
			Tabelhtml.Tagemach = strconv.Itoa(p.Tagemach)
			Tabelhtml.Hourmach = strconv.Itoa(p.Hourmach)
			Tabelhtml.Oberhour = strconv.Itoa(p.Oberhour)
			Tabelhtml.Oplata = strconv.Itoa(p.Oplata)
			Tabelhtml.Bonusproz = strconv.Itoa(p.Bonusproz)
			Tabelhtml.Bonus = strconv.Itoa(p.Bonus)
			Tabelhtml.Total = strconv.Itoa(p.Total)
			Tabelhtml.MondId = strconv.Itoa(p.MondId)
			Tabelhtml.PersonId = strconv.Itoa(p.PersonId)
			Tabelhtml.Ready = "1"

			// добавление строки в таблицу Personalstab для personals_index.html
			Tabtable.Tabelstable = append(Tabtable.Tabelstable, Tabelhtml)
		}

		// sidebar
		Tabtable.Jetzyahre = mond.Jetzyahre
		Tabtable.Jetzmonat = mond.Jetzmonat
		Tabtable.Ready = "1"
		// загрузка tabel_index.html
		err = t.ExecuteTemplate(w, "base", Tabtable)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "tabels Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func ShowHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/tabel_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		idhtml := req.URL.Query().Get("Id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM tabels WHERE Id=$1", id)

		var p Tabfrombase
		err = row.Scan( // чтение строки из таблицы
			&p.Id,
			&p.Yahre,
			&p.Monat,
			&p.Title,
			&p.Forename,
			&p.Email,
			&p.Kadr,
			&p.Department,
			&p.Tarif,
			&p.Tagemach,
			&p.Hourmach,
			&p.Oberhour,
			&p.Oplata,
			&p.Bonusproz,
			&p.Bonus,
			&p.Total,
			&p.MondId,
			&p.PersonId,
		)
		if err != nil {
			fmt.Println("show tabels ошибка распаковки строки ")
			panic(err)
			return
		}
		var tabelhtml Tabel // переменная по сотруднику при вводе и отображении в personal.HTML
		tabelhtml.Id = strconv.Itoa(p.Id)
		tabelhtml.Yahre = strconv.Itoa(p.Yahre)
		tabelhtml.Monat = p.Monat
		tabelhtml.Title = p.Title
		tabelhtml.Forename = p.Forename
		tabelhtml.Email = p.Email
		tabelhtml.Kadr = p.Kadr
		tabelhtml.Department = p.Department
		tabelhtml.Tarif = strconv.Itoa(p.Tarif)
		tabelhtml.Tagemach = strconv.Itoa(p.Tagemach)
		tabelhtml.Hourmach = strconv.Itoa(p.Hourmach)
		tabelhtml.Oberhour = strconv.Itoa(p.Oberhour)
		tabelhtml.Oplata = strconv.Itoa(p.Oplata)
		tabelhtml.Bonusproz = strconv.Itoa(p.Bonusproz)
		tabelhtml.Bonus = strconv.Itoa(p.Bonus)
		tabelhtml.Total = strconv.Itoa(p.Total)
		tabelhtml.Ready = "1"    // "1" - ввод корректен
		tabelhtml.Errors = "0"   // "1" - ошибка при вводе полей
		tabelhtml.Empty = "0"    // "1" - остались пустые поля
		tabelhtml.ErrRange = "0" // "1" - выход за пределы диапазона

		Tabtable.Jetzyahre = mond.Jetzyahre
		Tabtable.Jetzmonat = mond.Jetzmonat

		err = t.ExecuteTemplate(w, "base", tabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func EditHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var tabelhtml Tabel // переменная по сотруднику при вводе и отображении в personal.HTML
		//var tabrow _tabelrow

		files := append(partials, "./static/tabel_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		Tabtable.Jetzyahre = mond.Jetzyahre
		Tabtable.Jetzmonat = mond.Jetzmonat

		err = t.ExecuteTemplate(w, "base", tabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func NewHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/tabel_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		var tabelhtml Tabel
		// выборка записи из monds по IdRealMond
		idrealmond := mond.IdRealMond
		if idrealmond > 0 {
			row := db.QueryRow("SELECT * FROM monds WHERE Id=$1", idrealmond)
			var p mond.Mondfrombase
			err = row.Scan( // чтение строки из таблицы
				&p.Id,
				&p.Yahre,
				&p.Nummonat,
				&p.Tag,
				&p.Hour,
				&p.Kf,
				&p.Blmond,
				&p.Bltime,
				&p.Bltabel,
				&p.Blbuch,
				&p.Blpers,
				&p.Monat,
			)
			if err != nil {
				fmt.Println("Tabel ошибка распаковки строки monds")
				panic(err)
			}
			tabelhtml.Tag = strconv.Itoa(p.Tag)
			tabelhtml.Monat = p.Monat
			tabelhtml.Nummonat = strconv.Itoa(p.Nummonat)
		}
		tabelhtml.Ready = "0"    // 1 - ввод успешный
		tabelhtml.Errors = "0"   // 1 - ошибки при вводе
		tabelhtml.ErrRange = "0" // 1 - выход за пределы диапазона
		if req.Method == "POST" {
			req.ParseForm()
			tabelhtml.Ready = "0" // 1 - ввод успешный
			tabelhtml.Errors = "0"
			tabelhtml.Title = req.Form["Title"][0]
			tabelhtml.Tagemach = req.Form["Tagemach"][0]
			tabelhtml.Hourmach = req.Form["Hourmach"][0]
			tabelhtml.Oberhour = req.Form["Oberhour"][0]

			//}
			//var p Tabfrombase
			//p.Yahre, _ = strconv.Atoi(tabelhtml.Yahre)
			//p.Nummonat, _ = strconv.Atoi(tabelhtml.Nummonat)
			//p.Tag, _ = strconv.Atoi(tabelhtml.Tag)
			//p.Hour, _ = strconv.Atoi(tabelhtml.Hour) // перевод в int для базы
			//p.Kf, _ = strconv.Atoi(tabelhtml.Kf)     // перевод в int для базы
			//p.Monat = monatArray[p.Nummonat-1]
			//
			//sqlStatement := `INSERT INTO monds (Yahre,Nummonat,Tag,Hour,Kf,blmond,bltime,bltabel,blbuch,blpers,Monat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
			//_, err2 := db.Exec(sqlStatement,
			//	&p.Yahre,
			//	&p.Nummonat,
			//	&p.Tag,
			//	&p.Hour,
			//	&p.Kf,
			//	&p.Monat,
			//)
			//if err2 != nil {
			//	fmt.Println("Ошибка записи новой строки в mondNew")
			//	panic(err2)
			//}
		}
		Tabtable.Jetzyahre = mond.Jetzyahre
		Tabtable.Jetzmonat = mond.Jetzmonat

		err1 := t.ExecuteTemplate(w, "base", tabelhtml)
		if err1 != nil {
			log.Println(err1.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

//f.select :personal_id, options_from_collection_for_select(Personal.where(personal_admin: $personal_admin), :Id, :Title)
