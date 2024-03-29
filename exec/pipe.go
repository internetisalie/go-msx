// Copyright © 2022, Cisco Systems Inc.
// Use of this source code is governed by an MIT-style license that can be
// found in the LICENSE file or at https://opensource.org/licenses/MIT.

package exec

import (
	"cto-github.cisco.com/NFV-BU/go-msx/log"
	"gopkg.in/pipe.v2"
	"io"
	"net/http"
	"os"
	"strings"
)

var logger = log.NewLogger("msx.exec")

func RemoveAll(dir string) pipe.Pipe {
	return func(s *pipe.State) error {
		return os.RemoveAll(s.Path(dir))
	}
}

func Exec(name string, args []string, moreArgs ...[]string) pipe.Pipe {
	for _, moreArg := range moreArgs {
		args = append(args, moreArg...)
	}
	logger.Infof("cmd: %s %v", name, args)
	return pipe.Exec(name, args...)
}

func ExecQuiet(name string, args []string, moreArgs ...[]string) pipe.Pipe {
	for _, moreArg := range moreArgs {
		args = append(args, moreArg...)
	}
	return pipe.Exec(name, args...)
}

func ExecSimple(command ...string) pipe.Pipe {
	return Exec(command[0], command[1:])
}

func Info(template string, args ...interface{}) pipe.Pipe {
	return func(s *pipe.State) error {
		logger.Infof(template, args...)
		return nil
	}
}

// ExecutePipes executes the pipes, and on failure, sends output/error directly to our stderr
func ExecutePipes(pipes ...pipe.Pipe) error {
	for _, p := range pipes {
		if outputBytes, err := pipe.CombinedOutput(WithOutput(p)); err != nil {
			_, _ = os.Stderr.Write(outputBytes)
			return err
		}
	}

	return nil
}

// ExecutePipesStderr executes the pipes and sends output/error directly to our stdout
func ExecutePipesStderr(pipes ...pipe.Pipe) error {
	for _, p := range pipes {
		p = WithOutput(p)

		s := pipe.NewState(os.Stdout, os.Stdout)
		err := p(s)
		if err == nil {
			err = s.RunTasks()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func WithEnv(env map[string]string, p pipe.Pipe) pipe.Pipe {
	var pipes []pipe.Pipe
	for k, v := range env {
		k = strings.ToUpper(k)
		logger.Infof("env: %s=`%s`", k, v)
		pipes = append(pipes, pipe.SetEnvVar(k, v))
	}
	pipes = append(pipes, p)
	return pipe.Line(pipes...)
}

func WithDir(directory string, p pipe.Pipe) pipe.Pipe {
	if directory == "" {
		return p
	}

	return pipe.Line(
		pipe.ChDir(directory),
		p)
}

// WithOutput sends the output of p to our stdout
func WithOutput(p pipe.Pipe) pipe.Pipe {
	return pipe.Line(
		p,
		pipe.Write(os.Stdout))
}

// Run executes the pipe and sends output/error directly to our stdout/stderr
func Run(p pipe.Pipe) error {
	s := pipe.NewState(os.Stdout, os.Stderr)
	if err := p(s); err != nil {
		return err
	}
	return s.RunTasks()
}

// ReadUrl reads data from url and writes it to the pipe's stdout.
func ReadUrl(url string) pipe.Pipe {
	return pipe.TaskFunc(func(s *pipe.State) error {
		response, err := http.DefaultClient.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		_, err = io.Copy(s.Stdout, response.Body)
		return err
	})
}
