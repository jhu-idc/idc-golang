// Provides methods for accessing resources of the Drupal JSON API
package jsonapi

import (
	"encoding/json"
	"fmt"
	"github.com/jhu-idc/idc-golang/drupal/env"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// Encapsulates the Entity type and bundle of a Drupal resource.
//
// DrupalType is parsed from the JSONAPI response, where type is represented, e.g. as:
//   "type": "taxonomy_term--person"
type DrupalType string

// The entity (e.g. taxonomy_term, node, etc) encapsulated by this type
func (t DrupalType) Entity() string {
	return strings.Split(string(t), "--")[0]
}

// The bundle (e.g. 'person', 'islandora_object', etc) encapsulated by this type
func (t DrupalType) Bundle() string {
	// TODO: some entities (like User) do not have a bundle type
	return strings.Split(string(t), "--")[1]
}

// Default HTTP client
var httpClient = &http.Client{}

// Encapsulates the relevant components of a URL which executes a JSON API request against Drupal; the typical
// entrypoint into the JSON API for making queries and retrieving results.
//
// Filter and Value are used to match an entity, e.g.: `title` and `The Adventures of Sherlock Holmes` respectively.
// For any filters more complex than a single field with a matching value, use RawFilter to specify the entire filter.
// If the RawFilter is present, Filter and Value are ignored, and the RawFilter is appended to the JSON API url as-is.
// An example RawFilter value might be: `filter[name-group][condition][operator]=ENDS_WITH&filter[name-group][condition][path]=name&filter[name-group][condition][value]=Thumbnail Image.jpg&filter[of-group][condition][path]=field_media_of.title&filter[of-group][condition][value]=Derivative Image 04`
type JsonApiUrl struct {
	T            assert.TestingT
	BaseUrl      string
	DrupalEntity string
	DrupalBundle string
	// Filter is the name of the field to match on, e.g. `title`, `name`, or `id`.
	// If RawFilter is supplied, this field is ignored.
	Filter string
	// Value is the value that the Filter field must match, e.g. `The Adventures of Sherlock Holmes`,
	// `Ansel Adams Images`, or `329c57a2-97f2-4350-8b54-439237c68311`.  If RawFilter is supplied, this field is
	// ignored.
	Value string
	// RawFilter is supplied by the caller and is used as-is.  In that case, Filter and Value are not used.
	RawFilter string
	// The username to use when authenticating to Drupal's JSONAPI endpoint.  If this value is empty, no `Authorization` header will be sent, otherwise Basic authentication is used.
	Username  string
	// The password to use when authenticating to Drupal's JSONAPI endpoint.
	Password  string
}

// Get the JSON API content from the URL and unmarshal the response into the supplied interface (which must be a
// pointer).  This method asserts that there is a single object in the `data` element of the JSON response.
func (jar *JsonApiUrl) GetSingle(v interface{}) {
	var res *http.Response
	var body[]byte

	// retrieve json of the migrated entity from the jsonapi and unmarshal the single response
	if len(strings.TrimSpace(jar.Username)) == 0 {
		res, body = GetResource(jar.T.(*testing.T), jar.String())
	} else {
		res, body = GetResourceWithBasicAuth(jar.T.(*testing.T), jar.String(), jar.Username, jar.Password)
	}
	defer func() { _ = res.Close }()
	UnmarshalSingleResponse(jar.T.(*testing.T), body, res, &JsonApiResponse{}).To(v)
}

// Get the JSON API content from the URL and unmarshal the response into the supplied interface (which must be a
// pointer).
func (jar *JsonApiUrl) Get(v interface{}) {
	var res *http.Response
	var body[]byte

	// retrieve json of the migrated entity from the jsonapi and unmarshal the single response
	if len(strings.TrimSpace(jar.Username)) == 0 {
		res, body = GetResource(jar.T.(*testing.T), jar.String())
	} else {
		res, body = GetResourceWithBasicAuth(jar.T.(*testing.T), jar.String(), jar.Username, jar.Password)
	}
	defer func() { _ = res.Close }()
	UnmarshalResponse(jar.T.(*testing.T), body, res, &JsonApiResponse{}, nil).To(v)
}

// Encapsulates a generic JSON API response
type JsonApiResponse struct {
	// The 'data' element(s) of the response
	Data []map[string]interface{}
}

// Handles the case where the 'data' key contains an array of objects, or a single object.
func (jar *JsonApiResponse) UnmarshalJSON(b []byte) error {
	fullRes := make(map[string]interface{})

	if err := json.Unmarshal(b, &fullRes); err != nil {
		return err
	}

	if e, ok := fullRes["data"]; !ok {
		return fmt.Errorf("missing 'data' key when unmarshaling JSONAPI response: %v", e)
	} else {
		switch e.(type) {
		case []interface{}:
			jar.Data = make([]map[string]interface{}, len(e.([]interface{})))
			for i, v := range e.([]interface{}) {
				jar.Data[i] = v.(map[string]interface{})
			}
		case map[string]interface{}:
			jar.Data = make([]map[string]interface{}, 1)
			jar.Data[0] = e.(map[string]interface{})
		default:
			return fmt.Errorf("unable to determine type of JSONAPI key 'data': %v", e)
		}
	}
	return nil
}

// Adapts the generic JsonApiResponse to a higher-fidelity type
func (jar *JsonApiResponse) To(v interface{}) {
	if b, e := json.Marshal(jar); e != nil {
		log.Fatalf("Unable to marshal %v as json: %s", jar, e)
	} else {
		json.Unmarshal(b, v)
	}
}

// Compose and return a string representation of the JSONAPI URL
func (moo *JsonApiUrl) String() string {
	var u *url.URL
	var err error

	assert.NotEmpty(moo.T, moo.BaseUrl, "error generating a JsonAPI URL from %v: %s", moo, "base url must not be empty")
	assert.NotEmpty(moo.T, moo.DrupalEntity, "error generating a JsonAPI URL from %v: %s", moo, "drupal entity must not be empty")
	assert.NotEmpty(moo.T, moo.DrupalBundle, "error generating a JsonAPI URL from %v: %s", moo, "drupal bundle must not be empty")

	baseUrl := env.BaseUrlOr(moo.BaseUrl)
	if strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl[:len(baseUrl) - 1]
	}
	u, err = url.Parse(fmt.Sprintf("%s", strings.Join([]string{baseUrl, "jsonapi", moo.DrupalEntity, moo.DrupalBundle}, "/")))
	assert.Nil(moo.T, err, "error generating a JsonAPI URL from %v: %s", moo, err)

	// If a raw filter is supplied, use it as-is, otherwise use the .Filter and .Value
	if moo.RawFilter != "" {
		u, err = url.Parse(fmt.Sprintf("%s?%s", u.String(), moo.RawFilter))
	} else if moo.Filter != "" {
		u, err = url.Parse(fmt.Sprintf("%s?filter[%s]=%s", u.String(), moo.Filter, moo.Value))
	}

	assert.Nil(moo.T, err, "error generating a JsonAPI URL from %v: %s", moo, err)
	return u.String()
}

// Unmarshal a JSONAPI response body and assert that exactly one data element is present
func UnmarshalSingleResponse(t *testing.T, body []byte, res *http.Response, value *JsonApiResponse) *JsonApiResponse {
	UnmarshalResponse(t, body, res, value, func(value *JsonApiResponse) {
		assert.Equal(t, 1, len(value.Data), "Exactly one JSONAPI data element is expected in the response, but found %d element(s)", len(value.Data))
	})
	return value
}

// Unmarshal a JSONAPI response body and perform supplied assertions on the response
func UnmarshalResponse(t *testing.T, body []byte, res *http.Response, value *JsonApiResponse, responseAssertions func(res *JsonApiResponse)) *JsonApiResponse {
	err := json.Unmarshal(body, value)
	assert.Nil(t, err, "Error unmarshaling JSONAPI response body: %s", err)
	if responseAssertions != nil {
		responseAssertions(value)
	}
	return value
}

// GetResource returns the HTTP response and body from the supplied url.  It asserts that the HTTP status code is 200,
// and that no errors are encountered reading the response body.  The requeest will be unauthenticated
func GetResource(t *testing.T, u string) (*http.Response, []byte) {
	return GetResourceWithBasicAuth(t, u, "", "")
}

// GetResourceWithAuthn returns the HTTP response and body from the supplied url.  It asserts that the HTTP status code
// is 200, and that no errors are encountered reading the response body.  The supplied username and password are used to
// send a Basic Authorization header.  If the supplied username is empty, then the request will be sent without an
// Authorization header.
func GetResourceWithBasicAuth(t *testing.T, url, username, password string) (*http.Response, []byte) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if len(strings.TrimSpace(username)) > 0 {
		req.SetBasicAuth(username, password)
		log.Printf("Retrieving (with Authorization: basic) %s", url)
	} else {
		log.Printf("Retrieving %s", url)
	}
	res, err := httpClient.Do(req)
	assert.Nil(t, err, "encountered error requesting %s: %s", url, err)
	assert.Equal(t, 200, res.StatusCode, "%d status encountered when requesting %s", res.StatusCode, url)
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err, "error encountered reading response body from %s: %s", url, err)
	return res, body
}