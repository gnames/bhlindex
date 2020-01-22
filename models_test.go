package bhlindex_test

import (
	"path/filepath"

	"github.com/gnames/bhlindex/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type PageTest struct {
	FileName string
	Result   bool
}

var _ = BeforeEach(func() {
	truncateAll()
})

var _ = Describe("Models", func() {
	Describe("Item", func() {
		It("creates defaults", func() {
			t := models.Item{InternetArchiveID: "test"}
			Expect(t.InternetArchiveID).To(Equal("test"))
			Expect(t.Status).To(Equal(0))
		})
	})

	Describe("item.Insert()", func() {
		It("inserts item, returns id", func() {
			t := models.Item{InternetArchiveID: "test"}
			Expect(t.Insert(db)).To(BeNumerically(">", 0))
		})

		It("ignores duplicates", func() {
			t := models.Item{InternetArchiveID: "test"}
			id := t.Insert(db)
			Expect(id).To(BeNumerically(">", 0))
			t2 := models.Item{InternetArchiveID: "test"}
			id2 := t2.Insert(db)
			Expect(id2).To(Equal(0))
		})
	})

	Describe("ItemFind()", func() {
		It("finds a item in db", func() {
			t := models.Item{InternetArchiveID: "test"}
			id := t.Insert(db)
			t2 := models.ItemFind(db, id)
			Expect(t2.InternetArchiveID).To(Equal("test"))
		})
	})

	Describe("IsPageFile()", func() {
		It("determines if a file is a BHL page", func() {
			tests := []PageTest{
				{"/home/test_1234.txt", false},
				{filepath.Base("/home/test_1234.txt"), true},
				{"something.txt", false},
				{"something_0000.txt", true},
				{"som_ething_0543.txt", true},
				{"smt_0.txt", false},
				{"smt_00000.txt", false},
				{"smt_smt-1234.txt", false},
				{"smt_smt_1234.csv", false}}
			for _, t := range tests {
				res := models.IsPageFile(t.FileName)
				Expect(res).To(Equal(t.Result))
			}
		})
	})

	Describe("PageID()", func() {
		It("strips extension from filename", func() {
			tests := [][2]string{
				{"name_1234.txt", "name_1234"},
				{"name", "name"},
				{"name123.pdf", "name123"},
			}
			for _, t := range tests {
				Expect(models.PageID(t[0])).To(Equal(t[1]))
			}
		})
	})
})
