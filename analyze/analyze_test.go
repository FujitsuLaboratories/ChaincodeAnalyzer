package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCheckImports(t *testing.T) {
	testsrc := `package main

import (
	"fmt"
	"github.com/yamashita/cc_analyzer/analyze"
	r "math/rand"
)

func main() {
	fmt.Println("test code")
}`

	si, f := createASTFromSrc(testsrc, t)
	f.checkImports(si)

	testData := []struct {
		full   string
		abbrev string
	}{
		{"fmt", "fmt"},
		{"github.com/yamashita/cc_analyzer/analyze", "analyze"},
		{"math/rand", "r"},
	}
	got := si.importMap
	for _, test := range testData {
		if got[test.full] != test.abbrev {
			t.Error("checkImports is fail: got", got)
		}
	}
}

func TestReceiverType(t *testing.T) {
	fset := token.NewFileSet()
	testsrc := `package main

import (
	"fmt"
	"github.com/yamashita/cc_analyzer/analyze"
	r "math/rand"
)

type st struct {}

func (s *st) test() {
	fmt.Println("test")
}`

	p, err := parser.ParseFile(fset, "", testsrc, 0)
	if err != nil {
		t.Error("parse failed")
	}

	for _, n := range p.Decls {
		if m, ok := n.(*ast.FuncDecl); ok {
			got := receiverType(m)
			if got != "st" {
				t.Error("receiver type is not match", got)
			}
		}
	}
}

func TestCheckGlobalVar(t *testing.T) {
	testsrc := `package main

import (
	"fmt"
)

var a string

func main() {
	a = "test"
	fmt.Println("a")
}`

	si, f := createASTFromSrc(testsrc, t)
	f.checkGlobalVar(si)

	got := si.problems[0]
	if got.Category != "Global Variable" || got.VarName != "" || got.Function != "Global Space" {
		t.Error("checkGlobalVar failed")
	}
}

func TestCheckMapIter(t *testing.T) {
	testsrc := `package main

import (
	"fmt"
)

func main() {
	type test struct {
		a map[int]int
	}
	a := map[int]int{1:1, 2:2, 3:3}
	for k, v := range a {
		fmt.Println(k, v)
	}
	n := []map[int]int{a, a}
	for k, v := range n[0] {
		fmt.Println(k, v)
	}
	o := test{a: a}
	for k, v := range o.a {
		fmt.Println(k, v)
	}
}`

	var a Analyzer
	problems, _ := a.Analyze(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile), "t.go", []byte(testsrc))
	//si, f := createASTFromSrc(testsrc, t)
	//f.storeInfo(si)

	if len(problems) != 3 {
		t.Error("Failed to detect problems")
	} else {
		got := problems[0]
		if got.Category != "MapIter" || got.VarName != "a" || got.Function != "main" {
			t.Error("checkMapIter: category, varname, or function is different:\n expected: MapIter a main\n got: ", got.Category, got.VarName, got.Function)
		}
		got = problems[1]
		if got.Category != "MapIter" || got.VarName != "n[0]" || got.Function != "main" {
			t.Error("checkMapIter: category, varname, or function is different:\n expected: MapIter n[0] main\n got: ", got.Category, got.VarName, got.Function)
		}
		got = problems[2]
		if got.Category != "MapIter" || got.VarName != "o.a" || got.Function != "main" {
			t.Error("checkMapIter: category, varname, or function is different:\n expected: MapIter n main\n got: ", got.Category, got.VarName, got.Function)
		}
	}
}

func TestCheckExternalLibrary(t *testing.T) {
	testsrc := `package main

import (
	"github.com/kzhry/external/lib"
)

func main() {
	lib.Println("a")
}`

	si, f := createASTFromSrc(testsrc, t)
	f.checkImports(si)

	got := si.problems[0]
	if got.Category != "External Library" || got.VarName != "" || got.Function != "Imports" {
		t.Error("External Library failed", got.Category, got.VarName, got.Function)
	}
}

func TestRandVar(t *testing.T) {
	testsrc := `
	package main

	import (
	"math/rand"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	)
	func (t *SimpleChaincode) example(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]
	Aval := rand.Float32()
	
	err := stub.PutState(A, Aval)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
	`
	si, f := createASTFromSrc(testsrc, t)
	f.checkImports(si)
	f.storeInfo(si)
	f.detectProblems(si)
	if len(si.problems) != 1 {
		t.Error("checkRandVar failed: ", si.problems, si.targetVar)
	} else {
		got := si.problems[0]
		if got.Category != "Rand" || got.VarName != "Aval" || got.Function != "example" {
			t.Error("checkRandVar failed", got)
		}
	}
}

func TestTimestampVar(t *testing.T) {
	testsrc := `
	package main

	import (
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	)
	func (t *SimpleChaincode) example(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]
	Aval := time.Now()
	
	err := stub.PutState(A, Aval)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
	`
	si, f := createASTFromSrc(testsrc, t)
	f.checkImports(si)
	f.storeInfo(si)
	f.detectProblems(si)
	if len(si.problems) != 1 {
		t.Error("checkTimestampVar failed: ", si.problems)
	} else {
		got := si.problems[0]
		if got.Category != "Time" || got.VarName != "Aval" || got.Function != "example" {
			t.Error("checkTimestampVar failed")
		}
	}
}

func TestAPIUsage(t *testing.T) {
	testsrc := `
	package main

	import (
	"net/http"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	)

	func (t *SimpleChaincode) example(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	client := &http.Client{Transport: tr}
	resp := client.Get("https://example.com")
	err := stub.PutState(A, resp)

	return shim.Success(nil)
	}

	func main() {
err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
	}
	`
	si, f := createASTFromSrc(testsrc, t)
	f.checkImports(si)
	f.storeInfo(si)
	f.detectProblems(si)
	if len(si.problems) != 1 {
		t.Error("checkAPIVar failed: ", si.problems, si.targetVar)
	} else {
		got := si.problems[0]
		if got.Category != "API" || got.VarName != "resp" || got.Function != "example" {
			t.Error("checkAPIVar failed")
		}
	}
}

func TestGoroutine(t *testing.T) {
	testsrc := `
	package main

	import (
		"github.com/hyperledger/fabric/core/chaincode/shim"
		pb "github.com/hyperledger/fabric/protos/peer"
	)

	func example() bool {
		return true
	}

	func main() {
		err := shim.Start(new(SimpleChaincode))
		go example()
		if err != nil {
			logger.Errorf("Error starting Simple chaincode: %s", err)
		}
	}
	`
	si, f := createASTFromSrc(testsrc, t)
	f.checkImports(si)
	f.storeInfo(si)
	f.detectProblems(si)
	if len(si.problems) != 1 {
		t.Error("The number of problems is different to expected: ", si.problems, si.targetVar)
	} else {
		got := si.problems[0]
		if got.Category != "Goroutine" || got.VarName != "" || got.Function != "main" {
			t.Error("Different problem is included in the src")
		}
	}
}

func TestFieldDeclaration(t *testing.T) {
	// This example is from explanation of field declarations
	// in https://chaincode.chainsecurity.com
	testsrc := `
	package main

	import (
		"fmt"
		"github.com/hyperledger/fabric/core/chaincode/shim"
		"github.com/hyperledger/fabric/protos/peer"
	)

	type BadChaincode struct {
		globalValue string
	}

	func (t *BadChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
		return shim.Success(nil)
	}

	func (t BadChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
		fn, args := stub.GetFunctionAndParameters()

		if fn == "setValue" {
			t.globalValue = args[0]
			shim.PutState("key", []byte(t.globalValue))
			return shim.Success([]byte("success"))
		} else if fn == "getValue" {
			value, _ := shim.GetState("key")
			return shim.Success(value)
		}
		return shim.Error("not a valid function")
	}

	func main() {
		if err := shim.Start(new(BadChaincode)); err != nil {
			fmt.Printf("Error starting BadChaincode chaincode: %s", err)
		}
	}
	`
	si, f := createASTFromSrc(testsrc, t)
	f.checkImports(si)
	f.checkFieldDeclaration(si)
	if len(si.problems) != 1 {
		t.Error("The number of problems is different to expected: ", si.problems, si.targetVar)
	} else {
		got := si.problems[0]
		if got.Category != "FieldDeclaration" || got.VarName != "BadChaincode" || got.Function != "Declaration" {
			t.Error("Different problem is included in the src")
		}
	}
}

// TODO: Modify Test
func TestPointer(t *testing.T) {
}

func TestSysCom(t *testing.T) {
}

func TestCrossChan(t *testing.T) {
}

func TestReadFile(t *testing.T) {
}

func TestReadYourWrite(t *testing.T) {
	//testsrc := "../test/RiskQuery.go"
	//si, f := createASTFromFileName(testsrc, t)
	//f.checkGlobalVar(si)
	//f.storeInfo(si)
	//f.detectProblems(si)
	//if len(si.problems) > 0 {
	//	got := si.problems[0]
	//	if got.Category != "ReadYourWrite" || got.VarName != "" || got.Function != "ReadYourWrite" {
	//		t.Error("checkReadYourWrite failed: different problems are found", got)
	//	}
	//} else {
	//	t.Error("checkReadYourWrite failed: any problem is not found")
	//}
}

func TestRangeQuery(t *testing.T) {
	testsrc := "../test/RiskQuery.go"
	si, f := createASTFromFileName(testsrc, t)
	f.checkImports(si)
	f.checkGlobalVar(si)
	f.checkFieldDeclaration(si)
	f.checkImports(si)
	f.checkGlobalVar(si)
	f.checkFieldDeclaration(si)

	f.storeInfo(si)
	f.detectProblems(si)
	if len(si.problems) > 0 {
		got := si.problems[0]
		if got.Category != "RangeQuery" || got.VarName != "content" || got.Function != "initLedger" {
			t.Error("checkRangeQuery failed: different problems are found", got)
		}
	} else {
		t.Error("checkRangeQuery failed: any problem is not found")
	}
}

// createASTFromXXX creates storedInfo and file objects from each source
func createASTFromSrc(testsrc string, t *testing.T) (*storedInfo, *file) {
	fset := token.NewFileSet()
	p, err := parser.ParseFile(fset, "", testsrc, 0)
	if err != nil {
		t.Error("parse failed")
	}
	f := &file{f: p, fset: fset, src: []byte(testsrc), functions: make(map[string]function)}
	si := &storedInfo{
		file:      *f,
		targetVar: make(map[string]map[string]map[token.Pos]ast.Expr),
		mappings: mappings{
			opMap:        make(map[string]map[string][]ast.Expr),
			dotImportMap: make(map[string]bool),
			importMap:    make(map[string]string),
			pointerMap:   make(map[string]map[string]token.Pos),
		},
		isValidProblem: make(map[string]map[string]map[string]bool),
		flagNonDet: map[string]bool{"Rand": false, "Time": false, "API": false,
			"SysCom": false, "ReadFile": false, "RangeQuery": false, "CrossChan": false},
	}
	return si, f
}

func createASTFromFileName(filename string, t *testing.T) (*storedInfo, *file) {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		t.Error("no such file")
	}
	p, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		t.Error("parse failed")
	}
	f := &file{f: p, fset: fset, src: src, functions: make(map[string]function)}
	si := &storedInfo{
		file:      *f,
		targetVar: make(map[string]map[string]map[token.Pos]ast.Expr),
		mappings: mappings{
			opMap:        make(map[string]map[string][]ast.Expr),
			dotImportMap: make(map[string]bool),
			importMap:    make(map[string]string),
			pointerMap:   make(map[string]map[string]token.Pos),
		},
		isValidProblem: make(map[string]map[string]map[string]bool),
		flagNonDet: map[string]bool{"Rand": false, "Time": false, "API": false,
			"SysCom": false, "ReadFile": false, "RangeQuery": false, "CrossChan": false},
	}
	return si, f
}
