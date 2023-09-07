package main

import (
	"work-space/tools/db/pg"
)

func main() {
	pg.Init("162.14.115.114", "cill", "12345678", "test", "5432")
	pg.Client.AutoMigrate(model.)
}
