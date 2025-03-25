package utils

import (
	"encoding/json"
	"net/http"
)

func ErrorHandle(w http.ResponseWriter, err error, statusCode int) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		// แบบเก่า
		//errMessage := fmt.Sprintf(`{"error": "HTTP Status %d"}`, statusCode) // สร้าง JSON message
		//w.Write([]byte(errMessage))                                          // ส่งข้อความ error

		errResponse := struct {
			Error   string `json:"error"`             // เอาไว้ใส่ statusCode
			Message string `json:"message,omitempty"` // รายละเอียดข้อผิดพลาด
		}{
			Error: http.StatusText(statusCode),
		}

		errResponse.Message = err.Error()

		//วิธีที่สะอาดและมีประสิทธิภาพในการเขียน JSON response
		//แปลง struct เป็น JSON และเขียนลงใน response writer โดยตรง
		json.NewEncoder(w).Encode(errResponse)
		// return
	}
}
