package example

import "time"

type DemoStruct2 struct {
	Test1 string `json:"test1"`
	Test2 int    `json:"test2"`
}

type DemoStruct struct {
	ID1       uint32      `json:"id_1"`
	ID2       uint64      `json:"id_2"`
	ID3       uint8       `json:"id_3"`
	ID4       uint16      `json:"id_4"`
	Status    float32     `json:"status"`
	Name      string      `json:"name"`
	Age       int         `json:"age"`
	Address   string      `json:"address"`
	IsMarried bool        `json:"is_married"`
	Children  []string    `json:"children"`
	Salary    float64     `json:"salary"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt []time.Time `json:"deletedAt"`
	CreatedBy uint64      `json:"createdBy"`
	UpdatedBy uint64      `json:"updatedBy"`
	DeletedBy uint64      `json:"deletedBy"`
	Test      DemoStruct2 `cosy:"add:123" json:"test"`
}
