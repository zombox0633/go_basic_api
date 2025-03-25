package constraints

import "net/http"

// customResponseWriter เป็น wrapper รอบ http.ResponseWriter
// เพื่อช่วยในการจับ status code
type CustomResponseWriterType struct {
	http.ResponseWriter
	StatusCode int
}

// Method WriteHeader เพื่อบันทึก status code
func (crw *CustomResponseWriterType) WriteHeader(statusCode int) {
	crw.StatusCode = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
