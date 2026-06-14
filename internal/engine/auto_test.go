package engine

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubEngine is a controllable Engine for testing the auto composite.
type stubEngine struct {
	name   string
	res    Result
	err    error
	called *int
}

func (s *stubEngine) Name() string { return s.name }
func (s *stubEngine) Fetch(_ context.Context, _ string, _ Options) (Result, error) {
	if s.called != nil {
		*s.called++
	}
	if s.err != nil {
		return Result{}, s.err
	}
	return s.res, nil
}

func longText(n int) string { return strings.Repeat("essence content ", n) }

func TestAutoUsesFirstSufficientEngine(t *testing.T) {
	var c1, c2 int
	auto := NewAutoEngine(
		&stubEngine{name: "http", res: Result{HTML: "<p>" + longText(50) + "</p>", Engine: "http"}, called: &c1},
		&stubEngine{name: "lightpanda", res: Result{HTML: "x", Engine: "lightpanda"}, called: &c2},
	)
	res, err := auto.Fetch(context.Background(), "http://x", Options{})
	require.NoError(t, err)
	assert.Equal(t, "http", res.Engine)
	assert.Equal(t, 1, c1)
	assert.Equal(t, 0, c2, "should not escalate when first engine is sufficient")
}

func TestAutoEscalatesOnThinJSDependentPage(t *testing.T) {
	var c1, c2 int
	auto := NewAutoEngine(
		&stubEngine{name: "http", res: Result{
			HTML:   "<div id=app></div>",
			Raw:    `<html><body><div id=app></div><script>render()</script></body></html>`,
			Engine: "http",
		}, called: &c1},
		&stubEngine{name: "lightpanda", res: Result{HTML: "<p>" + longText(50) + "</p>", Engine: "lightpanda"}, called: &c2},
	)
	res, err := auto.Fetch(context.Background(), "http://x", Options{})
	require.NoError(t, err)
	assert.Equal(t, "lightpanda", res.Engine, "thin JS-dependent page should escalate")
	assert.Equal(t, 1, c1)
	assert.Equal(t, 1, c2)
}

func TestAutoDoesNotEscalateThinStaticPage(t *testing.T) {
	var c1, c2 int
	auto := NewAutoEngine(
		&stubEngine{name: "http", res: Result{
			HTML:   "<html><body>Short but no scripts here</body></html>",
			Raw:    "<html><body>Short but no scripts here</body></html>",
			Engine: "http",
		}, called: &c1},
		&stubEngine{name: "lightpanda", res: Result{HTML: "y", Engine: "lightpanda"}, called: &c2},
	)
	res, err := auto.Fetch(context.Background(), "http://x", Options{})
	require.NoError(t, err)
	assert.Equal(t, "http", res.Engine, "thin page without scripts should NOT escalate")
	assert.Equal(t, 0, c2)
}

func TestAutoEscalatesOnError(t *testing.T) {
	var c2 int
	auto := NewAutoEngine(
		&stubEngine{name: "http", err: errors.New("boom")},
		&stubEngine{name: "lightpanda", res: Result{HTML: "<p>" + longText(50) + "</p>", Engine: "lightpanda"}, called: &c2},
	)
	res, err := auto.Fetch(context.Background(), "http://x", Options{})
	require.NoError(t, err)
	assert.Equal(t, "lightpanda", res.Engine)
	assert.Equal(t, 1, c2)
}

func TestAutoAcceptsLastEngineEvenIfThin(t *testing.T) {
	auto := NewAutoEngine(
		&stubEngine{name: "http", res: Result{HTML: "<div></div>", Raw: "<script></script>", Engine: "http"}},
		&stubEngine{name: "chrome", res: Result{HTML: "<div>tiny</div>", Raw: "<script></script>", Engine: "chrome"}},
	)
	res, err := auto.Fetch(context.Background(), "http://x", Options{})
	require.NoError(t, err)
	assert.Equal(t, "chrome", res.Engine, "last engine is accepted even when thin")
}

func TestAutoAllFail(t *testing.T) {
	auto := NewAutoEngine(
		&stubEngine{name: "http", err: errors.New("e1")},
		&stubEngine{name: "lightpanda", err: errors.New("e2")},
	)
	_, err := auto.Fetch(context.Background(), "http://x", Options{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "all engines failed")
}

func TestVisibleTextDropsScripts(t *testing.T) {
	html := `<html><body><script>var x = "lots of code here that is not visible";</script><p>Hello</p></body></html>`
	assert.Equal(t, "Hello", visibleText(html))
}
