module github.com/1uvu/Fabric-Demo/api

go 1.15

require (
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 // indirect
	github.com/1uvu/Fabric-Demo/structures v0.0.0 // indirect
	github.com/1uvu/fabric-sdk-client v0.0.4
	github.com/1uvu/serve v0.0.1-beta
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 => ../crypt
	github.com/1uvu/Fabric-Demo/structures v0.0.0 => ../structures
)
