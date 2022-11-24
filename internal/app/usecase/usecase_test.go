package usecase

import (
	"context"
	"net/url"
	"testing"
	"time"

	"shorten/internal/app/entity"
	"shorten/internal/app/repository/db/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UsecaseTestSuite struct {
	suite.Suite
	us     *usecase
	mockDb *mock.MockRepositoryDb
}

func (suite *UsecaseTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.mockDb = mock.NewMockRepositoryDb(ctrl)
	suite.us = &usecase{
		db: suite.mockDb,
		cfg: Config{
			Scheme:   "http",
			Hostname: "localhost",
		},
	}
}

func (suite *UsecaseTestSuite) TearDownTest() {
}

func TestWithdrawTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite *UsecaseTestSuite) TestCreateWithdrawal() {
	suite.Run("success", func() {
		type TC struct {
			Param1 string
			Ans    string
		}
		testcases := []TC{
			{
				"http://google.com",
				"http://google.com",
			},
		}
		for _, tc := range testcases {
			suite.mockDb.EXPECT().
				SaveShortenUrl(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, shorten *entity.Url) (*entity.Url, error) {
					return shorten, nil
				})

			ans, err := suite.us.UploadUrl(context.TODO(), tc.Param1, time.Now())
			shorten, _ := url.Parse(ans.ShortUrl)

			suite.NoError(err)
			suite.Equal(tc.Ans, ans.Url)
			suite.Equal(9, len(shorten.Path))
		}
	})
}

func FuzzUploadUrl(f *testing.F) {
	u := usecase{
		cfg: Config{
			Scheme:   "http",
			Hostname: "localhost",
		},
	}
	f.Fuzz(func(t *testing.T, source string) {
		ctrl := gomock.NewController(t)
		mockDb := mock.NewMockRepositoryDb(ctrl)
		u.db = mockDb

		mockDb.EXPECT().
			SaveShortenUrl(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, shorten *entity.Url) (*entity.Url, error) {
				return shorten, nil
			})

		ans, _ := u.UploadUrl(context.TODO(), source, time.Now())
		shorten, _ := url.Parse(ans.ShortUrl)
		assert.Equal(t, 9, len(shorten.Path))
	})
}
