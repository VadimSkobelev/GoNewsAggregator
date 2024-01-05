// Пакет для работы с БД PostgreSQL.
package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Storage struct {
	db *pgxpool.Pool
}

// Публикация, получаемая из RSS.
type Post struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

// Подключения к БД.
func New() (*Storage, error) {
	constr := "postgres://postgres:postgres@localhost:5432/news"
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

// News возвращает n последних новостей из БД.
func (s *Storage) News(n int) ([]Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			title,
			content,
			pub_time,
			link
		FROM news
		ORDER BY pub_time DESC
		LIMIT $1;
	`, n,
	)
	if err != nil {
		return nil, err
	}
	var news []Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		news = append(news, p)
	}
	return news, rows.Err()
}

// AddNews добовляет статьи в базу.
func (s *Storage) AddNews(p []Post) error {
	for _, post := range p {
		err := s.db.QueryRow(context.Background(), `
		INSERT INTO news (title, content, pub_time, link)
		VALUES ($1, $2, $3, $4) RETURNING id;
		`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		).Scan(&post.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
