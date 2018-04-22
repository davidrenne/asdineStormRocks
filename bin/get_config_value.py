import json, os, sys
gopath = os.getenv("GOPATH")
base = gopath + "/src/github.com/davidrenne/asdineStormRocks/bin/globalcache/"
config = json.load(open(gopath + '/src/github.com/davidrenne/asdineStormRocks/webConfig.json', 'r'))
ptrOne = sys.argv[1]
ptrTwo = sys.argv[2]
print config[ptrOne][ptrTwo]
