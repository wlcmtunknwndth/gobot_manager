package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/wlcmtunknwndth/gobot_manager/lib/error_handler"
	"github.com/wlcmtunknwndth/gobot_manager/storage"
)

const (
	// usr
	HelpCmd  = "/help"
	StartCmd = "/start"

	RndCmd    = "/rnd"
	RemoveCmd = "/remove"
	sendPages = "/pages"

	LinksCmd = "/links"
	//admin
	RmLink   = "/removeLink"
	saveLink = "/saveLinks"
	// SaveCmd = "/saveLink"
)

// doCmd —— checks if there a command sent by usr.
func (p *Processor) doCmd(ctx context.Context, text string, chatID int) error {
	// text = strings.TrimSpace(text)
	commands := strings.Split(text, " ")
	if len(commands) == 0 {
		return nil
	}

	log.Printf("got new command '%s' from '%d'", text, chatID)

	chatid := strconv.Itoa(chatID)

	if isAddCmd(text) { //is URL
		return p.savePage(ctx, chatID, storage.PagesTable, text, chatid)
	}

	switch commands[0] {
	case RndCmd:
		return p.sendRandom(ctx, chatID, chatid)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case RemoveCmd:
		if len(commands) != 2 {
			return p.tg.SendMessage(chatID, msgNotEnoughArg)
		}
		return p.Remove(ctx, storage.PagesTable, commands[1])
	case sendPages:
		return p.sendLinks(ctx, chatID, storage.PagesTable)

	case saveLink:
		if len(commands) != 3 {
			return p.tg.SendMessage(chatID, msgNotEnoughArg)
		}
		return p.savePage(ctx, chatID, storage.LinksTable, commands[1], commands[2])
	case LinksCmd:
		return p.sendLinks(ctx, chatID, storage.LinksTable)
	case RmLink:
		if len(commands) != 2 {
			return p.tg.SendMessage(chatID, msgNotEnoughArg)
		}
		return p.Remove(ctx, storage.LinksTable, commands[1])

	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) Remove(ctx context.Context, table, id string) error {
	if table = p.storage.CheckTable(table); table == "" {
		return fmt.Errorf("can't remove id: %s", id)
	}
	numID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return p.storage.Remove(ctx, table, numID)
}

func (p *Processor) isAdmin(chatID int) bool {
	strAdminID, _ := os.LookupEnv("ChatID_admin")
	adminID, _ := strconv.Atoi(strAdminID)
	if chatID == adminID {
		return true
	} else {
		return false
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, table, pageURL, name string) (err error) {
	defer func() { err = error_handler.WrapIfErr("can't do command: save page", err) }()

	if table == storage.LinksTable && !p.isAdmin(chatID) {
		return nil
	}
	page := &storage.Page{
		URL:  pageURL,
		Name: name,
	}

	isExists, err := p.storage.IsExists(ctx, table, page)
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

	if err := p.storage.Save(ctx, table, page); err != nil {
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

func (p *Processor) sendRandom(ctx context.Context, chatID int, name string) (err error) {
	defer func() { err = error_handler.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(ctx, name)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

// sendHello —— Sends hello to chat
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func (p *Processor) sendLinks(ctx context.Context, chatID int, table string) (err error) {
	defer func() { err = error_handler.WrapIfErr("can't do command: can't send links", err) }()

	var links *[]storage.Page

	if p.storage.CheckTable(table) == storage.LinksTable {
		links, err = p.storage.SendLinks(ctx, table, "")
	} else if p.storage.CheckTable(table) == storage.PagesTable {
		links, err = p.storage.SendLinks(ctx, table, strconv.Itoa(chatID))
	} else {
		return fmt.Errorf("can't check a table: %s", err)
	}

	if err != nil && !errors.Is(err, storage.ErrNoSavedLinks) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedLinks) {
		return p.tg.SendMessage(chatID, msgNoSavedLinks)
	}

	for _, link := range *links {
		if err := p.tg.SendMessage(chatID, fmt.Sprintf("%d. %s —— %s", link.ID, link.Name, link.URL)); err != nil {
			return nil
		}
	}
	return nil
}

// isAddCmd —— checks if a command is for adding the page
func isAddCmd(text string) bool {
	return isURL(text)
}

// isURL —— checks if the string is URL
func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
