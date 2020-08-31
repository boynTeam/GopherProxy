package pkg

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
)

// Author:Boyn
// Date:2020/8/31

func TestLog(t *testing.T) {
	buf := bytes.Buffer{}
	logrus.SetOutput(&buf)
	logrus.Info("test msg")
	fmt.Println(buf.String())
}
