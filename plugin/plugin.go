package plugin

import (
	"encoding/json"
	"errors"
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type PluginManager struct {
	Commands  map[string]LuaPlugin
	Events    map[string]map[string]LuaPlugin
	Files     map[string]time.Time
	GoPlugins map[string]GoPlugin
	L         *lua.State
}

type LuaPlugin struct {
	name string
	fn   *luar.LuaObject
}

type GoPlugin interface {
	init(bot Bot)
}

func (PM *PluginManager) CallCommand(name string, ctx Context, bot Bot) (response string, err error) {
	fn := PM.Commands[name].fn
	if fn == nil {
		return "", errors.New("no command named " + name)
	}
	resp, err := fn.Call(ctx.Text, ctx, bot)
	if err != nil {
		return "", err
	}
	r, _ := resp.(string)
	return r, nil
}

func (PM *PluginManager) CallEvent(event string, ctx Context, bot Bot) (responses chan string, err error) {
	responses = make(chan string)
	go func() {
		for _, plug := range PM.Events[event] {
			fn := plug.fn
			if fn == nil {
				continue
			}
			resp, err := fn.Call(ctx.Text, ctx, bot)
			if err != nil {
				continue
			}
			if resp != nil {
				responses <- resp.(string)
			}
		}
	}()
	return responses, nil
}

func (PM *PluginManager) loadPlugin(path string) {
	log.Printf("loaded %s\n", path)
	err := PM.L.DoFile(path)
	if err != nil {
		log.Printf("!!! %s: %s\n", path, err.Error())
	}
}

func (PM *PluginManager) watchDirectory(directory string) {
	for {
		dir, err := os.Open(directory)
		if err != nil {
			log.Println(err)
			return
		}
		files, err := dir.Readdir(0)
		if err != nil {
			log.Println(err)
			return
		}

		for _, f := range files {
			if !strings.HasSuffix(f.Name(), "lua") {
				continue
			}
			path := filepath.Join(directory, f.Name())
			stat, err := os.Stat(path)
			if err != nil {
				log.Println(err)
				continue
			}
			mtime := stat.ModTime()
			if oldMtime, ok := PM.Files[path]; ok {
				if oldMtime.Before(mtime) {
					PM.loadPlugin(path)
					PM.Files[path] = mtime
				}
			} else {
				PM.loadPlugin(path)
				PM.Files[path] = mtime
			}
		}
		dir.Close()
		time.Sleep(5 * time.Second)
	}
}

func (PM *PluginManager) initLua() {
	PM.L = luar.Init()
	luar.RawRegister(PM.L, "", luar.Map{
		"RegisterCommand": func(L *lua.State) int {
			name := L.ToString(1)
			fn := luar.NewLuaObject(L, 2)
			PM.Commands[name] = LuaPlugin{name, fn}
			log.Printf("    %-10s command\n", name)
			return 0
		},
		"RegisterEvent": func(L *lua.State) int {
			name := L.ToString(1)
			event := L.ToString(2)
			fn := luar.NewLuaObject(L, 3)
			if _, ok := PM.Events[event]; !ok {
				PM.Events[event] = make(map[string]LuaPlugin)
			}
			PM.Events[event][name] = LuaPlugin{name, fn}
			log.Printf("    %-10s event\n", name)
			return 0
		},
	})

	luar.Register(PM.L, "go", luar.Map{
		"Split":  strings.Split,
		"SplitN": strings.SplitN,
		"PrintTable": func(table interface{}) {
			log.Printf("%#v\n", table)
		},
		"GetHTTP": func(url string) string {
			resp, err := http.Get(url)
			if err != nil {
				return ""
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return ""
			}
			return string(body)
		},
		"GetJSON": func(url string) luar.Map {
			resp, err := http.Get(url)
			if err != nil {
				return nil
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil
			}
			var js luar.Map
			json.Unmarshal(body, &js)
			return js
		},
	})
}

func (PM *PluginManager) initGo(bot Bot) {
	for _, v := range PM.GoPlugins {
		go v.init(bot)
	}
}

func NewPluginManager(directory string, bot Bot) *PluginManager {
	PM := new(PluginManager)

	PM.Commands = make(map[string]LuaPlugin)
	PM.Events = make(map[string]map[string]LuaPlugin)
	PM.Files = make(map[string]time.Time)
	PM.GoPlugins = make(map[string]GoPlugin)

	PM.GoPlugins["http"] = new(HttpServerPlugin)

	PM.initLua()
	PM.initGo(bot)

	go PM.watchDirectory(directory)
	return PM
}
