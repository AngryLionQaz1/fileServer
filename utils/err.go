package utils

import (
	"os"
	"path/filepath"
	"time"
)

func CheckErr(err error) {
	if err != nil {
		Logs(err.Error())
	}
}

//打印日志
func Logs(s string) {

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fd, _ := os.OpenFile(filepath.Join(dir, "logs.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd.WriteString(time.Now().Format("2006-01-02 15:04:05") + ":" + s + "\n")
	defer fd.Close()

}
