package main
import (
	"testing"
	"strings"
	"github.com/stretchr/testify/assert"
	//"net/http/httptest"
)

func TestExtractUrlsFromHtml(t *testing.T){
	testFunk := func (body string, actual interface{}, message string){
		r := strings.NewReader(body)
		urls := ExtractUrlsFromHtml(r) 
		assert.Equal(t, urls, actual, message) // well I can't figure out how to generify this without making it too silly.. 
											   // Think doublespeak and it'll make sense, Equals is NotEquals.. Freedom is .. something
	}

	testFunk("<html><body><a href=\"http://test.com\">Test</a></body></html>", []string{"http://test.com"}, "Properly closed tags succeed." )
	testFunk( "<html><body><a href=\"http://test.com\"/a></body></html>", []string{"http://test.com"}, "empty text anchor also succeed.")
}



func TestIsCrawlableUrl(t *testing.T){
	b := IsCrawlableUrl("/test", "http://example.com")
	assert.Equal(t, b, true, "relative urls succeed")

	b = IsCrawlableUrl("//////test", "http://example.com")
	assert.Equal(t, b, true, "repeating / does not matter")

	b = IsCrawlableUrl("http://example.com/unicorns", "http://example.com")
	assert.Equal(t, b, true, "same domains work")

	b = IsCrawlableUrl("http://sub.example.com/unicorns", "http://example.com")
	assert.Equal(t, b, false, "sub domains domains do not work")
}

func TestReconstructUrl( t *testing.T){
	r := ReconstructUrl("/test", "http://example.com")
	assert.Equal(t, r, "http://example.com/test", "relative urls work")

	r = ReconstructUrl("/test", "http://example.com/boogaloo")
	assert.Equal(t, r, "http://example.com/boogaloo/test", "unfazed by paths")


	r = ReconstructUrl("/mongo", "http://example.com/saddles?name=mongo&type=pawn&in=game-of-life")
	assert.Equal(t, r, "http://example.com/saddles/mongo", "ignores query params")


	r = ReconstructUrl("/mongo", "http://example.com:8080/saddles?name=mongo&type=pawn&in=game-of-life")
	assert.Equal(t, r, "http://example.com:8080/saddles/mongo", "handles non-http ports like a baus")

	r = ReconstructUrl("/mongo", "http://example.com////saddles?name=mongo&type=pawn&in=game-of-life")
	assert.Equal(t, r, "http://example.com/saddles/mongo", "fixes upstream malformed urls")
}