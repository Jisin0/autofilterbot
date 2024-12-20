package callbackdata_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/stretchr/testify/assert"
)

func TestCallbackData(t *testing.T) {
	assert := assert.New(t)

	{
		c := callbackdata.New()

		c.AddPath("dayumm")
		c.AddArg("foo")

		assert.Equal("dayumm|foo", c.ToString())
	}

	{
		c := callbackdata.FromString("paf1:paf2|arg1_arg2")

		assert.Equal([]string{"paf1", "paf2"}, c.Path)
		assert.Equal([]string{"arg1", "arg2"}, c.Args)
	}

	{
		c := callbackdata.FromString("p1|a1")

		assert.Equal([]string{"p1"}, c.Path)
		assert.Equal([]string{"a1"}, c.Args)
	}

	{
		c := callbackdata.FromString("justpath")

		assert.Equal([]string{"justpath"}, c.Path)
	}
}
