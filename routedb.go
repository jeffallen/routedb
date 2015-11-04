// Package routedb includes routines for maintaining and querying
// a database of transport network routes.
//
// Because this package is used in Android apps, this API is
// (and must remain) compatible with the limitations of gobind.
package routedb

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/google/flatbuffers/go"
	"github.com/jeffallen/routedb/route"
	"github.com/kellydunn/golang-geo"
	"github.com/rndz/gpx"
)

// A Stop is a place where a bus stops (or could be hailed).
type Stop struct {
	Lat, Lon float64
}

// A Box is a region defined by two latitudes (N, S) and two
// longitudes (E, W).
type Box struct {
	N, E, S, W float64
}

// A Db represents an in-memory copy of the transport database.
type Db struct {
	zip    *zip.Reader
	routes []*gpx.Gpx
	bounds Box
}

// Load loads a routedb, returning a Db that can be queried, or an
// error.
func Load(in []byte) (db *Db, err error) {
	db = &Db{}
	db.zip, err = zip.NewReader(bytes.NewReader(in), int64(len(in)))
	for _, zf := range db.zip.File {
		file, err := zf.Open()
		fn := zf.FileHeader.Name
		if err != nil {
			return nil, fmt.Errorf("Failed to read file %v: %v", fn, err)
		}
		gpx, err := gpx.Parse(file)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %v: %v", fn, err)
		}
		if len(gpx.Trk) != 1 {
			return nil, fmt.Errorf("In file %v expected 1 track, found %v", fn, len(gpx.Trk))
		}
		if len(gpx.Trk[0].Trkseg) != 1 {
			return nil, fmt.Errorf("In file %v expected 1 track segment, found %v", fn, len(gpx.Trk[0].Trkseg))
		}
		db.routes = append(db.routes, gpx)
	}

	// If we have any points at all, use the first one as the anchor for
	// the bounds, then expand the bounds by processing the rest.
	if len(db.routes) >= 1 && len(db.routes[0].Trk[0].Trkseg[0].Trkpt) >= 1 {
		pt0 := db.routes[0].Trk[0].Trkseg[0].Trkpt[0]
		db.bounds.N, db.bounds.E = pt0.Lat, pt0.Lon
		db.bounds.S, db.bounds.W = pt0.Lat, pt0.Lon
		for _, route := range db.routes {
			for _, pt := range route.Trk[0].Trkseg[0].Trkpt {
				if pt.Lat > db.bounds.N {
					db.bounds.N = pt.Lat
				}
				if pt.Lon > db.bounds.E {
					db.bounds.E = pt.Lon
				}
				if pt.Lat < db.bounds.S {
					db.bounds.S = pt.Lat
				}
				if pt.Lon < db.bounds.W {
					db.bounds.W = pt.Lon
				}
			}
		}
	} else {
		// No waypoints in our db, so leave the bounds at the zero
		// value.
	}

	return db, err
}

// This can't be global because gobind cannot handle it.
// TODO: File an issue on this bug.
//var ErrNoStop = errors.New("No stop found matching criteria.")

func (db *Db) Nearest(lat, lon float64) (stop *Stop, err error) {
	p1 := geo.NewPoint(lat, lon)
	err = errors.New("No stop found matching criteria.")
	minD := 1e10

	for _, route := range db.routes {
		for _, trkpt := range route.Trk[0].Trkseg[0].Trkpt {
			p2 := geo.NewPoint(trkpt.Lat, trkpt.Lon)
			d := p1.GreatCircleDistance(p2)
			if d < minD {
				minD = d
				stop = &Stop{Lat: p2.Lat(), Lon: p2.Lng()}
				err = nil
			}
		}
	}
	return
}

// Bounds returns the box bounding all the waypoints in all the routes
// in the database. It returns a *Box to be compatible with gobind.
func (db *Db) Bounds() *Box {
	return &db.bounds
}

func split_md(in string) (country, city, name string) {
	x := strings.SplitN(in, "-", 3)
	if len(x) == 3 {
		return x[0], x[1], x[2]
	}
	return
}

// Routes returns the number of routes.
func (db *Db) Routes() int {
	return len(db.routes)
}

// Route returns the selected route as a FlatBuffer.
func (db *Db) Route(i int) ([]byte, error) {
	if i >= len(db.routes) {
		return nil, errors.New("out of range")
	}

	gpx := db.routes[i]
	country, city, name := split_md(gpx.Metadata.Name)

	b := flatbuffers.NewBuilder(0)

	l1 := b.CreateString(country)
	l2 := b.CreateString(city)
	l3 := b.CreateString(name)
	route.RouteStartPathVector(b, len(gpx.Trk[0].Trkseg[0].Trkpt))
	for j := len(gpx.Trk[0].Trkseg[0].Trkpt) - 1; j >= 0; j-- {
		trkpt := gpx.Trk[0].Trkseg[0].Trkpt[j]
		lat := int32(trkpt.Lat * 1e6)
		lon := int32(trkpt.Lon * 1e6)
		route.CreateGeoPoint(b, lat, lon)
	}
	l4 := b.EndVector(len(gpx.Trk[0].Trkseg[0].Trkpt))

	route.RouteStart(b)
	route.RouteAddCountry(b, l1)
	route.RouteAddCity(b, l2)
	route.RouteAddName(b, l3)
	route.RouteAddPath(b, l4)
	b.Finish(route.RouteEnd(b))

	return b.Bytes[b.Head():], nil
}
