package mond

import (
	"database/sql"
	_ "errors"
	"fmt"
	//"github.com/jiliaevyp/web_yp_project/servfunc"
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

type Monat struct { // данные по месяцу при вводе и отображении в mond.HTML
	Id        string
	Yahre     string
	Nummonat  string
	Monat     string
	Tag       string
	Hour      string
	Kf        string
	Blmond    string
	Blpers    string
	Bltime    string
	Bltabel   string
	Blbuch    string
	Ready     string // "1" - ввод корректен
	Errors    string // "1" - ошибка при вводе полей
	Empty     string // "1" - остались пустые поля
	ErrRange  string // "1" - выход за пределы диапазона
	Jetzmonat string // если id==0 значит месяц не выбран
	Jetzyahre string
}

type Mondfrombase struct { // для чтения из таблицы monds
	Id       int
	Yahre    int
	Nummonat int
	Tag      int
	Hour     int
	Kf       int
	Blmond   int
	Bltime   int
	Bltabel  int
	Blbuch   int
	Blpers   int
	Monat    string
}

//============================================================================
var (
	IdRealMond int    // id текущего рабочего месяца для всех таблиц
	Jetzmonat  string // если id==0 значит месяц не выбран
	Jetzyahre  string
)

//============================================================================

var mondtable struct {
	Ready      string // флаг готовности
	IdRealMond string
	Jetzmonat  string
	Jetzyahre  string
	Mondstable []Monat // таблица по сотрудниам  в monds_index.html
}

// просмотр таблицы monds из geoplastdb
func Indexhandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/monds_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		// обработка key  del & id ---> удаление записи
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
		// выборка всей таблицы
		rows, err1 := db.Query(`SELECT * FROM monds ORDER BY nummonat`)
		if err1 != nil {
			fmt.Println(" table monds ошибка чтения ")
			panic(err1)
		}
		var p Mondfrombase
		// установка текущего рабочего месяца
		// обработка key=="idrealmond" при срабатывании в mond_show кнопки
		// <button class="button"> <a href= "/monds_index?idrealmond={{ .Id}}" >Назначить рабочим</a></button>
		idmondfromshow := req.URL.Query().Get("idrealmond")
		if idmondfromshow != "" {
			// если была нажата кнопка то идет выбор рабочего месяца
			//fmt.Println("idrealmond=", idmondfromshow)
			IdRealMond, err = strconv.Atoi(idmondfromshow)
			if err != nil {
				IdRealMond = 0
				mondtable.IdRealMond = "0"
				mondtable.Jetzyahre = "не выбран!"
				mondtable.Jetzmonat = "не выбран"
				Jetzyahre = "не выбран!"
				Jetzmonat = "не выбран"
			} else {
				row := db.QueryRow("SELECT * FROM monds WHERE id=$1", IdRealMond)
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
					fmt.Println("ошибка распаковки строки IdRealMond")
					panic(err)
				} else {
					mondtable.IdRealMond = strconv.Itoa(p.Id)
					mondtable.Jetzyahre = strconv.Itoa(p.Yahre)
					mondtable.Jetzmonat = p.Monat

					IdRealMond = p.Id // глобальная id для модулей personals tabels
					Jetzyahre = mondtable.Jetzyahre
					Jetzmonat = mondtable.Jetzmonat
				}
			}
		}
		defer rows.Close()
		for rows.Next() {
			err = rows.Scan( // пересылка  данных строки базы personals в "p"
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
				fmt.Println("index monds ошибка распаковки строки ")
				panic(err)
				return
			}
			var monathtml Monat
			monathtml.Id = strconv.Itoa(p.Id)
			monathtml.Ready = "0" // "1" - ввод корректен
			monathtml.Id = strconv.Itoa(p.Id)
			monathtml.Yahre = strconv.Itoa(p.Yahre)
			monathtml.Nummonat = strconv.Itoa(p.Nummonat)
			monathtml.Monat = p.Monat
			monathtml.Tag = strconv.Itoa(p.Tag)
			monathtml.Hour = strconv.Itoa(p.Hour)
			monathtml.Kf = strconv.Itoa(p.Kf)
			monathtml.Blmond = strconv.Itoa(p.Blmond)
			monathtml.Bltime = strconv.Itoa(p.Bltime)
			monathtml.Bltabel = strconv.Itoa(p.Bltabel)
			monathtml.Blbuch = strconv.Itoa(p.Blbuch)
			monathtml.Blpers = strconv.Itoa(p.Blpers)
			// добавление строки в таблицу Personalstab для personals_index.html
			mondtable.Mondstable = append(mondtable.Mondstable, monathtml)
		}
		mondtable.Ready = "1"
		mondtable.Jetzyahre = Jetzyahre
		mondtable.Jetzmonat = Jetzmonat

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
		// отработка key = "id"
		idhtml := req.URL.Query().Get("Id")
		id, _ := strconv.Atoi(idhtml)
		fmt.Println("id=", id)
		row := db.QueryRow("SELECT * FROM monds WHERE id=$1", id)
		var p Mondfrombase
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
			fmt.Println("ошибка распаковки строки show")
			panic(err)
		}
		var monathtml Monat
		monathtml.Ready = "0" // 1 - ввод успешный
		// подготовка HTML
		monathtml.Id = strconv.Itoa(p.Id)
		monathtml.Yahre = strconv.Itoa(p.Yahre)
		monathtml.Nummonat = strconv.Itoa(p.Nummonat)
		monathtml.Monat = p.Monat
		monathtml.Tag = strconv.Itoa(p.Tag)
		monathtml.Hour = strconv.Itoa(p.Hour)
		monathtml.Kf = strconv.Itoa(p.Kf)
		monathtml.Blmond = strconv.Itoa(p.Blmond)
		monathtml.Bltime = strconv.Itoa(p.Bltime)
		monathtml.Bltabel = strconv.Itoa(p.Bltabel)
		monathtml.Blbuch = strconv.Itoa(p.Blbuch)
		monathtml.Blpers = strconv.Itoa(p.Blpers)

		monathtml.Jetzmonat = Jetzmonat
		monathtml.Jetzyahre = Jetzyahre
		err = t.ExecuteTemplate(w, "base", monathtml)
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
		var monathtml Monat
		if req.Method == "POST" {
			req.ParseForm()
			monathtml.Ready = "0" // 1 - ввод успешный
			monathtml.Errors = "0"
			monathtml.Yahre = req.Form["yahre"][0]
			monathtml.Nummonat = req.Form["nummonat"][0]
			monathtml.Monat = req.Form["Monat"][0]
			monathtml.Tag = req.Form["tag"][0]
			monathtml.Hour = req.Form["hour"][0]
			monathtml.Kf = req.Form["kf"][0]
			monathtml.Blmond = req.Form["blmond"][0]
			monathtml.Bltime = req.Form["bltime"][0]
			monathtml.Bltabel = req.Form["bltabel"][0]
			monathtml.Blbuch = req.Form["blbuch"][0]
			monathtml.Blpers = req.Form["blpers"][0]
			// перевод в int для базы
			var p Mondfrombase
			p.Yahre, _ = strconv.Atoi(monathtml.Yahre)
			p.Nummonat, _ = strconv.Atoi(monathtml.Nummonat)
			p.Monat = monathtml.Monat //monatArray[p.nummonat-1]
			p.Tag, _ = strconv.Atoi(monathtml.Tag)
			p.Hour, _ = strconv.Atoi(monathtml.Hour)
			p.Kf, _ = strconv.Atoi(monathtml.Kf)
			p.Blmond, _ = strconv.Atoi(monathtml.Blmond)
			p.Bltime, _ = strconv.Atoi(monathtml.Bltime)
			p.Bltabel, _ = strconv.Atoi(monathtml.Bltabel)
			p.Blbuch, _ = strconv.Atoi(monathtml.Blbuch)
			p.Blpers, _ = strconv.Atoi(monathtml.Blpers)
			nummonat := p.Nummonat
			_, err1 := db.Exec("DELETE FROM monds WHERE nummonat = $1", nummonat)
			if err1 != nil {
				fmt.Println("Ошибка при удалении старой записи в monds nummonat = ", nummonat)
				panic(err1)
			}
			sqlStatement := `INSERT INTO monds (yahre,nummonat,tag,hour,kf,blmond,bltime,bltabel,blbuch,blpers,Monat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
			_, err2 := db.Exec(sqlStatement,
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
			if err2 != nil {
				fmt.Println("Ошибка записи новой строки в mondNew")
				panic(err2)
			} else {
				monathtml.Ready = "1"
				//row := db.QueryRow("returning id")
			}
		}
		monathtml.Jetzmonat = Jetzmonat
		monathtml.Jetzyahre = Jetzyahre
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
		var p Mondfrombase
		err = row.Scan(
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
		if err == nil {
			// подготовка HTML
			var monathtml Monat
			monathtml.Ready = "0" // "1" - ввод корректен
			monathtml.Id = strconv.Itoa(p.Id)
			monathtml.Yahre = strconv.Itoa(p.Yahre)
			monathtml.Nummonat = strconv.Itoa(p.Nummonat)
			monathtml.Monat = p.Monat
			monathtml.Tag = strconv.Itoa(p.Tag)
			monathtml.Hour = strconv.Itoa(p.Hour)
			monathtml.Kf = strconv.Itoa(p.Kf)
			monathtml.Blmond = strconv.Itoa(p.Blmond)
			monathtml.Bltime = strconv.Itoa(p.Bltime)
			monathtml.Bltabel = strconv.Itoa(p.Bltabel)
			monathtml.Blbuch = strconv.Itoa(p.Blbuch)
			monathtml.Blpers = strconv.Itoa(p.Blpers)

			if req.Method == "POST" {
				req.ParseForm()
				monathtml.Errors = "0" // 1 - ввод успешный
				monathtml.Yahre = req.Form["yahre"][0]
				monathtml.Nummonat = req.Form["nummonat"][0]
				monathtml.Monat = req.Form["Monat"][0]
				monathtml.Tag = req.Form["tag"][0]
				monathtml.Hour = req.Form["hour"][0]
				monathtml.Kf = req.Form["kf"][0]
				monathtml.Blmond = req.Form["blmond"][0]
				monathtml.Bltime = req.Form["bltime"][0]
				monathtml.Bltabel = req.Form["bltabel"][0]
				monathtml.Blbuch = req.Form["blbuch"][0]
				monathtml.Blpers = req.Form["blpers"][0]
				// перевод в int для базы
				var p Mondfrombase
				p.Id, _ = strconv.Atoi(monathtml.Id)
				p.Yahre, _ = strconv.Atoi(monathtml.Yahre)
				p.Nummonat, _ = strconv.Atoi(monathtml.Nummonat)
				p.Monat = monathtml.Monat //monatArray[p.nummonat-1]
				p.Tag, _ = strconv.Atoi(monathtml.Tag)
				p.Hour, _ = strconv.Atoi(monathtml.Hour)
				p.Kf, _ = strconv.Atoi(monathtml.Kf)
				p.Blmond, _ = strconv.Atoi(monathtml.Blmond)
				p.Bltime, _ = strconv.Atoi(monathtml.Bltime)
				p.Bltabel, _ = strconv.Atoi(monathtml.Bltabel)
				p.Blbuch, _ = strconv.Atoi(monathtml.Blbuch)
				p.Blpers, _ = strconv.Atoi(monathtml.Blpers)
				//nummonat := p.Nummonat

				//_, err1 := db.Exec("DELETE FROM monds WHERE id = $1", id)
				//if err1 != nil {
				//	fmt.Println("Ошибка при удалении старой записи в monds id = ", id)
				//	panic(err1)
				//}
				sqlStatement := `INSERT INTO monds (id,yahre,nummonat,tag,hour,kf,blmond,bltime,bltabel,blbuch,blpers,Monat) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
				_, err2 := db.Exec(sqlStatement,
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
				if err2 != nil {
					fmt.Println("Ошибка записи edit строки в mondNew")
					monathtml.Errors = "1"
					panic(err2)
				} else {
					//fmt.Println( "новая запись id=", p.id)
					monathtml.Ready = "1"
				}
			}
			monathtml.Jetzmonat = Jetzmonat
			monathtml.Jetzyahre = Jetzyahre
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
