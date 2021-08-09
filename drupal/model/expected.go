package model

// Expected entities must have a Drupal entity type and bundle.
type ExpectedEntity interface {
	// The Drupal entity type of the expected entity
	EntityType() string
	// The Drupal bundle of the expected entity
	EntityBundle() string
}

// NamedOrTitled entities are Drupal objects that have either a name or title. Most, if not all, expected entities will
// have either a name or title.
type NamedOrTitled interface {
	ExpectedEntity
	// The value of the name or title of the ExpectedEntity
	NameOrTitle() string
	// The name of the field, which ought to be 'name' or 'title'
	Field() string
}

type Expected struct {
	Type   string
	Bundle string
}

type ExpectedWithName struct {
	Expected
	Name string
}

type ExpectedWithTitle struct {
	Expected
	Title string
}

func (e Expected) EntityType() string {
	return e.Type
}

func (e Expected) EntityBundle() string {
	return e.Bundle
}

func (e ExpectedWithName) NameOrTitle() string {
	return e.Name
}

func (e ExpectedWithName) Field() string {
	return "name"
}

func (e ExpectedWithTitle) NameOrTitle() string {
	return e.Title
}

func (e ExpectedWithTitle) Field() string {
	return "title"
}

// Represents the expected results of a migrated person
type ExpectedPerson struct {
	ExpectedWithName
	PrimaryName string   `json:"primary_name"`
	RestOfName  []string `json:"rest_of_name"`
	FullerForm  []string `json:"fuller_form"`
	Prefix      []string
	Suffix      []string
	Number      []string
	AltName     []string `json:"alt_name"`
	Date        []string
	Knows       []string
	Authority   []struct {
		Uri  string
		Name string
		Type string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated repository object
type ExpectedRepoObj struct {
	ExpectedWithTitle
	Abstract         []LanguageString
	AccessRights     []string         `json:"access_rights"`
	AltTitle         []LanguageString `json:"alt_title"`
	CollectionNumber []string         `json:"collection_number"`
	CopyrightAndUse  string           `json:"copyright_and_use"`
	CopyrightHolder  []string         `json:"copyright_holder"`
	Contributor      []struct {
		RelType string `json:"rel_type"`
		Name    string
	}
	Creator []struct {
		RelType string `json:"rel_type"`
		Name    string
	}
	CustodialHistory  []LanguageString `json:"custodial_history"`
	DateAvailable     string           `json:"date_available"`
	DateCopyrighted   []string         `json:"date_copyrighted"`
	DateCreated       []string         `json:"date_created"`
	DatePublished     []string         `json:"date_published"`
	DigitalIdentifier []string         `json:"digital_identifier"`
	DigitalPublisher  []string         `json:"digital_publisher"`
	DisplayHint       string           `json:"display_hints"`
	DspaceIdentifier  string           `json:"dspace_identifier"`
	DspaceItemId      string           `json:"dspace_itemid"`
	Extent            []string
	FeaturedItem      bool `json:"featured_item"`
	FindingAid        []struct {
		Uri   string
		Title string
	} `json:"finding_aid"`
	Genre              []string
	GeoportalLink      string   `json:"geoportal_link"`
	AccessTerms        []string `json:"access_terms"`
	Issn               string
	IsPartOf           string   `json:"is_part_of"`
	ItemBarcode        []string `json:"item_barcode"`
	JhirUri            string   `json:"jhir"`
	LibraryCatalogLink []string `json:"catalog_link"`
	Model              struct {
		Name        string
		ExternalUri string `json:"external_uri"`
	}
	OclcNumber       []string `json:"oclc_number"`
	Publisher        []string
	PublisherCountry []string `json:"publisher_country"`
	ResourceType     []string `json:"resource_type"`
	SpatialCoverage  []string `json:"spatial_coverage"`
	Subject          []string
	TableOfContents  []LanguageString `json:"toc"`
	MemberOf         string           `json:"member_of"`
	LinkedAgent      []struct {
		Rel  string
		Name string
	}
	Description []struct {
		Value    string
		LangCode string `json:"language"`
	}
}

// Represents the expected results of a migrated Access Rights taxonomy term
type ExpectedAccessRights struct {
	ExpectedWithName
	Authority []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Islandora Access Terms taxonomy term
type ExpectedIslandoraAccessTerms struct {
	ExpectedWithName
	Parent      []string `json:"parent"`
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Copyright and Use taxonomy term
type ExpectedCopyrightAndUse struct {
	ExpectedWithName
	Authority []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Family taxonomy term
type ExpectedFamily struct {
	ExpectedWithName
	Date       []string
	FamilyName string `json:"family_name"`
	Title      string
	Authority  []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
	KnowsAbout []string `json:"knowsAbout"`
}

// Represents the expected results of a migrated Genre taxonomy term
type ExpectedGenre struct {
	ExpectedWithName
	Authority []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Geolocation taxonomy term
type ExpectedGeolocation struct {
	ExpectedWithName
	GeoAltName []string `json:"geo_alt_name"`
	Broader    []struct {
		Uri   string
		Title string
	}
	Authority []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Resource Types taxonomy term
type ExpectedResourceType struct {
	ExpectedWithName
	Authority []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Subject taxonomy term
type ExpectedSubject struct {
	ExpectedWithName
	Authority []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Language taxonomy term
type ExpectedLanguage struct {
	ExpectedWithName
	LanguageCode string `json:"language_code"`
	Authority    []struct {
		Uri    string
		Title  string
		Source string
	}
	Description struct {
		Value     string
		Format    string
		Processed string
	}
}

// Represents the expected results of a migrated Collection entity
type ExpectedCollection struct {
	ExpectedWithTitle
	TitleLangCode string `json:"title_language"`
	AltTitle      []struct {
		Value    string
		LangCode string `json:"language"`
	} `json:"alternative_title"`
	Description []struct {
		Value    string
		LangCode string `json:"language"`
	}
	ContactEmail     string   `json:"contact_email"`
	ContactName      string   `json:"contact_name"`
	CollectionNumber []string `json:"collection_number"`
	MemberOf         string   `json:"member_of"`
	AccessTerms      []string `json:"access_terms"`
	FindingAid       []struct {
		Uri   string
		Title string
	} `json:"finding_aid"`
}

// Represents the expected results of a migrated Corporate Body taxonomy term
type ExpectedCorporateBody struct {
	ExpectedWithName
	Description struct {
		Value     string
		Format    string
		Processed string
	}
	PrimaryName     string   `json:"primary_name"`
	SubordinateName []string `json:"subordinate_name"`
	DateOfMeeting   []string `json:"date_of_meeting_or_treaty"`
	Location        []string `json:"location_of_meeting"`
	NumberOrSection []string `json:"num_of_section_or_meet"`
	AltName         []string `json:"corporate_body_alternate_name"`
	Authority       []struct {
		Uri    string
		Title  string
		Source string
	}
	Date         []string
	Relationship []struct {
		Name string
		Rel  string `json:"rel_type"`
	} `json:"relationships"`
}

type LanguageString struct {
	Value    string
	LangCode string `json:"language"`
}

type ExpectedMediaGeneric struct {
	ExpectedWithName
	OriginalName string `json:"original_name"`
	Size         int
	MimeType     string   `json:"mime_type"`
	AccessTerms  []string `json:"access_terms"`
	MediaUse     []string `json:"use"`
	MediaOf      string   `json:"media_of"`
	Uri          struct {
		Url   string
		Value string
	}
	RestrictedAccess bool `json:"restricted_access"`
}

type ExpectedMediaImage struct {
	ExpectedMediaGeneric
	AltText string `json:"alt_text"`
	Height  int
	Width   int
}

type ExpectedMediaExtractedText struct {
	ExpectedMediaGeneric
	ExtractedText struct {
		Value     string
		Format    string
		Processed string
	} `json:"extracted_text"`
}

type ExpectedMediaRemoteVideo struct {
	ExpectedWithName
	EmbedUrl string `json:"embed_url"`
	MediaOf  string `json:"media_of"`
}
