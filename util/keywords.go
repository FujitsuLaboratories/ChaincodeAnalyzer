package util

// TODO: Cover all patterns for now, only cover the functions in math/rand and time
var Keywords = map[string]map[string]bool{
	"Rand": {
		"Float32":     true,
		"Float64":     true,
		"ExpFloat64":  true,
		"NormFloat64": true,
		"Int":         true,
		"Int31":       true,
		"Int31n":      true,
		"Int63":       true,
		"Int63n":      true,
		"Intn":        true,
		"Uint32":      true,
		"Uint64":      true,
	},
	"Time": {
		"Now":      true,
		"Date":     true,
		"Unix":     true,
		"Local":    true,
		"Location": true,
	},
	"API": {
		"Get":        true,
		"Head":       true,
		"Post":       true,
		"PostForm":   true,
		"NewRequest": true,
		// for client
		"Do": true,
	},
	"SysCom": {
		"Command": true,
	},
	"ReadFile": {
		// os
		"Open":     true,
		"OpenFile": true,
		// io/ioutil
		"ReadFile": true,
	},
	"RangeQuery": {
		"GetHistoryForKey": true,
		"GetQueryResult":   true,
	},
	"CrossChan": {
		"InvokeChaincode": true,
	},
}

var LibFullPath = map[string][]string{
	"Rand":       []string{"math/rand"},
	"API":        []string{"net/http"},
	"Time":       []string{"time"},
	"SysCom":     []string{"os/exec"},
	"ReadFile":   []string{"os", "io/ioutil"},
	"RangeQuery": []string{"github.com/hyperledger/fabric/core/chaincode/shim"},
	"CrossChan":  []string{"github.com/hyperledger/fabric/core/chaincode/shim"},
}
