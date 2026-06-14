package extractor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const pageHTML = `<!DOCTYPE html><html><head><title>Important Article Title</title></head><body>
<header><nav><a href="/home">Home</a></nav></header>
<aside class="sidebar"><div class="ad">Advertisement Content</div></aside>
<main><article>
<h1>Important Article Title</h1>
<p class="byline">By John Doe</p>
<p>This is the main article content that should be extracted. It contains important information that the user wants to read about lightweight browsers and content extraction.</p>
<h2>Section Header</h2>
<p>More content under the section header providing valuable insights worth keeping.</p>
</article></main>
<footer><p>Copyright 2024</p></footer></body></html>`

func TestExtractContentKeepsTitleAndContentDropsBoilerplate(t *testing.T) {
	md, err := New().ExtractContent(pageHTML)
	require.NoError(t, err)

	assert.Contains(t, md, "Important Article Title", "title must survive")
	assert.Contains(t, md, "main article content", "body must survive")
	assert.Contains(t, md, "Section Header")
	assert.Contains(t, md, "valuable insights")

	assert.NotContains(t, md, "Advertisement Content", "ads dropped")
	assert.NotContains(t, md, "Copyright 2024", "footer dropped")
	assert.NotContains(t, md, "Home", "nav dropped")
}

func TestExtractContentFallsBackOnMinimalPage(t *testing.T) {
	md, err := New().ExtractContent(`<html><body><div>Minimal content</div></body></html>`)
	require.NoError(t, err)
	assert.Contains(t, md, "Minimal content", "minimal pages should still produce output")
}

func TestExtractContentTitleNotDuplicated(t *testing.T) {
	md, err := New().ExtractContent(pageHTML)
	require.NoError(t, err)
	assert.Equal(t, 1, strings.Count(md, "Important Article Title"), "title should appear once")
}
