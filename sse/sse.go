package sse

import (
	"log"
	"video-streaming-server/shared"
	"video-streaming-server/types"
	"video-streaming-server/utils"

	"github.com/google/uuid"
)

func CreateNewSSEChannel(refererPath string, bufferSize int) types.SSEChannel {
	return types.SSEChannel{
		OriginatingPage: refererPath,
		EventChannel:    make(chan types.SSEType, bufferSize),
	}
}

func InitializeSSEConnection(userID types.UserID, path string) types.SessionID {
	sessionID := types.SessionID(uuid.New().String())
	log.Printf("userID is: %s", userID)
	log.Printf("originating page is: %s", path)
	userSSEChannelMap, userExists := shared.GlobalUserSSEConnectionsMap[userID]
	if !userExists {
		log.Printf("user %s not found in GlobalUserSSEConnectionsMap", userID)
		userSSEChannelMap = types.SessionSSEChannelMap{
			Sessions: make(map[types.SessionID]types.SSEChannel),
		}
	}
	userSSEChannelMap.Sessions[sessionID] = CreateNewSSEChannel(path, 10)
	shared.GlobalUserSSEConnectionsMap[userID] = userSSEChannelMap
	return sessionID
}

func RemoveSSEConnection(userID types.UserID, sessionID types.SessionID) {
	log.Printf("userID is: %s", userID)
	log.Printf("sessionID to be deleted is: %s", sessionID)
	userSSEChannelMap, userExists := shared.GlobalUserSSEConnectionsMap[userID]
	if !userExists {
		log.Printf("user %s not found in GlobalUserSSEConnectionsMap", userID)
		return
	}
	channel, channelExists := userSSEChannelMap.Sessions[sessionID]
	if !channelExists {
		log.Printf("session %s not found in userSSEChannelMap", sessionID)
		return
	}
	log.Printf("closing channel for user %s, session %s", userID, sessionID)
	close(channel.EventChannel)
	log.Printf("deleting sessionID %s from userSSEChannelMap of user %s", sessionID, userID)
	delete(userSSEChannelMap.Sessions, sessionID)
	if len(userSSEChannelMap.Sessions) == 0 {
		log.Printf("all sessions for user %s are closed", userID)
		log.Printf("deleting user %s from GlobalUserSSEConnectionsMap", userID)
		delete(shared.GlobalUserSSEConnectionsMap, userID)
	} else {
		shared.GlobalUserSSEConnectionsMap[userID] = userSSEChannelMap // clarity
	}

	utils.PrettyPrintMap(shared.GlobalUserSSEConnectionsMap, "GlobalUserSSEConnectionsMap")
}
