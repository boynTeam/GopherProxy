package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Author:Boyn
// Date:2020/9/8

func TestInitRedisCli(t *testing.T) {
	InitRedisCli()
	pong := RedisCli.Ping()
	result, err := pong.Result()
	assert.Nil(t, err)
	assert.Equal(t, result, "PONG")
}
