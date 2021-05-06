package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	md5 := MD5("test")
	assert.EqualValues(t, "098f6bcd4621d373cade4e832627b4f6", md5)

	sha1 := SHA1("test")
	assert.EqualValues(t, "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3", sha1)

	sha256 := SHA256("test")
	assert.EqualValues(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", sha256)

}
