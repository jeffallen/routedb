// Package routedb includes routines for maintaining and querying
// a database of transport network routes.
//
// Because this package is used in Android apps, this API is
// (and must remain) compatible with the limitations of gobind.
package routedb

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

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

type Route struct {
	Country, City, Name string
}

// read parses input of the form kg-osh-101 into r.
func (r *Route) read(in string) {
	x := strings.SplitN(in, "-", 3)
	if len(x) == 3 {
		r.Country = x[0]
		r.City = x[1]
		r.Name = x[2]
	}
}

// Routes returns the number of routes.
func (db *Db) Routes() int {
	return len(db.routes)
}

// Route returns the selected route description.
func (db *Db) Route(i int) (*Route, error) {
	if i >= len(db.routes) {
		return nil, errors.New("out of range")
	}

	route := &Route{}
	route.read(db.routes[i].Metadata.Name)
	return route, nil
}

// Points returns a []byte with the path of the specified route
// in it encoded as pairs of (int64(lat*1e6)),(int64(lon*1e6))
// in little endian format.
func (db *Db) Points(i int) ([]byte, error) {
	if i >= len(db.routes) {
		return nil, errors.New("out of range")
	}

	gpx := db.routes[i]
	b := &bytes.Buffer{}
	for _, trkpt := range gpx.Trk[0].Trkseg[0].Trkpt {
		lat := int64(trkpt.Lat * 1e6)
		lon := int64(trkpt.Lon * 1e6)
		binary.Write(b, binary.LittleEndian, lat)
		binary.Write(b, binary.LittleEndian, lon)
	}
	return b.Bytes(), nil
}
