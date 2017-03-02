package main

import (
	"flag"
	"log"

	"os"
	"os/signal"

	. "github.com/auvn/go.cache/commands"
	"github.com/auvn/go.cache/journal"
	"github.com/auvn/go.cache/server"
	"github.com/auvn/go.cache/session"
	"github.com/auvn/go.cache/storage"
	"github.com/auvn/go.cache/util/sync"
)

type Options struct {
	JournalFile string

	server.TelnetOptions
	server.HttpOptions

	Pass string
}

var (
	opts = &Options{}
)

func parseFlags() {
	flag.StringVar(&opts.JournalFile, "journal", "", "Journal file for cache")

	flag.StringVar(&opts.TelnetOptions.Addr, "telnet", "0.0.0.0:1234", "Addr to listen telnet on")
	flag.StringVar(&opts.HttpOptions.Addr, "http", "", "Addr to listen http on")

	flag.StringVar(&opts.Pass, "pass", "", "Password for cache auth")

	flag.Parse()
}

func initJournal(handler *Handler, s session.Session, group sync.ServeGroup) error {
	if opts.JournalFile == "" {
		return nil
	} else {
		log.Println("initializing journal file:", opts.JournalFile)
		journal := journal.MustInitFile(opts.JournalFile)
		journalAdapter := NewJournalAdapter(journal, s)
		if err := journalAdapter.Restore(handler); err != nil {
			return err
		}
		journalAdapter.AttachTo(handler)
		group.Serve(journalAdapter)
		return nil
	}
}

func initRegistry() Registry {
	options := new(RegistryOptions)
	options.Auth = opts.Pass
	return InitReflectRegistry(options)
}

func initHttp(handler *Handler, s session.Session, group sync.ServeGroup) {
	if opts.HttpOptions.Addr != "" {
		log.Println("serving http at:", opts.HttpOptions.Addr)
		group.Serve(server.Http(handler, s, &opts.HttpOptions))
	}
}

func initTelnet(handler *Handler, s session.Session, group sync.ServeGroup) {
	log.Println("serving telnet at:", opts.TelnetOptions.Addr)
	group.Serve(server.Telnet(handler, s, &opts.TelnetOptions))
}

func fatal(err error, group sync.ServeGroup) {
	group.Shutdown()
	log.Fatal("fail: ", err)
}

func initServeGroup() sync.ServeGroup {
	group := sync.NewServeGroup()
	go func() {
		err := <-group.Err()
		if err != nil {
			fatal(err, group)
		}
	}()
	return group
}

func waitForInterrupt(group sync.ServeGroup) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down")
	group.Shutdown()
}

func main() {
	parseFlags()
	group := initServeGroup()
	baseSession := session.WithStorage(session.New(), storage.New())

	registry := initRegistry()
	handler := NewHandler(registry)

	if err := initJournal(handler, baseSession, group); err != nil {
		fatal(err, group)
	}
	group.Serve(handler)

	initTelnet(handler, baseSession, group)
	initHttp(handler, baseSession, group)

	waitForInterrupt(group)
}
