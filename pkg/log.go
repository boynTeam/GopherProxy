package pkg

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Init log repository
// Author:Boyn
// Date:2020/8/31

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrusFormater{})
}

type logrusFormater struct {
}

func (s *logrusFormater) Format(entry *logrus.Entry) ([]byte, error) {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s] %s %+v \n", strings.ToUpper(entry.Level.String()), timeStr, entry.Message)
	return []byte(msg), nil
}
