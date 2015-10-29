package routedb

import (
	"io/ioutil"
	"testing"
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
