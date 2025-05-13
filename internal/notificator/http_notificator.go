package notificator

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/netip"

	"medods-test-task/internal/notificator/config"
)

type HTTPNotificator struct {
	cfg *config.Config
	lgr *slog.Logger
}

func NewHTTPNotificator(cfg *config.Config, lgr *slog.Logger) *HTTPNotificator {
	return &HTTPNotificator{
		lgr: lgr,
		cfg: cfg,
	}
}

func (h *HTTPNotificator) NewIp(guid, oldIp netip.Addr, newIp string) {
	data := map[string]string{"guid": guid.String(), "oldIp": oldIp.String(), "newIp": newIp}

	j, _ := json.Marshal(data)

	_, _ = http.Post(h.cfg.NewIpUrl, "application/json", bytes.NewBuffer(j))
}
