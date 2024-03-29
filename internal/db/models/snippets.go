package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Author  string
	Work    string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, author string, work string, content string, expires int) (id int, err error) {
	id = -1
	qry := `INSERT INTO snippets (title, author, work, content, created, expires) VALUES (
			 $1,
			 $2,
			 $3,
			 $4,
			 CURRENT_TIMESTAMP,
			 CURRENT_TIMESTAMP + $5 * INTERVAL '1 day'
			 ) RETURNING id;`

	row := m.DB.QueryRow(qry, title, author, work, content, expires)
	err = row.Scan(&id)
	if err != nil {
		return
	}

	return
}

func (m *SnippetModel) Get(id int) (s *Snippet, err error) {
	// query db
	qry := `SELECT *
	FROM snippets
	WHERE expires > CURRENT_TIMESTAMP AND id = $1; `
	row := m.DB.QueryRow(qry, id)

	// unmarshal
	s = &Snippet{}
	err = row.Scan(&s.ID, &s.Title, &s.Author, &s.Work, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrNoRecord
		}
		return nil, err
	}
	return
}

func (m *SnippetModel) List(limit int) (snippets []*Snippet, err error) {
	// query db
	qry := `SELECT *
	FROM snippets
	WHERE expires > CURRENT_TIMESTAMP
	ORDER BY created DESC
	LIMIT $1`
	rows, err := m.DB.Query(qry, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// unmarshal query results
	snippets = []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Author, &s.Work, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return
}
