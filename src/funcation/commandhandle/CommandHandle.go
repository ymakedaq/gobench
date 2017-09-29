package commandhandle

import (
	"errors"
	"fmt"
	"funcation/golog"
	"io/ioutil"
	"os/exec"
	"strconv"
)

func CommandExecResultBytes(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	defer stderr.Close()
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		golog.Error("command", "command", fmt.Sprintf("Execue %s Fail!", err), 0)
		return nil, err
	}

	opBytes, err := ioutil.ReadAll(stdout)
	opError, err := ioutil.ReadAll(stderr)
	if err != nil {
		golog.Error("command", "command", fmt.Sprintf("Execue %s Fail!", err), 0)
		return nil, err
	}
	if len(opError) > 0 {
		return nil, errors.New(string(opError))
	} else {
		return opBytes, nil
	}
}
