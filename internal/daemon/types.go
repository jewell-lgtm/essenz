package daemon

// Cookie represents a single HTTP cookie.
type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain,omitempty"`
	Path   string `json:"path,omitempty"`
}

// Request represents a client request to the daemon.
type Request struct {
	Action    string         `json:"action"`
	URL       string         `json:"url,omitempty"`
	Readiness *ReadinessSpec `json:"readiness,omitempty"`
	Cookies   []Cookie       `json:"cookies,omitempty"`
}

// Response represents the daemon's response.
type Response struct {
	Success bool   `json:"success"`
	Content string `json:"content,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ReadinessSpec conveys DOM readiness preferences to the daemon.
type ReadinessSpec struct {
	TimeoutMillis int      `json:"timeout_ms,omitempty"`
	UseLCP        bool     `json:"use_lcp,omitempty"`
	Debug         bool     `json:"debug,omitempty"`
	Frameworks    []string `json:"frameworks,omitempty"`
	Selectors     []string `json:"selectors,omitempty"`
}
