package specs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTLDRSummarizerSpec validates F6: TL;DR Article Summarization
//
// SPEC: sz tldr command generates concise article summaries using AI,
// leveraging the complete F1-F5 content extraction pipeline and adding
// intelligent summarization as the final processing step.
func TestTLDRSummarizerSpec(t *testing.T) {
	// Skip if no API key available for CI/testing
	if os.Getenv("OPENAI_API_KEY") == "" && os.Getenv("TEST_TLDR") == "" {
		t.Skip("Skipping TL;DR tests - set OPENAI_API_KEY or TEST_TLDR=1 to run")
	}

	// Build test binary
	szBinary := buildTLDRBinary(t)
	defer func() { _ = os.Remove(szBinary) }()

	t.Run("basic_tldr_functionality", func(t *testing.T) {
		t.Log("SPEC: Basic TL;DR Functionality")
		t.Log("GIVEN a URL with article content")
		t.Log("WHEN user runs `sz tldr https://example.com`")
		t.Log("THEN it should generate a concise summary")

		// Create a detailed article for summarization
		articleHTML := `<!DOCTYPE html>
<html>
<head>
    <title>The Future of Web Development</title>
</head>
<body>
    <header>
        <nav>
            <a href="/">Home</a>
            <a href="/about">About</a>
        </nav>
    </header>

    <main>
        <article>
            <h1>The Future of Web Development</h1>
            <p>Web development is rapidly evolving with new technologies and frameworks emerging constantly. In this comprehensive analysis, we explore the key trends that will shape the future of web development over the next decade.</p>

            <h2>Emerging Technologies</h2>
            <p>Several technologies are revolutionizing how we build web applications. WebAssembly (WASM) allows developers to run high-performance code in browsers, enabling applications that were previously impossible. Progressive Web Apps (PWAs) bridge the gap between web and native applications, providing offline functionality and native-like user experiences.</p>

            <h2>Framework Evolution</h2>
            <p>JavaScript frameworks continue to evolve, with React, Vue, and Angular leading the way. Server-side rendering (SSR) and static site generation (SSG) are becoming standard practices for performance optimization. The rise of full-stack frameworks like Next.js and Nuxt.js demonstrates the trend toward comprehensive development solutions.</p>

            <h2>Developer Experience</h2>
            <p>Developer experience (DX) has become a primary focus, with tools like Vite providing lightning-fast development builds. TypeScript adoption continues to grow, offering better code reliability and developer productivity. The integration of AI-powered development tools is beginning to transform how developers write and debug code.</p>

            <h2>Performance Optimization</h2>
            <p>Web performance remains critical, with Core Web Vitals becoming important ranking factors for search engines. Techniques like code splitting, lazy loading, and edge computing are essential for delivering fast user experiences. The shift toward edge computing brings computation closer to users, reducing latency significantly.</p>

            <h2>Conclusion</h2>
            <p>The future of web development is bright, with exciting technologies enabling more powerful and efficient applications. Developers who stay current with these trends will be well-positioned to build the next generation of web experiences.</p>
        </article>
    </main>

    <footer>
        <p>© 2024 Tech Blog. All rights reserved.</p>
    </footer>
</body>
</html>`

		// Create test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(articleHTML))
		}))
		defer server.Close()

		// Test basic tldr functionality
		cmd := exec.Command(szBinary, "tldr", server.URL)
		if os.Getenv("OPENAI_API_KEY") != "" {
			cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"))
		}
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "TL;DR command should succeed")

		outputStr := string(output)

		// Should generate a summary
		assert.NotEmpty(t, outputStr, "Should generate summary output")

		// Summary should be significantly shorter than original
		assert.True(t, len(outputStr) < len(articleHTML)/2, "Summary should be shorter than original content")

		// Should contain key concepts from the article
		assert.Contains(t, strings.ToLower(outputStr), "web development", "Should mention web development")
		// Check for technology-related terms (technology, technologies, or specific tech like webassembly)
		containsTech := strings.Contains(strings.ToLower(outputStr), "technology") ||
			strings.Contains(strings.ToLower(outputStr), "technologies") ||
			strings.Contains(strings.ToLower(outputStr), "webassembly") ||
			strings.Contains(strings.ToLower(outputStr), "javascript")
		assert.True(t, containsTech, "Should mention technology-related concepts")

		// Should not contain navigation/footer content
		assert.NotContains(t, outputStr, "Home", "Should not include navigation")
		assert.NotContains(t, outputStr, "© 2024", "Should not include footer")
	})

	t.Run("tldr_with_summary_length_options", func(t *testing.T) {
		t.Log("SPEC: TL;DR with Summary Length Options")
		t.Log("GIVEN a URL with article content")
		t.Log("WHEN user runs `sz tldr --summary-length short https://example.com`")
		t.Log("THEN it should generate a shorter summary")

		articleHTML := createLongArticle()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(articleHTML))
		}))
		defer server.Close()

		// Test short summary
		cmd := exec.Command(szBinary, "tldr", "--summary-length", "short", server.URL)
		if os.Getenv("OPENAI_API_KEY") != "" {
			cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"))
		}
		shortOutput, err := cmd.CombinedOutput()
		require.NoError(t, err, "Short summary should succeed")

		// Test medium summary
		cmd = exec.Command(szBinary, "tldr", "--summary-length", "medium", server.URL)
		if os.Getenv("OPENAI_API_KEY") != "" {
			cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"))
		}
		mediumOutput, err := cmd.CombinedOutput()
		require.NoError(t, err, "Medium summary should succeed")

		// Test long summary
		cmd = exec.Command(szBinary, "tldr", "--summary-length", "long", server.URL)
		if os.Getenv("OPENAI_API_KEY") != "" {
			cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"))
		}
		longOutput, err := cmd.CombinedOutput()
		require.NoError(t, err, "Long summary should succeed")

		// Verify length relationships
		shortLen := len(string(shortOutput))
		mediumLen := len(string(mediumOutput))
		longLen := len(string(longOutput))

		assert.True(t, shortLen < mediumLen, "Short summary should be shorter than medium")
		assert.True(t, mediumLen < longLen, "Medium summary should be shorter than long")
	})

	t.Run("tldr_with_api_key_flag", func(t *testing.T) {
		t.Log("SPEC: TL;DR with API Key Flag")
		t.Log("GIVEN a URL with article content and an API key")
		t.Log("WHEN user runs `sz tldr --api-key <key> https://example.com`")
		t.Log("THEN it should use the provided API key")

		if os.Getenv("OPENAI_API_KEY") == "" {
			t.Skip("Skipping API key test - no OPENAI_API_KEY available")
		}

		articleHTML := `<html><body><article><h1>Test Article</h1><p>This is a test article for API key functionality.</p></article></body></html>`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(articleHTML))
		}))
		defer server.Close()

		cmd := exec.Command(szBinary, "tldr", "--api-key", os.Getenv("OPENAI_API_KEY"), server.URL)
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "TL;DR with API key flag should succeed")

		outputStr := string(output)
		assert.NotEmpty(t, outputStr, "Should generate summary with API key flag")
	})

	t.Run("tldr_with_model_selection", func(t *testing.T) {
		t.Log("SPEC: TL;DR with Model Selection")
		t.Log("GIVEN a URL with article content")
		t.Log("WHEN user runs `sz tldr --model gpt-3.5-turbo https://example.com`")
		t.Log("THEN it should use the specified model")

		articleHTML := `<html><body><article><h1>Model Test</h1><p>Testing model selection functionality.</p></article></body></html>`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(articleHTML))
		}))
		defer server.Close()

		cmd := exec.Command(szBinary, "tldr", "--model", "gpt-3.5-turbo", server.URL)
		if os.Getenv("OPENAI_API_KEY") != "" {
			cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"))
		}
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "TL;DR with model selection should succeed")

		outputStr := string(output)
		assert.NotEmpty(t, outputStr, "Should generate summary with specified model")
	})

	t.Run("tldr_error_handling", func(t *testing.T) {
		t.Log("SPEC: TL;DR Error Handling")
		t.Log("GIVEN no API key is available")
		t.Log("WHEN user runs `sz tldr https://example.com`")
		t.Log("THEN it should show helpful error message")

		articleHTML := `<html><body><article><h1>Error Test</h1><p>Testing error handling.</p></article></body></html>`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(articleHTML))
		}))
		defer server.Close()

		cmd := exec.Command(szBinary, "tldr", server.URL)
		// Explicitly remove API key from environment
		cmd.Env = []string{}
		output, err := cmd.CombinedOutput()

		// Should exit with error
		assert.Error(t, err, "Should fail without API key")

		outputStr := string(output)
		assert.Contains(t, strings.ToLower(outputStr), "api", "Error message should mention API")
		assert.Contains(t, strings.ToLower(outputStr), "key", "Error message should mention key")
	})

	t.Run("tldr_with_complex_content", func(t *testing.T) {
		t.Log("SPEC: TL;DR with Complex Content")
		t.Log("GIVEN a URL with complex article structure (images, lists, quotes)")
		t.Log("WHEN user runs `sz tldr https://example.com`")
		t.Log("THEN it should handle complex content and generate coherent summary")

		complexHTML := `<!DOCTYPE html>
<html>
<body>
    <article>
        <h1>Complex Article Structure</h1>
        <p>This article demonstrates various content types that should be handled properly by the summarization system.</p>

        <img src="chart.jpg" alt="Performance metrics showing 40% improvement">

        <h2>Key Points</h2>
        <ul>
            <li>Performance improved by 40% through optimization</li>
            <li>User engagement increased significantly</li>
            <li>Cost reduction achieved through automation</li>
        </ul>

        <blockquote>
            "The results exceeded our expectations and demonstrated the value of systematic optimization." - Lead Engineer
        </blockquote>

        <p>The implementation involved multiple phases including analysis, development, testing, and deployment. Each phase contributed to the overall success of the project.</p>

        <h2>Technical Details</h2>
        <p>Technical implementation focused on several key areas:</p>
        <ol>
            <li>Database query optimization</li>
            <li>Caching layer implementation</li>
            <li>Frontend performance tuning</li>
        </ol>

        <p>These changes resulted in measurable improvements across all metrics.</p>
    </article>
</body>
</html>`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(complexHTML))
		}))
		defer server.Close()

		cmd := exec.Command(szBinary, "tldr", server.URL)
		if os.Getenv("OPENAI_API_KEY") != "" {
			cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"))
		}
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Complex content TL;DR should succeed")

		outputStr := string(output)
		assert.NotEmpty(t, outputStr, "Should generate summary for complex content")

		// Should capture key information from various content types
		assert.Contains(t, strings.ToLower(outputStr), "performance", "Should capture performance improvements")
		assert.Contains(t, strings.ToLower(outputStr), "optimization", "Should capture optimization theme")
	})

	t.Run("tldr_integration_with_f1_f5_pipeline", func(t *testing.T) {
		t.Log("SPEC: TL;DR Integration with F1-F5 Pipeline")
		t.Log("GIVEN a URL with JavaScript-heavy content requiring F1-F5 processing")
		t.Log("WHEN user runs `sz tldr https://example.com`")
		t.Log("THEN it should use the complete extraction pipeline before summarization")

		// Create content that would benefit from F1-F5 processing
		modernWebHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Modern Framework Article</title>
    <script>
        // Simulate framework hydration
        document.addEventListener('DOMContentLoaded', function() {
            const content = document.getElementById('dynamic-content');
            content.innerHTML = '<p>This content was loaded dynamically and should be captured by F1 DOM ready detection.</p>';
        });
    </script>
</head>
<body>
    <nav class="navigation">
        <a href="/">Home</a>
        <a href="/articles">Articles</a>
    </nav>

    <aside class="sidebar">
        <div class="ad-banner">Advertisement</div>
        <div class="social-links">Follow us</div>
    </aside>

    <main class="main-content">
        <article class="post-content">
            <h1>JavaScript Framework Evolution</h1>
            <p>JavaScript frameworks have evolved significantly over the past decade, with React, Vue, and Angular leading the charge in modern web development.</p>

            <div id="dynamic-content">
                <!-- This will be populated by JavaScript -->
            </div>

            <img src="framework-timeline.png" alt="Timeline showing the evolution of JavaScript frameworks from 2010 to 2024">

            <p>The adoption of TypeScript has accelerated, providing better developer experience and code reliability across these frameworks.</p>
        </article>
    </main>

    <footer class="site-footer">
        <p>© 2024 Developer Blog</p>
    </footer>
</body>
</html>`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(modernWebHTML))
		}))
		defer server.Close()

		cmd := exec.Command(szBinary, "tldr", server.URL)
		if os.Getenv("OPENAI_API_KEY") != "" {
			cmd.Env = append(os.Environ(), "OPENAI_API_KEY="+os.Getenv("OPENAI_API_KEY"))
		}
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "F1-F5 pipeline integration should succeed")

		outputStr := string(output)
		assert.NotEmpty(t, outputStr, "Should generate summary using F1-F5 pipeline")

		// Should include main content
		assert.Contains(t, strings.ToLower(outputStr), "javascript", "Should include main article content")
		assert.Contains(t, strings.ToLower(outputStr), "framework", "Should include framework discussion")

		// Should exclude navigation and ads (F3 content filtering)
		assert.NotContains(t, outputStr, "Home", "Should filter navigation via F3")
		assert.NotContains(t, outputStr, "Advertisement", "Should filter ads via F3")
		assert.NotContains(t, outputStr, "© 2024", "Should filter footer via F3")
	})
}

// createLongArticle creates a longer article for testing summary length variations
func createLongArticle() string {
	return `<!DOCTYPE html>
<html>
<head>
    <title>The Complete Guide to Modern Software Architecture</title>
</head>
<body>
    <article>
        <h1>The Complete Guide to Modern Software Architecture</h1>
        <p>Software architecture has undergone significant transformation in recent years, driven by cloud computing, microservices, and the need for scalable, maintainable systems. This comprehensive guide explores the fundamental principles and practices that define modern software architecture.</p>

        <h2>Microservices Architecture</h2>
        <p>Microservices architecture represents a fundamental shift from monolithic applications to distributed systems composed of small, independent services. Each microservice is responsible for a specific business capability and can be developed, deployed, and scaled independently. This approach offers numerous advantages including improved scalability, technology diversity, and fault isolation.</p>

        <h2>Cloud-Native Design Principles</h2>
        <p>Cloud-native architecture emphasizes building applications specifically for cloud environments. Key principles include containerization using Docker and Kubernetes, serverless computing with functions-as-a-service, and event-driven architectures that respond to changes in real-time. These principles enable organizations to leverage the full potential of cloud platforms.</p>

        <h2>Data Architecture Patterns</h2>
        <p>Modern data architecture requires careful consideration of data storage, processing, and access patterns. Popular approaches include data lakes for storing vast amounts of unstructured data, data warehouses for analytical workloads, and real-time streaming architectures for processing events as they occur. The choice of data architecture significantly impacts application performance and scalability.</p>

        <h2>Security and Compliance</h2>
        <p>Security must be integrated into every layer of modern software architecture. This includes implementing zero-trust security models, encrypting data both in transit and at rest, and ensuring compliance with regulations such as GDPR and HIPAA. Security considerations should be addressed from the initial design phase rather than added as an afterthought.</p>

        <h2>Observability and Monitoring</h2>
        <p>Modern applications require comprehensive observability to understand system behavior and performance. This includes logging, metrics collection, distributed tracing, and alerting systems. Observability enables teams to quickly identify and resolve issues, optimize performance, and ensure reliable operation of complex distributed systems.</p>

        <h2>DevOps Integration</h2>
        <p>Successful modern architecture requires tight integration with DevOps practices. This includes infrastructure as code, continuous integration and deployment pipelines, automated testing, and collaborative development practices. DevOps enables rapid, reliable delivery of software while maintaining high quality and stability.</p>

        <h2>Future Trends</h2>
        <p>The future of software architecture will likely be shaped by artificial intelligence, edge computing, and quantum computing. These emerging technologies will require new architectural patterns and approaches. Organizations should prepare for these changes by building flexible, adaptable architectures that can evolve with technological advancement.</p>

        <h2>Conclusion</h2>
        <p>Modern software architecture requires balancing multiple competing concerns including performance, scalability, maintainability, and security. Success depends on choosing the right architectural patterns for specific requirements and maintaining flexibility to adapt as needs evolve. Organizations that invest in sound architectural practices will be better positioned for future growth and innovation.</p>
    </article>
</body>
</html>`
}

// buildTLDRBinary builds the sz binary for testing TL;DR functionality
func buildTLDRBinary(t *testing.T) string {
	// Build the sz binary from project root
	cmd := exec.Command("go", "build", "-o", "/tmp/sz-tldr-test", "./cmd/essenz")
	// Set working directory to project root
	cmd.Dir = ".."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build binary for TL;DR testing")

	return "/tmp/sz-tldr-test"
}
