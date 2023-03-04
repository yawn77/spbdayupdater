package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gookit/slog"
	"github.com/yawn77/spbdayupdater/pkg/bdayupdater"
	"github.com/yawn77/spcontrol"
)

var version string

type Config struct {
	printVersion bool

	// args are the positional (non-flag) command-line arguments.
	args []string
}

// source: https://eli.thegreenplace.net/2020/testing-flag-parsing-in-go-programs/
// parseFlags parses the command-line arguments provided to the program.
// Typically os.Args[0] is provided as 'progname' and os.args[1:] as 'args'.
// Returns the Config in case parsing succeeded, or an error. In any case, the
// output of the flag.Parse is returned in output.
// A special case is usage requests with -h or -help: then the error
// flag.ErrHelp is returned and output will contain the usage message.
func parseFlags(progname string, args []string) (config *Config, output string, err error) {
	flags := flag.NewFlagSet(progname, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)

	var conf Config
	flags.BoolVar(&conf.printVersion, "version", false, "print version")

	err = flags.Parse(args)
	if err != nil {
		return nil, buf.String(), err
	}

	conf.args = flags.Args()
	return &conf, buf.String(), nil
}

func runUpdateJob() {
	s := gocron.NewScheduler(time.Local)
	// _, err := s.Every(1).Day().At("00:00").Do(func() {
	_, err := s.Every(5).Second().Do(func() {
		bdayupdater.Update(true)
	})
	if err != nil {
		slog.Error(err)
	}
	s.StartBlocking()
}

func subMain(conf *Config) int {
	if conf.printVersion {
		fmt.Println(version)
		return 0
	}
	slog.Info("started SP Bday Updater")
	// test if valid credentials are provided
	creds, err := spcontrol.GetCredentials()
	if err != nil {
		slog.Error(err)
		return 1
	}
	slog.Infof("user: %s", creds.Username)
	runUpdateJob()
	return 0
}

func main() {
	conf, output, err := parseFlags(os.Args[0], os.Args[1:])
	if err == flag.ErrHelp {
		slog.Error(output)
		os.Exit(1)
	}
	if err != nil {
		slog.Error(err)
		slog.Error(output)
		os.Exit(1)
	}
	exitCode := subMain(conf)
	os.Exit(exitCode)
}
