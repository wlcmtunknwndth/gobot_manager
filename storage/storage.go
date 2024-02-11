package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/wlcmtunknwndth/gobot_manager/lib/error_handler"
)

const (
	LinksTable = "links"
	PagesTable = "pages"
)

type Storage interface {
	//Read-adviser
	Save(ctx context.Context, table string, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, table string, id int) error

	IsExists(ctx context.Context, table string, p *Page) (bool, error)

	SendLinks(ctx context.Context, table, key string) (*[]Page, error)

	CheckTable(table string) string
}

var (
	ErrNoSavedPages = errors.New("no saved pages")
	ErrNoSavedLinks = errors.New("no saved links")
)

type Page struct {
	ID   uint32
	URL  string
	Name string
	// Created  time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", error_handler.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.Name); err != nil {
		return "", error_handler.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
