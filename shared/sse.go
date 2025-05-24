package shared

import (
	"log"
	"video-streaming-server/types"
)

// map of user ID and SessionSSEChannelMap
var GlobalUserSSEConnectionsMap = make(map[types.UserID]types.SessionSSEChannelMap)

func SendEventToUser(userID types.UserID, event string) {
	userSSEChannelMap, userExists := GlobalUserSSEConnectionsMap[userID]
	if !userExists {
		log.Printf("user %s not found in GlobalUserSSEConnectionsMap", userID)
		return
	}
	userSessions := userSSEChannelMap.Sessions
	for sessionID, channel := range userSessions {
		log.Printf("sending event %s to user %s on session %s", event, userID, sessionID)
		channel.EventChannel <- event
	}
}
