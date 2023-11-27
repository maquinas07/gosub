package utf8

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUTF8BomStrip(t *testing.T) {
	fd, err := os.Open("./test")
	if err != nil {
		t.Fatalf("Couldn't open test file")
	}

	content, err := io.ReadAll(fd)
	if err != nil {
		t.Fatalf("Fatal error in parsing %s", err)
	}
	fd.Close()
	assert.Subset(t, content, []byte{0xEF, 0xBB, 0xBF})
	assert.Subset(t, []byte(content), []byte("some text"))
	assert.NotEqual(t, []byte(content), []byte("some text"))
	content = StripUTF8BOM(content)
	assert.NotEmpty(t, content)
	assert.NotSubset(t, content, []byte{0xEF, 0xBB, 0xBF})
	assert.Equal(t, content, []byte("some text"))

}
