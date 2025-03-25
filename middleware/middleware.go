package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zombox0633/api/constraints"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		fmt.Println("Middleware is working! 😾")

		// Log ข้อมูลการร้องขอ
		fmt.Printf(
			"Incoming Request - [%s] %s %s - Remote Address: %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			r.RemoteAddr,
		)

		crw := &constraints.CustomResponseWriterType{
			ResponseWriter: w,
			StatusCode:     http.StatusOK, // Default status
		}

		// เรียกใช้ handler ถัดไป
		next.ServeHTTP(crw, r)

		duration := time.Since(startTime)

		// Log การตอบกลับ
		fmt.Printf(
			"Response - Status: %d, Duration: %v",
			crw.StatusCode,
			duration,
		)
	}
}
