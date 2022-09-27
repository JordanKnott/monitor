package jobs

import (
	"github.com/sirupsen/logrus"
)

// MachineryLogger is a customer logrusger for machinery worker
type MachineryLogger struct{}

// Print sends to logrusrus.Info
func (m *MachineryLogger) Print(args ...interface{}) {
	logrus.Info(args...)
}

// Printf sends to logrusrus.Infof
func (m *MachineryLogger) Printf(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Println sends to logrusrus.Info
func (m *MachineryLogger) Println(args ...interface{}) {
	logrus.Info(args...)
}

// Fatal sends to logrusrus.Fatal
func (m *MachineryLogger) Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf sends to logrusrus.Fatalf
func (m *MachineryLogger) Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// Fatalln sends to logrusrus.Fatal
func (m *MachineryLogger) Fatalln(args ...interface{}) {
	logrus.Fatal(args...)
}

// Panic sends to logrusrus.Panic
func (m *MachineryLogger) Panic(args ...interface{}) {
	logrus.Panic(args...)
}

// Panicf sends to logrusrus.Panic
func (m *MachineryLogger) Panicf(format string, args ...interface{}) {
	logrus.Panic(args...)
}

// Panicln sends to logrusrus.Panic
func (m *MachineryLogger) Panicln(args ...interface{}) {
	logrus.Panic(args...)
}
