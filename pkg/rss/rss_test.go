// Пакет для обработки RSS-потока.

package rss

import (
	"testing"
)

func TestReadRSS(t *testing.T) {

	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		{"habr_rss", args{"https://habr.com/ru/rss/news/?fl=ru"}},
		{"rbc_rss", args{"https://rssexport.rbc.ru/rbcnews/news/30/full.rss"}},
		{"rg_rss", args{"https://rg.ru/xml/index.xml"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadRSS(tt.args.url)
			if (err != nil) && len(got) == 0 {
				t.Errorf("ReadRSS() error = %v, RSS не раскодировано", err)
				return
			}
			t.Logf("получено %d новостей\n", len(got))
		})
	}
}
