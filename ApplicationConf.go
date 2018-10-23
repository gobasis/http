package http

import (
	"path/filepath"
	"github.com/BurntSushi/toml"
	"github.com/gobasis/log"
)

var ApplicationConf ApplicationConfType

type ApplicationConfType struct {
	Title   string `toml:"title"`
	Base    string
	BaseCss string
	BaseJs  string
	BaseImg string
	BaseFile string
	BaseApi string
	UserImg string
	WebImg string
	Domain  string

	//API APIType
	Sys SysType

	initData map[interface{}]interface{}
}


type APIType struct {
	Login string
	Logout string
}

type SysType struct {
	RootDir     string
	DoPagesPath string
	LoginToken string
}

/*
Description:
initialize parse conf file before application start up
 * Author: architect.bian
 * Date: 2018/10/23 13:38
 */
func init() {
	ApplicationConf.Sys.RootDir, _ = filepath.Abs("./")
	_, err := toml.DecodeFile("conf/application.conf", &ApplicationConf)
	if err != nil {
		log.Fatal("application conf file parse failed", "error", err)
	}
}