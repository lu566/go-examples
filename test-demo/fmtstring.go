package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadLine(fileName string, handler func(string)) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func Print(line string) {
	message := LogFormatting(line)
	if message == "" {
		return
	}
	fmt.Println(LogFormatting(line))
}

func LogFormatting(message string)string{

	if strings.Contains(message,"in namespace") || strings.Contains(message,"Server closed") {
		return ""
	}

	if strings.Contains(message,`{"stream"`) {
		if strings.Contains(message,"u003e") {
			return ""
		}
		message = strings.Replace(message,`{"stream":"`,"",-1)
		message = strings.Replace(message,`\n"}`,"",-1)
		message = strings.Replace(message,`\"`,"",-1)

	}

	if strings.Contains(message,`{"status"`) {
		if strings.Contains(message,"u003e") {
			return ""
		}
		message = strings.Replace(message,`"`,"",-1)
		message = strings.Replace(message,`{`,"",-1)
		message = strings.Replace(message,`}`,"",-1)
	}

	return message
}

func main() {
	ReadLine("/Users/wang/go/src/github.com/lu566/go-examples/save/20181112/file04", Print)

}