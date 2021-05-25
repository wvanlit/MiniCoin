/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 25-May-2021
 */

package messages

type ErrorType int

const (
	INVALID_MESSAGE = iota + 101
	INVALID_FORMAT
	NOT_AUTHENTICATED
	INVALID_AUTHENTICATION
)

func CreateErrorMessage(errorType ErrorType, errorMessage string) Message {
	return Message{ERROR, Error{errorType, errorMessage}}
}
