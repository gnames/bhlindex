package page

// Page represents a page of text.
type Page struct {
	// ID is the ID from BHL database assigned to a Page.
	ID int

	// FileNum is the number provided in the file name. This number is rarely
	// larger than 500.
	FileNum int

	// ItemID is the ID from BHL database assigned to an Item.
	ItemID int

	// FileName is taken from the filesystem.
	FileName string

	// Offset is the number of UTF-8 characters from the start of an Item
	// to the start of the Page.
	Offset int

	// OffsetNext is the offset for the next Page if it exists.
	OffsetNext int
}
