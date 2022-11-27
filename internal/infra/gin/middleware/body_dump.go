package middleware

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// BodyDump dump request body and response body
func BodyDump(disable bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if disable {
			return
		}
		var reqBody strings.Builder
		var resBody strings.Builder
		var isSkip bool

		//dump req body
		if !isSkip && ctx.Request.ContentLength > 0 {
			buf, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				log.Err(err).Msgf("%s", err.Error())
				return
			}
			reqDumpReader := ioutil.NopCloser(bytes.NewBuffer(buf))
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

			bts, err := ioutil.ReadAll(reqDumpReader)
			if err != nil {
				log.Err(err).Msgf("%s", err.Error())
				return
			}

			reqBody.WriteString(string(bts))
		}

		ctx.Writer = &bodyWriter{bodyCache: bytes.NewBuffer([]byte{}), ResponseWriter: ctx.Writer}
		ctx.Next()

		if isSkip {
			return
		}

		bw, _ := ctx.Writer.(*bodyWriter)

		if bw.bodyCache.Len() > 0 && bw.bodyCache.Len() < 1024 {
			s := bw.bodyCache.Bytes()
			resBody.WriteString(string(s))
		} else if bw.bodyCache.Len() > 1024 {
			s := bw.bodyCache.Bytes()
			resBody.WriteString(string(s[:1024]))
		}
		if len(ctx.Errors) != 0 {
			log.Error().
				Str("request_body", reqBody.String()).
				Str("response_body", resBody.String()).
				Msg("http body dump")
			return
		}
		log.Info().
			Str("request_body", reqBody.String()).
			Str("response_body", resBody.String()).
			Msg("http body dump")
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	bodyCache *bytes.Buffer
}

//rewrite Write()
func (w bodyWriter) Write(b []byte) (int, error) {
	w.bodyCache.Write(b)
	return w.ResponseWriter.Write(b)
}
