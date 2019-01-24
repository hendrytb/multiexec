package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	apps := make([]*App, 0)

	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error reading config.json:", err)
	}

	err = json.Unmarshal(b, &apps)
	if err != nil {
		fmt.Println("JSON error:", err)
	}

	c := make(chan AppResult)
	for i := range apps {
		go apps[i].Run(c)
	}

	sigs := make(chan os.Signal, 1)
	// done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case r := <-c:
		if r.Error != nil {
			fmt.Printf("App [%s] Error: %s\n", r.App.Name, r.Error)
		}
		for i := range apps {
			if apps[i] == r.App {
				continue
			}
			fmt.Printf("Stopping [%s]...\n", apps[i].Name)
			apps[i].Stop()
		}
	case <-sigs:
		fmt.Println("Termination signal received")
		for i := range apps {
			fmt.Printf("Stopping [%s]...\n", apps[i].Name)
			apps[i].Stop()
		}
	}
	fmt.Println("Process done.")
}

type App struct {
	Name      string   `json:"name"`
	Path      string   `json:"path"`
	Arguments []string `json:"args"`
	cmd       *exec.Cmd
}

type AppResult struct {
	App   *App
	Error error
}

func (a *App) Run(c chan AppResult) {
	a.cmd = exec.Command(a.Path, a.Arguments...)
	a.cmd.Stdout = Decorator(a.Name, os.Stdout)
	a.cmd.Stderr = Decorator(a.Name, os.Stderr)
	err := a.cmd.Run()
	if err != nil {
		_, _ = a.cmd.Stderr.Write([]byte(err.Error() + "\n"))
	}
	c <- AppResult{a, err}
}

func (a *App) Stop() {
	err := a.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		_, _ = a.cmd.Stderr.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
	}
	_, _ = a.cmd.Process.Wait()
}

func Decorator(name string, w io.Writer) io.Writer {
	return &decorator{name, w}
}

type decorator struct {
	name string
	w    io.Writer
}

func (d *decorator) Write(p []byte) (n int, err error) {
	_, _ = d.w.Write([]byte("[" + d.name + "] "))
	return d.w.Write(p)
}
