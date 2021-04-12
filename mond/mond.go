package mond

import (
	"database/sql"
	_ "errors"
	"fmt"
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
var monatArray = [12]string{
	"январь",
	"февраль",
	"март",
	"апрель",
	"май",
	"июнь",
	"июль",
	"август",
	"сентябрь",
	"октябрь",
	"ноябрь",
	"декабрь",
}

type monat struct { // данные по месяцу при вводе и отображении в mond.HTML
	Id       string
	Yahre    string
	Nummonat string
	Monat    string
	Tag      string
	Hour     string
	Kf       string
	Blmond   string
	Blpers   string
	Bltime   string
	Bltabel  string
	Blbuch   string
	Ready    string // "1" - ввод корректен
	Errors   string // "1" - ошибка при вводе полей
	Empty    string // "1" - остались пустые поля
	ErrRange string // "1" - выход за пределы диапазона
}

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

var mondtable struct {
	Ready      string  // флаг готовности
	Mondstable []monat // таблица по сотрудниам  в monds_index.html
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

// просмотр таблицы из personaldb
func Indexhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/monds_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		// обработка key  del & id
		del := req.URL.Query().Get("del")
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		if del == "del" {
			_, err = db.Exec("DELETE FROM monds WHERE id = $1", id)
			if err != nil { // удаление старой записи
				panic(err)
			}
		}
		mondtable.Mondstable = nil

		rows, err1 := db.Query(`SELECT * FROM monds ORDER BY nummonat`)
		if err1 != nil {
			fmt.Println(" table monds ошибка чтения ")
			panic(err1)
		}
		defer rows.Close()
		for rows.Next() {
			var p frombase
			err = rows.Scan( // пересылка  данных строки базы personals в "p"
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
			if err != nil {
				fmt.Println("index monds ошибка распаковки строки ")
				panic(err)
				return
			}
			var monatlhtml monat
			monatlhtml.Id = strconv.Itoa(p.id)
			monatlhtml.Yahre = strconv.Itoa(p.yahre)
			monatlhtml.Nummonat = strconv.Itoa(p.nummonat)
			monatlhtml.Monat = p.monat
			monatlhtml.Tag = strconv.Itoa(p.tag)
			monatlhtml.Hour = strconv.Itoa(p.hour) // int ---> string for HTML
			monatlhtml.Kf = strconv.Itoa(p.kf)
			monatlhtml.Blmond = strconv.Itoa(p.blmond)
			monatlhtml.Bltime = strconv.Itoa(p.bltime)
			monatlhtml.Bltabel = strconv.Itoa(p.bltabel)
			monatlhtml.Blbuch = strconv.Itoa(p.blbuch)
			monatlhtml.Blpers = strconv.Itoa(p.blpers)
			// добавление строки в таблицу Personalstab для personals_index.html
			mondtable.Mondstable = append(mondtable.Mondstable, monatlhtml)
		}
		mondtable.Ready = "1"
		err = t.ExecuteTemplate(w, "base", mondtable)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Index Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

// просмотр записи из personaldb
func Showhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/mond_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		var monathtml monat
		monathtml.Ready = "0" // 1 - ввод успешный
		//realIdMonat := req.URL.Query().Get("realIdMonat")
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM monds WHERE id=$1", id)
		var p frombase
		err = row.Scan( // чтение строки из таблицы
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
		if err != nil {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		}
		// подготовка HTML
		var monatlhtml monat
		monatlhtml.Id = strconv.Itoa(p.id)
		monatlhtml.Yahre = strconv.Itoa(p.yahre)
		monatlhtml.Nummonat = strconv.Itoa(p.nummonat)
		monatlhtml.Monat = p.monat
		monatlhtml.Tag = strconv.Itoa(p.tag)
		monatlhtml.Hour = strconv.Itoa(p.hour)
		monatlhtml.Kf = strconv.Itoa(p.kf)
		monatlhtml.Blmond = strconv.Itoa(p.blmond)
		monatlhtml.Bltime = strconv.Itoa(p.bltime)
		monatlhtml.Bltabel = strconv.Itoa(p.bltabel)
		monatlhtml.Blbuch = strconv.Itoa(p.blbuch)
		monatlhtml.Blpers = strconv.Itoa(p.blpers)
		err = t.ExecuteTemplate(w, "base", monatlhtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error Showhandler", http.StatusInternalServerError)
			return
		}
	}
}

// новая запись формы personal в базу personaldb
func Newhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/mond_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error personalNewhandler", http.StatusInternalServerError)
			return
		}
		var monathtml monat
		if req.Method == "POST" {
			req.ParseForm()
			monathtml.Ready = "0" // 1 - ввод успешный
			monathtml.Errors = "0"
			monathtml.Yahre = req.Form["yahre"][0]
			monathtml.Nummonat = req.Form["nummonat"][0]
			monathtml.Monat = req.Form["monat"][0]
			monathtml.Tag = req.Form["tag"][0]
			monathtml.Hour = req.Form["hour"][0]
			monathtml.Kf = req.Form["kf"][0]
			monathtml.Blmond = req.Form["blmond"][0]
			monathtml.Bltime = req.Form["bltime"][0]
			monathtml.Bltabel = req.Form["bltabel"][0]
			monathtml.Blbuch = req.Form["blbuch"][0]
			monathtml.Blpers = req.Form["blpers"][0]
			// перевод в int для базы
			var p frombase
			p.yahre, _ = strconv.Atoi(monathtml.Yahre)
			p.nummonat, _ = strconv.Atoi(monathtml.Nummonat)
			p.monat = monathtml.Monat //monatArray[p.nummonat-1]
			p.tag, _ = strconv.Atoi(monathtml.Tag)
			p.hour, _ = strconv.Atoi(monathtml.Hour)
			p.kf, _ = strconv.Atoi(monathtml.Kf)
			p.blmond, _ = strconv.Atoi(monathtml.Blmond)
			p.bltime, _ = strconv.Atoi(monathtml.Bltime)
			p.bltabel, _ = strconv.Atoi(monathtml.Bltabel)
			p.blbuch, _ = strconv.Atoi(monathtml.Blbuch)
			p.blpers, _ = strconv.Atoi(monathtml.Blpers)
			nummonat := p.nummonat
			_, err1 := db.Exec("DELETE FROM monds WHERE nummonat = $1", nummonat)
			if err1 != nil {
				fmt.Println("Ошибка при удалении старой записи в monds nummonat = ", nummonat)
				panic(err1)
			}
			sqlStatement := `INSERT INTO monds (yahre,nummonat,tag,hour,kf,blmond,bltime,bltabel,blbuch,blpers,monat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
			_, err2 := db.Exec(sqlStatement,
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
				fmt.Println("Ошибка записи новой строки в mondNew")
				panic(err2)
			} else {
				monathtml.Ready = "1"
				//row := db.QueryRow("returning id")
			}
		}
		err3 := t.ExecuteTemplate(w, "base", monathtml)
		if err3 != nil {
			log.Println(err.Error())
			http.Error(w, "Newmond Internal Server Execute Error", http.StatusInternalServerError)
			panic(err3)
			return
		}
	}
}

// редактирование формы personal и замена в базе personaldb
func Edithandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/mond_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		idhtml := req.URL.Query().Get("id")
		id, _ := strconv.Atoi(idhtml)
		row := db.QueryRow("SELECT * FROM monds WHERE id=$1", id)
		// выборка записи из таблицы
		var p frombase
		err = row.Scan(
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
		if err == nil {
			// подготовка HTML
			//id := p.id
			var monathtml monat
			monathtml.Ready = "0" // "1" - ввод корректен
			monathtml.Id = strconv.Itoa(p.id)
			monathtml.Yahre = strconv.Itoa(p.yahre)
			monathtml.Nummonat = strconv.Itoa(p.nummonat)
			monathtml.Monat = p.monat //string
			monathtml.Tag = strconv.Itoa(p.tag)
			monathtml.Hour = strconv.Itoa(p.hour) // int ---> string for HTML
			monathtml.Kf = strconv.Itoa(p.kf)
			monathtml.Blmond = strconv.Itoa(p.blmond)
			monathtml.Bltime = strconv.Itoa(p.bltime)
			monathtml.Bltabel = strconv.Itoa(p.bltabel)
			monathtml.Blbuch = strconv.Itoa(p.blbuch)
			monathtml.Blpers = strconv.Itoa(p.blpers)

			if req.Method == "POST" {
				req.ParseForm()
				monathtml.Errors = "0" // 1 - ввод успешный
				monathtml.Yahre = req.Form["yahre"][0]
				monathtml.Nummonat = req.Form["nummonat"][0]
				monathtml.Monat = req.Form["monat"][0]
				monathtml.Tag = req.Form["tag"][0]
				monathtml.Hour = req.Form["hour"][0]
				monathtml.Kf = req.Form["kf"][0]
				monathtml.Blmond = req.Form["blmond"][0]
				monathtml.Bltime = req.Form["bltime"][0]
				monathtml.Bltabel = req.Form["bltabel"][0]
				monathtml.Blbuch = req.Form["blbuch"][0]
				monathtml.Blpers = req.Form["blpers"][0]
				// перевод в int для базы
				var p frombase
				p.id, _ = strconv.Atoi(monathtml.Id)
				p.yahre, _ = strconv.Atoi(monathtml.Yahre)
				p.nummonat, _ = strconv.Atoi(monathtml.Nummonat)
				p.monat = monathtml.Monat
				p.tag, _ = strconv.Atoi(monathtml.Tag)
				p.hour, _ = strconv.Atoi(monathtml.Hour)
				p.kf, _ = strconv.Atoi(monathtml.Kf)
				p.blmond, _ = strconv.Atoi(monathtml.Blmond)
				p.bltime, _ = strconv.Atoi(monathtml.Bltime)
				p.bltabel, _ = strconv.Atoi(monathtml.Bltabel)
				p.blbuch, _ = strconv.Atoi(monathtml.Blbuch)
				p.blpers, _ = strconv.Atoi(monathtml.Blpers)

				//_, err1 := db.Exec("DELETE FROM monds WHERE id = $1", id)
				//if err1 != nil {
				//	fmt.Println("Ошибка при удалении старой записи в monds id = ", id)
				//	panic(err1)
				//}
				sqlStatement := `INSERT INTO monds (id,yahre,nummonat,tag,hour,kf,blmond,bltime,bltabel,blbuch,blpers,monat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				_, err2 := db.Exec(sqlStatement,
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
					fmt.Println("Ошибка записи edit строки в mondNew")
					monathtml.Errors = "1"
					panic(err2)
				} else {
					//fmt.Println( "новая запись id=", p.id)
					monathtml.Ready = "1"
				}
			}
			err3 := t.ExecuteTemplate(w, "base", monathtml)
			if err3 != nil {
				log.Println(err.Error())
				http.Error(w, "Newmond Internal Server Execute Error", http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		}
	}
}
