package main
import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
)

func TestExtractUrlsFromHtml(t *testing.T){
	// Succeeds on properly ended tags
	r := strings.NewReader("<html><body><a href=\"http://test.com\">Test</a></body></html>")
	urls := ExtractUrlsFromHtml(r) 
	assert.Equal(t, urls, []string{"http://test.com"}, "Properly closed tags succeed.")

	r2 := strings.NewReader("<html><body><a href=\"http://test.com\"/></body></html>")
	urls2 := ExtractUrlsFromHtml(r2) 
	assert.NotEqual(t, urls2, []string{"http://test.com"}, "empty a tags fail.")
}