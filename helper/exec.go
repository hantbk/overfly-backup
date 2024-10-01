// Copyright © 2024 Ha Nguyen <captainnemot1k60@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package helper

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hantbk/vtsbackup/logger"
)

var (
	spaceRegexp = regexp.MustCompile(`[\s]+`)
)

// Exec cli commands
func Exec(command string, args ...string) (output string, err error) {
	return ExecWithStdio(command, false, args...)
}

// ExecWithStdio cli commands with stdio
func ExecWithStdio(command string, stdout bool, args ...string) (output string, err error) {
	commands := spaceRegexp.Split(command, -1)
	command = commands[0]
	commandArgs := []string{}
	if len(commands) > 1 {
		commandArgs = commands[1:]
	}
	if len(args) > 0 {
		commandArgs = append(commandArgs, args...)
	}

	fullCommand, err := exec.LookPath(command)
	if err != nil {
		return "", fmt.Errorf("%s cannot be found", command)
	}

	cmd := exec.Command(fullCommand, commandArgs...)
	cmd.Env = os.Environ()

	var stdErr bytes.Buffer
	var stdOut bytes.Buffer
	cmd.Stderr = &stdErr

	if stdout {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = &stdOut
	}

	err = cmd.Run()
	if err != nil {
		logger.Debug(fullCommand, " ", strings.Join(commandArgs, " "))
		err = errors.New(stdErr.String())
	}
	output = strings.Trim(stdOut.String(), "\n")

	return
}
