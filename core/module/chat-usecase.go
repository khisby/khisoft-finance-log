package module

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/khisby/finance-log/core/entity"
	"github.com/khisby/finance-log/core/repository"
	"github.com/khisby/finance-log/utils"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type ChatUsecase struct {
	gsheetUsecase *repository.GSheetRepository
	waCli         *whatsmeow.Client
	whitelistUser map[string]bool
	financeLog    entity.FinanceLog
	mapCategory   map[string]entity.Category
}

func NewChatUsecase(gsheetUsecase *repository.GSheetRepository, waCli *whatsmeow.Client, whitelistUser []string) *ChatUsecase {
	whitelistUserMap := map[string]bool{}
	for _, d := range whitelistUser {
		whitelistUserMap[d] = true
	}

	var mapCategory = map[string]entity.Category{}
	mapCategory["makan"] = entity.Makan
	mapCategory["makanan"] = entity.Makan
	mapCategory["minum"] = entity.Makan
	mapCategory["minuman"] = entity.Makan
	mapCategory["minuman"] = entity.Makan
	mapCategory["jajan"] = entity.Jajan
	mapCategory["jalan"] = entity.Jalan
	mapCategory["jalan2"] = entity.Jalan
	mapCategory["jalan-jalan"] = entity.Jalan
	mapCategory["jalan²"] = entity.Jalan
	mapCategory["transport"] = entity.Transport
	mapCategory["transportasi"] = entity.Transport
	mapCategory["bensin"] = entity.Transport
	mapCategory["parkir"] = entity.Transport
	mapCategory["bulanan"] = entity.Bulanan
	mapCategory["belanja"] = entity.Belanja
	mapCategory["belanjaan"] = entity.Belanja
	mapCategory["project"] = entity.Project
	mapCategory["kerja"] = entity.Kerja
	mapCategory["gaji"] = entity.Kerja
	mapCategory["utang"] = entity.Hutang
	mapCategory["hutang"] = entity.Hutang
	mapCategory["pinjam"] = entity.Hutang
	mapCategory["pinjem"] = entity.Hutang
	mapCategory["sedekah"] = entity.Sedekah
	mapCategory["ngasih"] = entity.Sedekah
	mapCategory["infaq"] = entity.Sedekah

	return &ChatUsecase{gsheetUsecase, waCli, whitelistUserMap, entity.FinanceLog{}, mapCategory}
}

func (c *ChatUsecase) HandlerEvent(event interface{}) {
	switch v := event.(type) {
	case *events.Message:
		c.FinanceLog(v.Info.Sender.User, v.Message)
	}
}

func (c *ChatUsecase) FinanceLog(sender string, message *waProto.Message) {
	if !c.whitelistUser[sender] {
		return
	}

	messageText := message.GetConversation()
	if messageText == "" {
		messageText = message.GetExtendedTextMessage().GetText()
	}
	fmt.Printf(time.Now().Format("2006-01-02 15:04:05")+" message from: %s with content: %s\n", sender, messageText)
	checkSheetIfExist := c.gsheetUsecase.CheckSheetIfExist(sender)
	if !checkSheetIfExist {
		err := c.gsheetUsecase.CreateSheet(sender)
		if err != nil {
			fmt.Printf("Error creating sheet: %s", err)
			return
		}

		row := [][]string{}
		columnName := []string{"Time", "Category", "Amount", "Status", "Description"}
		row = append(row, columnName)
		err = c.gsheetUsecase.UpdateSheetData(sender, row)
		if err != nil {
			fmt.Printf("Error updating sheet: %s", err)
			return
		}

		go func() {
			err = c.SendMessage(sender, fmt.Sprintf("%s %s", entity.NewCommers, entity.MenuText))
			if err != nil {
				fmt.Printf("Error sending message: %s", err)
				return
			}
		}()
	}

	splitMessageText := strings.Split(strings.ToLower(messageText), " ")
	splitMessageTextTemplate2 := strings.Split(strings.ToLower(messageText), "\n")
	if len(splitMessageText) > 0 {
		if splitMessageText[0] == "masuk" || splitMessageText[0] == "keluar" || splitMessageTextTemplate2[0] == "debit" || splitMessageTextTemplate2[0] == "kredit" {
			c.Catat(sender, messageText)
		} else if splitMessageText[0] == "bantuin" || splitMessageText[0] == "bantuan" {
			err := c.SendMessage(sender, entity.MenuText)
			if err != nil {
				fmt.Printf("Error sending message: %s", err)
				return
			}
		} else if splitMessageText[0] == "report" || splitMessageText[0] == "laporan" {
			c.Report(sender, messageText)
		} else if splitMessageText[0] == "link" {
			link := c.gsheetUsecase.GetSheetLink()
			err := c.SendMessage(sender, fmt.Sprintf("Link spreadsheet kamu: %s", link))
			if err != nil {
				fmt.Printf("Error sending message: %s", err)
				return
			}
		} else if splitMessageText[0] == "hapus" {
			c.Delete(sender, messageText)
		}
	}

}

func (c *ChatUsecase) Delete(sender, message string) {
	splitMessageText := strings.Split(strings.ToLower(message), " ")
	if len(splitMessageText) > 1 {
		posisiHapus := 0
		teks := ""
		if splitMessageText[1] == "barusan" || splitMessageText[1] == "terakhir" {
			posisiHapus = 1
			teks = "terakhir"
		} else {
			err := c.SendMessage(sender, "Mau hapus apa kak? barusan/terakhir/-1/-2/-3/-4/-5")
			if err != nil {
				fmt.Printf("Error sending message: %s", err)
				return
			}
			return
		}

		rows, err := c.gsheetUsecase.GetSheetData(sender)
		if err != nil {
			fmt.Printf("Error getting sheet data: %s", err)
			return
		}

		if len(rows) <= posisiHapus || len(rows) == 1 {
			err := c.SendMessage(sender, "Datanya udh ga ada kak")
			if err != nil {
				fmt.Printf("Error sending message: %s", err)
				return
			}
			return
		}

		lastRow := rows[len(rows)-posisiHapus]
		rows = rows[:len(rows)-posisiHapus]

		go func() {
			err = c.gsheetUsecase.UpdateSheetData(sender, rows)
			if err != nil {
				fmt.Printf("Error updating sheet: %s", err)
				return
			}
		}()

		go func() {
			amount, err := strconv.ParseInt(lastRow[2], 10, 64)
			if err != nil {
				fmt.Printf("Error parsing string to int64: %s", err)
				return
			}

			err = c.SendMessage(sender, fmt.Sprintf(entity.DeletedText, teks, lastRow[3], utils.FormatRupiah(amount), lastRow[1], lastRow[4]))
			if err != nil {
				fmt.Printf("Error sending message: %s", err)
				return
			}
		}()
	}
}

func (c *ChatUsecase) Report(sender, message string) {
	splitMessageText := strings.Split(strings.ToLower(message), " ")
	jenisReport := ""
	waktuReport := ""

	if len(splitMessageText) > 1 {
		if splitMessageText[1] == "mingguan" || splitMessageText[1] == "minggu" {
			jenisReport = "minggu"
		} else if splitMessageText[1] == "bulanan" || splitMessageText[1] == "bulan" {
			jenisReport = "bulan"
		} else {
			jenisReport = "hari"
		}
	}

	if len(splitMessageText) > 2 {
		if splitMessageText[2] == "lalu" || splitMessageText[2] == "kemarin" {
			waktuReport = "kemarin"
		} else {
			waktuReport = "ini"
		}
	}

	rows, err := c.gsheetUsecase.GetSheetData(sender)
	if err != nil {
		fmt.Printf("Error getting sheet data: %s", err)
		return
	}

	if len(rows) == 1 {
		err := c.SendMessage(sender, entity.ReportTextNotFound)
		if err != nil {
			fmt.Printf("Error sending message: %s", err)
			return
		}
		return
	}

	pemasukan, pengeluaran, totalPemasukanCategory, totalPengeluaranCategory, err := c.countReport(rows, jenisReport, waktuReport)
	if err != nil {
		fmt.Printf("Error counting report: %s", err)
		return
	}

	text := fmt.Sprintf(entity.ReportTextHeader, jenisReport, waktuReport)
	text += fmt.Sprintf(entity.ReportTextPemasukan, utils.FormatRupiah(pemasukan))
	text += fmt.Sprintf(entity.ReportTextPengeluaran, utils.FormatRupiah(pengeluaran))
	text += fmt.Sprintf(entity.ReportTextCategoryHeader, "Pemasukan")
	for k, v := range totalPemasukanCategory {
		text += fmt.Sprintf(entity.ReportTextCategory, k, utils.FormatRupiah(v))
	}
	text += fmt.Sprintf(entity.ReportTextCategoryHeader, "Pengeluaran")
	for k, v := range totalPengeluaranCategory {
		text += fmt.Sprintf(entity.ReportTextCategory, k, utils.FormatRupiah(v))
	}

	err = c.SendMessage(sender, text)
	if err != nil {
		fmt.Printf("Error sending message: %s", err)
		return
	}

}

func (c *ChatUsecase) countReport(rows [][]string, jenisReport, waktuReport string) (int64, int64, map[string]int64, map[string]int64, error) {
	var totalPemasukan int64
	var totalPengeluaran int64
	totalPemasukanCategory := map[string]int64{}
	totalPengeluaranCategory := map[string]int64{}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}

	now := time.Now().In(location)

	weekStart := now
	for weekStart.Weekday() != time.Monday {
		weekStart = weekStart.AddDate(0, 0, -1)
	}
	weekStart = time.Date(weekStart.Year(), weekStart.Month(), weekStart.Day(), 0, 0, 0, 0, weekStart.Location())

	weekEnd := now
	for weekEnd.Weekday() != time.Sunday {
		weekEnd = weekEnd.AddDate(0, 0, 1)
	}
	weekEnd = time.Date(weekEnd.Year(), weekEnd.Month(), weekEnd.Day(), 23, 59, 59, 0, weekEnd.Location())

	for _, row := range rows {
		if len(row) == 0 || row[0] == "Time" {
			continue
		}

		timeRow, err := time.ParseInLocation("02-01-2006 15:04:05", row[0], location)
		if err != nil {
			fmt.Printf("Error parsing time: %s", err)
			return 0, 0, nil, nil, err
		}
		if jenisReport == "hari" && waktuReport == "ini" {
			if timeRow.Day() != now.Day() {
				continue
			}
			if timeRow.Month() != now.Month() {
				continue
			}
			if timeRow.Year() != now.Year() {
				continue
			}
		} else if jenisReport == "bulan" && waktuReport == "ini" {
			if timeRow.Month() != now.Month() {
				continue
			}
			if timeRow.Year() != now.Year() {
				continue
			}
		} else if jenisReport == "hari" && waktuReport == "kemarin" {
			if timeRow.Day() != now.AddDate(0, 0, -1).Day() && timeRow.Year() == now.Year() {
				continue
			}
			if timeRow.Month() != now.Month() {
				continue
			}
			if timeRow.Year() != now.Year() {
				continue
			}
		} else if jenisReport == "bulan" && waktuReport == "kemarin" {
			if timeRow.Month() != now.AddDate(0, -1, 0).Month() {
				continue
			}

			if now.Month() == time.January && timeRow.Year() != now.AddDate(-1, 0, 0).Year() {
				continue
			}

			if now.Month() != time.January && timeRow.Year() != now.Year() {
				continue
			}
		} else if jenisReport == "minggu" && waktuReport == "ini" {
			if !(timeRow.After(weekStart) && timeRow.Before(weekEnd)) {
				continue
			}
		} else if jenisReport == "minggu" && waktuReport == "kemarin" {
			if !(timeRow.After(weekStart.AddDate(0, 0, -7)) && timeRow.Before(weekEnd.AddDate(0, 0, -7))) {
				continue
			}
		}

		amount, err := strconv.ParseInt(row[2], 10, 64)
		if err != nil {
			fmt.Printf("Error parsing amount: %s", err)
			return 0, 0, nil, nil, err
		}

		if row[3] == "Debit" {
			totalPemasukan += amount
			totalPemasukanCategory[row[1]] += amount
		} else if row[3] == "Kredit" {
			totalPengeluaran += amount
			totalPengeluaranCategory[row[1]] += amount
		}
	}

	return totalPemasukan, totalPengeluaran, totalPemasukanCategory, totalPengeluaranCategory, nil
}

func (c *ChatUsecase) Catat(sender, message string) {
	financeLog, err := c.parseMessageToFinanceLog(message, sender)
	if err != nil {
		fmt.Printf("Error parsing message: %s", err)
		err = c.SendMessage(sender, err.Error())
		if err != nil {
			fmt.Printf("Error sending message: %s", err)
		}
		return
	}

	go func() {
		// string parse to int64
		amount, err := strconv.ParseInt(financeLog.Amount, 10, 64)
		if err != nil {
			fmt.Printf("Error parsing amount: %s", err)
			return
		}

		err = c.SendMessage(sender, fmt.Sprintf(entity.ReplyChatSaved, financeLog.Status, utils.FormatRupiah(amount), financeLog.Category, utils.CapitalizeFirstChar(financeLog.Description)))
		if err != nil {
			fmt.Printf("Error sending message: %s", err)
		}
	}()

	rows, err := c.gsheetUsecase.GetSheetData(sender)
	if err != nil {
		fmt.Printf("Error getting sheet data: %s", err)
		return
	}

	row := []string{
		financeLog.Time.Format("02-01-2006 15:04:05"),
		string(financeLog.Category),
		financeLog.Amount,
		financeLog.Status,
		financeLog.Description,
	}

	rows = append(rows, row)

	err = c.gsheetUsecase.UpdateSheetData(sender, rows)
	if err != nil {
		fmt.Printf("Error updating sheet data: %s", err)
		return
	}
}

func (c *ChatUsecase) parseMessageToFinanceLog(message, sender string) (entity.FinanceLog, error) {
	exampleFormat := []string{"keluar 2000 buat jajan pentol\natau gini\nkredit\n2000\njajan\npentol", "keluar 300k buat bulanan bayar indihome\natau gini\nkredit\n2000\njajan\npentol", "masuk 500ribu dari project Website Landing Page\natau gini\nkredit\n2000\njajan\npentol", "masuk 1jt dari kerja bulan januari 2023\natau gini\nkredit\n2000\njajan\npentol"}
	randomInt := rand.Intn(len(exampleFormat))
	example := exampleFormat[randomInt]

	financeLog := entity.FinanceLog{}
	financeLog.FillTime()
	financeLog.WhatsappNumber = sender

	splitString := strings.Split(strings.ToLower(message), "\n")
	status := ""
	jumlah := ""
	kategori := ""
	deskripsi := ""
	if len(splitString) == 4 {
		if len(splitString) < 3 {
			return financeLog, errors.New("Maaf, aku ga ngerti maksud pesanmu. Sepertinya templatemu salah \n\nContoh: " + example)
		}
		status = strings.TrimSpace(splitString[0])
		jumlah = strings.TrimSpace(splitString[1])
		kategori = strings.TrimSpace(splitString[2])
		deskripsi = strings.Join(splitString[3:], " ")
	} else {
		splitString := strings.Split(strings.ToLower(message), " ")
		status = strings.TrimSpace(splitString[0])
		jumlah = strings.TrimSpace(splitString[1])
		if splitString[2] != "untuk" && splitString[2] != "dari" && splitString[2] != "buat" && splitString[2] != "ke" {
			kategori = strings.TrimSpace(splitString[2])
			deskripsi = strings.Join(splitString[3:], " ")
		} else {
			kategori = strings.TrimSpace(splitString[3])
			deskripsi = strings.Join(splitString[4:], " ")
		}

	}

	if status == "keluar" || status == "kredit" {
		financeLog.Status = "Kredit"
	} else if status == "masuk" || status == "debit" {
		financeLog.Status = "Debit"
	} else {
		return financeLog, errors.New("Maaf, kamu maunya di catat sebagai keluar/kredit apa masuk/debit? \n\nContoh: " + example)
	}

	amount := strings.Replace(jumlah, ".", "", -1)
	if strings.Contains(amount, "k") || strings.Contains(amount, "ribu") {
		amount = strings.Replace(amount, "k", "", -1)
		amount = strings.Replace(amount, "ribu", "", -1)
		amount = amount + "000"
	} else if strings.Contains(amount, "jt") || strings.Contains(amount, "juta") {
		amount = strings.Replace(amount, "jt", "", -1)
		amount = strings.Replace(amount, "juta", "", -1)
		amount = amount + "000000"
	}

	if _, err := strconv.Atoi(amount); err != nil {
		return financeLog, errors.New("Maaf, duit harus angka kak, ga boleh ada huruf kecuali seperti contoh yak. \n\nContoh: " + example)
	}

	financeLog.Amount = amount

	category := strings.ToLower(kategori)
	if _, ok := c.mapCategory[category]; ok {
		financeLog.Category = c.mapCategory[category]
	} else {
		financeLog.Category = entity.Lainnya
	}

	financeLog.Description = deskripsi

	return financeLog, nil
}

func (c *ChatUsecase) SendMessage(to string, message string) error {
	fmt.Printf(time.Now().Format("2006-01-02 15:04:05")+" Sending message to: %s with content: %s\n", to, message)

	newJid := types.NewJID(to, "s.whatsapp.net")
	newMessage := &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(message),
		},
	}

	_, err := c.waCli.SendMessage(context.Background(), newJid, newMessage)
	if err != nil {
		fmt.Printf("Error sending message: %s", err)
		return err
	}

	return nil
}
