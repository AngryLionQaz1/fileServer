package main

import (
	"github.com/kardianos/service"
	"fileServer/utils"
	"fmt"
	"net/http"
	"os"
)

//配置信息
var program *utils.Program
//注册系统服务
var server service.Service
//http服务
var sHttp *http.Server
var err error
//优雅退出
var quit = make(chan os.Signal)

func main() {

    start()
	run()
}

func init()  {
	//读取配置信息
	program= utils.InitProgram()
	//初始化退出信号
	program.Quit=quit
	//初始化服务
	config := utils.InitService(program)
	server, err= service.New(program, config)
	utils.CheckErr(err)
}

func start()  {
	utils.HttpServer(program)
}

func run()  {

	if len(os.Args) < 2 {
		err = server.Run()
		utils.CheckErr(err)
		return
	}
	cmd := os.Args[1]
	if cmd == "install" {
		err = server.Install()
		utils.CheckErr(err)
		fmt.Println("安装成功")
	}
	if cmd == "uninstall" {
		err = server.Uninstall()
		utils.CheckErr(err)
		fmt.Println("卸载成功")
	}
}



