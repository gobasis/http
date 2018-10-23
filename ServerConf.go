package http

import (
	"github.com/BurntSushi/toml"
	"github.com/gobasis/log"
)

/*
Description: 
server config file
 * Author: architect.bian
 * Date: 2018/10/22 19:11
 */
type ServerConfType struct {
	Address string `toml:"address"`
	Port    int32
	Public string `toml:"public"`
}

var ServerConf ServerConfType

/*
Description:
initialize parse conf file before application start up
 * Author: architect.bian
 * Date: 2018/10/23 13:39
 */
func init() {
	_, err := toml.DecodeFile("conf/server.conf", &ServerConf)
	if err != nil {
		log.Fatal("server conf file parse failed", "error", err)
	}
}