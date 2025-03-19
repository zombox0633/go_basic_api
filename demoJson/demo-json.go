package demojson

import (
	"encoding/json"
	"fmt"
)

type AnimalType struct {
	Id         int //ชื่อตัวแปรต้องขึ้นต้นด้วย Uppercase เพื่อให้ json.Marshal ใช้งานได้
	Species    string
	SupSpecies string
}

func DemoJSon() {
	// & หมายถึงการส่งค่าแบบ pointer เพื่อลดการทำสำเนาข้อมูล (Efficiency) → ไม่ต้องคัดลอก struct ทั้งตัว
	animal := AnimalType{1, "Cat", "Siamese"}
	fmt.Println("😾 : ", animal)

	// แปลง struct เป็น JSON
	dataMarshal, _ := json.Marshal(&animal)
	fmt.Println("😸 : ", string(dataMarshal)) //การอ่านข้อมูลต้อง

	// แปลง JSON กลับเป็น struct
	newAnimal := AnimalType{}
	dataUnmarshal := json.Unmarshal(dataMarshal, &newAnimal)
	// dataUnmarshal := json.Unmarshal([]byte(`{"Id":1,"Species":"Cat","SupSpecies":"Siamese"}`), &newAnimal)

	if dataUnmarshal != nil {
		fmt.Println("😿 Unmarshal Error: ", dataUnmarshal)
	} else {
		fmt.Println("😿 : ", newAnimal)
	}
}
