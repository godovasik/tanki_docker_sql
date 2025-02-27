package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func SetupLogger() {
	file, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("cant open log file, using stdout only")
		Log.SetOutput(os.Stdout)
	} else {
		mw := io.MultiWriter(os.Stdout, file)
		Log.SetOutput(mw)
	}
	Log.SetLevel(logrus.DebugLevel)
	Log.SetFormatter(&logrus.TextFormatter{})
}
