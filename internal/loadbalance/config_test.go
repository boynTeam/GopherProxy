package loadbalance

import (
	"fmt"
	"testing"
	"time"

	"github.com/BoynChan/GopherProxy/pkg"
)

// Author:Boyn
// Date:2020/9/7

type testObserver struct {
}

var resultChan = make(chan int)

func (t *testObserver) Update(values []ConfigValue) error {
	for _, c := range values {
		if c.Key == "TEST_TEMP_NODE" && c.Value == "TTTT" {
			resultChan <- 0
		}
	}
	return nil
}

func TestNewZkConf(t *testing.T) {
	tO := &testObserver{}
	prefix := "/http_real/proxy"
	conf, err := NewZkConf("%s", prefix, []string{"127.0.0.1:2181"})
	if err != nil {
		t.Fatal(err)
	}
	configValues, err := conf.GetConf()
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range configValues {
		fmt.Printf("%+v\n", c)
	}
	conf.Attach(tO)
	manager := pkg.NewZkManager([]string{"127.0.0.1:2181"}...)
	_ = manager.GetConnect()
	_ = manager.RegistServerTmpNode(prefix, "TEST_TEMP_NODE", []byte("TTTT")...)
	timer := time.NewTimer(5 * time.Second)
	select {
	case <-timer.C:
		t.Fatal("timeout")
	case <-resultChan:
	}
}
