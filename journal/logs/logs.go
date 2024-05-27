package logs

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"hafidzresttemplate.com/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

func InitLog() *logrus.Logger {
	loggerInit := logrus.New()

	// Add this line for logging filename and line number!
	loggerInit.SetReportCaller(true)
	
	// Open the log file
	file, err := os.OpenFile(fmt.Sprintf("logs/daily/journal_%v.log", utils.CURRENT_DATE), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Set output to both file and standard output
		multiWriter := io.MultiWriter(os.Stdout, file)
		loggerInit.SetOutput(multiWriter)
	} else {
		loggerInit.Error("Failed to log to file, using default stderr")
	}

	// Set the formatter
	loggerInit.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s:%d", f.File, f.Line), ""
		},
	})

	// Add hook for writing to file
	loggerInit.AddHook(&writer.Hook{
		Writer: file,
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	return loggerInit
}
