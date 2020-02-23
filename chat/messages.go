package chat

import "fmt"

func newEnterMsg(userId string) string {
	return fmt.Sprintf("User_%v entered chatroom", userId)
}

func newLeaveMsg(userId string) string {
	return fmt.Sprintf("User_%v leaved chatroom", userId)
}

func newChatMsg(userId string, msg string) string {
	return fmt.Sprintf("User_%v: %v", userId, msg)
}
