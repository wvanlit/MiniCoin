package messages

const (
	AUTHENTICATION = 1
	DATA           = 2
)

type Message struct {
	MsgType int
	Data    interface{}
}

type Authentication struct {
	Id string
}

func CreateAuthenticationMessage(id string) Message {
	return Message{AUTHENTICATION, Authentication{id}}
}

type Data struct {
	Payload string
}

func CreateDataMessage(payload string) Message {
	return Message{DATA, Data{payload}}
}
