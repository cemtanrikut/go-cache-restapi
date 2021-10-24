package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

type LogMiddleware struct {
	logger *log.Logger
}

func NewLogMiddleware(logger *log.Logger) *LogMiddleware {
	return &LogMiddleware{logger: logger}
}

type LogData struct {
	Time     time.Time
	Duration string
	Status   int
	Body     string
}

func (m *LogMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			logRespWriter := NewLogResponseWriter(w)
			next.ServeHTTP(logRespWriter, r)

			file, err := os.OpenFile("tmp/"+"httpServerLog.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println("-- DOSYA AÇMA/YAZMA HATASI")
				log.Fatal(err)

			}

			logData := &LogData{
				Time:     time.Now(),
				Duration: time.Since(startTime).String(),
				Status:   logRespWriter.statusCode,
				Body:     logRespWriter.buf.String(),
			}

			b, er := json.MarshalIndent(logData, "", "")
			if er != nil {
				fmt.Println("-- JSON CONVERT HATASI")
				log.Fatal(er)
			}

			var jsonBlob = []byte(b)
			dataErr := json.Unmarshal(jsonBlob, &logData)
			if dataErr != nil {
				fmt.Println("-- UNMARSHALL JSON ERR")
				log.Fatal(dataErr)
			}

			jData, _ := json.MarshalIndent(logData, "", "")

			_, writeErr := file.WriteString(string(jData))
			if writeErr != nil {
				fmt.Println("-- DOSYAYA YAZMA HATASI")
				log.Fatal(writeErr)
			}

			defer file.Close()

			m.logger.Printf(
				"duration=%s status=%d body=%s",
				time.Since(startTime).String(),
				logRespWriter.statusCode,
				logRespWriter.buf.String())
		})
	}
}
