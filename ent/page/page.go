package page

type Page struct {
	ID         string `json:"id" gorm:"type:varchar(255);primary_key;auto_increment:false"`
	ItemID     int    `json:"item_id" gorm:"primary_key;auto_increment:false"`
	Offset     int    `json:"-"`
	OffsetNext int    `json:"-" gorm:"-"`
}
