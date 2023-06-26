package entity

import (
	"fmt"
	"time"
)

type FinanceLog struct {
	WhatsappNumber string
	Time           time.Time
	Category       Category
	Amount         string
	Status         string
	Description    string
}

func (f *FinanceLog) FillTime() {
	location, err := time.LoadLocation("Etc/GMT-7")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	f.Time = time.Now().In(location)
}
