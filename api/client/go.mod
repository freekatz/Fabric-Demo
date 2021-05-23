module client

go 1.15

require (
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 // indirect
	github.com/1uvu/Fabric-Demo/structures v0.0.0 // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23 // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190822125948-d2b42602e52e // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

replace (
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 => ../../crypt
	github.com/1uvu/Fabric-Demo/structures v0.0.0 => ../../structures
)
