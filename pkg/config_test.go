package pkg

import (
	"testing"

	"github.com/spf13/viper"
)

// Author:Boyn
// Date:2020/8/31

func TestConfig(t *testing.T) {
	keys := viper.AllKeys()
	if len(keys) == 0 {
		t.Fatal("no config key readed")
	}
}
