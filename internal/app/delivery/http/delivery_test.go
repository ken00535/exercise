package http

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shorten/internal/app/entity"
	"shorten/internal/app/usecase/mock"
	"shorten/internal/infra/clock"
	"shorten/internal/infra/config"
	"shorten/internal/infra/gin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type UsecaseTestSuite struct {
	suite.Suite
	d         *Delivery
	ts        *httptest.Server
	mockUs    *mock.MockUsecase
	mockClock *clock.MockClock
}

func (suite *UsecaseTestSuite) SetupTest() {
	config.Init()
	cfg := config.Get()
	cfg.Http.DisableAccessLog = true
	cfg.Http.DisableBodyDump = true
	cfg.Http.LimiterBucket = 1000
	cfg.Http.LimiterTimeout = 1000
	router := gin.New(cfg.Http)
	ctrl := gomock.NewController(suite.T())
	suite.mockUs = mock.NewMockUsecase(ctrl)
	suite.d = New(router, suite.mockUs)
	suite.mockClock = clock.NewMockClock(ctrl)
	suite.d.clock = suite.mockClock
	suite.ts = httptest.NewServer(router)
}

func (suite *UsecaseTestSuite) TearDownTest() {
	suite.ts.Close()
}

func TestUrlTestSuite(t *testing.T) {
	suite.Run(t, new(UsecaseTestSuite))
}

func (suite *UsecaseTestSuite) TestUploadUrl() {
	suite.Run("success", func() {
		type TC struct {
			ParamPath string
			ParamBody string
			AnsStatus int
		}
		testcases := []TC{
			{
				"api/v1/urls",
				`{"url": "https://www.google.com.tw/","expireAt": "2036-12-19T16:39:57+08:00"}`,
				http.StatusOK,
			},
		}
		for _, tc := range testcases {
			suite.mockUs.EXPECT().
				UploadUrl(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&entity.Url{ShortUrl: "some"}, nil)

			reader := bytes.NewBufferString(tc.ParamBody)
			res, err := http.Post(suite.ts.URL+"/"+tc.ParamPath, "application/json", reader)

			if suite.NoError(err) {
				suite.Equal(tc.AnsStatus, res.StatusCode)
			}
		}
	})
	suite.Run("invalid param", func() {
		type TC struct {
			ParamPath string
			ParamBody string
			AnsStatus int
		}
		testcases := []TC{
			{
				"api/v1/urls",
				`{"url": "https//www.google.com.tw/","expireAt": "2036-12-19T16:39:57+08:00"}`,
				http.StatusBadRequest,
			},
			{
				"api/v1/urls",
				`{"url": "https://www.google.com.tw/","expireAt": "203612-19T16:39:57+08:00"}`,
				http.StatusBadRequest,
			},
		}
		for _, tc := range testcases {

			reader := bytes.NewBufferString(tc.ParamBody)
			res, err := http.Post(suite.ts.URL+"/"+tc.ParamPath, "application/json", reader)

			if suite.NoError(err) {
				suite.Equal(tc.AnsStatus, res.StatusCode)
			}
		}
	})
	suite.Run("not found", func() {
		type TC struct {
			ParamPath string
			ParamBody string
			AnsStatus int
		}
		testcases := []TC{
			{
				"api/v1/wrong",
				`{"url": "https://www.google.com.tw/","expireAt": "2036-12-19T16:39:57+08:00"}`,
				http.StatusNotFound,
			},
		}
		for _, tc := range testcases {

			reader := bytes.NewBufferString(tc.ParamBody)
			res, err := http.Post(suite.ts.URL+"/"+tc.ParamPath, "application/json", reader)

			if suite.NoError(err) {
				suite.Equal(tc.AnsStatus, res.StatusCode)
			}
		}
	})
	suite.Run("internal error", func() {
		type TC struct {
			ParamPath string
			ParamBody string
			AnsStatus int
		}
		testcases := []TC{
			{
				"api/v1/urls",
				`{"url": "https://www.google.com.tw/","expireAt": "2036-12-19T16:39:57+08:00"}`,
				http.StatusInternalServerError,
			},
		}
		for _, tc := range testcases {
			suite.mockUs.EXPECT().
				UploadUrl(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&entity.Url{ShortUrl: "some"}, errors.New("some"))

			reader := bytes.NewBufferString(tc.ParamBody)
			res, err := http.Post(suite.ts.URL+"/"+tc.ParamPath, "application/json", reader)

			if suite.NoError(err) {
				suite.Equal(tc.AnsStatus, res.StatusCode)
			}
		}
	})
}

func (suite *UsecaseTestSuite) TestServeShortenUrl() {
	suite.Run("success", func() {
		type TC struct {
			ParamExpire time.Time
			AnsStatus   int
			AnsUrl      string
			AnsAge      int
		}
		testcases := []TC{
			{
				time.Date(2024, 1, 1, 0, 0, 1, 0, time.Local),
				http.StatusFound,
				"/www.google.com",
				1,
			},
			{
				time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
				http.StatusFound,
				"/www.google.com",
				24 * 60 * 60,
			},
		}
		for _, tc := range testcases {
			t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local)
			suite.mockUs.EXPECT().
				GetUrl(gomock.Any(), gomock.Any()).
				Return(&entity.Url{Url: tc.AnsUrl, ExpireAt: tc.ParamExpire}, nil)
			suite.mockClock.EXPECT().Now().Return(t)
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			res, err := client.Get(suite.ts.URL + "/" + "shorten")
			expire := fmt.Sprintf("public, max-age=%d", tc.AnsAge)

			if suite.NoError(err) {
				suite.Equal(tc.AnsStatus, res.StatusCode)
				suite.Equal(tc.AnsUrl, res.Header.Get("Location"))
				suite.Equal(expire, res.Header.Get("Cache-Control"))
			}
		}
	})
}
