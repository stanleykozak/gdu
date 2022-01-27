package analyze

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"github.com/dundee/gdu/v5/pkg/fs"
	"github.com/stretchr/testify/assert"
)

func TestEncDec(t *testing.T) {
	var d fs.Item = &StoredDir{
		Dir: &Dir{
			File: &File{
				Name: "xxx",
			},
			BasePath: "/yyy",
		},
	}

	b := &bytes.Buffer{}
	enc := gob.NewEncoder(b)
	err := enc.Encode(d)
	assert.NoError(t, err)

	var x fs.Item = &StoredDir{}
	dec := gob.NewDecoder(b)
	err = dec.Decode(x)
	assert.NoError(t, err)

	fmt.Println(d, x)
	assert.Equal(t, d.GetName(), x.GetName())
}
