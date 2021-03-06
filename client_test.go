package hobknob

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, "{\"action\":\"get\",\"node\":{\"key\":\"/v1/toggles/testApp\",\"dir\":true,\"nodes\":[{\"key\":\"/v1/toggles/testApp/mytoggle\",\"value\":\"true\",\"modifiedIndex\":78,\"createdIndex\":78}],\"modifiedIndex\":75,\"createdIndex\":75}}")
}

func initServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4001", nil)
}

func TestServer(t *testing.T) {
	go initServer()
}

func Setup(t *testing.T) (*Client, error) {
	c := NewClient([]string{"http://127.0.0.1:4001"}, "testApp", 1)
	err := c.Initialise()

	go func() {
		for {
			err := <-c.OnError
			t.Error(err)
		}
	}()

	go func() {
		for {
			<-c.OnUpdate
		}
	}()

	return c, err
}

func SetupBench(b *testing.B) (*Client, error) {
	c := NewClient([]string{"http://127.0.0.1:4001"}, "testApp", 1)
	err := c.Initialise()

	go func() {
		for {
			err := <-c.OnError
			b.Error(err)
		}
	}()

	go func() {
		for {
			<-c.OnUpdate
		}
	}()

	return c, err
}

func TestNew(t *testing.T) {
	c, _ := Setup(t)

	if c == nil {
		t.Fatalf("client was null")
	}

	if c.AppName != "testApp" {
		t.Fatalf("AppName not initialised: %v", c.AppName)
	}

	if c.CacheInterval != (time.Duration(1) * time.Second) {
		t.Fatalf("CacheInterval not initialised %v", c.CacheInterval)
	}
}

func TestInitialise(t *testing.T) {
	_, err := Setup(t)

	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	c, err := Setup(t)

	if err != nil {
		t.Error(err)
	}

	toggle, exists := c.Get("mytoggle")

	if toggle != true {
		t.Fatalf("expecting toggle 'mytoggle' to have value 'true' actual: '%v'", toggle)
	}

	if exists != true {
		t.Fatalf("expecting exists 'mytoggle' to have value 'true' actual: '%v'", toggle)
	}
}

func TestGetNonExistentToggle(t *testing.T) {
	c, err := Setup(t)

	if err != nil {
		t.Error(err)
	}

	toggle, exists := c.Get("unknowntoggle")

	if toggle != false {
		t.Fatalf("expecting toggle 'unknowntoggle' to have value 'false' actual: '%v'", toggle)
	}

	if exists != false {
		t.Fatalf("expecting exists 'unknowntoggle' to have value 'false' actual: '%v'", exists)
	}
}

func TestGetBadToggle(t *testing.T) {
	c, err := Setup(t)

	if err != nil {
		t.Error(err)
	}

	toggle, exists := c.Get("badtoggle")

	if toggle != false {
		t.Fatalf("expecting toggle 'badtoggle' to have value 'false' actual: '%v'", toggle)
	}

	if exists != false {
		t.Fatalf("expecting exists 'badtoggle' to have value 'false' actual: '%v'", exists)
	}
}

func TestGetOrDefault(t *testing.T) {
	c, err := Setup(t)

	if err != nil {
		t.Error(err)
	}

	toggle1 := c.GetOrDefault("mytoggle", true)

	if toggle1 != true {
		t.Fatalf("expecting toggle 'mytoggle' to have value 'true' actual: '%v'", toggle1)
	}

	toggle2 := c.GetOrDefault("unknowntoggle", true)

	if toggle2 != true {
		t.Fatalf("expecting toggle 'unknowntoggle' to have value 'true' actual: '%v'", toggle2)
	}
}

func TestSchedule(t *testing.T) {
	c, err := Setup(t)

	if err != nil {
		t.Error(err)
	}

	diffs := <-c.OnUpdate

	if diffs == nil {
		t.Fatalf("Got a nil update value: %v, was expecting: []Diffs{}", diffs)
	}
}

func BenchmarkGet(b *testing.B) {
	c, err := SetupBench(b)

	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		c.Get("mytoggle")
	}
}

func BenchmarkGetOrDefault(b *testing.B) {
	c, err := SetupBench(b)

	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		c.GetOrDefault("mytoggle", true)
	}
}
