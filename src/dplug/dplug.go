package goplug

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
	"sync"
)

var session *GoPlugSession

type GoPlugSession struct {
	plugins *map[string]Plugin
	rwlock  sync.RWMutex
}

func TestInitSession() error { //TODO
	session = &GoPlugSession{
		plugins: &map[string]Plugin{
			"test": Plugin{
				Name:        "test",
				MethodNames: []string{"doit"},
				Port:        1234,
			},
		},
	}
	return nil
}

type Plugin struct {
	Name        string
	MethodNames []string
	Port        int
}

type Method struct {
	plugin *Plugin
	name   string
}

type Parameters map[string]interface{}
type Results map[string]interface{}

type ExternalParameters struct {
	MethodName string
	Params     interface{}
}

type MethodHandler interface {
	Run(Parameters, *Results) error
}

type MethodHandlerFunc func(Parameters, *Results) error

func (f MethodHandlerFunc) Run(p Parameters, r *Results) error {
	return f(p, r)
}

type GoPlug struct { //TODO EFF THIS
    Server *GoPlugServer
}

type GoPlugServer struct {
	Self    Plugin
	Methods map[string]MethodHandler
}

func (gps *GoPlugServer) RegisterMethod(name string, handler MethodHandler) {
	if _, ok := gps.Methods[name]; ok {
		panic("Method name '" + name + "' already exists in plugin")
	}
	gps.Methods[name] = handler
}

func (gp *GoPlug) HandleMethod(p ExternalParameters, r *Results) error {
	handler, ok := gp.Server.Methods[p.MethodName]
	if !ok {
		return fmt.Errorf(
			"Method '%v' is not registered for plugin '%v'",
			p.MethodName,
			gp.Server.Self.Name,
		)
	}
	return handler.Run(p.Params, r)
}

func (gp *GoPlug) ListMethods(_ interface{}, methods *[]string) error {
	for name, _ := range gp.Server.Methods {
		*methods = append(*methods, name)
	}
	return nil
}

func (gps *GoPlugServer) Serve() error {
    gp := &GoPlug{gps}
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

func CallPluginMethod(pluginName, methodName string, p Parameters, r *Results) error {
	//TODO error handling
	if session == nil {
		return fmt.Errorf("Must initialize GoPlug session before calling plugins")
	}
	if session.plugins == nil {
		return fmt.Errorf("Must initialize GoPlug session before calling plugins")
		//TODO panic?
	}

	plugin, ok := (*session.plugins)[pluginName]
    if !ok {
        return fmt.Errorf("plugin '%v' does not exist in session", pluginName)
    }

    client, err := rpc.DialHTTP("tcp", "127.0.0.1:"+strconv.Itoa(plugin.Port))
	if err != nil {
        return fmt.Errorf("error connecting: %", err)
	}

	// Synchronous call
	err = client.Call(
		"GoPlug.HandleMethod",
		ExternalParameters{
			methodName,
			p,
		},
		r,
	)
	if err != nil {
		return fmt.Errorf("plugin error:", err)
	}

	return nil
}