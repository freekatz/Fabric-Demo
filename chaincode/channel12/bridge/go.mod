module bridge

go 1.15

require (
	"github.com/1uvu/Fabric-Demo/structures" v0.0.0 // indirect
	github.com/hyperledger/fabric-contract-api-go v1.1.1
)

replace "github.com/1uvu/Fabric-Demo/structures" v0.0.0 => ../../../structures
