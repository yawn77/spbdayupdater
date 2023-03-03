package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	version = "0.0.1"
	conf := Config{printVersion: true, args: []string{}}

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	exitCode := subMain(&conf)

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if exitCode != 0 {
		t.Errorf("exitCode got %d, want 0", exitCode)
	}
	if string(out) != fmt.Sprintln(version) {
		t.Errorf("version got %s, want %s", out, version)
	}
}

func TestParseFlags(t *testing.T) {
	var tests = []struct {
		args []string
		conf Config
	}{
		{[]string{"--version"},
			Config{printVersion: true, args: []string{}}},
		{[]string{},
			Config{printVersion: false, args: []string{}}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			conf, output, err := parseFlags("prog", tt.args)
			if err != nil {
				t.Errorf("err got %v, want nil", err)
			}
			if output != "" {
				t.Errorf("output got %q, want empty", output)
			}
			if !reflect.DeepEqual(*conf, tt.conf) {
				t.Errorf("conf got %+v, want %+v", *conf, tt.conf)
			}
		})
	}
}
