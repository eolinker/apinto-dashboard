package shared

import "github.com/hashicorp/go-plugin"

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "HTTP_HANDLER_PLUGIN",
	MagicCookieValue: "Apinto_Dashboard_Plugin",
}
