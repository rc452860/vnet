package eventbus

import (
	evbus "github.com/asaskevich/EventBus"
)

var eventBus evbus.Bus

func init() {
	eventBus = evbus.New()
}
func GetEventBus() evbus.Bus {
	return eventBus
}
