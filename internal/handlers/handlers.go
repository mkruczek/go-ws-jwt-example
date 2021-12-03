package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/mkruczek/go-ws-jwt-example/internal/auth"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	wsChan = make(chan WsPayload)

	views = jet.NewSet(
		jet.NewOSFileSystemLoader("./html"),
		jet.InDevelopmentMode(),
	)
)

func Home(ctx echo.Context) error {
	err := renderPage(ctx.Response(), "home.jet", nil)
	if err != nil {
		log.Println("ERROR: ", err.Error())
	}
	return nil
}

type WebSocketConnection struct {
	*websocket.Conn
	mu sync.Mutex
}

func (w *WebSocketConnection) Write(v interface{}) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.WriteJSON(v)
}

type WsJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type WsJwtResponse struct {
	Action string `json:"action"`
	Token  string `json:"token"`
}

type WsPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

func WsEndpoint(ctx echo.Context) error {

	upgradeConnection := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	ws, err := upgradeConnection.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Println("ERROR: ", err.Error())
	}

	conn := WebSocketConnection{Conn: ws}

	response := WsJsonResponse{
		Action:  "welcome",
		Message: "connection to the server",
	}

	conn.Write(response)

	go listenForWs(&conn)

	go listenForWsChannel(&conn)

	go authForWs(&conn)

	go displayData(&conn)

	return nil
}

func displayData(conn *WebSocketConnection) {
	t := time.NewTicker(2 * time.Second)

	data := []string{"drzewo", "kwiatek", "krzesło", "biurko", "słuchawki", "drapaczka"}
	rand.Seed(time.Now().Unix())

	for range t.C {
		request := WsPayload{
			Action:  "display_data",
			Message: data[rand.Intn(len(data))],
		}
		conn.Write(request)
	}
}

func authForWs(conn *WebSocketConnection) {
	t := time.NewTicker(15 * time.Second)

	for range t.C {
		request := WsPayload{
			Action:  "checked",
			Message: "validate token for: " + time.Now().String(),
		}
		conn.Write(request)
	}
}

func listenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(fmt.Sprintf("recovery... : %v", r))
		}
	}()

	var body WsPayload

	for {

		err := conn.ReadJSON(&body)
		if err != nil {
			log.Println(fmt.Sprintf("error for reading body from WS : %s", err.Error()))
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Println("connection to ws has been closed")
				break
			}
			continue
		}
		wsChan <- body
	}
}

func listenForWsChannel(con *WebSocketConnection) {

	for {
		e := <-wsChan

		switch e.Action {
		case "initial_connection":
			jwt, err := auth.CreateJWT()
			if err != nil {
				panic(err)
			}
			r := WsJwtResponse{
				Action: "initial_jwt",
				Token:  jwt,
			}
			con.Write(r)
		case "auth":
			valid := auth.CheckJWT(e.Message)
			if valid != nil {
				con.Close()
			} else {
				log.Println("token still valid")
			}
		default:
			log.Println(fmt.Sprintf("error, unknown action: %s", e.Action))
		}
	}
}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {

	view, err := views.GetTemplate(tmpl)
	if err != nil {
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		return err
	}

	return nil
}
