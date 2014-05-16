package dplug

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

type MethodHandler interface {
	Run(Parameters, *Results) error
}

type MethodHandlerFunc func(Parameters, *Results) error

func (f MethodHandlerFunc) Run(p Parameters, r *Results) error {
	return f(p, r)
}

type DPlug struct { //TODO EFF THIS
	Server *DPlugServer
}

type DPlugServer struct {
	Self    Plugin
	Methods map[string]MethodHandler
}

func (gps *DPlugServer) RegisterMethod(name string, handler MethodHandler) {
	if _, ok := gps.Methods[name]; ok {
		panic("Method name '" + name + "' already exists in plugin")
	}
	gps.Methods[name] = handler
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
	return handler.Run(p.Params, r)
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
