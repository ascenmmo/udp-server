// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"github.com/ascenmmo/udp-server/pkg/api/types"
	"github.com/google/uuid"
)

type requestServerSettingsGetConnectionsNum struct {
	Token string `json:"token"`
}

type responseServerSettingsGetConnectionsNum struct {
	CountConn int  `json:"countConn"`
	Exists    bool `json:"exists"`
}

type requestServerSettingsHealthCheck struct {
	Token string `json:"token"`
}

type responseServerSettingsHealthCheck struct {
	Exists bool `json:"exists"`
}

type requestServerSettingsGetServerSettings struct {
	Token string `json:"token"`
}

type responseServerSettingsGetServerSettings struct {
	Settings types.Settings `json:"settings"`
}

type requestServerSettingsCreateRoom struct {
	Token      string                  `json:"token"`
	CreateRoom types.CreateRoomRequest `json:"createRoom"`
}

// Formal exchange type, please do not delete.
type responseServerSettingsCreateRoom struct{}

type requestServerSettingsSetNotifyServer struct {
	Token string    `json:"token"`
	Id    uuid.UUID `json:"id"`
	Url   string    `json:"url"`
}

// Formal exchange type, please do not delete.
type responseServerSettingsSetNotifyServer struct{}
