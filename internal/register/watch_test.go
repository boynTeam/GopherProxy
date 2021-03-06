package register

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/BoynChan/GopherProxy/pkg"
)

// Author:Boyn
// Date:2020/9/4

func TestWatch(t *testing.T) {
	zkManager := pkg.NewZkManager([]string{"127.0.0.1:2181"}...)
	zkManager.GetConnect()
	defer zkManager.Close()

	zlist, err := zkManager.GetServerListByPath("/http_real_server/TEST_HTTP_PROXY")
	fmt.Println("server node:")
	fmt.Println(zlist)
	if err != nil {
		log.Println(err)
	}

	//动态监听节点变化
	chanList, chanErr := zkManager.WatchServerListByPath("/http_real_server/TEST_HTTP_PROXY")
	go func() {
		for {
			select {
			case changeErr := <-chanErr:
				fmt.Println("changeErr")
				fmt.Println(changeErr)
			case changedList := <-chanList:
				fmt.Println("watch node changed")
				fmt.Println(changedList)
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
