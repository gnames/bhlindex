package page

type Page struct {
	ID         int    `json:"id"`
	FileNum    int    `json:"fileNum"`
	ItemID     int    `json:"itemId"`
	FileName   string `json:"-" gorm:"varchar(255)"`
	Offset     int    `json:"-"`
	OffsetNext int    `json:"-" gorm:"-"`
}
