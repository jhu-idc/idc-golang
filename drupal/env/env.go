// Provides access to environment variables used by IDC
package env

import (
	"fmt"
	"os"
	"strconv"
)

const (
	drupalBaseUrl = "DRUPAL_BASE_URL"
	testBasedir   = "DRUPAL_TEST_BASEDIR"
	assetsBaseUrl = "BASE_ASSETS_URL"
)

// Answers the base url of Drupal from the environment variable 'DRUPAL_BASE_URL', or panics
func BaseUrl() string {
	return requireEnv(drupalBaseUrl)
}

// Answers the base url of Drupal from the environment variable 'DRUPAL_BASE_URL', or returns the default value if unset
func BaseUrlOr(defaultValue string) string {
	return GetEnvOr(drupalBaseUrl, defaultValue)
}

// Answers the name (not path) of the base directory for the test suite from the environment variable
// 'DRUPAL_TEST_BASEDIR', or panics
func TestBasedir() string {
	return requireEnv(testBasedir)
}

// Answers the base URL to the test assets docker container from the environment variable
// 'DRUPAL_TEST_BASEDIR', or returns the default value if unset
func TestBasedirOr(defaultValue string) string {
	return GetEnvOr(testBasedir, defaultValue)
}

// Answers the base URL to the test assets docker container from the environment variable
// 'BASE_ASSETS_URL', or panics
func AssetsBaseUrl() string {
	return requireEnv(assetsBaseUrl)
}

// Answers the base URL to the test assets docker container from the environment variable
// 'BASE_ASSETS_URL', or returns the default value if unset
func AssetsBaseUrlOr(defaultValue string) string {
	return GetEnvOr(assetsBaseUrl, defaultValue)
}

// Answers the value of the supplied environment variable, or the default value if unset
func GetEnvOr(envVar, defValue string) string {
	if val, ok := getEnv(envVar, false); ok {
		return val
	} else {
		return defValue
	}
}

// Answers the value of the supplied environment variable as an integer, or the default value if unset.  This function
// will panic if the value of the environment variable cannot be parsed as an integer.
func GetEnvOrInt(envVar string, defValue int) int {
	if val, ok := getEnv(envVar, false); ok {
		if intval, err := strconv.Atoi(val); err != nil {
			panic(fmt.Errorf("env: error formatting the value of environment variable '%s' as an integer: %w", envVar, err))
		} else {
			return intval
		}
	} else {
		return defValue
	}
}

// Answers the value of the supplied environment variable, or the default value if unset.  This function
// will panic if the value of the environment variable cannot be parsed as a bool.
func GetEnvOrBool(envVar string, defValue bool) bool {
	if val, ok := getEnv(envVar, false); ok {
		if boolval, err := strconv.ParseBool(val); err != nil {
			panic(fmt.Errorf("env: error formatting the value of environment variable '%s' as a bool: %w", envVar, err))
		} else {
			return boolval
		}
	} else {
		return defValue
	}
}

// Answers the value for the supplied environment variable, or panics
func requireEnv(envVar string) string {
	val, _ := getEnv(envVar, true)
	return val
}

// Answers the value for the supplied environment variable, or panics if `require` is true
func getEnv(envVar string, require bool) (val string, ok bool) {
	if val, ok = os.LookupEnv(envVar); !ok {
		if require {
			panic(fmt.Sprintf("env: missing required environment variable: %s", envVar))
		}
	}

	return
}
