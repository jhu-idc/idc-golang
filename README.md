# About

Provides rudimentary access to the Drupal JSON API, and supports test assertions between JSON API resources and their expected values.

While there are a number of JSON API packages for Go, this little library was meant to be quick and simple solution for the singular use case of accessing Drupal content, and assessing the correctness of the content using structured data provided by a well-behaved API.

The functions and structures originated for the initial IDC migration test suite, but were found to be useful for derivative testing as well.  To keep these packages in a single place, they were cobbled together into this module.

# Usage

The use case for this module is to query for Drupal resources using the JSON API, and check - or assert - that the JSON objects carry the expected values.

Typically, the caller will begin by ingesting objects into Drupal using some out-of-band means, e.g. using Drupal Migrations.  Each resource, if successfully ingested, will have a JSON API representation.

The caller will use this module to search for and retrieve the JSON API representations, and compare them to expected values.  The expected values have a simple representation, distinct and separate from the JSON API.  Therefore, this module supports two different data models: one for JSON API resources, and a simpler, less-cumbersome 'Expected' data model.

An example Expected Subject Taxonomy may read as follows, noting that the Expected model has zero to do with JSON API.  The field names and values should be familiar, as they correspond one-to-one with the business properties of the Drupal objects we're concerned with.
```json
{
  "type": "taxonomy_term",
  "bundle": "subject",
  "name": "Analog Photography",
  "authority": [
    {
      "uri": "http://www.google.com?q=Analog%20Photography",
      "title": "Google",
      "source": "other"
    },
    {
      "uri": "http://www.ford.com",
      "title": "Nothing to do with Analog Photography But Works For a Test",
      "source": "iso19115"
    }
  ],
  "description": {
    "value": "<p>Analog photography description.</p>",
    "format": "basic_html",
    "processed": "<p>Analog photography description.</p>"
  }
}
```

A simple test that retrieves a Subject Taxonomy entry and compares it to expected values may read:

```go
func Test_VerifyTaxonomySubject(t *testing.T) {
	// Instantiate the expected object and unmarshal its representation from the filesystem
	// (note: callers ought to consider using an embedded filesystem instead)
	expectedJson := model.ExpectedSubject{}
	unmarshalExpectedJson(t, "taxonomy-subject.json", &expectedJson)

	// Sanity check the expected json
	assert.Equal(t, "taxonomy_term", expectedJson.Type)
	assert.Equal(t, "subject", expectedJson.Bundle)

	// Craft a JSON API url by setting values for various parameters
	u := &jsonapi.JsonApiUrl{
		T:            t,
		BaseUrl:      DrupalBaseurl,
		DrupalEntity: expectedJson.Type,
		DrupalBundle: expectedJson.Bundle,
		Filter:       "name",
		Value:        expectedJson.Name,
	}

	// Retrieve the JSON of the migrated entity from the JSON API and unmarshal the single response
	// Remember that prior to executing this test, some out-of-band process has populated Drupal
	// with the particular Subject taxonomy we're retrieving here.
	res := &model.JsonApiSubject{}
	u.GetSingle(res)

	// The "meat" of the response is contained in the JSON API 'data' element.  We can safely assume
	// that a single value exists because the `GetSingle` method was used to retrieve the Subject.
	actual := res.JsonApiData[0]
	
	// Make our assertions between the JSON API response and the expected values
	assert.Equal(t, expectedJson.Type, actual.Type.Entity())
	assert.Equal(t, expectedJson.Bundle, actual.Type.Bundle())
	assert.Equal(t, expectedJson.Name, actual.JsonApiAttributes.Name)
	assert.Equal(t, expectedJson.Description.Format, actual.JsonApiAttributes.Description.Format)
	assert.Equal(t, expectedJson.Description.Value, actual.JsonApiAttributes.Description.Value)
	assert.Equal(t, expectedJson.Description.Processed, actual.JsonApiAttributes.Description.Processed)
	assert.Equal(t, len(expectedJson.Authority), len(actual.JsonApiAttributes.Authority))
	assert.Equal(t, 2, len(actual.JsonApiAttributes.Authority))
	for i, v := range actual.JsonApiAttributes.Authority {
		assert.Equal(t, expectedJson.Authority[i].Source, v.Source)
		assert.Equal(t, expectedJson.Authority[i].Uri, v.Uri)
	}
}
```

Resolving relationships can be performed using the `Resolve(...)` method of `JsonApiData`.  An example Expected Collection may read as follows, noting that the parent collection is not referenced, but simply has its name present:
```json
{
  "type": "node",
  "bundle": "collection_object",
  "title": "Test Collection One",
  "member_of": "Parent Collection"
}

```

To obtain the parent collection referenced by a member_of relationship, you can:
```go
    childCol := &model.JsonApiCollection{}
    
    // (elided) retrieve the parent using JsonApiUrl
    
    // create a convenient reference to the relationships in the response
    relData := childCol.JsonApiData[0].JsonApiRelationships
    
    // (elided) make assertions about the presence of the member collection relationship
    
    // create the parent member collection structure and resolve the relationship.  
    // The JSON API collection will be retrieved by its Drupal identifier, and unmarshalled
    // into the supplied struct.
    parentCol := model.JsonApiCollection{}
    relData.MemberOf.Data.Resolve(t, &parentCol)
```

## Raw Filters

Since version `0.0.2`

`JsonApiUrl.RawFilter` was added in version `0.0.2` to support arbitrary JSON API filter expressions.  Normally executing `JsonApiUrl.Get(...)` will result in a query that expects exactly one result.

Prior to version `0.0.2` the JSON API filter was derived from `JsonApiUrl.Filter` and `JsonApiUrl.Value`, e.g.:
```go
u := &jsonapi.JsonApiUrl{
	...
	DrupalEntity: "media",
	DrupalBundle: "document",
	Filter:       "id",
	Value:        "329c57a2-97f2-4350-8b54-439237c68311",
	...
}
```

Sometimes a simple key/value pair is not sufficient for matching a single result; a more complex filter is required.  In that case, set a value for `JsonApiUrl.RawFilter`, and leave `JsonApiUrl.Filter` and `JsonApiUrl.Value` empty. For example, the derivative tests use a complex filter to match exactly one resource that was ingested using a combination of file name and parent media:

```go
u := &jsonapi.JsonApiUrl{
	...
	DrupalEntity: "media",
	DrupalBundle: "document",
	RawFilter:    "filter[name-group][condition][operator]=ENDS_WITH&filter[name-group][condition][path]=name&filter[name-group][condition][value]=Thumbnail Image.jpg&filter[of-group][condition][path]=field_media_of.title&filter[of-group][condition][value]=Derivative Image 04"
	...
}
```

## Authenticated Requests

Since version `0.0.5`

Set the `JsonApiUrl.Username` and `JsonApiUrl.Password` to execute an authenticated request.  In order to trigger HTTP Basic Auth, `JsonApiUrl.Username` must be a non-empty string.

Authenticated requests may be useful when access to the resource is denied to the anonymous user, e.g. by a restricted access flag on the media.