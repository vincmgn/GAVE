package groupie

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AllArtists struct {
	Artists     []Artists   `json:"artists"`
	Pays        []Pays      `json:"pays"`
	User        User        `json:"user"`
	SpotifyUser SpotifyUser `json:"spotifyUser"`
}

type MapData struct {
	Artists     []Artists `json:"artists"`
	Pays        []Pays    `json:"pays"`
	Coordonnees string    `json:"coordonnees"`
	CooName     string    `json:"cooName"`
	CooId       string    `json:"cooId"`
	CooArtist   string    `json:"cooArtist"`
	CooDates    string    `json:"cooDates"`
	User        User      `json:"user"`
}

type Artists struct {
	Id           int           `json:"id"`
	Image        string        `json:"image"`
	Name         string        `json:"name"`
	Members      []string      `json:"members"`
	CreationDate int           `json:"creationDate"`
	FirstAlbum   string        `json:"firstAlbum"`
	Relations    Relation      `json:"relations"`
	SpotifyData  SpotifyArtist `json:"spotifyData"`
}

type Pays struct {
	Pays string
}

type Locations struct {
	Index []struct {
		Id        int      `json:"id"`
		Locations []string `json:"locations"`
	} `json:"index"`
}

type Location struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func GetArtistsData() []Artists {
	var myDataList []Artists
	req, err := http.NewRequest("GET", "https://groupietrackers.herokuapp.com/api/artists", bytes.NewBufferString(""))
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	json.Unmarshal(reponsebytes, &myDataList)

	for i := 0; i < len(myDataList); i++ {
		myDataList[i].Relations = GetOneRelation(strconv.Itoa(i + 1))
	}

	return myDataList
}

func GetOneArtistById(id string) Artists {

	var myartist Artists
	req, err := http.NewRequest("GET", "https://groupietrackers.herokuapp.com/api/artists/"+id, bytes.NewBufferString(""))
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	json.Unmarshal(reponsebytes, &myartist)

	myartist.Relations = GetOneRelation(id)
	return myartist
}

func GetOneLocation(id string) Location {
	var myLocation Location
	req, err := http.NewRequest("GET", "https://groupietrackers.herokuapp.com/api/locations/"+id, bytes.NewBufferString(""))
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	json.Unmarshal(reponsebytes, &myLocation)

	return myLocation
}

func GetOneRelation(id string) Relation {
	var myRelation Relation
	req, err := http.NewRequest("GET", "https://groupietrackers.herokuapp.com/api/relation/"+id, bytes.NewBufferString(""))
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	json.Unmarshal(reponsebytes, &myRelation)

	return myRelation
}

func GetArtistsByName(letters string) []Artists {
	letters = strings.ToLower(letters)
	var tab []Artists
	data := GetArtistsData()
	for i := 0; i < len(data); i++ {
		var bool bool
		if strings.Contains(strings.ToLower(data[i].Name), letters) {
			bool = true
		}
		if bool {
			tab = append(tab, data[i])
		}
	}
	return tab
}

func GetArtistsByCreationDate(dates string) []Artists {
	var tab []Artists
	data := GetArtistsData()
	for i := 0; i < len(data); i++ {
		if strconv.Itoa(data[i].CreationDate) == dates {
			tab = append(tab, data[i])
		}
	}
	return tab
}

func GetArtistsByFirstAlbumDate(dates string) []Artists {
	var tab []Artists
	data := GetArtistsData()
	for i := 0; i < len(data); i++ {
		if data[i].FirstAlbum[6:] == dates {
			tab = append(tab, data[i])
		}
	}
	return tab
}

func GetArtistsByConcertDate(dates string) []Artists {
	var tab []Artists
	data := GetArtistsData()
	for i := 0; i < len(data); i++ {
		boolean := false
		adress := GetPaysById([]string{strconv.Itoa(data[i].Id)})
		for k := 0; k < len(adress); k++ {
			for _, l := range data[i].Relations.DatesLocations[adress[k].Pays] {
				if l[6:] == dates {
					tab = append(tab, data[i])
					boolean = true
					break
				}
			}
			if boolean {
				break
			}
		}
	}
	return tab
}

func GetArtistsByMembers(nb string) []Artists {
	var tab []Artists
	data := GetArtistsData()
	result, _ := strconv.Atoi(nb)
	for i := 0; i < len(data); i++ {
		if len(data[i].Members) == result {
			tab = append(tab, data[i])
		}
	}
	return tab
}

func GetArtistsByLocation(location string) []Artists {
	var tab []Artists
	data := GetArtistsData()
	bool := false
	for i := 0; i < len(data); i++ {
		locationdata := GetOneLocation(strconv.Itoa(data[i].Id))
		for _, j := range locationdata.Locations {
			index := 0
			for ind, k := range j {
				if k == '-' {
					index = ind
				}
			}
			if j[index+1:] == location {
				bool = true
			}
			if bool {
				tab = append(tab, data[i])
				bool = false
				break
			}
		}
	}
	return tab
}

func GetPays() []Pays {
	var listloca []Pays
	var nb int
	var myDataList Locations
	req, err := http.NewRequest("GET", "https://groupietrackers.herokuapp.com/api/locations", bytes.NewBufferString(""))
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	json.Unmarshal(reponsebytes, &myDataList)

	for i := 0; i < len(myDataList.Index); i++ {
		for j := 0; j < len(myDataList.Index[i].Locations); j++ {
			nb = 0
			index := 0
			for g := 0; g < len(myDataList.Index[i].Locations[j]); g++ {
				if myDataList.Index[i].Locations[j][g] == '-' {
					index = g + 1
				}
			}
			for k := 0; k < len(listloca); k++ {
				if myDataList.Index[i].Locations[j][index:] == listloca[k].Pays {
					nb++
					break
				}
			}
			if nb == 0 {
				listloca = append(listloca, Pays{myDataList.Index[i].Locations[j][index:]})
			}
		}
	}
	return listloca
}

func GetLocationWithAdress(adress string) string {

	var newadress string
	for i := 0; i < len(adress); i++ {
		if adress[i] == ' ' {
			newadress += "%20"
		} else {
			newadress += string(adress[i])
		}
	}

	req, err := http.NewRequest("GET", "https://maps.googleapis.com/maps/api/geocode/json?address="+newadress+"&key=AIzaSyBeq2hPzggF1UpPuHe8A3qao36F9OuAlvc", bytes.NewBufferString(""))
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	var geocodeResp struct {
		Results []struct {
			Geometry struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"results"`
	}

	json.Unmarshal(reponsebytes, &geocodeResp)

	var coo string

	coo = "[" + strconv.FormatFloat(geocodeResp.Results[0].Geometry.Location.Lat, 'f', 6, 64) + "," + strconv.FormatFloat(geocodeResp.Results[0].Geometry.Location.Lng, 'f', 6, 64) + "],"

	return coo
}

func GetPaysById(ids []string) []Pays {

	var listloca []Pays
	var nb int

	var myDataList Location

	for _, l := range ids {
		req, err := http.NewRequest("GET", "https://groupietrackers.herokuapp.com/api/locations/"+l, bytes.NewBufferString(""))
		if err != nil {
			log.Fatalf("error creating HTTP request: %v", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("error creating HTTP request: %v", err)
		}

		reponsebytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("error creating HTTP request: %v", err)
		}

		json.Unmarshal(reponsebytes, &myDataList)
		for j := 0; j < len(myDataList.Locations); j++ {
			nb = 0
			for k := 0; k < len(listloca); k++ {
				if myDataList.Locations[j] == listloca[k].Pays {
					nb++
					break
				}
			}
			if nb == 0 {
				listloca = append(listloca, Pays{myDataList.Locations[j]})
			}
		}
	}

	return listloca
}

func LocationByParams(data MapData) (string, string) {
	var list []string
	var CooString string
	var CooName string
	for _, j := range data.Artists {
		for k := range j.Relations.DatesLocations {
			list = append(list, k)
		}
	}
	for _, k := range list {
		CooString += GetLocationWithAdress(k)
		CooName += k + ","
	}
	return CooString, CooName
}

func GetArtistsLiked(data User) []Artists {
	var id string
	var result []Artists
	for i := 0; i < len(data.LikedArtists); i++ {
		if data.LikedArtists[i] != ',' {
			id += string(data.LikedArtists[i])
		} else {
			result = append(result, GetOneArtistById(id))
			id = ""
		}
	}
	return result
}

func GetArtistsByConcertLocation(location string) []Artists {
	var newTab []Artists
	tab := GetArtistsByLocation(location)
	for _, i := range tab {
		var tempLoc Relation
		tempLoc.Id = i.Id
		concertRelation := GetOneRelation(strconv.Itoa(i.Id))
		tempLoc.DatesLocations = make(map[string][]string)
		for key, value := range concertRelation.DatesLocations {
			index := 0
			for g := 0; g < len(key); g++ {
				if key[g] == '-' {
					index = g + 1
				}
			}
			if key[index:] == location {
				tempLoc.DatesLocations[key] = value
			}
		}
		i.Relations = tempLoc
		newTab = append(newTab, i)
	}
	return newTab
}
