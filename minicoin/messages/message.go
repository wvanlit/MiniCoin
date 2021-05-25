/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package messages

import "encoding/json"

type MessageType int

const (
	DATA = iota
	AUTHENTICATION_REQUEST
	AUTHENTICATION_RESULT
	ERROR
)

type Message struct {
	MsgType MessageType `json:"msg_type"`
	Payload interface{} `json:"payload"`
}

func UnmarshalMessageJSON(msgBytes []byte) (*Message, error) {
	message := Message{}
	err := json.Unmarshal(msgBytes, &message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

type Data struct {
	Data string `json:"data"`
}

func CreateDataMessage(payload string) Message {
	return Message{DATA, Data{payload}}
}

type Error struct {
	ErrorType ErrorType `json:"error_type"`
	ErrorMsg  string    `json:"error_msg" json:"error_msg"`
}
