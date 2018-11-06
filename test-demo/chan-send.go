package main

import (
	"bytes"
	"fmt"
	"html/template"
)

func main() {
	var tmpl = template.Must(template.New("").Parse(`{{.Person}}{{.Timeframe}}`))

	buf := new(bytes.Buffer)
	tmpl.Execute(buf, map[string]interface{}{
		"Person":    "Bob",
		"Timeframe": 1,
	})

	fmt.Println(buf.String(), "=====")

}