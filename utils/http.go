package utils

import (
	"bytes"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var fileServer http.Handler

func HttpServer(program *Program) {

	signal.Notify(program.Quit, os.Interrupt)
	mux := http.NewServeMux()
	fileServer = http.FileServer(http.Dir(program.Path))
	mux.Handle("dir", fileServer)
	mux.HandleFunc("/", sx)
	mux.Handle("/"+program.Upload, uploadHandler(program))
	program.Server = &http.Server{
		Addr:         ":" + strconv.Itoa(program.Port),
		WriteTimeout: time.Second * 10,
		Handler:      mux,
	}
	go func() {
		//接受退出信号
		<-program.Quit
		if err := program.Server.Close(); err != nil {
			CheckErr(err)
		}
	}()
	err := program.Server.ListenAndServe()
	if err != nil {
		// 正常退出
		if err == http.ErrServerClosed {
			CheckErr(err)
		} else {
			CheckErr(err)
		}
		log.Fatal("Server exited")
		Logs("Server exited")
	}

}

func sx(writer http.ResponseWriter, request *http.Request) {
	if filter(request.URL) {
		writer.WriteHeader(http.StatusNotFound)
		io.WriteString(writer, "404 page not found!\n")
		return
	}
	fileServer.ServeHTTP(writer, request)

}

func uploadHandler(program *Program) http.Handler {

	maxUploadSize := int64(program.Size * 1024 * 1024) //M

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		request.Body = http.MaxBytesReader(writer, request.Body, maxUploadSize)
		if err := request.ParseMultipartForm(maxUploadSize); err != nil {
			renderError(writer, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		request.ParseMultipartForm(-1)
		files := request.MultipartForm.File["file"]

		for _, v := range files {
			if CheckFileType(program.Types, path.Ext(v.Filename)) {
				renderError(writer, "INVALID_FILE_TYPE", http.StatusBadRequest)
				break
			}

		}

		file := SaveFile(request.MultipartForm.File, program)

		writer.Write([]byte(file))
	})

}

//保存数据到磁盘
func SaveFile(headers map[string][]*multipart.FileHeader, p *Program) string {
	var buffer bytes.Buffer
	files := headers["file"]
	dir, dir2 := CreateDateDir(p.Path)
	v4 := strings.Replace(uuid.NewV4().String(), "-", "", -1)
	for k, v := range files {
		str := v4 + strconv.Itoa(k) + path.Ext(v.Filename)
		str2 := filepath.Join(dir2, str)
		newPath := filepath.Join(dir, str)
		go func(v *multipart.FileHeader) {
			f, _ := v.Open()
			defer f.Close()
			fileBytes, _ := ioutil.ReadAll(f)
			newFile, _ := os.Create(newPath)
			newFile.Write(fileBytes)

		}(v)
		buffer.WriteString(str2)
		if k != len(files)-1 {
			buffer.WriteString(",")
		}
	}
	return buffer.String()
}

//拦截文件夹
func filter(url *url.URL) bool {
	s := url.Path
	split := strings.Split(s, ".")
	if len(split) == 2 {
		return false
	}
	return true

}

//获取文件类型
func CheckFileType(s, s2 string) bool {

	return strings.Contains(s, strings.Replace(s2, ".", "", -1))

}

func renderError(writer http.ResponseWriter, s string, i int) {
	writer.WriteHeader(i)
	io.WriteString(writer, s)
}

// CreateDateDir 根据当前日期来创建文件夹
func CreateDateDir(basePath string) (string, string) {
	folderName := time.Now().Format("20060102")
	folderPath := filepath.Join(basePath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步
		// 先创建文件夹
		os.Mkdir(folderPath, 0777)
		// 再修改权限
		os.Chmod(folderPath, 0777)
	}
	return folderPath, folderName
}
