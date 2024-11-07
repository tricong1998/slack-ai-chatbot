package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectAction(t *testing.T) {
	action := DetectAction("abc (action: onboard_nhan_vien)")
	assert.Equal(t, "onboard_nhan_vien", action)
}
