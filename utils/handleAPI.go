package utils

import (
	"fmt"
	"net/http"
)

func HandleWriteHeader(w http.ResponseWriter, err error, statusCode int) {
	if err != nil {
		w.WriteHeader(statusCode)

		errMessage := fmt.Sprintf(`{"error": "HTTP Status %d"}`, statusCode) // สร้าง JSON message
		w.Write([]byte(errMessage))                                          // ส่งข้อความ error

		return
	}
}
