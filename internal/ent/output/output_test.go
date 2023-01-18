package output_test

import (
	"testing"

	"github.com/gnames/bhlindex/internal/ent/output"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeVerbatim(t *testing.T) {
	assert := assert.New(t)
	tests := []struct{ msg, inp, out string }{
		{"norm", "Bubo bubo", "Bubo bubo"},
		{"char1", "(Bubo bubo", "Bubo bubo"},
		{"char2", "[Bubo bubo/", "Bubo bubo"},
		{"space1", "Bubo 		  bubo", "Bubo bubo"},
		{"space2", "Bubo-  ␤␤␤ bubo", "Bubobubo"},
		{"space3", "Bubo-␤␤␤bubo", "Bubobubo"},
		{"caps", "BuBo-␤␤␤bubO", "Bubobubo"},
		{"caps2", "POMATomus Saltator", "Pomatomus saltator"},
		{"caps3", "Pomatomus (Salt) saltator", "Pomatomus (Salt) saltator"},
		{"caps4", "Pomatomus (salt) saltator", "Pomatomus (salt) saltator"},
		{"charint1", "P0matom^s s=ltar", "P0matom^s s=ltar"},
	}

	for _, v := range tests {
		oo := output.OutputOccurrence{
			DetectedVerbatim: v.inp,
		}
		oo.NormalizeVerbatim()
		assert.Equal(v.out, oo.DetectedVerbatim)
	}
}
