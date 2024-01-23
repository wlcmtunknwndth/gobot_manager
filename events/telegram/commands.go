package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/wlcmtunknwndth/gobot_manager/lib/error_handler"
	"github.com/wlcmtunknwndth/gobot_manager/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) { //is URL
		return p.savePage(ctx, chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
	// add page: http://...
	// rnd page: /rnd
	// help: /help
	// start: /start hi + help
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL, username string) (err error) {
	defer func() { err = error_handler.WrapIfErr("can't do command: save page", err) }()

	// send := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}

	// closure
	// if isExists{
	// 	return send(msgAlreadyExists)
	// }

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

// closure
// func NewMessageSender(chatID int, tg *telegram.Client) func(string) error{
// 	return func(msg string) error{
// 		return tg.SendMessage(chatID, msg)
// 	}
// }

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = error_handler.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)

}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
