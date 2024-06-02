package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	db "github.com/phelipperibeiro/technical-challenges-rate-limiter/adapter/db/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebServerTestSuite struct {
	suite.Suite
	router           http.Handler
	maxIPRequests    int
	maxTokenRequests int
	blockDuration    int
	redisHost        string
	redisPort        string
	cache            *db.RedisCache
}

func (suite *WebServerTestSuite) SetupSuite() {
	suite.maxIPRequests = 10
	suite.maxTokenRequests = 100
	suite.blockDuration = 5
	suite.redisHost = "localhost"
	suite.redisPort = "6379"

	log.Println("Initializing Redis cache...")
	redisCache, err := db.NewRedis(fmt.Sprintf("%s:%s", suite.redisHost, suite.redisPort))
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	suite.cache = redisCache

	router := NewWebServer(suite.maxIPRequests, suite.maxTokenRequests, suite.blockDuration, redisCache)
	suite.router = router
}

func (suite *WebServerTestSuite) SetupTest() {
	suite.cache.Clear()
}

func (suite *WebServerTestSuite) TestWebServerIsRunning() {
	server := httptest.NewServer(suite.router)
	defer server.Close()

	resp, err := http.Get(server.URL)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *WebServerTestSuite) TestIPRateLimiting() {
	server := httptest.NewServer(suite.router)
	defer server.Close()

	var okResponses, blockedResponses int
	const numberOfRequests = 100

	for i := 0; i < numberOfRequests; i++ {
		resp, err := http.Get(server.URL)
		assert.NoError(suite.T(), err)

		if resp.StatusCode == http.StatusOK {
			okResponses++
		} else if resp.StatusCode == http.StatusTooManyRequests {
			blockedResponses++
		}
	}

	assert.Equal(suite.T(), suite.maxIPRequests, okResponses)
	assert.Equal(suite.T(), numberOfRequests-suite.maxIPRequests, blockedResponses)
}

func (suite *WebServerTestSuite) TestTokenRateLimiting() {
	server := httptest.NewServer(suite.router)
	defer server.Close()

	var okResponses, blockedResponses int
	const numberOfRequests = 1000

	for i := 0; i < numberOfRequests; i++ {
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		assert.NoError(suite.T(), err)
		req.Header.Set("API_KEY", "1245487875445487844545454fdf4d5f4")
		resp, err := http.DefaultClient.Do(req)
		assert.NoError(suite.T(), err)

		if resp.StatusCode == http.StatusOK {
			okResponses++
		} else if resp.StatusCode == http.StatusTooManyRequests {
			blockedResponses++
		}
	}

	assert.Equal(suite.T(), suite.maxTokenRequests, okResponses)
	assert.Equal(suite.T(), numberOfRequests-suite.maxTokenRequests, blockedResponses)
}

func TestWebServerSuite(t *testing.T) {
	suite.Run(t, new(WebServerTestSuite))
}
