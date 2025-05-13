package notificator

import (
	"net/netip"
)

type Notificator interface {
	NewIp(guid, oldIp netip.Addr, newIp string)
}
