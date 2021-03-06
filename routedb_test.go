package routedb

import (
	"io/ioutil"
	"testing"

	"github.com/jeffallen/routedb/route"
)

var db *Db

func init() {
	bytes, err := ioutil.ReadFile("testdata/routedb.zip")
	if err != nil {
		panic("read db")
	}

	db, err = Load(bytes)
	if err != nil {
		panic("parse db")
	}
}

func TestRoutes(t *testing.T) {
	buf, err := db.Route(0)
	if err != nil {
		t.Fatal(err)
	}
	route := route.GetRootAsRoute(buf, 0)
	if string(route.Country()) != "kg" {
		t.Errorf("route is not expected: %v", route)
	}

	explen := 477
	if route.PathLength() != explen {
		t.Errorf("path len is %v", route.PathLength())
	}
}

func TestNearest(t *testing.T) {
	// a known point is: lat 40.50263 lon 72.821976
	// so we ask for a point near that and expect it to come back
	explat, explon := 40.50263, 72.821976

	n, err := db.Nearest(40.50265, 72.821978)
	if err != nil {
		t.Fatal("nearest err:", err)
	}

	if n.Lat != explat || n.Lon != explon {
		t.Errorf("expected %v/%v, got %v/%v)",
			explat, explon, n.Lat, n.Lon)
	}
}

func TestBounds(t *testing.T) {
	b := db.Bounds()
	// These expected values were checked by putting the .xml file
	// into Excel and sorting for mins/maxes.
	if b.N != 40.5432 || b.S != 40.501026 {
		t.Error("bounds n/s wrong")
	}
	if b.E != 72.822586 || b.W != 72.796295 {
		t.Error("bounds e/w wrong")
	}
}
