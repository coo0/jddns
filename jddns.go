package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/kardianos/service"
	"jddns/api"
	"jddns/lib/install"
	"os"
	"runtime"
)

var confPath = flag.String("c", "/etc/jddns/conf/config.yml", "配置文件路径")

var conf struct {
	IpURL     string     `yaml:"ip-url"`
	UpdateApi []string   `yaml:"update-api"`
	Hostname  [][]string `yaml:"hostname"`
	Token     []string   `yaml:"Token"`
	Cron      string     `yaml:"cron"`
}

func main() {
	NCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(NCPU)
	serConfig := &service.Config{
		Name:        "jddns",
		DisplayName: "jddns",
		Description: "jddns 开机自启动",
	}
	logs.Reset()
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logPath := "/var/log/jddns.log"
	level := "7"
	_ = logs.SetLogger("file", `{"level":`+level+`,"filename":"`+logPath+`","daily":false,"maxlines":100000,"color":true}`)
	pro := &Program{}
	s, err := service.New(pro, serConfig)
	if err != nil {
		fmt.Println(err, "service.New() err")
	}
	if len(os.Args) > 1 {
		fmt.Println("actions: ", os.Args[1])
		switch os.Args[1] {
		case "install":
			_ = service.Control(s, "stop")
			_ = service.Control(s, "uninstall")
			binPath := install.Install()
			serConfig.Executable = binPath
			s, err := service.New(pro, serConfig)
			if err != nil {
				logs.Error(err)
				return
			}
			err = s.Install()
			if err != nil {
				fmt.Println("install err", err)
				logs.Error("install err", err)
			} else {
				logs.Info("install success")
				fmt.Println("install success")
			}
			return
		case "start", "restart", "stop":
			err := service.Control(s, os.Args[1])
			if err != nil {
				logs.Error("Valid actions: %q\n%s", service.ControlAction, err.Error())
			}
			return
		case "uninstall":
			install.Uninstall("jddns")
			err = s.Uninstall()
			if err != nil {
				fmt.Println("Uninstall  err", err)
				logs.Error("Uninstall  err", err)
			} else {
				logs.Info("Uninstall  success")
				fmt.Println("Uninstall  success")
			}
			return
		}
	}

	err = s.Run() // 运行服务
	if err != nil {
		logs.Error("s.Run err", err)
		fmt.Println("s.Run err", err)
	}
}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	_, _ = s.Status()
	fmt.Println("server start")
	logs.Info("server start")
	go p.run()
	return nil
}
func (p *Program) Stop(s service.Service) error {
	_, _ = s.Status()
	fmt.Println("server stop")
	logs.Info("server stop")
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *Program) run() {
	c := cron.New()
	update()
	_, err := c.AddFunc(conf.Cron, update)
	if err != nil {
		logs.Error("启动计划任务失败：", err)
		return
	} else {
		logs.Info("已启动计划任务")
	}
	c.Start()

	fmt.Println("开机自启动服务")
	logs.Info("开机自启动服务")
}

func readConf(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		logs.Error("in file %q: %v", filename, err)
		return err
	}
	return nil
}

func update() {
	err := readConf(*confPath)
	if err != nil {
		logs.Error("更新配置文件失败：", err)
		return
	}
	logs.Info("启动更新")
	api.DomainNameBinding(conf.IpURL, conf.UpdateApi, conf.Hostname, conf.Token)
}
