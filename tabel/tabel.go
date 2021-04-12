package tabel

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	_ "strconv"
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

type tabel struct {
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
	Tag        string
	Hour       string
	Kf         string
	Tagemach   string
	Hourmach   string
	Oberhour   string
	Oplata     string
	Bonusproz  string
	Bonus      string
	Total      string
	Ready      string // "1" - ввод корректен
	Errors     string // "1" - ошибка при вводе полей
	Empty      string // "1" - остались пустые поля
	ErrRange   string // "1" - выход за пределы диапазона
}

type frombase struct { // строка  при чтении/записи из/в базы personaldb
	id          int
	yahre       int
	monat       string
	nummonat    string
	title       string
	forename    string
	email       string
	kadr        string
	department  string
	tarif       int // почасовая руб
	tag         int
	hour        int
	kf          int
	tagemach    int
	hourmach    int
	oberhour    int
	oplata      int
	bonusproz   int
	bonus       int
	total       int
	mondsID     int
	personalsID int
}

var tabtable struct {
	Ready       string
	Jetzyahre   string
	Jetzmonat   string
	Department  string
	Tabelstable []tabel // таблица по сотрудниам  в personals_index.html
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

func IndexHandler(db *sql.DB, jetzYahre int, jetzMonat string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/tabels_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "tabels Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		del := req.URL.Query().Get("del")
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		if del == "del" {
			_, err = db.Exec("DELETE FROM tabels WHERE id = $1", id)
			if err != nil { // удаление старой записи
				panic(err)
			}
		}
		tabtable.Tabelstable = nil

		rows, err1 := db.Query(`SELECT * FROM tabels ORDER BY title`)
		if err1 != nil {
			fmt.Println(" table tabels ошибка чтения ")
			panic(err1)
		}
		defer rows.Close()

		for rows.Next() {
			var p frombase
			err = rows.Scan( // пересылка  данных строки базы personals в "p"
				&p.id,
				&p.yahre,
				&p.monat,
				&p.title,
				&p.forename,
				&p.email,
				&p.kadr,
				&p.department,
				&p.tarif,
				&p.tag,
				&p.hour,
				&p.kf,
				&p.tagemach,
				&p.hourmach,
				&p.oberhour,
				&p.oplata,
				&p.bonusproz,
				&p.bonus,
				&p.total,
				&p.mondsID,
				&p.personalsID,
			)
			if err != nil {
				fmt.Println("index tabels ошибка распаковки строки ")
				panic(err)
				return
			}
			//fmt.Println("id=",p.id)
			var tabelhtml tabel
			tabelhtml.Id = strconv.Itoa(p.id)
			tabelhtml.Yahre = strconv.Itoa(p.yahre)
			tabelhtml.Monat = p.monat
			tabelhtml.Title = p.title
			tabelhtml.Forename = p.forename
			tabelhtml.Email = p.email
			tabelhtml.Kadr = p.kadr
			tabelhtml.Department = p.department
			tabelhtml.Tarif = strconv.Itoa(p.tarif)
			tabelhtml.Tag = strconv.Itoa(p.tag)
			tabelhtml.Hour = strconv.Itoa(p.hour)
			tabelhtml.Kf = strconv.Itoa(p.kf)
			tabelhtml.Tagemach = strconv.Itoa(p.tagemach)
			tabelhtml.Hourmach = strconv.Itoa(p.hourmach)
			tabelhtml.Oberhour = strconv.Itoa(p.oberhour)
			tabelhtml.Oplata = strconv.Itoa(p.oplata)
			tabelhtml.Bonusproz = strconv.Itoa(p.bonusproz)
			tabelhtml.Bonus = strconv.Itoa(p.bonus)
			tabelhtml.Total = strconv.Itoa(p.total)
			tabelhtml.Ready = "1"    // "1" - ввод корректен
			tabelhtml.Errors = "0"   // "1" - ошибка при вводе полей
			tabelhtml.Empty = "0"    // "1" - остались пустые поля
			tabelhtml.ErrRange = "0" // "1" - выход за пределы диапазона

			// добавление строки в таблицу Personalstab для personals_index.html
			tabtable.Tabelstable = append(tabtable.Tabelstable, tabelhtml)
		}
		tabtable.Ready = "1"
		tabtable.Jetzyahre = strconv.Itoa(jetzYahre)
		tabtable.Jetzmonat = jetzMonat
		err = t.ExecuteTemplate(w, "base", tabtable)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "tabels Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func ShowHandler(db *sql.DB, jetzYahre int, jetzMonat string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/tabel_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM tabels WHERE id=$1", id)

		var p frombase
		err = row.Scan( // чтение строки из таблицы
			&p.id,
			&p.yahre,
			&p.monat,
			&p.title,
			&p.forename,
			&p.email,
			&p.kadr,
			&p.department,
			&p.tarif,
			&p.tag,
			&p.hour,
			&p.kf,
			&p.tagemach,
			&p.hourmach,
			&p.oberhour,
			&p.oplata,
			&p.bonusproz,
			&p.bonus,
			&p.total,
			&p.mondsID,
			&p.personalsID,
		)
		if err != nil {
			fmt.Println("show tabels ошибка распаковки строки ")
			panic(err)
			return
		}
		var tabelhtml tabel // переменная по сотруднику при вводе и отображении в personal.HTML
		tabelhtml.Id = strconv.Itoa(p.id)
		tabelhtml.Yahre = strconv.Itoa(p.yahre)
		tabelhtml.Monat = p.monat
		tabelhtml.Title = p.title
		tabelhtml.Forename = p.forename
		tabelhtml.Email = p.email
		tabelhtml.Kadr = p.kadr
		tabelhtml.Department = p.department
		tabelhtml.Tarif = strconv.Itoa(p.tarif)
		tabelhtml.Tag = strconv.Itoa(p.tag)
		tabelhtml.Hour = strconv.Itoa(p.hour)
		tabelhtml.Kf = strconv.Itoa(p.kf)
		tabelhtml.Tagemach = strconv.Itoa(p.tagemach)
		tabelhtml.Hourmach = strconv.Itoa(p.hourmach)
		tabelhtml.Oberhour = strconv.Itoa(p.oberhour)
		tabelhtml.Oplata = strconv.Itoa(p.oplata)
		tabelhtml.Bonusproz = strconv.Itoa(p.bonusproz)
		tabelhtml.Bonus = strconv.Itoa(p.bonus)
		tabelhtml.Total = strconv.Itoa(p.total)
		tabelhtml.Ready = "1"    // "1" - ввод корректен
		tabelhtml.Errors = "0"   // "1" - ошибка при вводе полей
		tabelhtml.Empty = "0"    // "1" - остались пустые поля
		tabelhtml.ErrRange = "0" // "1" - выход за пределы диапазона
		tabtable.Jetzyahre = strconv.Itoa(jetzYahre)
		tabtable.Jetzmonat = jetzMonat
		err = t.ExecuteTemplate(w, "base", tabelhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func EditHandler(db *sql.DB, jetzYahre int, jetzMonat string) func(w http.ResponseWriter, req *http.Request) {
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

func NewHandler(db *sql.DB, jetzYahre int, jetzMonat string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/tabel_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		var tabelhtml tabel
		tabelhtml.Ready = "0"    // 1 - ввод успешный
		tabelhtml.Errors = "0"   // 1 - ошибки при вводе
		tabelhtml.ErrRange = "0" // 1 - выход за пределы диапазона

		if req.Method == "POST" {
			req.ParseForm()
			//makeReadyHtml(&personalhtml) // подготовка значений для web
			//readFromHtml(&personalhtml, req)  	// ввод значений из web
			tabelhtml.Ready = "0" // 1 - ввод успешный
			tabelhtml.Errors = "0"

			tabelhtml.Yahre = req.Form["yahre"][0]
			tabelhtml.Nummonat = req.Form["nummonat"][0]

			nummonat, _ := strconv.Atoi(tabelhtml.Nummonat)
			_, err1 := db.Exec("DELETE FROM monds WHERE nummonat = $1", nummonat)
			if err1 != nil {
				fmt.Println("Ошибка при удалении старой записи в monds nummonat = ", nummonat)
				panic(err1)
			}
			//}
			//var p frombase
			//p.yahre, _ = strconv.Atoi(tabelhtml.Yahre)
			//p.nummonat, _ = strconv.Atoi(tabelhtml.Nummonat)
			//p.tag, _ = strconv.Atoi(tabelhtml.Tag)
			//p.hour, _ = strconv.Atoi(tabelhtml.Hour) // перевод в int для базы
			//p.kf, _ = strconv.Atoi(tabelhtml.Kf)     // перевод в int для базы
			//p.monat = monatArray[p.nummonat-1]
			//
			//sqlStatement := `INSERT INTO monds (yahre,nummonat,tag,hour,kf,blmond,bltime,bltabel,blbuch,blpers,monat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
			//_, err2 := db.Exec(sqlStatement,
			//	&p.yahre,
			//	&p.nummonat,
			//	&p.tag,
			//	&p.hour,
			//	&p.kf,
			//	&p.monat,
			//)
			//if err2 != nil {
			//	fmt.Println("Ошибка записи новой строки в mondNew")
			//	panic(err2)
			//}
		}
		tabtable.Jetzyahre = strconv.Itoa(jetzYahre)
		tabtable.Jetzmonat = jetzMonat
		err1 := t.ExecuteTemplate(w, "base", tabelhtml)
		if err1 != nil {
			log.Println(err1.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

//f.select :personal_id, options_from_collection_for_select(Personal.where(personal_admin: $personal_admin), :id, :title)
