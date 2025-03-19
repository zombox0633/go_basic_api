package constraints

import (
	"encoding/json"
	"log"
)

type PetType struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Weight  float64 `json:"weight"`
	Species string  `json:"species"`
}

var PetList []PetType

func init() {
	PetJson := `[
		{
			"id":1,
			"name":"Tigger",
			"weight":2.7,
			"species":"Cat"
		},
		{
			"id":2,
			"name":"Toufu",
			"weight":2.4,
			"species":"Cat"
		},
		{
			"id":3,
			"name":"Chinchan",
			"weight":0.7,
			"species":"GuiniePig"
		},
		{
			"id":4,
			"name":"LittleDragon",
			"weight":0.5,
			"species":"Hamster"
		},
		{
			"id":5,
			"name":"Yellow",
			"weight":0.3,
			"species":"Hamster"
		}		
	]`

	err := json.Unmarshal([]byte(PetJson), &PetList)
	if err != nil {
		log.Fatal(err) // พิมพ์ข้อความ error ไปที่ stderr(Standard Error) ปิดโปรแกรมทันที ด้วย os.Exit(1) (ไม่รัน defer)
		// panic(err) // สร้าง panic และพิมพ์ stack trace
	}
}
