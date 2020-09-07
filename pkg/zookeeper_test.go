package pkg

import "testing"

// Author:Boyn
// Date:2020/9/7

func TestZkManager_RegisterServerNode(t *testing.T) {
	zk := NewZkManager([]string{"127.0.0.1:2181"}...)
	zk.GetConnect()
	err := zk.doRegister("/abc", false)
	if err != nil {
		t.Fatal(err)
	}
	exist, err := zk.NodeExist("/abc")
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatalf("node /abc not exists")
	}
	err = zk.doRegister("/a11/bsad/cdwcd1", false)
	if err != nil {
		t.Fatal(err)
	}
	exist, err = zk.NodeExist("/a11/bsad")
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatalf("node /a11/bsad not exists")
	}
	exist, err = zk.NodeExist("/a11/bsad/cdw")
	if err != nil {
		t.Fatal(err)
	}
	if !exist {
		t.Fatalf("node /a11/bsad/cdw not exists")
	}
}
