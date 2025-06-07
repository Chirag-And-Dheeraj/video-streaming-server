package shared

import (
	"encoding/json"
	"log"
	"video-streaming-server/types"
)

// map of user ID and SessionSSEChannelMap
var GlobalUserSSEConnectionsMap = make(map[types.UserID]types.SessionSSEChannelMap)

func SendEventToUser(userID types.UserID, event string, data any) {
	userSSEChannelMap, userExists := GlobalUserSSEConnectionsMap[userID]
	if !userExists {
		log.Printf("user %s not found in GlobalUserSSEConnectionsMap", userID)
		return
	}
	userSessions := userSSEChannelMap.Sessions
	for sessionID, channel := range userSessions {
		eventData, err := json.Marshal(data)
		if err != nil {
			log.Printf("failed marshalling data for event %s to user %s on session %s", event, userID, sessionID)
			return
		}
		log.Printf("sending event %s data %s to user %s on session %s", event, eventData, userID, sessionID)

		channel.EventChannel <- types.SSEType{
			Event: event,
			Data:  data,
		}
	}
}
