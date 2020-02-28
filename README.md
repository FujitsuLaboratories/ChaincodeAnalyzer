# Chaincode Analyzer
Chaincode Analyzer is a CLI tool to detect the codes which can be risks potentially 
such as nondeterminism in Chaincode (i.e., smart contract in Hyperledger Fabric) written in Golang.

## How to use
1. Clone this repository
1. `go build ccanalyzer.go`
1. `./ccanalyzer [file | directory]`

## How to read the output
If Chaincode Analyzer find any risk, it outputs the followings.
1. Category
	- The type of risk
	- e.g., Rand
1. Function
	- The function name which includes the risk
	- e.g., init => `func init() {}`
1. VarName
	- The name of variable which related to the risk
	- e.g., Aval 
1. Position
	- The position of the code related to the risk
	- e.g., `example.go:122:14 Aval = rand.Float32()`
1. Affected Position
	- The position of the code which is affected by the risk
	- e.g., `example.go:151:25 err = stub.PutState(A, Aval)`

## What can this tool detect
Currently, the tool can detect following risks.
For more information about risks, please refer the [paper](https://ieeexplore.ieee.org/abstract/document/8666486).

- Random value 
- Timestamp 
- Iteration on map object
- Calling external API
- File access
- Pointer
- Global variable
- External library
- System commands
- Goroutine
- Range query risk
- Field declaration 
- Read your write
- Cross channel Invocation 

## License
This tool is distributed under the Apache License Version 2.0, see [LICENSE](./LICENSE) file.
