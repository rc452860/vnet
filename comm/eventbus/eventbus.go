package eventbus

import (
	evbus "github.com/asaskevich/EventBus"
	"github.com/rc452860/vnet/utils"
)

var eventBus evbus.Bus

func GetEventBus() evbus.Bus {
	utils.Lock("eventbus:eventBus")
	defer utils.UnLock("eventbus:eventBus")
	if eventBus != nil {
		eventBus = evbus.New()
	}
	return eventBus
}
