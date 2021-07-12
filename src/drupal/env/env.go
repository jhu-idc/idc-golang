package env

import "os"

// BaseUrl Answers the base url of Drupal from the environment
func BaseUrl() string {
	return os.Getenv("DRUPAL_BASE_URL")
}

func TestBasedir() string {
	return os.Getenv("DRUPAL_TEST_BASEDIR")
}
