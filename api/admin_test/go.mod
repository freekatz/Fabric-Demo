module admin_test

go 1.15

require (
	github.com/1uvu/Fabric-Demo/api/app v0.0.0
	github.com/1uvu/Fabric-Demo/api/admin v0.0.0
	github.com/1uvu/Fabric-Demo/crypt v0.0.0
	github.com/1uvu/Fabric-Demo/structures v0.0.0
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/pkg/errors v0.9.1
)

replace (
	github.com/1uvu/Fabric-Demo/api/app v0.0.0 => ../app
	github.com/1uvu/Fabric-Demo/api/admin v0.0.0 => ../admin
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 => ../../crypt
	github.com/1uvu/Fabric-Demo/structures v0.0.0 => ../../structures
)
