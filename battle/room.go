package battle

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sr8e/mellow-ir/db"
)

type Room struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	Name       string
	toRoom     chan Message
	fromRoom   chan Message
	closed     bool
}

type Player struct {
	Id         string
	Name       string
	alive      bool
	aliveTimer *time.Timer
}

func CreateRoom(parent context.Context, name string) *Room {
	ctx, cancel := context.WithCancel(parent)
	to := make(chan Message)
	from := make(chan Message)

	return &Room{
		ctx:        ctx,
		cancelFunc: cancel,
		Name:       name,
		toRoom:     to,
		fromRoom:   from,
	}
}

func createPlayer(id, name string) *Player {
	p := &Player{
		Id:    id,
		Name:  name,
		alive: true,
	}
	p.aliveTimer = time.AfterFunc(30*time.Second, p.kill)
	return p
}

func (r *Room) CloseRoom() {
	log.Printf("closeRoom %s called", r.Name)
	r.cancelFunc()
	if !r.closed {
		close(r.toRoom)
		close(r.fromRoom)
		r.closed = true
	}
}

func (r *Room) ToRoom() (c chan<- Message, ok bool) {
	return r.toRoom, !r.closed
}
func (r *Room) FromRoom() <-chan Message {
	return r.fromRoom
}
func (r *Room) Execute() {
	log.Println("room routine executed")
	players := make(map[string]*Player)

	roomTimer := time.AfterFunc(time.Minute, func() {
		log.Printf("room %s closed as timeout", r.Name)
		r.CloseRoom()
	})
	roomTicker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-roomTicker.C:
			// every 10 seconds
			// checks any player is alive
			anyAlive := false
			for _, p := range players {
				anyAlive = anyAlive || p.alive
			}
			if anyAlive {
				roomTimer.Reset(time.Minute)
			}
			// also send thump to clients
			r.fromRoom <- Message{
				RoomId: r.Name,
				Cmd:    "THUMP_PUB",
			}
		case <-r.ctx.Done():
			if !r.closed {
				r.CloseRoom()
			}
			log.Printf("room routine %s closed", r.Name)
			return

		case msg, ok := <-r.toRoom:
			if !ok {
				//closed
				return
			}
			p, ok := players[msg.User]
			if !ok {
				if msg.Cmd == "JOIN_REQ" {
					dbUser := db.User{Id: msg.User}
					if ok, err := dbUser.Get(); !ok {
						log.Printf("user %s not in db is trying to connect room %s", msg.User, r.Name)
						continue
					} else if err != nil {
						log.Printf("error on registering user %s in room %s", msg.User, r.Name)
						continue
					}
					p = createPlayer(msg.User, dbUser.DisplayName)
					players[msg.User] = p
					r.fromRoom <- Message{
						User: msg.User,
						Cmd:  "JOIN_ACK",
					}
					log.Printf("user %s registered to room %s", msg.User, r.Name)
				} else {
					r.fromRoom <- Message{
						User: msg.User,
						Cmd:  "ERR_NOT_REGISTERED",
					}
				}
				continue
			}
			p.Refresh()
			if msg.Cmd == "THUMP_REQ" {
				r.fromRoom <- Message{
					User: msg.User,
					Cmd:  "THUMP_ACK",
				}
			} else if msg.Cmd == "POLL_OPPONENT" {
				if len(players) == 1 {
					r.fromRoom <- Message{
						User: msg.User,
						Cmd:  "OPPONENT_NOT_FOUND",
					}
				} else if len(players) == 2 {
					for k, v := range players {
						if k != msg.User {
							r.fromRoom <- Message{
								User:    msg.User,
								Cmd:     "OPPONENT_FOUND",
								Content: fmt.Sprintf("%s:%s", k, v.Name),
							}
							break
						}
					}
				} else {
					log.Printf("error on room %s: >2 players?", r.Name)
					r.CloseRoom()
				}
			}
		default:
		}

	}
}

func (p *Player) kill() {
	p.alive = false
	log.Printf("player %s (%s) killed as its inactivity", p.Name, p.Id)
}
func (p *Player) Refresh() {
	p.aliveTimer.Reset(30 * time.Second)
}
