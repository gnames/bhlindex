package page

type Page struct {
	ID         string `gorm:"type:varchar(255);primary_key;auto_increment:false"`
	ItemID     int    `gorm:"primary_key;auto_increment:false"`
	Offset     int
	OffsetNext int `gorm:"-"`
}
