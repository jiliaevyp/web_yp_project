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

type monatHtml struct { // данные по месяцу при вводе и отображении в mond.HTML
	Id             string
	Yahre          string
	Monat          string
	Nummonat       string
	Tag            string
	Hour           string
	Kf             string
	Blockmonat     string
	Nalog          string
	Blockpersonal  string
	Blocktimetabel string
	Blocktabel     string
	Blockbuchtabel string
	Timestamp      string
	Ready          string // "1" - ввод корректен
	Errors         string // "1" - ошибка при вводе полей
	Empty          string // "1" - остались пустые поля
	Range          string // "1" - выход за пределы диапазона
}

type monatrow struct { // для чтения из таблицы monds
	Yahre          int
	Monat          string
	Nummonat       int
	Tag            int
	Hour           int
	Kf             int
	Blockmonat     int
	Nalog          int
	Blockpersonal  int
	Blocktimetabel int
	Blocktabel     int
	Blockbuchtabel int
	Id             int
	Timestamp      int
}

type monattab struct { // данные по месяцу при отображении строки в monds_index.html
	Id       string
	Yahre    string
	Monat    string
	Nummonat string
	Tag      string
	Hour     string
	Kf       string
	Nalog    string
}

var mondtable struct {
	Ready    string     // флаг готовности
	Mondstab []monattab // таблица по сотрудниам  в monds_index.html
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

func MondIndexHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		files := append(partials, "./static/mond_index.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}

		mondtable.Mondstab = nil
		rows, err1 := db.Query("SELECT * FROM monds;")
		if err1 != nil {
			fmt.Println(" Monds ошибка чтения таблицы")
			panic(err)
		} else {
			fmt.Println(" Прочитали monds таблицу")
		}
		defer rows.Close()

		//var _monatrow monatrow
		rowtable := monattab{} // строка  в формате string для monds_index.html

		for rows.Next() {
			p := monattab{}
			fmt.Println(" начали распаковку строк monds таблицы")
			err2 := rows.Scan( // пересылка  данных строки таблицы monds
				&p.Id,
				&p.Monat,
				&p.Nummonat,
				&p.Tag,
				&p.Hour,
				&p.Kf,
				//&p.Blockmonat,
				&p.Nalog,
				//&p.Blockpersonal,
				//&p.Blocktimetabel,
				//&p.Blocktabel,
				//&p.Blockbuchtabel,
				//&_monatrow.Id,
				//&_monatrow.Timestamp,
			)
			fmt.Println(p)
			if err2 != nil {
				fmt.Println("monds_index ошибка распаковки строки ")
				panic(err)
				return
			} else {
				fmt.Println(" распаковали строку из monds_index")
			}
			//rowtable.Yahre 		= strconv.Itoa(p.Yahre)
			//rowtable.Monat 		= p.Monat
			//rowtable.Nummonat 	= strconv.Itoa(p.Nummonat)
			//rowtable.Tag 		= strconv.Itoa(p.Tag)
			//rowtable.Hour 		= strconv.Itoa(p.Hour)
			//rowtable.Kf 		= strconv.Itoa(p.Kf)
			////rowtable.BlockMonat 		= strconv.Itoa(p.BlockMonat)
			//rowtable.Nalog 		= strconv.Itoa(p.Nalog)
			//rowtable.BlockPersonal 	= strconv.Itoa(_monatrow.BlockPersonal)
			//rowtable.BlockTimetabel 	= strconv.Itoa(_monatrow.BlockTimetabel)
			//rowtable.BlockTabel 		= strconv.Itoa(_monatrow.BlockTabel)
			//rowtable.BlockBuchtabel 	= strconv.Itoa(_monatrow.BlockBuchtabel)
			//rowtable.Timestamp 		= strconv.Itoa(_monatrow.Timestamp)
			//rowtable.Id 				= strconv.Itoa(_monatrow.Id)
			mondtable.Mondstab = append(mondtable.Mondstab, rowtable) // добавление строки в таблицу Personalstab для personals_index.html
		}
		fmt.Println("after unpack rows in index")
		mondtable.Ready = "1"
		//monat := req.URL.Query().Get("monat")
		//del := req.URL.Query().Get("del")
		//
		//if del == "del" {
		//	_, err = db.Exec("delete from monds where monat = $2", monat)
		//	if err != nil { // удаление старой записи
		//		panic(err)
		//	}
		//}
		err = t.ExecuteTemplate(w, "base", mondtable)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func MondShowHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		var _monatHtml monatHtml
		var _monatrow monatrow

		files := append(partials, "./static/mond_show.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		monat := req.URL.Query().Get("monat")
		row := db.QueryRow("SELECT * FROM monds WHERE monat=$2", monat) // выборка строки
		//if err1 != nil  {
		//	panic(err1)
		//}
		err2 := row.Scan(&_monatrow.Yahre, // чтение переменных из полей
			&_monatrow.Monat,
			&_monatrow.Nummonat,
			&_monatrow.Tag,
			&_monatrow.Hour,
			&_monatrow.Kf,
			&_monatrow.Blockmonat,
			&_monatrow.Nalog,
			&_monatrow.Blockpersonal,
			&_monatrow.Blocktimetabel,
			&_monatrow.Blocktabel,
			&_monatrow.Blockbuchtabel,
			//&_monatrow.Id,
			&_monatrow.Timestamp,
		)
		if err2 != nil {
			fmt.Println("ошибка распаковки строки monds в show")
			panic(err)
		} else {
			_monatHtml.Yahre = strconv.Itoa(_monatrow.Yahre)
			_monatHtml.Monat = _monatrow.Monat
			_monatHtml.Nummonat = strconv.Itoa(_monatrow.Nummonat)
			_monatHtml.Tag = strconv.Itoa(_monatrow.Tag)
			_monatHtml.Hour = strconv.Itoa(_monatrow.Hour)
			_monatHtml.Kf = strconv.Itoa(_monatrow.Kf)
			_monatHtml.Blockmonat = strconv.Itoa(_monatrow.Blockmonat)
			_monatHtml.Nalog = strconv.Itoa(_monatrow.Nalog)
			_monatHtml.Blockpersonal = strconv.Itoa(_monatrow.Blockpersonal)
			_monatHtml.Blocktimetabel = strconv.Itoa(_monatrow.Blocktimetabel)
			_monatHtml.Blocktabel = strconv.Itoa(_monatrow.Blocktabel)
			_monatHtml.Blockbuchtabel = strconv.Itoa(_monatrow.Blockbuchtabel)
			_monatHtml.Timestamp = strconv.Itoa(_monatrow.Timestamp)
			_monatHtml.Id = strconv.Itoa(_monatrow.Id)
		}
		err = t.ExecuteTemplate(w, "base", _monatHtml)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}

func MondEditHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var _monatHtml monatHtml
		var _monatrow monatrow

		files := append(partials, "./static/mond_edit.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}

		monat := req.URL.Query().Get("monat")
		fmt.Println("monat = ", monat)
		row := db.QueryRow("SELECT * FROM monds WHERE monat=$2", monat) // выборка строки
		//if err1 != nil {
		//	fmt.Println("ошибка чтения из базы ")
		//	panic(err)
		//}
		err = row.Scan(&_monatrow.Yahre, // чтение переменных из полей
			&_monatrow.Monat,
			&_monatrow.Nummonat,
			&_monatrow.Tag,
			&_monatrow.Hour,
			&_monatrow.Kf,
			&_monatrow.Blockmonat,
			&_monatrow.Nalog,
			&_monatrow.Blockpersonal,
			&_monatrow.Blocktimetabel,
			&_monatrow.Blocktabel,
			&_monatrow.Blockbuchtabel,
			&_monatrow.Id,
			&_monatrow.Timestamp,
		)
		if err != nil {
			fmt.Println("ошибка распаковки строки monds в edit")
			panic(err)
		} else {
			_monatHtml.Yahre = strconv.Itoa(_monatrow.Yahre)
			_monatHtml.Monat = _monatrow.Monat
			_monatHtml.Nummonat = strconv.Itoa(_monatrow.Nummonat)
			_monatHtml.Tag = strconv.Itoa(_monatrow.Tag)
			_monatHtml.Hour = strconv.Itoa(_monatrow.Hour)
			_monatHtml.Kf = strconv.Itoa(_monatrow.Kf)
			_monatHtml.Blockmonat = strconv.Itoa(_monatrow.Blockmonat)
			_monatHtml.Nalog = strconv.Itoa(_monatrow.Nalog)
			_monatHtml.Blockpersonal = strconv.Itoa(_monatrow.Blockpersonal)
			_monatHtml.Blocktimetabel = strconv.Itoa(_monatrow.Blocktimetabel)
			_monatHtml.Blocktabel = strconv.Itoa(_monatrow.Blocktabel)
			_monatHtml.Blockbuchtabel = strconv.Itoa(_monatrow.Blockbuchtabel)
			_monatHtml.Timestamp = strconv.Itoa(_monatrow.Timestamp)
			_monatHtml.Id = strconv.Itoa(_monatrow.Id)

			if req.Method == "POST" {
				_monatHtml.Ready = "0"  // 1 - ввод успешный
				_monatHtml.Errors = "0" // 1 - ошибки при вводе
				_monatHtml.Empty = "0"  // 1 - есть пустые поля
				_monatHtml.Range = "0"  // 1 - выход за пределы диапазона
				req.ParseForm()
				_monatHtml.Yahre = req.Form["yahre"][0]
				_monatHtml.Monat = req.Form["monat"][0]
				_monatHtml.Nummonat = req.Form["num_monat"][0]
				_monatHtml.Tag = req.Form["tag"][0]
				_monatHtml.Hour = req.Form["hour"][0]
				_monatHtml.Kf = req.Form["kf_oberhour"][0]
				_monatHtml.Nalog = req.Form["nalog"][0]
				//_monatHtml.Blockmonat		= req.Form["block_monat"][0]
				//_monatHtml.Blockpersonal 	= req.Form["block_personal"][0]
				//_monatHtml.Blocktimetabel 	= req.Form["block_timetabel"][0]
				//_monatHtml.Blocktabel 		= req.Form["block_tabel"][0]
				//_monatHtml.Blockbuchtabel 	= req.Form["block_buchtabel"][0]
				err := 0
				err = err + checknum(_monatHtml.Yahre, 2040, 2021)
				err = err + checknum(_monatHtml.Monat, 12, 1)
				err = err + checknum(_monatHtml.Tag, 30, 1)
				err = err + checknum(_monatHtml.Hour, 176, 8)
				err = err + checknum(_monatHtml.Kf, 3, 1)
				err = err + checknum(_monatHtml.Nalog, 40, 13)
				if err > 0 {
					_monatHtml.Range = "1"
					_monatHtml.Errors = "1"
				}
				if _monatHtml.Yahre == "" || _monatHtml.Monat == "" || _monatHtml.Nummonat == "" || _monatHtml.Tag == "" || _monatHtml.Hour == "" || _monatHtml.Kf == "" || _monatHtml.Nalog == "" {
					_monatHtml.Empty = "1"
					_monatHtml.Errors = "1"
				}
				if _monatHtml.Errors == "0" {
					_monatHtml.Ready = "1"
				}
				//addMonat(_monatHtml) //добавление записи в базу
				// запись в базу
				_, err1 := db.Exec("delete from monds where monat = $2", monat)
				if err1 != nil { // удаление старой записи
					panic(err)
				}
				yahre, _ := strconv.Atoi(_monatHtml.Yahre)
				_monat := _monatHtml.Monat
				numMonat, _ := strconv.Atoi(_monatHtml.Nummonat)
				hour, _ := strconv.Atoi(_monatHtml.Hour)
				tag, _ := strconv.Atoi(_monatHtml.Tag)
				kfOberhour, _ := strconv.Atoi(_monatHtml.Kf)
				nalog, _ := strconv.Atoi(_monatHtml.Nalog)
				blockMonat := 0
				blockPersonal := 0
				blockTimetabel := 0
				blockTabel := 0
				blockBuchtabel := 0

				_, err2 := db.Exec("INSERT INTO monds VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
					yahre,
					_monat,
					numMonat,
					tag,
					hour,
					kfOberhour,
					blockMonat,
					nalog,
					blockPersonal,
					blockTimetabel,
					blockTabel,
					blockBuchtabel,
				)
				if err2 != nil {
					fmt.Println("Ошибка записи измененной строки в monds")
				}
			}

			err = t.ExecuteTemplate(w, "base", _monatHtml)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
				return
			}
		}
	}
}

func MondNewHandler(db *sql.DB) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var _monatHtml monatHtml

		files := append(partials, "./static/mond_new.html")
		t, err := template.ParseFiles(files...) // Parse template file.
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Mond Internal Server ParseFiles Error", http.StatusInternalServerError)
			return
		}
		if req.Method == "POST" {
			_monatHtml.Ready = "0"  // 1 - ввод успешный
			_monatHtml.Errors = "0" // 1 - ошибки при вводе
			_monatHtml.Empty = "0"  // 1 - есть пустые поля
			_monatHtml.Range = "0"  // 1 - выход за пределы диапазона
			req.ParseForm()
			_monatHtml.Yahre = req.Form["yahre"][0]
			_monatHtml.Monat = req.Form["monat"][0]
			_monatHtml.Nummonat = req.Form["nummonat"][0]
			_monatHtml.Tag = req.Form["tag"][0]
			_monatHtml.Hour = req.Form["hour"][0]
			_monatHtml.Kf = req.Form["kf"][0]
			_monatHtml.Nalog = req.Form["nalog"][0]
			//_monatHtml.Blockmonat 		= req.Form["blockmonat"][0]
			//_monatHtml.Blockpersonal 	= req.Form["blockpersonal"][0]
			//_monatHtml.Blocktimetabel 	= req.Form["blocktimetabel"][0]
			//_monatHtml.Blocktabel 		= req.Form["blocktabel"][0]
			//_monatHtml.Blockbuchtabel 	= req.Form["blockbuchtabel"][0]
			err := 0
			err = err + checknum(_monatHtml.Yahre, 2040, 2021)
			err = err + checknum(_monatHtml.Nummonat, 12, 1)
			err = err + checknum(_monatHtml.Tag, 30, 1)
			err = err + checknum(_monatHtml.Hour, 176, 8)
			err = err + checknum(_monatHtml.Kf, 3, 1)
			err = err + checknum(_monatHtml.Nalog, 40, 10)
			if err > 0 {
				_monatHtml.Range = "1"
				_monatHtml.Errors = "1"
			}
			if _monatHtml.Yahre == "" || _monatHtml.Monat == "" || _monatHtml.Tag == "" || _monatHtml.Hour == "" || _monatHtml.Kf == "" || _monatHtml.Nalog == "" {
				_monatHtml.Empty = "1"
				_monatHtml.Errors = "1"
			}
			if _monatHtml.Errors == "0" {
				_monatHtml.Ready = "1"
			}
			//addMonat(_monatHtml) //добавление записи в базу
			// запись в базу
			// запись в базу
			//monat := _monatHtml.Monat
			//_, err1 := db.Exec("delete from monds where monat = $2", monat)
			//if err1 != nil { // удаление старой записи
			//	fmt.Println("Ошибка удаления старой строки в monds")
			//	panic(err)
			//}
			//yahre, _ 		:= strconv.Atoi(_monatHtml.Yahre)
			//_monat 		:= _monatHtml.Monat
			//numMonat, _ 	:= strconv.Atoi(_monatHtml.Nummonat)
			//hour, _ 		:= strconv.Atoi(_monatHtml.Hour)
			//tag, _ 		:= strconv.Atoi(_monatHtml.Tag)
			//kfOberhour, _ := strconv.Atoi(_monatHtml.Kf)
			//nalog, _      := strconv.Atoi(_monatHtml.Nalog)
			yahre := 2021
			_monat := "zyd"
			numMonat := 1
			hour := 168
			tag := 21
			kfOberhour := 2
			nalog := 13
			blockMonat := 0
			blockPersonal := 0
			blockTimetabel := 0
			blockTabel := 0
			blockBuchtabel := 0
			fmt.Println("befor record  row in new monds")
			_, err2 := db.Exec("INSERT INTO monds VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
				yahre,
				_monat,
				numMonat,
				tag,
				hour,
				kfOberhour,
				blockMonat,
				nalog,
				blockPersonal,
				blockTimetabel,
				blockTabel,
				blockBuchtabel,
			)
			if err2 != nil {
				fmt.Println("Ошибка записи новой строки в monds")
			}
		}
		err1 := t.ExecuteTemplate(w, "base", _monatHtml)
		if err1 != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Execute Error", http.StatusInternalServerError)
			return
		}
	}
}
