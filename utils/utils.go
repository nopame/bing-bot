package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// ✅ แปลง `sql.NullInt64` เป็น `*int` รองรับ `NULL`
func ToNullableInt(value sql.NullInt64) *int {
	if value.Valid {
		v := int(value.Int64)
		return &v
	}
	return nil
}

// ✅ แปลง `struct` เป็น JSON String พร้อม Format
func ToJSON(v interface{}) string {
	jsonData, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(jsonData)
}

// ✅ Sleep พร้อม Debug
func Sleep(duration time.Duration, message string) {
	fmt.Printf("⏳ Waiting %v... (%s)\n", duration, message)
	time.Sleep(duration)
}

// ✅ **พิมพ์เส้นแบ่ง**
func PrintDivider() {
	fmt.Println("--------------------------------------------------")
}

// ✅ ตัดข้อความให้สั้นลง โดยเพิ่ม "..." ท้ายข้อความ
func TruncateString(s string, length int) string {
	if len(s) > length {
		return s[:length] + "..."
	}
	return s
}