package utils

import (
	"fmt"
	"github.com/kardianos/service"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const initFile = `E:\GoProjects\fileServer\main\init.yaml`

type Program struct {
	//服务名
	Name string `yaml:name`
	//服务器显示
	DisplayName string `yaml:displayName`
	//服务描述
	Description string `yaml:description`
	//路径
	Path string `yaml:path`
	//端口号
	Port int `yaml:port`
	//文件大小
	Size int `yaml:size`
	//上传地址
	Upload string `yaml:upload`
	//限制哪些文件不能上传
	Types string `yaml:types`
	//服务
	Server *http.Server
	//退出
	Quit chan os.Signal

}


func (p *Program) run() {

	fmt.Println("sssssssssssssssssssssss")

}

func (p *Program) Start(s service.Service) error {
	log.Println("开始服务")
	Logs("开始服务")
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	log.Println("停止服务")
	Logs("停止服务")
	err := p.Server.Shutdown(nil)
	if err != nil {
		log.Fatal([]byte("shutdown the server err"))
	}
	return nil
}

//初始化服务
func InitService(p *Program) *service.Config {

	var serviceConfig = &service.Config{
		Name:        p.Name,
		DisplayName: p.DisplayName,
		Description: p.Description,
	}

	return serviceConfig
}


//读取配置文件
func InitProgram() *Program {
	i := new(Program)
	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//yamlFile, err := ioutil.ReadFile(filepath.Join(dir, initFile))
	yamlFile, err := ioutil.ReadFile(initFile)
	CheckErr(err)
	err = yaml.Unmarshal(yamlFile, i)
	CheckErr(err)
	return i
}
