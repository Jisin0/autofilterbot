package configpanel_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/pkg/configpanel"
	"github.com/stretchr/testify/assert"
)

func TestCallbackData(t *testing.T) {
	rawData := "path1:path2|arg1_arg2_arg3"

	d := configpanel.CallbackDataFromString(rawData)
	assert := assert.New(t)

	assert.Equal([]string{"path1", "path2"}, d.Path)
	assert.Equal([]string{"arg1", "arg2", "arg3"}, d.Args)

	assert.Equal("path1:path2:path3|arg1_arg2_arg3", d.AddPath("path3").ToString())
	assert.Equal("path1:path2|arg1_arg2_arg3_arg4", d.AddArg("arg4").ToString())
	assert.Equal("path1|arg1_arg2_arg3", d.RemoveLastPath().ToString())
}
