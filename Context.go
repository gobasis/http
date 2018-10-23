package http

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)

/*
Description:
Context Http request context struct including Input, Output, http.Request and http.ResponseWriter.
Input and Output provides some api to operate request and response more easily.
 * Author: architect.bian
 * Date: 2018/10/23 12:15
 */
type Context struct {
	Input          *Input
	Output         *Output
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

/*
Description:
Input operates the http request header, Data, cookie and body.
it also contains router PathParams and current session.
 * Author: architect.bian
 * Date: 2018/10/23 12:16
 */
type Input struct {
	PathParams httprouter.Params
	Data       map[interface{}]interface{} // store some values in this context when calling context in filter or controller.
}

/*
Description:
Output does work for sending response header.
 * Author: architect.bian
 * Date: 2018/10/23 12:19
 */
type Output struct {
	Status     int
	EnableGzip bool
}