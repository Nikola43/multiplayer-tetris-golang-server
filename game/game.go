package game

import (
	"fmt"
	"github.com/antoniodipinto/ikisocket"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/nikola43/tetrisMultiplayer/middleware"
	"github.com/nikola43/tetrisMultiplayer/models"
	"github.com/nikola43/tetrisMultiplayer/websockets"
	"log"
	"math/rand"
	"time"
)

var httpServer *fiber.App

type Game struct {
	Players []*models.Player
}

func (a *Game) Initialize(port string) {
	a.Players = make([]*models.Player, 0)
	a.InitializeHttpServer(port)
}

func HandleRoutes(api fiber.Router) {
	//app.Use(middleware.Logger())

	//routes.ClientRoutes(api)
	

}

func (a *Game) FindOpponent(player models.Player) int {
	max := len(a.Players) - 1
	min := 1
	rand.Seed(time.Now().UnixNano())

	index := models.PlayerExists(a.Players, player.UUID)
	opponent := rand.Intn(max - min + 1) + min

	for ok := true; ok; ok = index != opponent {
		opponent = rand.Intn(max - min + 1) + min
	}

	return opponent
}
func (a *Game) InitializeHttpServer(port string) {
	httpServer = fiber.New(fiber.Config{
		BodyLimit: 2000 * 1024 * 1024, // this is the default limit of 4MB
	})
	//httpServer.Use(middlewares.XApiKeyMiddleware)
	httpServer.Use(cors.New(cors.Config{}))

	ws := httpServer.Group("/ws")

	// Setup the middleware to retrieve the data sent in first GET request
	ws.Use(middleware.WebSocketUpgradeMiddleware)

	// Pull out in another function
	// all the ikisocket callbacks and listeners
	setupSocketListeners(a)

	ws.Get("/:walletAddress", ikisocket.New(func(kws *ikisocket.Websocket) {
		websockets.SocketInstance = kws

		// Retrieve the user id from endpoint
		userId := kws.Params("walletAddress")

		// Add the connection to the list of the connected clients
		// The UUID is generated randomly and is the key that allow
		// ikisocket to manage Emit/EmitTo/Broadcast
		websockets.SocketClients[userId] = kws.UUID

		// Every websocket connection has an optional session key => value storage
		kws.SetAttribute("walletAddress", userId)

		index := models.PlayerExists(a.Players, kws.UUID)

		if index > 0 {
			a.Players[index].UUID = kws.UUID
			//ep.Kws.Emit([]byte(fmt.Sprintf("Player reconnected" + ep.Kws.UUID)))
			color.Yellow(fmt.Sprintf(fmt.Sprintf("Player Connected " + kws.UUID)))
		} else {
			a.Players = append(a.Players, &models.Player{
				UUID:          kws.UUID,
				WalletAddress: userId,
				IsPlaying: false,
			})
			color.Yellow(fmt.Sprintf(fmt.Sprintf("Player Reconnected " + kws.UUID)))
			//fmt.Println(a.Players[len(a.Players) - 1])

			//ep.Kws.Emit([]byte(fmt.Sprintf("New Player Connected " + ep.Kws.UUID)))
		}



		//Broadcast to all the connected users the newcomer
		// kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID)), true)
		//Write welcome message
		kws.Emit([]byte(fmt.Sprintf("Socket connected")))
	}))

	api := httpServer.Group("/api") // /api
	v1 := api.Group("/v1")          // /api/v1
	HandleRoutes(v1)

	err := httpServer.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}



// Setup all the ikisocket listeners
// pulled out main function
func setupSocketListeners(a *Game) {

	// Multiple event handling supported
	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {

	})

	// On message event
	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		fmt.Println(fmt.Sprintf("Message socket event - User: %s", ep.Kws.GetStringAttribute("walletAddress")))
	})

	// On disconnect event
	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		// Remove the user from the local clients
		delete(websockets.SocketClients, ep.Kws.GetStringAttribute("user_id"))
		color.Red(fmt.Sprintf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("walletAddress")))
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		// Remove the user from the local clients
		delete(websockets.SocketClients, ep.Kws.GetStringAttribute("walletAddress"))
		fmt.Println(fmt.Sprintf("Close event - User: %s", ep.Kws.GetStringAttribute("walletAddress")))
	})

	// On error event
	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		fmt.Println(fmt.Sprintf("Error event - User: %s", ep.Kws.GetStringAttribute("walletAddress")))
	})
}
