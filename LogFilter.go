package http

import (
	"github.com/gobasis/log"
	"github.com/gobasis/utils"
)

const requestUid = "requestUid"

/*
Description:
logBefore write a log before a request handler
 * Author: architect.bian
 * Date: 2018/10/23 13:32
 */
func logBefore(ctx *Context) {
	ctx.Input.Data[requestUid] = utils.UUID.Get()
	log.Info("A new requesting", "host", ctx.Request.Host, "method", ctx.Request.Method,
		"url", ctx.Request.RequestURI, "remoteAddr", ctx.Request.RemoteAddr, requestUid, ctx.Input.Data[requestUid])
}

/*
Description:
logAfter write a log after a request handler
 * Author: architect.bian
 * Date: 2018/10/23 13:32
 */
func logAfter(ctx *Context) {
	log.Info("request finished", requestUid, ctx.Input.Data[requestUid])
}

/*
Description: 
initialize before application start up
 * Author: architect.bian
 * Date: 2018/10/23 13:33
 */
func init() {
	Router.InitFilter(".*", logBefore, logAfter)
}
