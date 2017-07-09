// +build !appengine
package app

//
// common
//

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/pata/", handlePata)
	http.HandleFunc("/route/", handleRouteSearch)
	http.HandleFunc("/stainfo/", handleStationInfo)
}

// https://mrekucci.blogspot.jp/2015/07/dont-abuse-mathmax-mathmin.html
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

//
// handle station info
//

func handleStationInfo(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get("http://fantasy-transit.appspot.com/net?format=json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// http://m-shige1979.hatenablog.com/entry/2016/01/29/080000
	// http://golang-jp.org/pkg/text/template/
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		panic(err3)
	}
	var tracks []Track
	if err := json.Unmarshal(body, &tracks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// gen stations
	stations := map[string]Station{}
	for ti, t := range tracks {
		for si, s := range t.Stations {
			_, ok := stations[s]
			if ok == false {
				stations[s] = Station{Key: s, Name: s + " @"}
			}
			newTIndexList := append(stations[s].TrackIndexList, ti)
			newSIndexList := append(stations[s].StationIndexInTrackList, si)
			stations[s] = Station{
				Key:                     stations[s].Key,
				Name:                    stations[s].Name + " " + t.Name,
				TrackIndexList:          newTIndexList,
				StationIndexInTrackList: newSIndexList,
			}
		}
	}

	for _, sta := range stations {
		adjStaList := findAdjStations(sta, tracks)
		stations[sta.Key] = Station{
			Key:                     sta.Key,
			Name:                    sta.Name,
			TrackIndexList:          sta.TrackIndexList,
			StationIndexInTrackList: sta.StationIndexInTrackList,
			AdjStations:             adjStaList,
		}

	}

	//s := fmt.Sprintf("%#v", stations)
	r.ParseForm()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	station := stations[r.Form.Get("key")]

	tpl, err1 := template.ParseFiles("templates/stationinfo.gtpl")
	if err1 != nil {
		panic(err1)
	}
	err2 := tpl.Execute(w, struct {
		Sta      Station
		Stations map[string]Station
	}{
		Sta:      station,
		Stations: stations,
	})
	if err2 != nil {
		panic(err2)
	}

}

//
// route search
//

type Track struct {
	Name     string   `json:"Name"`
	Stations []string `json:"Stations"`
}

type Station struct {
	Key                     string
	Name                    string
	TrackIndexList          []int
	StationIndexInTrackList []int
	AdjStations             []string
}

func findAdjStations(station Station, tracks []Track) []string {
	var adjStations []string
	for tii, ti := range station.TrackIndexList {
		si := station.StationIndexInTrackList[tii]
		if si != 0 {
			adjStations = append(adjStations, tracks[ti].Stations[si-1])
		}
		if si != len(tracks[ti].Stations)-1 {
			adjStations = append(adjStations, tracks[ti].Stations[si+1])
		}
	}
	return adjStations
}

func shiftRoute(slice [][]string) ([]string, [][]string) {
	ans := slice[0]
	slice = slice[1:]
	return ans, slice
}

// import fmt
func findRoute(from, to string, stations map[string]Station) []string {
	_, ok := stations[to]
	if !ok {
		return []string{}
	}
	var checkingRoutes [][]string
	checkingRoutes = append(checkingRoutes, []string{from})

	var route []string
	for i := 0; ; i++ {
		if len(checkingRoutes) == 0 {
			return []string{}
		}
		route, checkingRoutes = shiftRoute(checkingRoutes)
		lastSta := route[len(route)-1]
		if lastSta == to {
			return route
		}
		for _, adjSta := range stations[lastSta].AdjStations {
			found := false
			for _, s := range route {
				if s == adjSta {
					found = true
				}
			}
			if !found {
				checkingRoutes = append(checkingRoutes, append(route, adjSta))
			}
		}
	}

	return []string{}
}

func handleRouteSearch(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get("http://fantasy-transit.appspot.com/net?format=json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//
	// http://m-shige1979.hatenablog.com/entry/2016/01/29/080000
	// http://golang-jp.org/pkg/text/template/
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		panic(err3)
	}
	var tracks []Track
	if err := json.Unmarshal(body, &tracks); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// gen stations
	stations := map[string]Station{}
	for ti, t := range tracks {
		for si, s := range t.Stations {
			_, ok := stations[s]
			if ok == false {
				stations[s] = Station{Key: s, Name: s + " @"}
			}
			newTIndexList := append(stations[s].TrackIndexList, ti)
			newSIndexList := append(stations[s].StationIndexInTrackList, si)
			stations[s] = Station{
				Key:                     stations[s].Key,
				Name:                    stations[s].Name + " " + t.Name,
				TrackIndexList:          newTIndexList,
				StationIndexInTrackList: newSIndexList,
			}
		}
	}

	for _, sta := range stations {
		adjStaList := findAdjStations(sta, tracks)
		stations[sta.Key] = Station{
			Key:                     sta.Key,
			Name:                    sta.Name,
			TrackIndexList:          sta.TrackIndexList,
			StationIndexInTrackList: sta.StationIndexInTrackList,
			AdjStations:             adjStaList,
		}

	}

	//s := fmt.Sprintf("%#v", stations)
	r.ParseForm()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	stationFrom := stations[r.Form.Get("stationFrom")]
	stationTo := stations[r.Form.Get("stationTo")]

	route := findRoute(stationFrom.Key, stationTo.Key, stations)

	tpl, err1 := template.ParseFiles("templates/routesearch.gtpl")
	if err1 != nil {
		panic(err1)
	}
	err2 := tpl.Execute(w, struct {
		Stations    map[string]Station
		Result      []string
		StationFrom Station
		StationTo   Station
	}{
		Stations:    stations,
		Result:      route,
		StationFrom: stationFrom,
		StationTo:   stationTo,
	})
	if err2 != nil {
		panic(err2)
	}

}

//
// pata
//

func handlePata(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	a := r.Form.Get("a")
	b := r.Form.Get("b")
	ans := generatePatatokakushi(a, b)

	// http://m-shige1979.hatenablog.com/entry/2016/01/29/080000
	// http://golang-jp.org/pkg/text/template/
	tpl, err1 := template.ParseFiles("templates/pata.gtpl")
	if err1 != nil {
		panic(err1)
	}
	err2 := tpl.Execute(w, struct {
		Result string
	}{
		Result: ans,
	})
	if err2 != nil {
		panic(err2)
	}
}

func generatePatatokakushi(a, b string) string {
	var ans string
	ra := []rune(a)
	rb := []rune(b)
	maxLen := Max(len(ra), len(rb))
	for i := 0; i < maxLen; i++ {
		if i < len(ra) {
			ans += string(ra[i])
		}
		if i < len(rb) {
			ans += string(rb[i])
		}
	}
	return ans
}
