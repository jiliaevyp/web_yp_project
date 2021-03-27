package servfunc

import (
	"net/mail"
	"strconv"
)

// проверка корректности емайл адреса nameAddress --> "имя <email@mail.com>
func InpMailAddress(nameAddress string) (err int, email string, title string) {
	e, err1 := mail.ParseAddress(nameAddress)
	if err1 != nil {
		return 1, e.Address, e.Name //"?", "?"
	}
	return 0, e.Address, e.Name
}

// валидация  числовых вводов и диапазонов
func Checknum(checknum string, min int, max int) int {
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
