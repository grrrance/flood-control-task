package floodctrl

import "net/http"

type Handlers interface {
	TriggerUser() http.HandlerFunc
}
