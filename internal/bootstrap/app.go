package bootstrap

import (
	"go-fiber-svelte/internal/config"
	"go-fiber-svelte/internal/db"
)

func Init() {
	config.InitConfigApp()
	config.I18n()
	db.Init()
}
