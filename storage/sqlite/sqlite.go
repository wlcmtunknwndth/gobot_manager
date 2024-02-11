package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wlcmtunknwndth/gobot_manager/storage"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) CheckTable(table string) string {
	switch table {
	case storage.LinksTable:
		return storage.LinksTable
	case storage.PagesTable:
		return storage.PagesTable
	default:
		return ""
	}
}

// New creates a new storage
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}

	return &Storage{db: db}, nil
}

// Save saves the page to the storage
func (s *Storage) Save(ctx context.Context, table string, p *storage.Page) error {
	// var q string
	// if table == storage.LinksTable {
	// 	q = `INSERT INTO links(url, name) VALUES (?, ?)`
	// } else if table == storage.PagesTable {
	// 	q = `INSERT INTO pages(url, name) VALUES (?, ?)`
	// 	//   INSERT INTO pages(url, name) VALUE (?, ?)
	// } else {
	// 	return fmt.Errorf("no such table")
	// }
	q := fmt.Sprintf("INSERT INTO %s(url, name) VALUES (?, ?)", s.CheckTable(table))

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Name); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// PickRandom picks a random page from storage
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE name = ? ORDER BY RANDOM() LIMIT 1`

	var url string
	err := s.db.QueryRowContext(ctx, q, userName).Scan(&url)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}

	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:  url,
		Name: userName,
	}, nil
}

// Remove removes page from the storage.
func (s *Storage) Remove(ctx context.Context, table string, id int) error {

	q := fmt.Sprintf(`DELETE FROM %s WHERE id = ?`, s.CheckTable(table))

	if _, err := s.db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("can't remove page: %w", err)
	}

	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, table string, page *storage.Page) (bool, error) {
	// first "?" is to choose from links or pages
	q := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE url = ? AND name = ?`, s.CheckTable(table))

	var count int
	// table_name
	if err := s.db.QueryRowContext(ctx, q, page.URL, page.Name).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check if page exists: %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, name TEXT)`
	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table pages: %w", err)
	}

	q = `CREATE TABLE IF NOT EXISTS links (id INTEGER PRIMARY KEY AUTOINCREMENT, url TEXT, name TEXT)`
	_, err = s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table links: %w", err)
	}
	return nil
}

// func (s *Storage) SaveLink(ctx context.Context, link *storage.Page) error {
// 	q := `INSERT INTO links(name, value) VALUE (?, ?)`

// 	if _, err := s.db.ExecContext(ctx, q, link.Name, link.URL); err != nil {
// 		return fmt.Errorf("can't save link: %w", err)
// 	}

// 	return nil
// }

func (s *Storage) SendLinks(ctx context.Context, table, key string) (*[]storage.Page, error) {
	// q := fmt.Sprintf(`SELECT * FROM %s WHERE name = ?`, s.CheckTable(table))
	table = s.CheckTable(table)
	var q string
	if table == storage.LinksTable {
		q = fmt.Sprintf(`SELECT * FROM %s`, table)
	} else if table == storage.PagesTable {
		q = fmt.Sprintf(`SELECT * FROM %s WHERE name = ?`, table)
	} else {
		return nil, fmt.Errorf("wrong table: %s", table)
	}

	rows, err := s.db.QueryContext(ctx, q, key)
	if err != nil {
		return nil, err
	}

	var links []storage.Page = make([]storage.Page, 0)

	for rows.Next() {
		link := storage.Page{}
		err = rows.Scan(&link.ID, &link.URL, &link.Name)
		if err == sql.ErrNoRows {
			return nil, storage.ErrNoSavedLinks
		}

		if err != nil {
			return nil, fmt.Errorf("can't get a page: %w", err)
		}
		links = append(links, link)
	}

	return &links, nil
}

// func (s *Storage) sendGit(ctx context.Context) error{

// }
