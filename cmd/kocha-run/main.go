package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/woremacx/kocha/util"
	"github.com/naoina/miyabi"
	"gopkg.in/fsnotify.v1"
)

type runCommand struct {
	option struct {
		Help bool `short:"h" long:"help"`
	}
}

func (c *runCommand) Name() string {
	return "kocha run"
}

func (c *runCommand) Usage() string {
	return fmt.Sprintf(`Usage: %s [OPTIONS]

Run the your application.

Options:
    -h, --help        display this help and exit

`, c.Name())
}

func (c *runCommand) Option() interface{} {
	return &c.option
}

func (c *runCommand) Run(args []string) error {
	basedir, err := os.Getwd()
	if err != nil {
		return err
	}
	execName := filepath.Base(basedir)
	if runtime.GOOS == "windows" {
		execName += ".exe"
	}
	if err := util.PrintEnv(); err != nil {
		return err
	}
	fmt.Println("Starting...")
	var cmd *exec.Cmd
	for {
		if cmd != nil {
			if err := cmd.Process.Signal(miyabi.ShutdownSignal); err != nil {
				cmd.Process.Kill()
			}
			if err := cmd.Wait(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		newCmd, err := runApp(basedir, execName)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		fmt.Println()
		cmd = newCmd
		if err := watchApp(basedir, execName); err != nil {
			if err := cmd.Process.Signal(miyabi.ShutdownSignal); err != nil {
				cmd.Process.Kill()
			}
			return err
		}
		fmt.Println("\nRestarting...")
	}
}

func runApp(basedir, execName string) (*exec.Cmd, error) {
	execArgs := []string{"build", "-o", execName}
	if runtime.GOARCH == "amd64" {
		execArgs = append(execArgs, "-race")
	}
	c, err := execCmd("go", execArgs...)
	if err != nil {
		return nil, err
	}
	if err := c.Wait(); err != nil {
		c.Process.Kill()
		return nil, err
	}
	c, err = execCmd(filepath.Join(basedir, execName))
	if err != nil {
		c.Process.Kill()
	}
	return c, err
}

func watchApp(basedir, execName string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()
	watchFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name()[0] == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if err := watcher.Add(path); err != nil {
			return err
		}
		return nil
	}
	for _, path := range []string{
		"app", "config", "main.go",
	} {
		if err := filepath.Walk(filepath.Join(basedir, path), watchFunc); err != nil {
			return err
		}
	}
	select {
	case <-watcher.Events:
	case err := <-watcher.Errors:
		return err
	}
	return nil
}

func execCmd(name string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func main() {
	util.RunCommand(&runCommand{})
}
