package properties

import (
	"bufio"
	"math/rand"
	"os"
	"properties/plugins"
	"testing"
)

var (
	location      = "example.properties"
	configuration = NewProperties()
)

//goland:noinspection GoUnhandledErrorResult
func TestNewProperties(t *testing.T) {
	// Add resolvers
	configuration.PutResolver(plugins.EnvironmentDetector)
	configuration.PutResolver(plugins.EnvironmentDetector)
	// Load
	file, err := os.Open(location)
	if err != nil {
		t.Log(err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		err = configuration.Load(reader)
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Log(configuration)
	t.Log(configuration.GetPropertyOrDefault("example.env", ""))
	t.Log(configuration.GetPropertyOrDefault("example.env.miss", ""))
}

func TestFileProperties_Store(t *testing.T) {
	randSys := rand.Int()
	old := configuration.Put("example.random", randSys)

	// Info
	t.Logf("Old -> %v \n New -> %v", old, randSys)
	// Store
	out, err := os.Create(location)
	if err != nil {
		t.Fatal(err)
	}

	defer out.Close()
	w := bufio.NewWriter(out)
	err = configuration.Store(w)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(configuration)
}
