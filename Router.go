package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/gobasis/log"
	"net/http"
	"regexp"
	"reflect"
	"encoding/json"
	"compress/gzip"
	"strings"

)

/*
Description: 
RouterType wrap with httprouter.RouterType, and some defined-self methods.
 * Author: architect.bian
 * Date: 2018/10/23 11:00
 */
type RouterType struct {
	rawRouter *httprouter.Router
	filters map[string]filterHandler
}

/*
Description:
Handler a handle function for the request, accept a param context
 * Author: architect.bian
 * Date: 2018/10/23 12:26
 */
type Handler func(*Context) interface{}

/*
Description:
Filter A filter function for the request
 * Author: architect.bian
 * Date: 2018/10/23 13:52
 */
type Filter func(*Context)

/*
Description:
filterHandler a filter type with before or after handler
 * Author: architect.bian
 * Date: 2018/10/23 11:38
 */
type filterHandler struct {
	before Filter
	after  Filter
}

var Router = RouterType{
	rawRouter: httprouter.New(),
	filters: make(map[string]filterHandler),
}

func (r *RouterType) GetRawRouter() *httprouter.Router {
	return r.rawRouter
}

/*
Description:
mapping a handler with get method
 * Author: architect.bian
 * Date: 2018/10/22 20:27
 */
func (r *RouterType) InitGet(path string, handler Handler) {
	r.rawRouter.GET(path, proxyFun(handler))
}

/*
Description:
mapping a handler with HEAD method
 * Author: architect.bian
 * Date: 2018/10/22 20:31
 */
func (r *RouterType) InitHEAD(path string, handler Handler) {
	r.rawRouter.GET(path, proxyFun(handler))
}

/*
Description:
mapping a handler with OPTIONS method
 * Author: architect.bian
 * Date: 2018/10/22 20:31
 */
func (r *RouterType) InitOPTIONS(path string, handler Handler) {
	r.rawRouter.OPTIONS(path, proxyFun(handler))
}

/*
Description:
mapping a handler with POST method
 * Author: architect.bian
 * Date: 2018/10/22 20:31
 */
func (r *RouterType) InitPOST(path string, handler Handler) {
	r.rawRouter.POST(path, proxyFun(handler))
}

/*
Description:
mapping a handler with PUT method
 * Author: architect.bian
 * Date: 2018/10/22 20:31
 */
func (r *RouterType) InitPUT(path string, handler Handler) {
	r.rawRouter.PUT(path, proxyFun(handler))
}

/*
Description:
mapping a handler with PATCH method
 * Author: architect.bian
 * Date: 2018/10/22 20:31
 */
func (r *RouterType) InitPATCH(path string, handler Handler) {
	r.rawRouter.PATCH(path, proxyFun(handler))
}

/*
Description:
mapping a handler with DELETE method
 * Author: architect.bian
 * Date: 2018/10/22 20:31
 */
func (r *RouterType) InitDELETE(path string, handler Handler) {
	r.rawRouter.DELETE(path, proxyFun(handler))
}

/*
Description:
mapping a handler with special method
 * Author: architect.bian
 * Date: 2018/10/22 20:31
 */
func (r *RouterType) InitHandler(method, path string, handler Handler) {
	r.rawRouter.Handle(method, path, proxyFun(handler))
}

/*
Description:
filter some request as interceptor
 * Author: architect.bian
 * Date: 2018/10/23 10:31
 */
func (r *RouterType) InitFilter(pattern string, handlerBefore Filter, handlerAfter Filter) {
	r.filters[pattern] = filterHandler {
		before: handlerBefore,
		after:  handlerAfter,
	}
}

/*
Description:
proxyFun a proxy for any handler
 * Author: architect.bian
 * Date: 2018/10/23 11:38
 */
func proxyFun(handler Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var filter filterHandler
		for key, value := range Router.filters {
			ok, err := regexp.MatchString(key, r.URL.RawPath)
			if err == nil && ok {
				filter = value
				break
			}
		}
		var ctxNew = &Context{
			Input: &Input{
				PathParams: ps,
				Data:       make(map[interface{}]interface{}),
			},
			Output: &Output{
				Status: 200,
				EnableGzip: *ServerConf.EnableGzip,
			},
			Request: r,
			ResponseWriter: w,
		}
		if filter.before != nil {
			filter.before(ctxNew) //invoke before filter
		}
		result := handler(ctxNew) //invoke handler
		if result != nil {
			var bytes []byte
			if reflect.TypeOf(result).String() == "string" {
				w.Header().Set("Content-Type", "text/html")
				bytes = []byte(result.(string)) //response raw result
			} else {
				bs, err := json.Marshal(result)
				if err != nil {
					log.Error("json marshal failed", "result", result, "error", err)
					w.WriteHeader(500)
				} else {
					w.Header().Set("Content-Type", "application/json")
					bytes = bs
				}
			}
			if len(bytes) > 0 {
				if ctxNew.Output.EnableGzip && strings.Contains(ctxNew.Request.Header.Get("Accept-Encoding"), "gzip") {
					w.Header().Set("Content-Encoding", "gzip")
					gzipWriter := gzip.NewWriter(w)
					gzipWriter.Write(bytes)
					defer gzipWriter.Close()
				} else {
					w.Write(bytes)
				}
			}
		}
		if filter.after != nil {
			filter.after(ctxNew) //invoke after filter
		}
	}
}