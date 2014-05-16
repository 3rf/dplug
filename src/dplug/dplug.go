package dplug

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var session *DPlugSession

type DPlugSession struct {
	plugins *map[string]Plugin
	rwlock  sync.RWMutex
}

func Init(conf Config) error {
	session = &DPlugSession{
		plugins: &map[string]Plugin{},
	}
	for _, c := range conf.PluginConfigs {
		plugin, err := startPluginFromConfig(c.Path, c.Port)
		if err != nil {
			return fmt.Errorf("dplug: %v", err)
		}
		(*session.plugins)[plugin.Name] = *plugin
	}
	return nil
}

func ShutDown() error {
	failed := 0
	for _, plugin := range *session.plugins {
		err := plugin.Process.Kill()
		if err != nil {
			fmt.Printf("ERROR TERMINATING PLUGIN [%v]", plugin.Name) //TODO make a list?
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("failed to terminate %v plugins", failed)
	}
	return nil
}

func startPluginFromConfig(path string, port int) (*Plugin, error) {
	cmd := exec.Command(path, "-port", strconv.Itoa(port))
	err := cmd.Start() //non-blocking
	if err != nil {
		return nil, fmt.Errorf("error starting '%v' on port %v: %v", path, port, err)
	}

	time.Sleep(100 * time.Millisecond) //TODO FIXME??? CONFIG????

	name, err := getPluginNameFromPort(port)
	if err != nil {
		return nil, fmt.Errorf("error getting plugin name on port %v: %v", port, err)
	}
	methods, err := getPluginMethodsFromPort(port)
	if err != nil {
		return nil, fmt.Errorf("error getting plugin '%v' methods on port %v: %v", name, port, err)
	}

	return &Plugin{name, methods, port, cmd.Process}, nil
}

type PluginRoute struct {
	Path string
	Port int
}

type Config struct {
	PluginConfigs []PluginRoute
}

type Plugin struct {
	Name        string
	MethodNames []string
	Port        int
	Process     *os.Process
}

type Method struct {
	plugin *Plugin
	name   string
}

type Parameters map[string]interface{}
type Results map[string]interface{}

type ExternalParameters struct {
	MethodName string
	Params     Parameters
}

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

func getPluginNameFromPort(port int) (string, error) {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return "", fmt.Errorf("error connecting: %", err)
	}

	name := ""
	err = client.Call(
		"DPlug.Name",
		0,
		&name,
	)

	if err != nil {
		return "", fmt.Errorf("error getting plugin name: %v", err)
	}
	if name == "" {
		return "", fmt.Errorf("name not set for plugin")
	}

	return name, nil
}

func getPluginMethodsFromPort(port int) ([]string, error) {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("error connecting: %", err)
	}

	methods := []string{}
	err = client.Call(
		"DPlug.Methods",
		0,
		&methods,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting plugin methods: %v", err)
	}
	if methods == nil {
		return nil, fmt.Errorf("no methods registered for plugin")
	}

	return methods, nil
}

func CallPluginMethod(pluginName, methodName string, p Parameters, r *Results) error {
	//TODO error handling
	if session == nil {
		return fmt.Errorf("Must initialize DPlug session before calling plugins")
	}
	if session.plugins == nil {
		return fmt.Errorf("Must initialize DPlug session before calling plugins")
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
		"DPlug.HandleMethod",
		ExternalParameters{
			methodName,
			p,
		},
		r,
	)
	if err != nil {
		return fmt.Errorf("plugin error: %v", err)
	}

	return nil
}
