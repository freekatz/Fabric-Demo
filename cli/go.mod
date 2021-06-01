module cli

go 1.15

require (
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 // indirect
	github.com/1uvu/Fabric-Demo/structures v0.0.0
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23 // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric v0.0.0-20190822125948-d2b42602e52e // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/Kubuxu/go-ipfs-swarm-key-gen v0.0.0-20170218193930-0ee739ec6d32 // indirect
	github.com/ipfs/go-ipfs v0.8.0
	github.com/ipfs/go-ipfs-api v0.2.0
	github.com/ipfs/go-ipfs-config v0.14.0
	github.com/ipfs/go-ipfs-files v0.0.8
	github.com/ipfs/go-ipfs-util v0.0.2
	github.com/ipfs/interface-go-ipfs-core v0.4.0
	github.com/klauspost/cpuid/v2 v2.0.6 // indirect
	github.com/libp2p/go-libp2p-core v0.8.0
	github.com/libp2p/go-libp2p-peerstore v0.2.7
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/multiformats/go-multihash v0.0.15 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/sys v0.0.0-20210531080801-fdfd190a6549 // indirect
)

replace (
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 => ../crypt
	github.com/1uvu/Fabric-Demo/structures v0.0.0 => ../structures
)
