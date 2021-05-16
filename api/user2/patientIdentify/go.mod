module patientIdentify

go 1.15

require (
    github.com/hyperledger/fabric-sdk-go v1.0.0
    github.com/SWU-Blockchain/mol-server/chaincode/structures v0.0.0
)

replace github.com/SWU-Blockchain/mol-server/chaincode/structures v0.0.0 => ../../../structures
