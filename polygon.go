// Also added other functions and some tests related to geo based polygons.

package geo

import (
	"math"
)

// A Polygon is carved out of a 2D plane by a set of (possibly disjoint) contours.
// It can thus contain holes, and can be self-intersecting.
type Polygon struct {
	points []*Point
}

// Creates and returns a new pointer to a Polygon
// composed of the passed in points.  Points are
// considered to be in order such that the last point
// forms an edge with the first point.
func NewPolygon(points []*Point) *Polygon {
	return &Polygon{points: points}
}

// Returns the points of the current Polygon.
func (p *Polygon) Points() []*Point {
	return p.points
}

// Appends the passed in contour to the current Polygon.
func (p *Polygon) Add(point *Point) {
	p.points = append(p.points, point)
}

// Returns whether or not the polygon is closed.
// TODO:  This can obviously be improved, but for now,
//        this should be sufficient for detecting if points
//        are contained using the raycast algorithm.
func (p *Polygon) IsClosed() bool {
	if len(p.points) < 3 {
		return false
	}

	return true
}

// Returns whether or not the current Polygon contains the passed in Point.
func (p *Polygon) Contains(point *Point) bool {
	if !p.IsClosed() {
		return false
	}

	start := len(p.points) - 1
	end := 0

	contains := false

	// This assumes polygon is non-intersecting and contains the same point to
	// begin and end the polygon. If a point is repeated, it is assumed to
	// complete the polygon. The first polygon is assumed to be the bounding
	// polygon. All polygons that follow the bounding polygon are assumed to
	// define a hole within the polygon and must have an ending point that
	// equals the beginning point. If the first point does not match the last
	// point, the polygon contains at least one hole.
	//
	// if there are no holes, we can use simpler and faster logic
	if p.points[start] == p.points[end] {
		for i := 1; i < len(p.points); i++ {
			if p.intersectsWithRaycast(point, p.points[i-1], p.points[i]) {
				contains = !contains
			}
		}
	} else {
		pt := p.points[0]
		for i := 1; i < len(p.points); i++ {

			// if true, we are transitioning to a hole, skip
			if i != 1 && pt == p.points[i-1] {
				pt = p.points[i-1]
			} else {
				if p.intersectsWithRaycast(point, p.points[i-1], p.points[i]) {
					contains = !contains
				}
			}
		}
	}

	return contains
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func (p *Polygon) intersectsWithRaycast(point *Point, start *Point, end *Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.lng > end.lng {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.lng == start.lng || point.lng == end.lng {
		newLng := math.Nextafter(point.lng, math.Inf(1))
		point = NewPoint(point.lat, newLng)
	}

	// If we are outside of the polygon, indicate so.
	if point.lng < start.lng || point.lng > end.lng {
		return false
	}

	if start.lat > end.lat {
		if point.lat > start.lat {
			return false
		}
		if point.lat < end.lat {
			return true
		}

	} else {
		if point.lat > end.lat {
			return false
		}
		if point.lat < start.lat {
			return true
		}
	}

	raySlope := (point.lng - start.lng) / (point.lat - start.lat)
	diagSlope := (end.lng - start.lng) / (end.lat - start.lat)

	return raySlope >= diagSlope
}
