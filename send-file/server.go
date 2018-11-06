package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/upload", uploadHandle)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {

	// 根据字段名获取表单文件
	formFile, _, err := r.FormFile("uploadfile")
	if err != nil {
		log.Printf("Get form file failed: %s\n", err)
		return
	}
	defer formFile.Close()

	// 创建保存文件
	destFile, err := os.Create("/Users/wang/go/src/github.com/lu566/go-examples/send-file/upload/t.jpg")//"." + r.URL.Path + "/" + header.Filename)
	if err != nil {
		log.Printf("Create failed: %s\n", err)
		return
	}
	defer destFile.Close()

	// 读取表单文件，写入保存文件
	_, err = io.Copy(destFile, formFile)
	if err != nil {
		log.Printf("Write file failed: %s\n", err)
		return
	}
}
