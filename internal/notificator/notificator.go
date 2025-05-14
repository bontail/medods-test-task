package notificator

import (
	"net/netip"
)

type Notificator interface {
	NewIp(guid string, oldIp, newIp netip.Addr)
}
