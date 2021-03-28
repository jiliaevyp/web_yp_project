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

type monat struct { // данные по месяцу при вводе и отображении в mond.HTML
	Id       string
	Yahre    string
	Nummonat string
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

		files := append(partials, "./static/mond_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		del := req.URL.Query().Get("del")
		nummonat := req.URL.Query().Get("nummonat")
		num, _ := strconv.Atoi(nummonat)
		if del == "del" {
			_, err = db.Exec("DELETE FROM personals WHERE nummonat = $2", num)
			if err != nil { // удаление старой записи
				panic(err)
			}
		}
		mondtable.Mondstable = nil

		rows, err1 := db.Query(`SELECT * FROM monds`)
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
			)
			if err != nil {
				fmt.Println("index monds ошибка распаковки строки ")
				panic(err)
				return
			}
			var monatlhtml monat
			monatlhtml.Yahre = strconv.Itoa(p.yahre)
			monatlhtml.Nummonat = strconv.Itoa(p.nummonat)
			monatlhtml.Tag = strconv.Itoa(p.tag)
			monatlhtml.Hour = strconv.Itoa(p.hour) // int ---> string for HTML
			monatlhtml.Kf = strconv.Itoa(p.kf)
			monatlhtml.Blmond = strconv.Itoa(p.blmond)
			monatlhtml.Bltime = strconv.Itoa(p.bltime)
			monatlhtml.Bltabel = strconv.Itoa(p.bltabel)
			monatlhtml.Blbuch = strconv.Itoa(p.blbuch)
			monatlhtml.Blpers = strconv.Itoa(p.blpers)
			monatlhtml.Ready = "1"    // "1" - ввод корректен
			monatlhtml.Errors = "0"   // "1" - ошибка при вводе полей
			monatlhtml.Empty = "0"    // "1" - остались пустые поля
			monatlhtml.ErrRange = "0" // "1" - выход за пределы диапазона

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
		monathtml.Ready = "1"    // 1 - ввод успешный
		monathtml.Errors = "0"   // 1 - ошибки при вводе
		monathtml.Empty = "0"    // 1 - есть пустые поля
		monathtml.ErrRange = "0" // 1 - выход за пределы диапазона

		nummonat := req.URL.Query().Get("nummonat")
		row := db.QueryRow("SELECT * FROM monds WHERE nummonat=$2", nummonat)

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
		)
		if err != nil {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		} // подготовка HTML
		var monatlhtml monat
		monatlhtml.Yahre = strconv.Itoa(p.yahre)
		monatlhtml.Nummonat = strconv.Itoa(p.nummonat)
		monatlhtml.Tag = strconv.Itoa(p.tag)
		monatlhtml.Hour = strconv.Itoa(p.hour) // int ---> string for HTML
		monatlhtml.Kf = strconv.Itoa(p.kf)
		monatlhtml.Blmond = strconv.Itoa(p.blmond)
		monatlhtml.Bltime = strconv.Itoa(p.bltime)
		monatlhtml.Bltabel = strconv.Itoa(p.bltabel)
		monatlhtml.Blbuch = strconv.Itoa(p.blbuch)
		monatlhtml.Blpers = strconv.Itoa(p.blpers)
		// проверка корректности ввода
		monatlhtml.Ready = "1"    // "1" - ввод корректен
		monatlhtml.Errors = "0"   // "1" - ошибка при вводе полей
		monatlhtml.Empty = "0"    // "1" - остались пустые поля
		monatlhtml.ErrRange = "0" // "1" - выход за пределы диапазона

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
		monathtml.Ready = "0"    // 1 - ввод успешный
		monathtml.Errors = "0"   // 1 - ошибки при вводе
		monathtml.ErrRange = "0" // 1 - выход за пределы диапазона

		if req.Method == "POST" {
			req.ParseForm()
			//makeReadyHtml(&personalhtml) // подготовка значений для web
			//readFromHtml(&personalhtml, req)  	// ввод значений из web
			monathtml.Ready = "0" // 1 - ввод успешный
			monathtml.Errors = "0"

			monathtml.Yahre = req.Form["yahre"][0]
			if checknum(monathtml.Yahre, 2021, 2031) != 0 {
				monathtml.Yahre = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Nummonat = req.Form["nummonat"][0]
			if checknum(monathtml.Nummonat, 0, 12) != 0 {
				monathtml.Nummonat = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Tag = req.Form["tag"][0]
			if checknum(monathtml.Tag, 10, 24) != 0 {
				monathtml.Tag = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Hour = req.Form["hour"][0]
			if checknum(monathtml.Hour, 50, 180) != 0 {
				monathtml.Hour = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Kf = req.Form["kf"][0]
			if checknum(monathtml.Kf, 1, 2) != 0 {
				monathtml.Kf = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Blmond = req.Form["blmond"][0]
			if checknum(monathtml.Blmond, 0, 1) != 0 {
				monathtml.Blmond = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Bltime = req.Form["bltime"][0]
			if checknum(monathtml.Bltime, 0, 1) != 0 {
				monathtml.Bltime = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Bltabel = req.Form["bltabel"][0]
			if checknum(monathtml.Bltabel, 0, 1) != 0 {
				monathtml.Bltabel = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Blbuch = req.Form["blbuch"][0]
			if checknum(monathtml.Blbuch, 0, 1) != 0 {
				monathtml.Blbuch = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Blpers = req.Form["blpers"][0]
			if checknum(monathtml.Blpers, 0, 1) != 0 {
				monathtml.Blpers = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			if monathtml.Errors == "0" {
				monathtml.Ready = "1"
				//добавление записи в базу
				nummonat := monathtml.Nummonat
				// удаление старой записи
				row := db.QueryRow("SELECT * FROM monds WHERE nummonat=$2", nummonat)
				if row != nil { // если запись есть удаляем
					_, err1 := db.Exec("DELETE FROM monds WHERE nummonat = $2", nummonat)
					if err1 != nil {
						fmt.Println("Ошибка при удалении старой записи в monds nummonat = ", nummonat)
						panic(err)
					}
				}
				var p frombase
				p.yahre, _ = strconv.Atoi(monathtml.Yahre)
				p.nummonat, _ = strconv.Atoi(monathtml.Nummonat)
				p.tag, _ = strconv.Atoi(monathtml.Tag)
				p.hour, _ = strconv.Atoi(monathtml.Hour) // перевод в int для базы
				p.kf, _ = strconv.Atoi(monathtml.Kf)     // перевод в int для базы
				p.blmond, _ = strconv.Atoi(monathtml.Blmond)
				p.bltime, _ = strconv.Atoi(monathtml.Bltime)
				p.bltabel, _ = strconv.Atoi(monathtml.Bltabel)
				p.blbuch, _ = strconv.Atoi(monathtml.Blbuch)
				p.blpers, _ = strconv.Atoi(monathtml.Blpers)

				sqlStatement := `INSERT INTO monds (yahre,nummonat,tag,hour,kf,blmond,bltime,bltabel,blbuch,blpers) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
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
				)
				if err2 != nil {
					fmt.Println("Ошибка записи новой строки в mondNew")
					panic(err2)
				}
			}
		}
		err = t.ExecuteTemplate(w, "base", monathtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Newmond Internal Server Execute Error", http.StatusInternalServerError)
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

		nummonat := req.URL.Query().Get("nummonat")
		num, _ := strconv.Atoi(nummonat)
		row := db.QueryRow("SELECT * FROM monds WHERE nummonat=$2", num)

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
		)
		if err != nil {
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		}
		// подготовка HTML
		var monathtml monat
		monathtml.Empty = "1"
		monathtml.Ready = "0"    // "1" - ввод корректен
		monathtml.Errors = "0"   // "1" - ошибка при вводе полей
		monathtml.ErrRange = "0" // "1" - выход за пределы диапазона
		monathtml.Yahre = strconv.Itoa(p.yahre)
		monathtml.Nummonat = strconv.Itoa(p.nummonat)
		monathtml.Tag = strconv.Itoa(p.tag)
		monathtml.Hour = strconv.Itoa(p.hour) // int ---> string for HTML
		monathtml.Kf = strconv.Itoa(p.kf)
		monathtml.Blmond = strconv.Itoa(p.blmond)
		monathtml.Bltime = strconv.Itoa(p.bltime)
		monathtml.Bltabel = strconv.Itoa(p.bltabel)
		monathtml.Blbuch = strconv.Itoa(p.blbuch)
		monathtml.Blpers = strconv.Itoa(p.blpers)
		// проверка корректности ввода

		if req.Method == "POST" {
			req.ParseForm()
			monathtml.Ready = "0"  // "1" - ввод корректен
			monathtml.Errors = "0" // "1" - ошибка при вводе полей
			//makeReadyHtml(&personalhtml) // подготовка значений для web
			//readFromHtml(&personalhtml, req)  	// ввод значений из web
			monathtml.Yahre = req.Form["yahre"][0]
			if checknum(monathtml.Yahre, 2021, 2031) != 0 {
				monathtml.Yahre = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Nummonat = req.Form["nummonat"][0]
			if checknum(monathtml.Nummonat, 0, 12) != 0 {
				monathtml.Nummonat = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Tag = req.Form["tag"][0]
			if checknum(monathtml.Tag, 10, 24) != 0 {
				monathtml.Tag = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Hour = req.Form["hour"][0]
			if checknum(monathtml.Hour, 50, 180) != 0 {
				monathtml.Hour = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Kf = req.Form["kf"][0]
			if checknum(monathtml.Kf, 1, 2) != 0 {
				monathtml.Kf = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Blmond = req.Form["blmond"][0]
			if checknum(monathtml.Blmond, 0, 1) != 0 {
				monathtml.Blmond = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Bltime = req.Form["bltime"][0]
			if checknum(monathtml.Bltime, 0, 1) != 0 {
				monathtml.Bltime = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Bltabel = req.Form["bltabel"][0]
			if checknum(monathtml.Bltabel, 0, 1) != 0 {
				monathtml.Bltabel = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Blbuch = req.Form["blbuch"][0]
			if checknum(monathtml.Blbuch, 0, 1) != 0 {
				monathtml.Blbuch = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"
			}
			monathtml.Blpers = req.Form["blpers"][0]
			if checknum(monathtml.Blpers, 0, 1) != 0 {
				monathtml.Blpers = "???"
				monathtml.ErrRange = "1"
				monathtml.Errors = "1"

			}
			if monathtml.Errors == "0" {
				monathtml.Ready = "1"
				monathtml.Empty = "0"
				//удаление старой записи в базе

				// удаление старой записи
				row := db.QueryRow("SELECT * FROM monds WHERE nummonat=$2", num)
				if row != nil { // если запись есть удаляем
					_, err1 := db.Exec("DELETE FROM monds WHERE nummonat = $2", num)
					if err1 != nil {
						fmt.Println("Ошибка при удалении старой записи в monds nummonat = ", num)
						panic(err)
					}
				}
				//добавление записи в базу
				var p frombase
				p.yahre, _ = strconv.Atoi(monathtml.Yahre)
				p.nummonat, _ = strconv.Atoi(monathtml.Nummonat)
				p.tag, _ = strconv.Atoi(monathtml.Tag)
				p.hour, _ = strconv.Atoi(monathtml.Hour) // перевод в int для базы
				p.kf, _ = strconv.Atoi(monathtml.Kf)     // перевод в int для базы
				p.blmond, _ = strconv.Atoi(monathtml.Blmond)
				p.bltime, _ = strconv.Atoi(monathtml.Bltime)
				p.bltabel, _ = strconv.Atoi(monathtml.Bltabel)
				p.blbuch, _ = strconv.Atoi(monathtml.Blbuch)
				p.blpers, _ = strconv.Atoi(monathtml.Blpers)

				sqlStatement := `INSERT INTO monds (yahre,nummonat,tag,hour,kf,blmond,bltime,bltabel,blbuch,blpers) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
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
				)
				if err2 != nil {
					fmt.Println("Ошибка записи новой строки в mondNew")
					panic(err2)
				}
			}
		}
		err = t.ExecuteTemplate(w, "base", monathtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}

}
