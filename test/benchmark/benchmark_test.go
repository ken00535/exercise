package benchmark

import (
	"net/http"
	"testing"
)

func BenchmarkRequest(b *testing.B) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	for i := 0; i < b.N; i++ {
		client.Get("http://localhost/qiI5wXYJ")
	}
}
