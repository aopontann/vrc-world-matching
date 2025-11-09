package vrc_world_matching

import "fmt"

var VRChatAPINotExistWorldError = fmt.Errorf("not exist world error")
var VRChatAPIError = fmt.Errorf("VRChatAPI error")
var AlreadyRegisteredError = fmt.Errorf("already registered")
var NotRegisteredError = fmt.Errorf("not registered")
var NotFoundError = fmt.Errorf("not found")
