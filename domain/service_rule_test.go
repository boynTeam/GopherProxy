package domain

import (
	"testing"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/stretchr/testify/assert"
)

// Author:Boyn
// Date:2020/9/12

func TestInitHttpRuleTable(t *testing.T) {
	pkg.InitDB()
	err := pkg.DefaultDB.AutoMigrate(&HttpRule{})
	assert.Nil(t, err)
}

func TestInitGrpcRuleTable(t *testing.T) {
	pkg.InitDB()
	err := pkg.DefaultDB.AutoMigrate(&GrpcRule{})
	assert.Nil(t, err)
}
