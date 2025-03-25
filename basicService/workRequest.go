package basicService

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/zombox0633/api/constraints"
	"github.com/zombox0633/api/middleware"
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

// Func สำหรับการหาข้อมูลจาก id ที่ส่งมา
// return ออก 2 ค่าคือ ตัวแรกคือ ข้อมูลที่ใช้ pointer ข้อมูลหนึ่งตัวใน PetList และ index ของ pat
func findPetById(id int) (*constraints.PetType, int) { // *😺 ของน้องใช้ส่วนของ return เป็น Pet (ค่าชนิดข้อมูลหนึ่งตัว) และ PetList (หลายตัว)
	for i, pat := range constraints.PetList {
		if pat.Id == id {
			return &pat, i
		}
	}
	return nil, -1 //ใส่ -1 เพื่อบอกว่าไม่มีข้อมูลอยู่ใน PetList
}

func deletePetByID(id int) error {
	_, index := findPetById(id)
	if index == -1 {
		return fmt.Errorf("pet not found")
	}

	//constraints.PetList[:index] เป็นตำแหน่งที่ต้องการลบ และ constraints.PetList[index+1:] เป็นส่วนหลังตำแหน่งที่ต้องการลบ
	constraints.PetList = append(constraints.PetList[:index], constraints.PetList[index+1:]...)
	return nil
}

func petByIDHandler(w http.ResponseWriter, r *http.Request) {
	//ทำการตัด url แยกเป็นส่วน ๆ จาก http://localhost:5000/pet/id จะได้ค่า ["http://localhost:5000/","id"]
	segment := strings.Split(r.URL.Path, "pet/")

	// urlPathSegment[len(urlPathSegment)-1] เลือกตำแหน่งสุดใน array จะได้ตัว id
	//จากนั้นเอา id ที่เป็น string ให้เป็น int จากคำสั่ง strconv.Atoi()
	id, err := strconv.Atoi(segment[len(segment)-1])

	if err != nil {
		log.Print(err)
		//ทำการส่งไปให้ HandleWriteHeader ส่ง error ให้โดยตัว utils เอิททำเอง
		utils.ErrorHandle(w, err, http.StatusNotFound)
		return
	}
	pet, index := findPetById(id)
	if index == -1 {
		utils.ErrorHandle(w, fmt.Errorf("invalid URL"), http.StatusNotFound)

		return
	}

	switch r.Method {
	// get by id
	case http.MethodGet:
		// แบบเก่า
		// petJson, err := json.Marshal(pet) //นำค่า pet มาแปลงเป็น json
		// utils.HandleWriteHeader(w, err, http.StatusNotFound)
		// w.Header().Set("Content-Type", "application/json")
		// w.Write(petJson)

		w.Header().Set("Content-Type", "application/json")

		//json.NewEncoder() ที่ใช้แปลง Go struct เป็น JSON และ สร้าง encoder ที่สามารถเขียน JSON ลงใน io.Writer
		// เมธอดที่แปลง Go struct (pet) เป็น JSON
		if err := json.NewEncoder(w).Encode(pet); err != nil {
			utils.ErrorHandle(w, err, http.StatusInternalServerError)
		}

	//แก้ไขข้อมูล pet โดยใช้ id
	case http.MethodPut:
		var updatePet constraints.PetType

		//ตัว ioutil.ReadAll ถูก deprecated  จึงใช้  io.ReadAll()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ErrorHandle(w, err, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		//แปลงค่า json ให้เป็น struct
		if err := json.Unmarshal(body, &updatePet); err != nil {
			utils.ErrorHandle(w, err, http.StatusBadRequest)
			return
		}

		//ถ้าค่าที่รับ ไม่เท่ากับ id ที่ใช้ใน URL จะ error
		if updatePet.Id != id {
			utils.ErrorHandle(w, fmt.Errorf("ID mismatch"), http.StatusBadRequest)
			return
		}

		constraints.PetList[index] = updatePet

		w.Header().Set("Content-Type", "application/json")
		// json.Marshal() + w.Write()  เขียน JSON ลง w ได้โดยตรง
		if err := json.NewEncoder(w).Encode(updatePet); err != nil {
			utils.ErrorHandle(w, err, http.StatusInternalServerError)
		}

	case http.MethodDelete:
		if err := deletePetByID(id); err != nil {
			utils.ErrorHandle(w, err, http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Pet with ID %d deleted successfully", id),
		})

	default:
		//405 Method Not Allowed
		utils.ErrorHandle(w, fmt.Errorf("method not allowed"), http.StatusMethodNotAllowed)

	}

}

// Func Handler สำหรับเรียกใช้งาน Method CRUD
func petHandler(w http.ResponseWriter, r *http.Request) {
	//แปลง PetList ให้อยู่ในรูปแบบ Json
	petJson, err := json.Marshal(constraints.PetList)

	switch r.Method {
	//method GET
	case http.MethodGet:
		utils.ErrorHandle(w, err, http.StatusInternalServerError)

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
		utils.ErrorHandle(w, err, http.StatusBadRequest)

		//แปลงค่า Json ที่รับมาจาก body เป็น struct และเพิ่มเข้าไปใน newPet
		err = json.Unmarshal(body, &newPet)
		utils.ErrorHandle(w, err, http.StatusBadRequest)

		//กรณีเกิด error ก็จะโยน error กลับมาให้
		if newPet.Id != 0 {
			utils.ErrorHandle(w, err, http.StatusBadRequest)
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
	fmt.Println("Registering Handlers...")

	//ประกาศ path และเรียกใช้งาน handler(Function การทำงานของ method ต่างๆ)
	//ประกาศใช้งาน middleware
	http.HandleFunc("/pet", middleware.LoggingMiddleware(petHandler))

	//การต่อ path by id
	http.HandleFunc("/pet/", middleware.LoggingMiddleware(petByIDHandler))
}
