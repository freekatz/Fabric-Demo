module patient

go 1.15

require (
	github.com/SWU-Blockchain/mol-server/chaincode/structures v0.0.0
	github.com/hyperledger/fabric-contract-api-go v1.1.1
)

replace github.com/SWU-Blockchain/mol-server/chaincode/structures v0.0.0 => ../../../structures
