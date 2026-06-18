package main

import (
	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/db"
	"go-fiber-svelte/internal/db/seed"
)

func main() {
	config.InitConfigApp()
	db.Init()
	seed.Run()
}
