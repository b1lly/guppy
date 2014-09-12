package main

import "flag"

const (
	guppyFile = ".guppy"
)

var (
	guppyCfg = NewGuppyConfig()
	commands = map[string]func(args []string) error{
		"install":  Install,
		"register": Register,
		"add":      Add,
	}
	glog = GuppyLog{}
)

func main() {
	// cmd line interface
	flag.Parse()
	args := flag.Args()

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

	cmd([]string{})
}

func Install(args []string) error {
	// - look for local gup.json package requirement
	// - convert deps into request
	// - send get request to registry server
	// - use response remote/commit hash to run a git clone to destination dir (perhaps use .gup.conf)

	return nil
}

func Register(args []string) error {
	// - look for local gup.json setting
	// - convert settings into request
	// - confirm remote pkg repo exists
	// - send register request to gup registry server
	return nil
}

func Add(args []string) error {
	// - convert [pkg] to get request
	// - send get request to registry server
	// - use response remote/commit hash to run a git clone to destination dir
	// - add pkg to gup.json dependencies
	return nil
}
