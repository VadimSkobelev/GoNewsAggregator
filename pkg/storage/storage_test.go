// Пакет для работы с БД PostgreSQL.

package storage

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorage_AddNews(t *testing.T) {

	db, err := New()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		p []Post
	}

	tests := []struct {
		name    string
		s       *Storage
		args    args
		wantErr bool
	}{
		{"AddPostInDB-1", db, args{[]Post{{Title: "Test post-01", Content: "Testing_1", Link: strconv.Itoa(rand.Intn(1_000_000_000))}}}, false},
		{"AddPostInDB-2", db, args{[]Post{{Title: "Test post-02", Content: "Testing_2", Link: strconv.Itoa(rand.Intn(1_000_000_000))}}}, false},
		{"AddPostInDB-3", db, args{[]Post{{Title: "Test post-03", Content: "Testing_3", Link: "https://test.test"}}}, false},
		{"AddPostInDB-4", db, args{[]Post{{Title: "Test post-04", Content: "Testing_4", Link: "https://test.test"}}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.AddNews(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("Storage.AddNews() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		news, err := db.News(2)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v", news)
	}

	// Очищаем БД от тестовых данных.
	ss, _ := New()
	for _, tt := range tests {
		for _, pp := range tt.args.p {
			_, _ = ss.db.Exec(context.Background(), `DELETE FROM news WHERE title=$1;`, pp.Title)
		}
	}
}
