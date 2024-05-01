package battle

import "github.com/go-zeromq/zmq4"

type Message struct {
	RoomId  string
	User    string
	Content string
	Cmd     string
}

func fromZmqMsg(msg zmq4.Msg) Message {
	user := string(msg.Frames[0])
	room := string(msg.Frames[2])
	cmd := string(msg.Frames[3])
	buf := make([]byte, 0)
	for _, f := range msg.Frames[4:] {
		buf = append(buf, f...)
	}

	return Message{
		RoomId:  room,
		User:    user,
		Cmd:     cmd,
		Content: string(buf),
	}
}
func (m *Message) toZmqMsg() zmq4.Msg {
	if m.User == "" {
		return zmq4.NewMsgFrom([]byte(m.RoomId), []byte(m.Cmd), []byte(m.Content))
	} else {
		return zmq4.NewMsgFrom([]byte(m.User), []byte(""), []byte(m.Cmd), []byte(m.Content))
	}
}
