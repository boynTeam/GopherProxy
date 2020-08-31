package pkg

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Author:Boyn
// Date:2020/8/31

func init() {
	wd, _ := os.Getwd()
	path := []string{wd, "conf"}
	viper.AddConfigPath(strings.Join(path, string(os.PathSeparator)))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("load config error:%v\n", err))
	}
}
