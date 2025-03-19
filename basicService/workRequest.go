package basicService

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/zombox0633/api/constraints"
	"github.com/zombox0633/api/utils"
)

// Func สำหรับการ auto id
// โดยการกำหนดค่าเริ่มต้นเท่ากับ -1 เนื่องด้วย id ไม่มีทางเป็นค่าติดลบ
// จากนั้นนำมา loop for โดยกำหนดการวนอยู่ใน PetList
// ซึ่งถ้า highestId มีค่าน้อยกว่า pet.Id ก็จะให้ highestId มีค่าเท่ากับ pet.Id
// และจะวนจนกว่าค่า highestId มีค่าไม่น้อยกว่า pet.Id
// จากนั้นจะให้ highestId + 1
func getNextID() int {
	highestId := -1
	for _, pet := range constraints.PetList { //range ใช้วนลูป ผ่าน slice, array, map, channel หรือ string เพื่อดึงค่าหรือ key ออกมาใช้งาน
		if highestId < pet.Id {
			highestId = pet.Id
		}
	}
	return highestId + 1
}

// Func Handler สำหรับเรียกใช้งาน Method CRUD
func petHandler(w http.ResponseWriter, r *http.Request) {
	//แปลง PetList ให้อยู่ในรูปแบบ Json
	petJson, err := json.Marshal(constraints.PetList)

	switch r.Method {
	//method GET
	case http.MethodGet:
		utils.HandleWriteHeader(w, err, http.StatusInternalServerError)

		//ทำการ set Header รูปแบบการจัดเรียงข้อมูลแบบ json
		w.Header().Set("Content-Type", "application/json")
		//ใช้ Write เพื่อแสดงข้อมูล Json ออกมา
		w.Write(petJson)
		return

		//method POST
	case http.MethodPost:
		//ประกาศตัวแปรสำหรับเพิ่มข้อมูล
		var newPet constraints.PetType

		//สร้างตัวแปล body มารับค่าที่อ่านผ่าน request ที่อยู่ใน body ทั้งหมด
		body, err := io.ReadAll(r.Body)
		utils.HandleWriteHeader(w, err, http.StatusBadRequest)

		//แปลงค่า Json ที่รับมาจาก body เป็น struct และเพิ่มเข้าไปใน newPet
		err = json.Unmarshal(body, &newPet)
		utils.HandleWriteHeader(w, err, http.StatusBadRequest)

		//กรณีเกิด error ก็จะโยน error กลับมาให้
		if newPet.Id != 0 {
			utils.HandleWriteHeader(w, err, http.StatusBadRequest)
		}

		//ถ้าไม่ error ก็จะทำการเพิ่ม newPet เข้าไปไว้ใน PetList
		//ซึ่งจะมีการกำหนด auto id ผ่าน func getNextID
		newPet.Id = getNextID()
		constraints.PetList = append(constraints.PetList, newPet)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Pet added successfully!"))
		return
	}

}

func WorkRequest() {
	//ประกาศ path และเรียกใช้งาน handler(Function การทำงานของ method ต่างๆ)
	http.HandleFunc("/pet", petHandler)
}
