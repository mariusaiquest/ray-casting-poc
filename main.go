package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
)

type xy struct {
	x, y float64
}

type closedPoly struct {
	name string
	vert []xy
}

func inside(pt xy, pg closedPoly) bool {
	if len(pg.vert) < 3 {
		return false
	}
	in := rayIntersectsSegment(pt, pg.vert[len(pg.vert)-1], pg.vert[0])
	for i := 1; i < len(pg.vert); i++ {
		if rayIntersectsSegment(pt, pg.vert[i-1], pg.vert[i]) {
			in = !in
		}
	}
	return in
}

func rayIntersectsSegment(p, a, b xy) bool {
	return (a.y > p.y) != (b.y > p.y) &&
		p.x < (b.x-a.x)*(p.y-a.y)/(b.y-a.y)+a.x
}

var tpg = []closedPoly{

	{"poly", []xy{
		{42.160359, 18.985600},
		{42.138387, 19.006795},
		{42.133652, 19.011903},
		{42.138360, 19.016029},
		{42.140007, 19.017325},
		{42.140213, 19.019176},
		{42.140041, 19.019963},
		{42.140504, 19.020750},
		{42.141191, 19.021861},
		{42.140539, 19.023666},
		{42.139852, 19.024777},
		{42.138868, 19.024370},
		{42.137990, 19.024859},
		{42.137457, 19.025400},
		{42.136765, 19.025773},
		{42.135825, 19.024954},
		{42.127649, 19.014226},
		{42.126663, 18.994998},
		{42.138566, 18.983621},
		{42.160359, 18.985600},
	}},
	{"poly2", []xy{
		{42.134953, 18.852932},
		{42.145176, 18.973353},
		{42.136866, 18.972514},
		{42.099694, 18.854581},
		{42.134953, 18.852932},
	}},
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("index.html")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		defer f.Close()
		io.Copy(w, f)
	})

	http.HandleFunc("/onwater", func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query()
		lat, _ := strconv.ParseFloat(param.Get("lat"), 64)
		lng, _ := strconv.ParseFloat(param.Get("lng"), 64)

		reqPoint := xy{lat, lng}

		res := struct {
			Inp1 bool `json:"inP1"`
			Inp2 bool `json:"inP2"`
		}{inside(reqPoint, tpg[1]), inside(reqPoint, tpg[2])}
		json.NewEncoder(w).Encode(&res)
	})

	http.ListenAndServe(":8000", nil)
}
