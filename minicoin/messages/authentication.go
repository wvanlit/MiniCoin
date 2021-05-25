/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package messages

import "encoding/json"

type AuthenticationRequest struct {
	Id string `json:"id"`
}

func CreateAuthenticationRequestMessage(id string) Message {
	return Message{AUTHENTICATION_REQUEST, AuthenticationRequest{id}}
}

func UnmarshalAuthenticationRequest(payload map[string]interface{}) (*AuthenticationRequest, error) {
	jsonString, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request := AuthenticationRequest{}
	err = json.Unmarshal(jsonString, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

type AuthenticationResult struct {
	Result bool   `json:"result"`
	Reason string `json:"reason,omitempty"`
}

func CreateAuthenticationResultMessage(result bool, reason string) Message {
	return Message{AUTHENTICATION_RESULT, AuthenticationResult{result, reason}}
}

func UnmarshalAuthenticationResult(payload map[string]interface{}) (*AuthenticationResult, error) {
	jsonString, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	request := AuthenticationResult{}
	err = json.Unmarshal(jsonString, &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}
