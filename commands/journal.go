package commands

import (
	"fmt"
	"log"

	"github.com/auvn/go.cache/journal"
	"github.com/auvn/go.cache/session"
	"github.com/auvn/go.cache/util/sync"
)

var (
	_ (sync.Server) = (*JournalAdapter)(nil)

	NonJournalableFlag = (TDFlag | RFlag)
)

type journalCmd struct {
	flag int
	body [][]byte
}

type JournalAdapter struct {
	cmds    chan *journalCmd
	journal journal.Journal
	session session.Session
}

func (self *JournalAdapter) successCommand(flag int, body [][]byte) {
	go func() { self.cmds <- &journalCmd{flag: flag, body: body} }()
}

func (self *JournalAdapter) loopCommands(quit sync.Quit) {
	for {
		select {
		case <-quit:
			return
		case cmd := <-self.cmds:
			if !CheckFlag(cmd.flag, NonJournalableFlag) {
				if err := self.journal.Write(cmd.body); err != nil {
					log.Println("cannot write to journal:", err)
				}
				if err := self.journal.Commit(); err != nil {
					log.Println("cannot commit the journal:", err)
				}
			}
		}
	}
}

func (self *JournalAdapter) Restore(h *Handler) error {
	resp := make(chan interface{}, 1)
	for {
		entry, err := self.journal.NextEntry()
		if err != nil {
			if err == journal.ErrEmpty {
				return nil
			} else {
				return err
			}
		} else {
			h.Handle(self.session, entry, resp)
			select {
			case ret := <-resp:
				if err, ok := ret.(error); ok {
					return fmt.Errorf("cannot restore command: %s", err.Error())
				}
			}
		}
	}
}

func (self *JournalAdapter) Serve(quit sync.Quit) error {
	self.loopCommands(quit)
	return nil
}

func (self *JournalAdapter) AttachTo(h *Handler) {
	h.AddSuccessHook(self.successCommand)
}

func NewJournalAdapter(journal journal.Journal, s session.Session) *JournalAdapter {
	return &JournalAdapter{
		cmds:    make(chan *journalCmd, 1000),
		journal: journal,
		session: s,
	}
}
