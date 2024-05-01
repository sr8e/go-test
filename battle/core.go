package battle

import (
	"context"
	"log"

	"github.com/go-zeromq/zmq4"
)

func BattleMainRoutine(ctx context.Context, cancel context.CancelFunc, roomReq <-chan string) {
	log.Println("battle main routine executed")

	routerSock := zmq4.NewRouter(ctx)
	defer routerSock.Close()
	pubSock := zmq4.NewPub(ctx)
	defer pubSock.Close()

	if err := routerSock.Listen("tcp://*:5555"); err != nil {
		log.Fatal("error listening router socket 5555 port")
	}
	repChan := make(chan zmq4.Msg)
	// receiver
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				log.Print("router recv routine canceled")
				return
			default:
				// recv() will unblock when ctx canceled
				msg, err := routerSock.Recv()
				if err != nil {
					log.Printf("error at router.recv, %s", err)
					continue
				}
				repChan <- msg
			}
		}
	}(ctx)

	if err := pubSock.Listen("tcp://*:5563"); err != nil {
		log.Fatal("error listening pub socket 5563 port")
	}

	rooms := make(map[string]*Room)

	aggrChan := make(chan Message, 10)

	// aggregating channel closure
	go func(ctx context.Context, aggrChan chan<- Message) {
		for {
			select {
			case <-ctx.Done():
				log.Print("aggrChan@battleMain canceled")
				return
			default:
			}

			for k, v := range rooms {
				select {
				case msg, ok := <-v.FromRoom():
					if !ok {
						// room has been closed
						delete(rooms, k)
					} else {
						aggrChan <- msg
					}
				default:
				}
			}
		}

	}(ctx, aggrChan)

	for {
		select {
		case <-ctx.Done():
			log.Println("battle main routine closed")
			return
		case newRoomName, ok := <-roomReq:
			if !ok {
				// channel closed
				cancel()
				log.Fatal("room open request channel closed unexpectedly")
			}
			log.Printf("new Room! %s", newRoomName)
			room := CreateRoom(ctx, newRoomName)
			rooms[newRoomName] = room
			go room.Execute()
		case fanout, ok := <-aggrChan:
			if !ok {
				log.Print("aggrChan closed")
				cancel()
				break
			}
			if fanout.User == "" {
				// pub message
				err := pubSock.Send(fanout.toZmqMsg())
				if err != nil {
					log.Printf("error sending pub: room %s content %s", fanout.RoomId, fanout.Content)
				}
			} else {
				err := routerSock.Send(fanout.toZmqMsg())

				if err != nil {
					log.Printf("error sending router: to %s, content %s", fanout.User, fanout.Content)
				}
			}
		case fanin := <-repChan:
			msg := fromZmqMsg(fanin)
			log.Printf("incoming rep: %s, %s,%s, %s", msg.RoomId, msg.User, msg.Cmd, msg.Content)
			room, ok := rooms[msg.RoomId]
			if !ok {
				repMsg := Message{
					User: msg.User,
					Cmd:  "ERR_INVALID_ROOM",
				}
				routerSock.Send(repMsg.toZmqMsg())
				continue
			}
			toRoomChan, ok := room.ToRoom()
			if !ok {
				repMsg := Message{
					User: msg.User,
					Cmd:  "ERR_INVALID_ROOM",
				}
				routerSock.Send(repMsg.toZmqMsg())
				room.CloseRoom()
				delete(rooms, msg.RoomId)
				continue
			}
			toRoomChan <- msg
		default:
		}

	}
}
