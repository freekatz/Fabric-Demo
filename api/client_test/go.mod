module client_test

go 1.15

require (
    "github.com/hyperledger/fabric-sdk-go" v1.0.0
	"github.com/1uvu/Fabric-Demo/api/client" v0.0.0
    "github.com/1uvu/Fabric-Demo/crypt" v0.0.0
	"github.com/1uvu/Fabric-Demo/structures" v0.0.0
)

replace (
	"github.com/1uvu/Fabric-Demo/api/client" v0.0.0 => ../client
	"github.com/1uvu/Fabric-Demo/crypt" v0.0.0 => ../../crypt
	"github.com/1uvu/Fabric-Demo/structures" v0.0.0 => ../../structures
)