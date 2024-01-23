package telegram

import (
	"context"
	"errors"

	"github.com/wlcmtunknwndth/gobot_manager/clients/telegram"
	"github.com/wlcmtunknwndth/gobot_manager/events"
	"github.com/wlcmtunknwndth/gobot_manager/lib/error_handler"
	"github.com/wlcmtunknwndth/gobot_manager/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEvent    = errors.New("unknown event type")
	ErrUnknownMetaType = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage, offset int) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
		offset:  offset,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, error_handler.Wrap("can't get events:", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	default:
		return error_handler.Wrap("can't process message", ErrUnknownEvent)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return error_handler.Wrap("can't procces message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return error_handler.Wrap("can't process message", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)

	if !ok {
		return Meta{}, error_handler.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}
