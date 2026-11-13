package routes

import (
	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/presentation/core/controllers"
)

// RegisterCronRoutes mounts every /api/v1/crons endpoint (10 endpoints total).
func RegisterCronRoutes(group fiber.Router, perm func(string) fiber.Handler, controller *controllers.CronController) {
	crons := group.Group("/crons")

	// Static / specific runs routes MUST be registered before parameterized :id routes
	crons.Get("/runs", perm("get-cron-runs"), controller.ListAllRuns)
	crons.Get("/runs/:run_id", perm("get-cron-run-detail"), controller.ShowRun)

	// Crons CRUD and actions
	crons.Get("/", perm("get-crons"), controller.Index)
	crons.Post("/", perm("create-crons"), controller.Store)
	crons.Get("/:id", perm("get-cron-detail"), controller.Show)
	crons.Put("/:id", perm("update-crons"), controller.Update)
	crons.Delete("/:id", perm("delete-crons"), controller.Delete)
	crons.Post("/:id/toggle", perm("toggle-crons"), controller.Toggle)
	crons.Post("/:id/trigger", perm("trigger-crons"), controller.Trigger)
	crons.Get("/:id/runs", perm("get-cron-specific-runs"), controller.ListRuns)
}
