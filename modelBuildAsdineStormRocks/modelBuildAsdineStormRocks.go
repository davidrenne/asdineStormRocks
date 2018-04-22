package main

import (
	"flag"
	"github.com/DanielRenne/GoCore/modelBuild"
)

func main() {
	// allow -configFile=test.json to be passed to build different configs other than webConfig.json
	configFile := flag.String("configFile", "webConfig.json", "Configuration File Name.  Ex...  webConfig.json")
	flag.Parse()
	modelBuild.Initialize("src/github.com/davidrenne/asdineStormRocks", *configFile)
}
