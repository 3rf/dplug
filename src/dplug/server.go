package dplug

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"reflect"
	"strconv"
)

var port *int

func init() {
	port = flag.Int("dplugport", 0, "port to run plugin on")
}

type DPlug struct { //TODO EFF THIS
	Server *DPlugServer
}

type DPlugServer struct {
	Self    Plugin
	Methods map[string]reflect.Value
}

func StartDPlugServer(name string) *DPlugServer {
	if !flag.Parsed() {
		flag.Parse()
	}
	dps := DPlugServer{
		Self: Plugin{
			Name: name,
			Port: *port},
		Methods: map[string]reflect.Value{},
	}
	return &dps
}

func validateHandler(hType reflect.Type) {
	if hType.Kind() == reflect.Func {
		if hType.NumIn() == 2 {
			if hType.In(1).Kind() == reflect.Ptr {
				if hType.NumOut() == 1 {
					var err error
					if hType.Out(0).Implements(reflect.TypeOf(err)) {
						return
					} else {
						panic("handler function must return error")
					}
				} else {
					panic("handler function must return one value")
				}
			} else {
				panic("handler function must take a reference as the second parameter")
			}
		} else {
			panic("handler function must take 2 input parameters")
		}
	} else {
		panic("handler must be Func")
	}
}

func (gps *DPlugServer) RegisterMethod(name string, handler interface{}) {
	if _, ok := gps.Methods[name]; ok {
		panic("Method name '" + name + "' already exists in plugin")
	}

	validateHandler(reflect.TypeOf(handler))
	gps.Methods[name] = reflect.ValueOf(handler)
}

func (gp *DPlug) HandleMethod(p ExternalParameters, r *Results) error {
	handler, ok := gp.Server.Methods[p.MethodName]
	if !ok {
		return fmt.Errorf(
			"Method '%v' is not registered for plugin '%v'",
			p.MethodName,
			gp.Server.Self.Name,
		)
	}
	retVals := handler.Call([]reflect.Value{
		reflect.ValueOf(p.Params),
		reflect.ValueOf(r),
	})
	err := retVals[0]
	if err.IsNil() {
		return nil
	} else {
		return err.Interface().(error)
	}
}

func (gp *DPlug) Methods(_ int, methods *[]string) error {
	for name, _ := range gp.Server.Methods {
		*methods = append(*methods, name)
	}
	return nil
}

func (dp *DPlug) Name(_ int, name *string) error {
	*name = dp.Server.Self.Name
	return nil
}

func (gps *DPlugServer) Serve() error {
	gp := &DPlug{gps}
	rpc.Register(gp)
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(gps.Self.Port))
	if err != nil {
		return fmt.Errorf("listener error:", err)
	}
	http.Serve(listener, nil)
	// Serve blocks, this return should not be reached
	return nil
}
