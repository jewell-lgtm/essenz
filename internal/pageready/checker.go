// Package pageready provides DOM readiness detection for reliable content extraction.
package pageready

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// ReadinessChecker manages DOM readiness detection for web pages.
type ReadinessChecker struct {
	MaxWaitTime     time.Duration
	FrameworkHints  []string
	CustomSelectors []string
	Debug           bool
	UseLCP          bool
}

// ReadinessResult contains information about page readiness detection.
type ReadinessResult struct {
	IsReady   bool
	EventType string
	WaitTime  time.Duration
	Error     error
	DebugInfo string
}

// NewReadinessChecker creates a new readiness checker with default settings.
func NewReadinessChecker() *ReadinessChecker {
	return &ReadinessChecker{
		MaxWaitTime:     5 * time.Second,
		FrameworkHints:  []string{},
		CustomSelectors: []string{},
		Debug:           false,
		UseLCP:          true, // default: LCP-based readiness
	}
}

// WithTimeout sets the maximum wait time for readiness detection.
func (r *ReadinessChecker) WithTimeout(timeout time.Duration) *ReadinessChecker {
	r.MaxWaitTime = timeout
	return r
}

// WithFrameworkHints sets framework-specific detection patterns.
func (r *ReadinessChecker) WithFrameworkHints(hints []string) *ReadinessChecker {
	r.FrameworkHints = hints
	return r
}

// WithCustomSelectors sets custom CSS selectors to wait for.
func (r *ReadinessChecker) WithCustomSelectors(selectors []string) *ReadinessChecker {
	r.CustomSelectors = selectors
	return r
}

// WithDebug enables debug information collection.
func (r *ReadinessChecker) WithDebug(debug bool) *ReadinessChecker {
	r.Debug = debug
	return r
}

// WithLCP toggles Largest Contentful Paint readiness.
func (r *ReadinessChecker) WithLCP(enabled bool) *ReadinessChecker {
	r.UseLCP = enabled
	return r
}

// WaitForReady waits for the page to be ready according to configured criteria.
func (r *ReadinessChecker) WaitForReady(ctx context.Context, chromeCtx context.Context) (*ReadinessResult, error) {
	start := time.Now()

	// Create a timeout context
	timeoutCtx, cancel := context.WithTimeout(ctx, r.MaxWaitTime)
	defer cancel()

	result := &ReadinessResult{
		IsReady:   false,
		EventType: "unknown",
		WaitTime:  0,
		Error:     nil,
		DebugInfo: "",
	}

	var err error

	// Prefer LCP readiness when enabled
	if r.UseLCP {
		if r.Debug {
			result.DebugInfo += "Waiting for LCP; "
		}
		if err = r.waitForLCP(timeoutCtx, chromeCtx, result); err != nil {
			if r.Debug {
				result.DebugInfo += fmt.Sprintf("LCP wait failed: %v; falling back; ", err)
			}
			// fall back to basic DOM ready
			err = r.waitForBasicDOMReady(timeoutCtx, chromeCtx, result)
		}
	} else {
		// Start with basic DOM ready detection
		err = r.waitForBasicDOMReady(timeoutCtx, chromeCtx, result)
	}
	if err != nil {
		result.Error = err
		result.WaitTime = time.Since(start)
		return result, err
	}

	// If we have custom selectors, wait for them
	if len(r.CustomSelectors) > 0 {
		err = r.waitForCustomSelectors(timeoutCtx, chromeCtx, result)
		if err != nil {
			result.Error = err
			result.WaitTime = time.Since(start)
			return result, err
		}
	}

	// If we have framework hints, try to detect framework readiness
	if len(r.FrameworkHints) > 0 {
		err = r.waitForFrameworkReady(timeoutCtx, chromeCtx, result)
		if err != nil {
			// Framework detection failure is not fatal - continue with basic readiness
			if r.Debug {
				result.DebugInfo += fmt.Sprintf("Framework detection failed: %v; ", err)
			}
		}
	}

	result.IsReady = true
	result.WaitTime = time.Since(start)

	if r.Debug {
		result.DebugInfo += fmt.Sprintf("Ready after %v", result.WaitTime)
	}

	return result, nil
}

// waitForLCP injects an observer and waits until the first LCP entry is observed.
func (r *ReadinessChecker) waitForLCP(_ context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	// Inject observer early
	var _ignored bool
	inject := `(function(){
      try{
        if (window.__essenzLCPInjected__) return true;
        window.__essenzLCPInjected__=true;window.__essenzLCP__=null;window.__essenzLCP_READY__=false;
        function S(e){var el=e.element||null;var o={startTime:e.startTime,url:e.url||null,size:e.size||null,tag:null,id:null,class:null,text:null,src:null};
          if(el){o.tag=(el.tagName||'').toLowerCase();o.id=el.id||null;o.class=el.className||null; if(o.tag==='img'&&el.currentSrc){o.src=el.currentSrc;} if(el.textContent&&o.tag!=='img'){o.text=(el.textContent||'').trim().slice(0,200);} }
          return o; }
        const po=new PerformanceObserver(list=>{ for(const e of list.getEntries()){ window.__essenzLCP__=S(e); window.__essenzLCP_READY__=true; }});
        po.observe({type:'largest-contentful-paint', buffered:true});
        document.addEventListener('visibilitychange',function(){ if(document.visibilityState==='hidden'){ try{po.disconnect();}catch(_){}}});
        return true;
      }catch(_){return false}
    })();`
	if err := chromedp.Run(chromeCtx, chromedp.EvaluateAsDevTools(inject, &_ignored)); err != nil {
		return fmt.Errorf("failed to inject LCP observer: %w", err)
	}
	// Poll for readiness flag
	deadline := time.Now().Add(r.MaxWaitTime)
	for time.Now().Before(deadline) {
		var ready bool
		if err := chromedp.Run(chromeCtx, chromedp.EvaluateAsDevTools("Boolean(window.__essenzLCP_READY__)||false", &ready)); err == nil && ready {
			result.EventType = "LCP"
			if r.Debug {
				result.DebugInfo += "LCP observed; "
			}
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for LCP")
}

// waitForBasicDOMReady waits for the basic DOM to be ready.
func (r *ReadinessChecker) waitForBasicDOMReady(_ context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	// Wait for DOMContentLoaded equivalent
	err := chromedp.Run(chromeCtx,
		chromedp.WaitReady("body"),
	)

	if err != nil {
		return fmt.Errorf("DOM ready timeout: %w", err)
	}

	result.EventType = "DOMContentLoaded"

	if r.Debug {
		result.DebugInfo += "DOM content loaded; "
	}

	return nil
}

// waitForCustomSelectors waits for custom CSS selectors to appear.
func (r *ReadinessChecker) waitForCustomSelectors(_ context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	for _, selector := range r.CustomSelectors {
		err := chromedp.Run(chromeCtx,
			chromedp.WaitVisible(selector),
		)

		if err != nil {
			return fmt.Errorf("custom selector '%s' not found: %w", selector, err)
		}

		if r.Debug {
			result.DebugInfo += fmt.Sprintf("Custom selector '%s' found; ", selector)
		}
	}

	result.EventType = "CustomSelector"
	return nil
}

// waitForFrameworkReady attempts to detect JavaScript framework readiness.
func (r *ReadinessChecker) waitForFrameworkReady(ctx context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	for _, hint := range r.FrameworkHints {
		switch strings.ToLower(hint) {
		case "react":
			if err := r.waitForReactReady(ctx, chromeCtx, result); err == nil {
				result.EventType = "ReactReady"
				return nil
			}
		case "vue":
			if err := r.waitForVueReady(ctx, chromeCtx, result); err == nil {
				result.EventType = "VueReady"
				return nil
			}
		case "angular":
			if err := r.waitForAngularReady(ctx, chromeCtx, result); err == nil {
				result.EventType = "AngularReady"
				return nil
			}
		case "nextjs":
			if err := r.waitForNextJSReady(ctx, chromeCtx, result); err == nil {
				result.EventType = "NextJSReady"
				return nil
			}
		}
	}

	return fmt.Errorf("no supported frameworks detected")
}

// waitForReactReady waits for React app to be ready.
func (r *ReadinessChecker) waitForReactReady(_ context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	var isReady bool

	// Try multiple approaches to detect React readiness
	err := chromedp.Run(chromeCtx,
		// Check if React is loaded
		chromedp.EvaluateAsDevTools(`
			(function() {
				// Check for React in global scope
				if (window.React || window.ReactDOM) {
					return true;
				}

				// Check for React dev tools
				if (window.__REACT_DEVTOOLS_GLOBAL_HOOK__) {
					return true;
				}

				// Check for common React patterns
				const reactElements = document.querySelectorAll('[data-reactroot], [data-react-]');
				if (reactElements.length > 0) {
					return true;
				}

				// Check for React Fiber
				const rootElement = document.querySelector('#root, [id*="react"], [class*="react"]');
				if (rootElement && rootElement._reactInternalFiber) {
					return true;
				}

				return false;
			})();
		`, &isReady),
	)

	if err != nil {
		return fmt.Errorf("React detection failed: %w", err)
	}

	if !isReady {
		return fmt.Errorf("React not detected")
	}

	// Wait a bit for React hydration to complete
	time.Sleep(500 * time.Millisecond)

	if r.Debug {
		result.DebugInfo += "React framework detected; "
	}

	return nil
}

// waitForVueReady waits for Vue.js app to be ready.
func (r *ReadinessChecker) waitForVueReady(_ context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	var isReady bool

	err := chromedp.Run(chromeCtx,
		chromedp.EvaluateAsDevTools(`
			(function() {
				// Check for Vue in global scope
				if (window.Vue) {
					return true;
				}

				// Check for Vue dev tools
				if (window.__VUE__) {
					return true;
				}

				// Check for Vue instances in DOM
				const vueElements = document.querySelectorAll('[data-v-]');
				if (vueElements.length > 0) {
					return true;
				}

				return false;
			})();
		`, &isReady),
	)

	if err != nil {
		return fmt.Errorf("Vue detection failed: %w", err)
	}

	if !isReady {
		return fmt.Errorf("Vue not detected")
	}

	if r.Debug {
		result.DebugInfo += "Vue framework detected; "
	}

	return nil
}

// waitForAngularReady waits for Angular app to be ready.
func (r *ReadinessChecker) waitForAngularReady(_ context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	var isReady bool

	err := chromedp.Run(chromeCtx,
		chromedp.EvaluateAsDevTools(`
			(function() {
				// Check for Angular in global scope
				if (window.ng || window.angular) {
					return true;
				}

				// Check for Angular elements
				const ngElements = document.querySelectorAll('[ng-app], [data-ng-app], app-root');
				if (ngElements.length > 0) {
					return true;
				}

				return false;
			})();
		`, &isReady),
	)

	if err != nil {
		return fmt.Errorf("Angular detection failed: %w", err)
	}

	if !isReady {
		return fmt.Errorf("Angular not detected")
	}

	if r.Debug {
		result.DebugInfo += "Angular framework detected; "
	}

	return nil
}

// waitForNextJSReady waits for Next.js app to be ready.
func (r *ReadinessChecker) waitForNextJSReady(_ context.Context, chromeCtx context.Context, result *ReadinessResult) error {
	var isReady bool

	err := chromedp.Run(chromeCtx,
		chromedp.EvaluateAsDevTools(`
			(function() {
				// Check for Next.js specific indicators
				if (window.__NEXT_DATA__) {
					return true;
				}

				// Check for Next.js build ID
				if (window.__BUILD_MANIFEST) {
					return true;
				}

				// Check for Next.js router
				if (window.next && window.next.router) {
					return true;
				}

				// Check for common Next.js patterns
				const nextElements = document.querySelectorAll('#__next');
				if (nextElements.length > 0) {
					return true;
				}

				return false;
			})();
		`, &isReady),
	)

	if err != nil {
		return fmt.Errorf("Next.js detection failed: %w", err)
	}

	if !isReady {
		return fmt.Errorf("Next.js not detected")
	}

	// Wait for hydration to complete
	time.Sleep(1 * time.Second)

	if r.Debug {
		result.DebugInfo += "Next.js framework detected; "
	}

	return nil
}
