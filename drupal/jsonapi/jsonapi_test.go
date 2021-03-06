package jsonapi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"html"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"testing"
)

// host that the test http server listens on
var testServerHost string

// random open high port that the test http server listens on
var testServerPort int

// A simple JSON API response that satisfies the constraints of JsonApiUrl.Get(...): namely that that response be
// unmarshaled into a JsonApiResponse, and that JsonApiResponse.Data has one element.
var stubResponse = `
{
  "jsonapi": {
    "version": "1.0",
    "meta": {
      "links": {
        "self": {
          "href": "http://jsonapi.org/format/1.0/"
        }
      }
    }
  },
  "data": [
    {
      "type": "media--document",
      "id": "fd0b8969-ecc9-4a0d-81d3-537ba95bd5a8"
    }
  ]
}`

// calledHandlerFunc wraps a http.HandlerFunc that keeps track of whether the handler has been invoked
type calledHandlerFunc struct {
	// called indicates whether or not the handler has been invoked
	called bool
	// handler the http.HandlerFunc that may or may not have been invoked
	handler http.HandlerFunc
}

// wasCalled checks to see if the handler has been invoked, then resets the flag
func (chf *calledHandlerFunc) wasCalled() bool {
	result := chf.called
	chf.called = false
	return result
}

// TestMain allocates an unused TCP port for the HTTP server or panics, then runs tests
func TestMain(m *testing.M) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Panicf("Unable to resolve TCP address: %s", err.Error())
	} else {
		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Panicf("Unable to listen on TCP address %s: %s", addr.String(), err.Error())
		}
		testServerHost = l.Addr().(*net.TCPAddr).IP.String()
		testServerPort = l.Addr().(*net.TCPAddr).Port
		log.Printf("Allocated unused port %d for test HTTP server", testServerPort)
		if err := l.Close(); err != nil {
			log.Panicf("Error closing listener %s: %s", l.Addr().String(), err.Error())
		}
	}

	os.Exit(m.Run())
}

// Insures that a JsonApiUrl with a non-empty Username will result in basic authentication being used.
//
// There is not a full suite of integration tests against the JSON API.  Empirically we know that unauthenticated
// JSON API requests do not result in a 4xx or 5xx error, even if the data being requested is partially restricted.
//
// HTTP basic authentication sent on the HTTP request are used to *authorize* access to data provided by the JSON API.
// For example, an unauthenticated JSON API request for a resource that you are not authorized for will *not* result in
// a 401 or a 403; you'll get a 200 from the JSON API, but the response will not contain unauthorized data.  This may
// result in a well-formed, but empty, JSON API response devoid of meaningful data. However, if you execute the same
// JSON API request *with* valid HTTP basic auth credentials, Drupal will authorize the access to the underlying data,
// and the response will contain what you expect.
//
// All this test does is ensure that when a non-empty JsonApiUrl.Username is provided, HTTP basic authentication is
// attempted.
func Test_GetResourceWithBasicAuthn(t *testing.T) {
	const (
		// the expected user name and password expected on authenticated HTTP requests to the dummy HTTP server
		expectedUser = "admin"
		expectedPass = "moo"

		// values used to coerce the JsonApiUrl URL request path to apply to different handlers
		request  = "request"
		withAuth = "withauth"
		noAuth   = "noauth"
	)

	// Url request paths, the first will expect a basic authorization header, the second will expect no authorization
	withAuthHandlerPath := fmt.Sprintf("/jsonapi/%s/%s", request, withAuth)
	noAuthHandlerPath := fmt.Sprintf("/jsonapi/%s/%s", request, noAuth)

	// Maps Url request paths to handlers.  Each handler will perform assertions then return a stub response.  Generally
	// the caller does not care about the response, it just needs to be valid JSON and able to be marshaled into a
	// JsonApiResponse.
	handlers := map[string]*calledHandlerFunc{
		withAuthHandlerPath: {
			false,
			func(writer http.ResponseWriter, request *http.Request) {
				user, pass, ok := request.BasicAuth()
				require.Equal(t, expectedUser, user)
				require.Equal(t, expectedPass, pass)
				require.True(t, ok)
				writer.Write([]byte(stubResponse))
			},
		},
		noAuthHandlerPath: {
			false,
			func(writer http.ResponseWriter, request *http.Request) {
				user, pass, ok := request.BasicAuth()
				require.Equal(t, "", user)
				require.Equal(t, "", pass)
				require.False(t, ok)
				writer.Write([]byte(stubResponse))
			},
		},
	}

	// Composite handler function which looks up and invokes the correct handler by exactly matching the url path from
	// the handlers map.
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		path := html.EscapeString(request.URL.Path)
		log.Printf("handling '%s'", path)
		handlerFunc := handlers[path]
		handlerFunc.handler(writer, request)
		handlerFunc.called = true
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	// Start the test http server on a high numbered port
	go func() {
		addr := fmt.Sprintf("%s:%d", testServerHost, testServerPort)
		server := http.Server{Addr: addr}
		ln, err := net.Listen("tcp4", addr)
		if err != nil {
			log.Panicf("Error starting HTTP servier on %s: %s", addr, err)
		} else {
			log.Printf("Listening on %s:%d", testServerHost, testServerPort)
			wg.Done()
		}

		err = server.Serve(ln)

		if err != nil {
			log.Panicf("Error running HTTP server on port %d: %s", testServerPort, err.Error())
		}
	}()

	wg.Wait() // wait for the server to begin listening for connections

	// Generic response object which we don't really care about.
	result := &JsonApiResponse{}

	// The JsonApiUrl forms the request, and is sent to our test server started earlier.  Note the values for
	// DrupalEntity, DrupalBundle are used to coerce the URL that is requested so that they match our handlers.
	u := &JsonApiUrl{
		T:            t,
		BaseUrl:      fmt.Sprintf("http://%s:%d", testServerHost, testServerPort),
		DrupalEntity: request,
		DrupalBundle: withAuth,
		Filter:       "name",
		Value:        "moo",
		Username:     expectedUser,
		Password:     expectedPass,
	}

	// Get a JsonApiResponse with authentication (the JsonApiUrl.Username is not empty)
	u.Get(result)
	assert.True(t, handlers[withAuthHandlerPath].wasCalled())

	// Get a JsonApiResponse without authentication (the JsonApiUrl.Username is the zero-length string)
	u.Username = ""
	u.Password = ""
	u.DrupalBundle = noAuth
	u.Get(result)
	assert.True(t, handlers[noAuthHandlerPath].wasCalled())

	// Get a JsonApiResponse without authentication (the JsonApiUrl.Username is the empty string)
	u.Username = "  "
	u.Get(result)
	assert.True(t, handlers[noAuthHandlerPath].wasCalled())

	// Get a JsonApiResponse without authentication (the JsonApiUrl.Username is the empty string; setting the
	// JsonApiUrl.Password on its own does not invoke basic authorization)
	u.Password = "foo"
	u.Get(result)
	assert.True(t, handlers[noAuthHandlerPath].wasCalled())
}
