package autodelete_test

import (
	"testing"
	"time"

	"github.com/Jisin0/autofilterbot/pkg/autodelete"
	"github.com/stretchr/testify/assert"
)

func TestAutoDelete(t *testing.T) {
	assert := assert.New(t)

	m, err := autodelete.NewManager(nil)
	assert.NoError(err)

	assert.NoError(m.Save(1234567, 69, time.Minute*5))
	assert.NoError(m.Remove(1234567, 69))
}
