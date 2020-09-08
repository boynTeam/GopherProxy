package domain

import (
	"testing"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/stretchr/testify/assert"
)

// Author:Boyn
// Date:2020/9/8

func TestInitAdminTable(t *testing.T) {
	pkg.InitDB()
	err := pkg.DefaultDB.AutoMigrate(&Admin{})
	assert.Nil(t, err)
}

func TestHashPassword(t *testing.T) {
	password := hashPassword("123456")
	assert.Equal(t, password, "MTIzNDU21B2M2Y8AsgTpgAmY7PhCfg==")
}

func TestAdmin_Save(t *testing.T) {
	pkg.InitDB()
	password := "123456"
	admin := Admin{
		UserName: "test01",
		Password: hashPassword(password),
	}
	err := admin.Save(nil, pkg.DefaultDB)
	if err != pkg.DuplicateRegisterError && err != nil {
		assert.Nil(t, err)
	}
}

func TestAdmin_Find(t *testing.T) {
	pkg.InitDB()
	admin := Admin{
		UserName: "test01",
	}
	result, err := admin.Find(nil, pkg.DefaultDB)
	assert.Nil(t, err)
	assert.Equal(t, result.UserName, "test01")
}

func TestAdmin_LoginCheck(t *testing.T) {
	pkg.InitDB()
	password := "123456"
	admin := Admin{
		UserName: "test01",
		Password: hashPassword(password),
	}
	check, err := admin.LoginCheck(nil, pkg.DefaultDB)
	assert.Nil(t, err)
	assert.NotNil(t, check)
}
