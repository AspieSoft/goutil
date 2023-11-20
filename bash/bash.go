package bash

import (
	"io"
	"os"
	"os/exec"
)

/*
	Run will run a bash command based on the given args

	note: stdin is piped to the os logs

		@dir: a directory to run the command in (set to an empty string to disable)
		@env: an optional list of environment variables (set to nil to disable)

	[optional]
		@liveOutput[0]: set to true to pipe stdout and stderr to the os
		@liveOutput[1]: set to false to only pipe stdout to the os, and keep stderr hidden
*/
func Run(args []string, dir string, env []string, liveOutput ...bool) (output []byte, err error) {
	arg1 := args[0]
	args = args[1:]
	cmd := exec.Command(arg1, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	if env != nil {
		cmd.Env = append(cmd.Env, env...)
	}
	cmd.Stdin = os.Stdin
	if len(liveOutput) != 0 && liveOutput[0] == true {
		cmd.Stdout = os.Stdout
		if len(liveOutput) <= 1 || liveOutput[1] == true {
			cmd.Stderr = os.Stderr
		}
		return []byte{}, cmd.Run()
	}
	return cmd.CombinedOutput()
}

/*
	RunRaw will run an unescaped (unquoted) bash command as a different user

	this method uses `bash -c` to get around the auto quotes added by golang

	note: user input is Not recommended for this method

	note: stdin is piped to the os logs

		@cmdStr: the command to run
		@dir: a directory to run the command in (set to an empty string to disable)
		@env: an optional list of environment variables (set to nil to disable)

	[optional]
		@liveOutput[0]: set to true to pipe stdout and stderr to the os
		@liveOutput[1]: set to false to only pipe stdout to the os, and keep stderr hidden
*/
func RunRaw(cmdStr string, dir string, env []string, liveOutput ...bool) (output []byte, err error) {
	cmd := exec.Command(`bash`, `-c`, cmdStr)
	if dir != "" {
		cmd.Dir = dir
	}
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	if env != nil {
		cmd.Env = append(cmd.Env, env...)
	}
	cmd.Stdin = os.Stdin
	if len(liveOutput) != 0 && liveOutput[0] == true {
		cmd.Stdout = os.Stdout
		if len(liveOutput) <= 1 || liveOutput[1] == true {
			cmd.Stderr = os.Stderr
		}
		return []byte{}, cmd.Run()
	}
	return cmd.CombinedOutput()
}

/*
	RunUser will run an unescaped (unquoted) bash command as a specified user

	this method uses `runuser -l [user] -c`

	note: user input is Not recommended for this method

	note: stdin is piped to the os logs

		@cmdStr: the command to run
		@user: the username to run the command as
		@dir: a directory to run the command in (set to an empty string to disable)
		@env: an optional list of environment variables (set to nil to disable)

	[optional]
		@liveOutput[0]: set to true to pipe stdout and stderr to the os
		@liveOutput[1]: set to false to only pipe stdout to the os, and keep stderr hidden
*/
func RunUser(cmdStr string, user string, dir string, env []string, liveOutput ...bool) (output []byte, err error) {
	cmd := exec.Command(`runuser`, `-l`, user, `-c`, cmdStr)
	if dir != "" {
		cmd.Dir = dir
	}
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	if env != nil {
		cmd.Env = append(cmd.Env, env...)
	}
	cmd.Stdin = os.Stdin
	if len(liveOutput) != 0 && liveOutput[0] == true {
		cmd.Stdout = os.Stdout
		if len(liveOutput) <= 1 || liveOutput[1] == true {
			cmd.Stderr = os.Stderr
		}
		return []byte{}, cmd.Run()
	}
	return cmd.CombinedOutput()
}

/*
	RunUserFile will run a bash file as a specified user

	this method uses `pkexec --user [user]` to simulate a user in a normal desktop environment

	note: user input is Not recommended for this method

	note: stdin is piped to the os logs

		@file: the file to run
		@user: the username to run the command as
		@dir: a directory to run the command in (set to an empty string to disable)
		@env: an optional list of environment variables (set to nil to disable)

	[optional]
		@liveOutput[0]: set to true to pipe stdout and stderr to the os
		@liveOutput[1]: set to false to only pipe stdout to the os, and keep stderr hidden
*/
func RunUserFile(file string, args []string, user string, dir string, env []string, liveOutput ...bool) (output []byte, err error) {
	cmd := exec.Command(`pkexec`, append([]string{`--user`, user, file}, args...)...)
	if dir != "" {
		cmd.Dir = dir
	}
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	if env != nil {
		cmd.Env = append(cmd.Env, env...)
	}
	cmd.Stdin = os.Stdin
	if len(liveOutput) != 0 && liveOutput[0] == true {
		cmd.Stdout = os.Stdout
		if len(liveOutput) <= 1 || liveOutput[1] == true {
			cmd.Stderr = os.Stderr
		}
		return []byte{}, cmd.Run()
	}
	return cmd.CombinedOutput()
}

/*
	Pipe allows you to pipe multiple bash commands

	[example (bash)]
		echo "test" | tee -a "./test.txt"

	[example (go)]
		bash.Pipe(".", []string{"echo", "test"}, []string{"tee", "-a", "./test.txt"})

		@dir: a directory to run the command in (set to an empty string to disable)
*/
func Pipe(dir string, args ...[]string){
	if len(args) == 1 {
		arg1 := args[0][0]
		args1 := args[0][1:]
		c1 := exec.Command(arg1, args1...)
		c1.Stdout = os.Stdout
	}

	cmd := []*exec.Cmd{}

	arg0 := args[0][0]
	args0 := args[0][1:]
	cmd = append(cmd, exec.Command(arg0, args0...))
	cmd[0].Stdin = os.Stdin
	if dir != "" {
		cmd[0].Dir = dir
	}

	if cmd[0].Env == nil {
		cmd[0].Env = os.Environ()
	}

	for i := 1; i < len(args); i++ {
		arg0 = args[i][0]
		args0 = args[i][1:]
		cmd = append(cmd, exec.Command(arg0, args0...))

		pr, pw := io.Pipe()
		cmd[i-1].Stdout = pw
		cmd[i].Stdin = pr
		if dir != "" {
			cmd[i].Dir = dir
		}

		if cmd[i].Env == nil {
			cmd[i].Env = os.Environ()
		}

		cmd[i-1].Start()

		go func(i int){
			defer pw.Close()

			cmd[i-1].Wait()
		}(i)
	}

	cmd[len(cmd)-1].Stdout = os.Stdout

	cmd[len(cmd)-1].Start()
	cmd[len(cmd)-1].Wait()
}

/*
	PipeMultiDir allows you to pipe multiple bash commands with a different directory for each of them

	note: the first arg is the directory

	[example]
		bash.PipeMultiDir([]string{"/dir1", "cat", "test.txt"}, []string{"./dir2", "tee", "-a", "./test.txt"})
*/
func PipeMultiDir(args ...[]string){
	if len(args) == 1 {
		arg1 := args[0][0]
		args1 := args[0][1:]
		c1 := exec.Command(arg1, args1...)
		c1.Stdout = os.Stdout
	}

	cmd := []*exec.Cmd{}

	dir := args[0][0]
	arg0 := args[0][1]
	args0 := args[0][2:]
	cmd = append(cmd, exec.Command(arg0, args0...))
	cmd[0].Stdin = os.Stdin
	if dir != "" {
		cmd[0].Dir = dir
	}

	if cmd[0].Env == nil {
		cmd[0].Env = os.Environ()
	}

	for i := 1; i < len(args); i++ {
		dir = args[i][0]
		arg0 = args[i][1]
		args0 = args[i][2:]
		cmd = append(cmd, exec.Command(arg0, args0...))

		pr, pw := io.Pipe()
		cmd[i-1].Stdout = pw
		cmd[i].Stdin = pr
		if dir != "" {
			cmd[i].Dir = dir
		}

		if cmd[i].Env == nil {
			cmd[i].Env = os.Environ()
		}

		cmd[i-1].Start()

		go func(i int){
			defer pw.Close()

			cmd[i-1].Wait()
		}(i)
	}

	cmd[len(cmd)-1].Stdout = os.Stdout

	cmd[len(cmd)-1].Start()
	cmd[len(cmd)-1].Wait()
}

/*
	PipeMultiDirEnv is just like the 'PipeMultiDir' method, but it also allows you to add custom envirronment vars

	note: the first arg is the directory

	[example]
		bash.PipeMultiDirEnv([]string{`MyEnvVar=CustomValue`}, []string{"/dir1", "cat", "test.txt"}, []string{"./dir2", "tee", "-a", "./test.txt"})
*/
func PipeMultiDirEnv(env []string, args ...[]string){
	if len(args) == 1 {
		arg1 := args[0][0]
		args1 := args[0][1:]
		c1 := exec.Command(arg1, args1...)
		c1.Stdout = os.Stdout
	}

	cmd := []*exec.Cmd{}

	dir := args[0][0]
	arg0 := args[0][1]
	args0 := args[0][2:]
	cmd = append(cmd, exec.Command(arg0, args0...))
	cmd[0].Stdin = os.Stdin
	if dir != "" {
		cmd[0].Dir = dir
	}

	if cmd[0].Env == nil {
		cmd[0].Env = os.Environ()
	}
	if env != nil {
		cmd[0].Env = append(cmd[0].Env, env...)
	}

	for i := 1; i < len(args); i++ {
		dir = args[i][0]
		arg0 = args[i][1]
		args0 = args[i][2:]
		cmd = append(cmd, exec.Command(arg0, args0...))

		pr, pw := io.Pipe()
		cmd[i-1].Stdout = pw
		cmd[i].Stdin = pr
		if dir != "" {
			cmd[i].Dir = dir
		}

		if cmd[i].Env == nil {
			cmd[i].Env = os.Environ()
		}
		if env != nil {
			cmd[i].Env = append(cmd[i].Env, env...)
		}

		cmd[i-1].Start()

		go func(i int){
			defer pw.Close()

			cmd[i-1].Wait()
		}(i)
	}

	cmd[len(cmd)-1].Stdout = os.Stdout

	cmd[len(cmd)-1].Start()
	cmd[len(cmd)-1].Wait()
}
