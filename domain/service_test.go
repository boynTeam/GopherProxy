package domain

import (
	"testing"

	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Author:Boyn
// Date:2020/9/12

func TestInitServiceInfoTable(t *testing.T) {
	pkg.InitDB()
	err := pkg.DefaultDB.AutoMigrate(&ServiceInfo{})
	assert.Nil(t, err)
}

func TestServiceInfo_Save(t *testing.T) {
	pkg.InitDB()
	info := &ServiceInfo{
		ServiceType: HttpLoadType,
		ServiceName: "TEST_SERVICE",
		ServicePort: 1,
		ServiceDesc: "测试Service",
		RoundType:   RandomRound,
	}
	err := info.Save(nil, pkg.DefaultDB)
	assert.Nil(t, err)

}

func TestServiceInfo_Find(t *testing.T) {
	pkg.InitDB()
	info := &ServiceInfo{
		ServiceName: "TEST_SERVICE",
	}
	search, err := info.Find(nil, pkg.DefaultDB)
	require.Nil(t, err)
	assert.Equal(t, "TEST_SERVICE", search.ServiceName)
}

func TestServiceInfo_GetServiceDetail(t *testing.T) {
	pkg.InitDB()
	info := &ServiceInfo{
		ServiceName: "TEST_SERVICE",
	}
	search, err := info.Find(nil, pkg.DefaultDB)
	require.Nil(t, err)
	rule := HttpRule{
		ServiceId:      search.ID,
		Prefix:         "TEST_PREFIX",
		NeedHttps:      1,
		NeedWebsocket:  1,
		HeaderTransfor: "",
	}
	err = rule.Save(nil, pkg.DefaultDB)
	require.Nil(t, err)

	manager := pkg.NewZkManager([]string{"127.0.0.1:2181"}...)
	err = manager.GetConnect()
	require.Nil(t, err)
	err = manager.RegisterServerNode("/http_real_server", "TEST_SERVICE")
	require.Nil(t, err)

	detail, err := info.GetServiceDetail(nil, pkg.DefaultDB)

	require.Nil(t, err)
	assert.Equal(t, "127.0.0.1:1", detail.ServiceAddress)
	assert.Equal(t, 0, len(detail.ServerList))
}

func TestServiceInfo_Delete(t *testing.T) {
	pkg.InitDB()
	info := &ServiceInfo{
		ServiceName: "TEST_SERVICE",
	}
	err := info.Delete(nil, pkg.DefaultDB)
	require.Nil(t, err)
}
