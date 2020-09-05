package register

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-zookeeper/zk"
)

// Author:Boyn
// Date:2020/9/4

var (
	host = []string{"localhost:2181"}
)

func TestCRUD(t *testing.T) {
	conn, _, err := zk.Connect(host, 5*time.Second)

	if err != nil {

		panic(err)
	}

	_, dStat, _ := conn.Get("/test_tree")
	if err := conn.Delete("/test_tree", dStat.Version); err != nil {
		fmt.Println("Delete err", err)
		//return
	}

	if _, err := conn.Create("/test_tree", []byte("tree_content"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		t.Fatalf("create fail %v", err)
	}
	nodeValue, dStat, err := conn.Get("/test_tree")
	if err != nil {
		t.Fatalf("get fail %v", err)
	}
	fmt.Println(string(nodeValue))
	//改
	if _, err := conn.Set("/test_tree", []byte("new_content"),
		dStat.Version); err != nil {
		t.Fatalf("update fail %v", err)
	}

	//_, dStat, _ = conn.Get("/test_tree")
	//if err := conn.Delete("/test_tree", dStat.Version); err != nil {
	//	t.Fatalf("Delete err %v", err)
	//	//return
	//}

	//设置子节点
	if _, err := conn.Create("/test_tree/subnode", []byte("node_content"),
		0, zk.WorldACL(zk.PermAll)); err != nil {
		t.Fatalf("create err %v", err)
	}

	//获取子节点列表
	childNodes, _, err := conn.Children("/test_tree")
	if err != nil {
		t.Fatalf("Children err %v", err)
	}
	fmt.Println("childNodes", childNodes)
}
