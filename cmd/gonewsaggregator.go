package main

import (
	"GoNewsAggregator/pkg/api"
	"GoNewsAggregator/pkg/rss"
	"GoNewsAggregator/pkg/storage"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

// Структура конфигурационного файла.
type config struct {
	UrlList []string `json:"rss"`
	Period  int      `json:"request_period"`
}

func main() {

	// Создаём канал для публикаций.
	chPosts := make(chan []storage.Post)

	// Создаём канал для агрегации ошибок.
	chErrs := make(chan error)

	// Обработка потока ошибок.
	go func() {
		for err := range chErrs {
			log.Println("ошибка:", err)
		}
	}()

	// Читаем файл конфигурации со списком RSS URLs и периодом опроса.
	fileIn, err := os.ReadFile("./config.json")
	if err != nil {
		chErrs <- err
	}

	// Формируем структуру конфигурации из считанного файла конфигурации.
	var config config
	err = json.Unmarshal(fileIn, &config)
	if err != nil {
		chErrs <- err
	}

	// Реляционная БД PostgreSQL.
	db, err := storage.New()
	if err != nil {
		chErrs <- err
	}

	api := api.New(db)

	// Проходим по списку RSS ссылок.
	// Для каждого RSS-канала запускается своя горутина.
	for _, url := range config.UrlList {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	// запись потока новостей в БД
	go func() {
		for posts := range chPosts {
			err := db.AddNews(posts)
			// Исключаем логирование ожидаемой ошибки записи дубликата поста в БД
			// в соответствии с правилом schema.sql
			// link TEXT NOT NULL UNIQUE -- UNIQUE для link позволяет избежать дублирования новостей в БД.
			var ErrDuplicate = errors.New("ERROR: duplicate key value violates unique constraint \"news_link_key\" (SQLSTATE 23505")
			if err != nil && errors.Is(err, ErrDuplicate) {
				chErrs <- err
			}
		}
	}()

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		chErrs <- err
	}
}

// Асинхронное чтение потока RSS. Раскодированные
// новости и ошибки пишутся в каналы.
func parseURL(url string, db *storage.Storage, posts chan<- []storage.Post, errs chan<- error, period int) {
	for {
		news, err := rss.ReadRSS(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
