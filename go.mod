module mapcleaner

go 1.20

replace github.com/minetest-go/mtdb => github.com/ronoaldo/mtdb v1.1.45-rc1

require (
	github.com/hashicorp/golang-lru/v2 v2.0.7
	github.com/minetest-go/areasparser v1.0.0
	github.com/minetest-go/mapparser v0.1.8
	github.com/minetest-go/mtdb v1.1.35
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/klauspost/compress v1.15.4 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	golang.org/x/sys v0.11.0 // indirect
)
