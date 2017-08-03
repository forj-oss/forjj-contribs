package main

import (
	"fmt"
	"github.com/forj-oss/forjj-modules/trace"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"crypto/md5"
)

func Copy(src, dst string) (written int64, err error, md5sum []byte) {
	var src_file *os.File

	src_file, err = os.Open(src)
	if err != nil {
		return
	}
	defer src_file.Close()

	src_file_stat, err := src_file.Stat()
	if err != nil {
		return
	}

	if !src_file_stat.Mode().IsRegular() {
		err = fmt.Errorf("%s is not a regular file", src)
		return
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dst_file.Close()

	md5_file := md5.New()
	tee_file := io.TeeReader(src_file, md5_file)
	written, err = io.Copy(dst_file, tee_file)
	md5sum = md5_file.Sum(nil)
	return
}

func md5sum(src string) ([]byte, error) {
	src_file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer src_file.Close()
	md5_file := md5.New()
	io.Copy(md5_file, src_file)
	return md5_file.Sum(nil), nil
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
