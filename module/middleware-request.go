package apinto_module

import (
	"encoding/json"
	"golang.org/x/net/http/httpguts"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"
)

type MiddlewareRequest struct {
	Header  http.Header    `json:"header"`
	Keys    map[string]any `json:"keys"`
	Module  string         `json:"module"`
	Method  string         `json:"method"`
	Url     string         `json:"url"`
	FulPath string         `json:"fulpath"`
}

func Read(r *http.Request) (*MiddlewareRequest, error) {

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	m := new(MiddlewareRequest)
	err = json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (m *MiddlewareRequest) Get(key string) (value any, exists bool) {

	value, exists = m.Keys[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (m *MiddlewareRequest) MustGet(key string) any {
	if value, exists := m.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (m *MiddlewareRequest) GetString(key string) (s string) {
	if val, ok := m.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (m *MiddlewareRequest) GetBool(key string) (b bool) {
	if val, ok := m.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (m *MiddlewareRequest) GetInt(key string) (i int) {
	if val, ok := m.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (m *MiddlewareRequest) GetInt64(key string) (i64 int64) {
	if val, ok := m.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetUint returns the value associated with the key as an unsigned integer.
func (m *MiddlewareRequest) GetUint(key string) (ui uint) {
	if val, ok := m.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

// GetUint64 returns the value associated with the key as an unsigned integer.
func (m *MiddlewareRequest) GetUint64(key string) (ui64 uint64) {
	if val, ok := m.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (m *MiddlewareRequest) GetFloat64(key string) (f64 float64) {
	if val, ok := m.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (m *MiddlewareRequest) GetTime(key string) (t time.Time) {
	if val, ok := m.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (m *MiddlewareRequest) GetDuration(key string) (d time.Duration) {
	if val, ok := m.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (m *MiddlewareRequest) GetStringSlice(key string) (ss []string) {
	if val, ok := m.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (m *MiddlewareRequest) GetStringMap(key string) (sm map[string]any) {
	if val, ok := m.Get(key); ok && val != nil {
		sm, _ = val.(map[string]any)
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (m *MiddlewareRequest) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := m.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (m *MiddlewareRequest) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := m.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}
func (m *MiddlewareRequest) GetCookie(name string) (string, error) {
	cookie, err := m.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}
func (m *MiddlewareRequest) Cookie(name string) (*http.Cookie, error) {
	for _, c := range readCookies(m.Header, name) {
		return c, nil
	}
	return nil, http.ErrNoCookie
}

// Cookies parses and returns the HTTP cookies sent with the request.
func (m *MiddlewareRequest) Cookies() []*http.Cookie {
	return readCookies(m.Header, "")
}

// readCookies  parses all "Cookie" values from the header h and
// returns the successfully parsed Cookies.
//
// if filter isn't empty, only cookies of that name are returned
func readCookies(h http.Header, filter string) []*http.Cookie {
	lines := h["Cookie"]
	if len(lines) == 0 {
		return []*http.Cookie{}
	}

	cookies := make([]*http.Cookie, 0, len(lines)+strings.Count(lines[0], ";"))
	for _, line := range lines {
		line = textproto.TrimString(line)

		var part string
		for len(line) > 0 { // continue since we have rest
			part, line, _ = strings.Cut(line, ";")
			part = textproto.TrimString(part)
			if part == "" {
				continue
			}
			name, val, _ := strings.Cut(part, "=")
			if !isCookieNameValid(name) {
				continue
			}
			if filter != "" && filter != name {
				continue
			}
			val, ok := parseCookieValue(val, true)
			if !ok {
				continue
			}
			cookies = append(cookies, &http.Cookie{Name: name, Value: val})
		}
	}
	return cookies
}

func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

func isCookieNameValid(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.IndexFunc(raw, isNotToken) < 0
}
func isNotToken(r rune) bool {
	return !httpguts.IsTokenRune(r)
}
func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}
