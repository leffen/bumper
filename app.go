package main

import (
	"os"

	"github.com/leffen/cruncy"
	"github.com/sirupsen/logrus"
)

type app struct {
	Name          string
	FileName      string
	Part          string
	Input         string
	Format        string
	Extract       bool
	MetricsPort   int
	Verbose       bool
	JSONFormatter bool
}

// VERSION of the application
const VERSION = "0.0.1"

var (
	// CommitHash is the revision hash of the build's git repository
	CommitHash string
	// BuildTime is the build's time
	BuildTime string
)

func (a *app) Assign(name string) {
	a.Name = name

	o := cruncy.NewCliOption("bumper")

	o.MakeString("filename", "f", "FILE_NAME", "VERSION", "Name of file containing version")
	o.MakeString("input", "i", "INPUT", "", "Optionally use input pattern to bump. Instead of filename")
	o.MakeString("format", "F", "FORMAT", "", "either M for major, M-m for major-minor or M-m-p")
	o.MakeString("part", "p", "PART", "p", "M for major, m for minor or p for patch")

	o.MakeBool("extract", "e", "EXTRACT", false, "Only extract version and displayit")

	o.MakeBool("json_formatter", "j", "JSON_LOG", true, "JSON logging")
	o.MakeBool("verbose", "V", "VERBOSE", false, "Verbose logging")

	o.ReadConfig()

	a.FileName = o.GetString("filename")
	a.Input = o.GetString("input")
	a.Extract = o.GetBool("extract")
	a.Format = o.GetString("format")
	a.Verbose = o.GetBool("verbose")
	a.JSONFormatter = o.GetBool("json_formatter")

	a.StartLog()

}

func (a *app) StartLog() {
	logLevel := logrus.InfoLevel

	if a.Verbose {
		logLevel = logrus.DebugLevel
	}

	logrus.SetLevel(logLevel)
	logrus.SetOutput(os.Stdout)
	if a.JSONFormatter {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
