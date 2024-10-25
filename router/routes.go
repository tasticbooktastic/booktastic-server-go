package router

import (
	"booktastic-server-go/address"
	"booktastic-server-go/chat"
	"booktastic-server-go/config"
	"booktastic-server-go/group"
	"booktastic-server-go/isochrone"
	"booktastic-server-go/job"
	"booktastic-server-go/location"
	"booktastic-server-go/message"
	"booktastic-server-go/misc"
	"booktastic-server-go/notification"
	"booktastic-server-go/shelf"
	"booktastic-server-go/user"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// We have two groups because of how the API is used in the old and new clients.
	api := app.Group("/api")
	apiv2 := app.Group("/apiv2")

	for _, rg := range []fiber.Router{api, apiv2} {
		rg.Get("/address", address.ListForUser)
		rg.Get("/address/:id", address.GetAddress)
		rg.Get("/chat", chat.ListForUser)
		rg.Get("/chat/:id/message", chat.GetChatMessages)
		rg.Post("/chat/:id/message", chat.CreateChatMessage)
		rg.Get("/chat/:id", chat.GetChat)
		rg.Get("/config/:key", config.Get)
		rg.Get("/group", group.ListGroups)
		rg.Get("/group/:id", group.GetGroup)
		rg.Get("/group/:id/message", group.GetGroupMessages)
		rg.Get("/isochrone", isochrone.ListIsochrones)
		rg.Get("/isochrone/message", isochrone.Messages)
		rg.Get("/job", job.GetJobs)
		rg.Get("/job/:id", job.GetJob)
		rg.Get("/location/:id", location.GetLocation)
		rg.Get("/user/:id?", user.GetUser)
		rg.Get("/user/:id/publiclocation", user.GetPublicLocation)
		rg.Get("/user/:id/message", message.GetMessagesForUser)
		rg.Get("/user/:id/search", user.GetSearchesForUser)
		rg.Get("/notification/count", notification.Count)
		rg.Get("/notification", notification.List)
		rg.Get("/online", misc.Online)
		rg.Put("/shelf", shelf.Create)
		rg.Get("/shelf/:id", shelf.Single)
		rg.Get("/shelf/:id/books", shelf.Books)
	}
}
