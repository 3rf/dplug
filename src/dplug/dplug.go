package dplug

import (
	"fmt"
	"net/rpc"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var session *DPlugSession

type DPlugSession struct {
	plugins *map[string]Plugin
	rwLock  sync.RWMutex
}

func Initialize(conf Config) error {
	session = &DPlugSession{
		plugins: &map[string]Plugin{},
	}

	session.rwLock.Lock()
	defer session.rwLock.Unlock()

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

	session.rwLock.Lock()
	defer session.rwLock.Unlock()

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
	cmd := exec.Command(path, "-dplugport", strconv.Itoa(port))
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

type ExternalResults struct {
	Results interface{}
}
type ExternalParameters struct {
	MethodName string
	Params     interface{}
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

func PluginMethods(pluginName string) ([]string, error) {
	if session == nil {
		return nil, fmt.Errorf("Must initialize DPlug session before looking up plugin methods")
	}

	session.rwLock.RLock()
	defer session.rwLock.RUnlock()

	plugin, ok := (*session.plugins)[pluginName]
	if !ok {
		return nil, fmt.Errorf("plugin '%' not registered", pluginName)
	}
	return plugin.MethodNames, nil
}

func PluginsWithMethod(methodName string) ([]string, error) {
	if session == nil {
		return nil, fmt.Errorf("Must initialize DPlug session before looking up plugin methods")
	}

	session.rwLock.RLock()
	defer session.rwLock.RUnlock()

	matches := []string{}
	for name, plugin := range *session.plugins {
		for _, method := range plugin.MethodNames {
			if method == methodName {
				matches = append(matches, name)
			}
		}
	}

	return matches, nil
}

func CallPluginMethod(pluginName, methodName string, p interface{}, r interface{}) error {
	//TODO error handling
	if session == nil {
		return fmt.Errorf("Must initialize DPlug session before calling plugins")
	}
	if session.plugins == nil {
		return fmt.Errorf("Must initialize DPlug session before calling plugins")
		//TODO panic?
	}

	//verify results interface is a reference
	if rKind := reflect.TypeOf(r).Kind(); rKind != reflect.Ptr {
		return fmt.Errorf(
			"results argument must be reference, but got: %v",
			rKind)
	}

	session.rwLock.RLock()
	defer session.rwLock.RUnlock()

	plugin, ok := (*session.plugins)[pluginName]
	if !ok {
		return fmt.Errorf("plugin '%v' does not exist in session", pluginName)
	}

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:"+strconv.Itoa(plugin.Port))
	if err != nil {
		return fmt.Errorf("error connecting: %", err)
	}

	results := &ExternalResults{r}

	// Synchronous call
	err = client.Call(
		"DPlug.HandleMethod",
		ExternalParameters{
			methodName,
			p,
		},
		results,
	)
	if err != nil {
		return fmt.Errorf("plugin error: %v", err)
	}

	// Set value pointed to by r to be results
	reflect.Indirect(reflect.ValueOf(r)).Set(reflect.ValueOf(results.Results))

	return nil
}
