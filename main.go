package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func init() {
	ioutil.WriteFile("./output.txt", []byte(""), 0644)
}

func main() {
	path := Pwd()
	goModBytes, err := ioutil.ReadFile(path + "go.mod")
	if err != nil {
		panic(err)
	}

	executableName := strings.Split(strings.Split(string(goModBytes), "\n")[0], " ")[1]

	var cmd *exec.Cmd
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd = exec.Command("go", "build")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	compileStart := time.Now()
	if err := cmd.Run(); err != nil {
		panic(stderr.String())
	}
	fmt.Println(out.String())
	fmt.Println("go build took " + time.Since(compileStart).String())

	cmd = exec.Command("./" + executableName)

	var streamStdout, streamStderr []byte
	var streamErrStdout, streamErrStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		streamStdout, streamErrStdout = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	streamStderr, streamErrStderr = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if streamErrStdout != nil || streamErrStderr != nil {
		log.Fatal("gorun -> failed to capture stdout or stderr\n")
	}
	_, _ = string(streamStdout), string(streamStderr)
	// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)

}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			original, _ := ioutil.ReadFile("./output.txt")
			newStr := string(original) + "\n" + string(out)
			ioutil.WriteFile("./output.txt", []byte(newStr), 0644)
			return out, err
		}
	}
}
func Pwd() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	return dir + "/"
}
