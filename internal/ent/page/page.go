package page

type Page struct {
	ID         int    `json:"id"`
	FileID     int    `json:"fileId"`
	ItemID     int    `json:"itemId"`
	FileName   string `json:"-" gorm:"varchar(255)"`
	Offset     int    `json:"-"`
	OffsetNext int    `json:"-" gorm:"-"`
}
