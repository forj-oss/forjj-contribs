package main

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Copy(src, dst string) (int64, error) {
	src_file, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return 0, err
	}

	if !src_file_stat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dst_file.Close()
	return io.Copy(dst_file, src_file)
}

// Simple function to call a shell command and display to stdout
// stdout is displayed as is when it arrives, while stderr is displayed in Red, line per line.
func run_cmd(command string, env []string, args ...string) (cmdlog []byte, err error) {

	cmd := exec.Command(command, args...)
	cmd.Env = env
	gotrace.Trace("RUNNING: %s %s", command, strings.Join(args, " "))

	// Execute command
	if cmdlog, err = cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("ERROR could not spawn command. %s.", err.Error())
		return
	}

	gotrace.Trace("Command done")
	if status := cmd.ProcessState.Sys().(syscall.WaitStatus); status.ExitStatus() != 0 {
		err = fmt.Errorf("\n%s ERROR: Unable to get process status - %s: %s", command, cmd.ProcessState.String())
	}
	return
}
