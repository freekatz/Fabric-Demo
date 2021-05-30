module api

go 1.15

require (
	github.com/1uvu/Fabric-Demo v0.0.0-20210529091904-4ea8b2dd4986 // indirect
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 // indirect
	github.com/1uvu/Fabric-Demo/serve v0.0.0 // indirect
	github.com/1uvu/Fabric-Demo/structures v0.0.0 // indirect
)

replace (
	github.com/1uvu/Fabric-Demo/crypt v0.0.0 => ../crypt
	github.com/1uvu/Fabric-Demo/serve v0.0.0 => ../serve
	github.com/1uvu/Fabric-Demo/structures v0.0.0 => ../structures
)
