package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/khisby/finance-log/config"
	"github.com/khisby/finance-log/core/module"
	"github.com/khisby/finance-log/core/repository"
	"github.com/khisby/finance-log/pkg/whatsapp"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

//go:embed credential.json
var credential string

func main() {
	cfg := config.Get()

	creds := fmt.Sprint(credential)
	sheetClient, err := sheets.NewService(context.Background(), option.WithCredentialsJSON([]byte(creds)))
	if err != nil {
		log.Fatal(err)
	}

	gsheetUsecase := repository.NewGSheetRepository(sheetClient, cfg.GSheetID)
	client := whatsapp.InitWhatsapp()
	chatUsecase := module.NewChatUsecase(gsheetUsecase, client, cfg.WhitelistUser)
	whatsapp.EventHandler(client, chatUsecase.HandlerEvent)
	whatsapp.StartWhatsapp(client)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	whatsapp.StopWhatsapp(client)

	fmt.Printf("khisby debug = %s\n", "DONE")
}
