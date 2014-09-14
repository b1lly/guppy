package main

import "flag"

const (
	guppyFile = ".guppy"
)

var (
	guppyCfg = NewGuppyConfig()
	glog     = GuppyLog{}
)

func main() {
	flag.Parse()
	args := flag.Args()

	project, err := NewProject()
	if err != nil {
		glog.Println(err)
		return
	}

	commands := map[string]func(args []string) error{
		"install":  project.Install,
		"register": project.Register,
	}

	if len(args) <= 0 {
		glog.Println("commands: install, register, add")
		return
	}

	var (
		cmd func(args []string) error
		ok  bool
	)

	if cmd, ok = commands[args[0]]; !ok {
		glog.Println("...list help...")
		return
	}

	cmd(args[1:])
}
