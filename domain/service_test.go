package domain

import (
	"testing"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/stretchr/testify/assert"
)

// Author:Boyn
// Date:2020/9/12

func TestInitServiceInfoTable(t *testing.T) {
	pkg.InitDB()
	err := pkg.DefaultDB.AutoMigrate(&ServiceInfo{})
	assert.Nil(t, err)
}
