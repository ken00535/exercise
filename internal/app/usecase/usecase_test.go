package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func FuzzUploadUrl(f *testing.F) {
	u := usecase{}
	f.Fuzz(func(t *testing.T, url string) {
		ans, _ := u.UploadUrl(context.TODO(), url, time.Now())
		assert.Equal(t, 8, len(ans.ShortUrl))
	})
}
