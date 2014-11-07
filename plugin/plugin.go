package plugin

import (
	"fmt"
	lua "github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
	"os"
	"path/filepath"
	"time"
)

type PluginManager struct {
	Plugins map[string]*luar.LuaObject
	Files   map[string]time.Time
	L       *lua.State
}

func (PM *PluginManager) Call(name string) (response string, err error) {
	fn := *PM.Plugins[name]
	resp, err := fn.Call()
	fmt.Printf("PLUGIN CALL %#v\n", resp)
	if err != nil {
		return "", err
	}
	r, _ := resp.(string)
	return r, nil
}

func (PM *PluginManager) loadPlugin(path string) {
	fmt.Printf("loaded %s\n", path)
	err := PM.L.DoFile(path)
	if err != nil {
		fmt.Print("LUA ")
		fmt.Println(err)
	}
}

func (PM *PluginManager) watchDirectory(directory string) {
	for {
		dir, err := os.Open(directory)
		if err != nil {
			fmt.Println(err)
			return
		}
		files, err := dir.Readdir(0)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, f := range files {
			path := filepath.Join(directory, f.Name())
			fh, err := os.Open(path)
			if err != nil {
				fmt.Println(err)
				continue
			}
			stat, _ := fh.Stat()
			mtime := stat.ModTime()
			if oldMtime, ok := PM.Files[path]; ok {
				if oldMtime.Before(mtime) {
					go PM.loadPlugin(path)
					PM.Files[path] = mtime
				}
			} else {
				go PM.loadPlugin(path)
				PM.Files[path] = mtime
			}
		}
		dir.Close()
		time.Sleep(5 * time.Second)
	}
}

func NewPluginManager(directory string) *PluginManager {
	PM := new(PluginManager)

	PM.Plugins = make(map[string]*luar.LuaObject)
	PM.Files = make(map[string]time.Time)
	PM.L = luar.Init()

	luar.RawRegister(PM.L, "", luar.Map{
		"RegisterPlugin": func(L *lua.State) int {
			name := L.ToString(1)
			fn := luar.NewLuaObject(L, 2)
			PM.Plugins[name] = fn
			fmt.Println("  " + name + " registered")
			return 0
		},
	})

	go PM.watchDirectory(directory)
	return PM
}
