package entity

import "time"

type FinanceLog struct {
	WhatsappNumber string
	Time           time.Time
	Category       Category
	Amount         string
	Status         string
	Description    string
}

func (f *FinanceLog) FillTime() {
	f.Time = time.Now()
}
