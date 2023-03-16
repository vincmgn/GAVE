package groupie

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AccesTokenResponse struct {
	AccesToken string `json:"access_token"`
	TokenType  string `json:"token_type"`
	ExpiresIn  int    `json:"expires_in"`
	Scope      string `json:"scope"`
}

type SpotifyUser struct {
	DisplayName  string `json:"display_name"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Id        string         `json:"id"`
	URI       string         `json:"uri"`
	Playlists []UserPlaylist `json:"playlists"`
}

type SpotifyArtist struct {
	External_urls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Total int `json:"total"`
	} `json:"followers"`
	Genres     []string       `json:"genres"`
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Popularity int            `json:"popularity"`
	URI        string         `json:"uri"`
	Albums     []SpotifyAlbum `json:"albums"`
}

type SpotifyAlbum struct {
	AlbumGroup   string `json:"album_group"`
	AlbumType    string `json:"album_type"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
	Images      []struct {
		Url string `json:"url"`
	} `json:"images"`
	Tracks []Track `json:"tracks"`
}

type UserPlaylist struct {
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	TracksNb struct {
		Total int `json:"total"`
	} `json:"tracksNb"`
	Images []struct {
		Url string `json:"url"`
	} `json:"images"`
	Uri    string  `json:"uri"`
	Tracks []Track `json:"tracks"`
}

type Track struct {
	DurationMs     int    `json:"duration_ms"`
	DurationMinute string `json:"duration_minute"`
	ExternalUrls   struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Images []struct {
		Url string `json:"url"`
	} `json:"images"`
	Album       ImagePlaylist `json:"album"`
	TrackNumber int           `json:"track_number"`
	URI         string        `json:"uri"`
}

type ImagePlaylist struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
}

type ArtistsResponse struct {
	Items []SpotifyArtist `json:"items"`
}

type AlbumsResponse struct {
	Items []SpotifyAlbum `json:"items"`
}

type PlaylistsResponse struct {
	Items []UserPlaylist `json:"items"`
}

type SpotifySearchResult struct {
	Artists ArtistsResponse `json:"artists"`
}

type TracksResponse struct {
	Items []Track `json:"items"`
}

type PlaylistTracksResponse struct {
	Items []PlaylistsTrack `json:"items"`
}

type PlaylistsTrack struct {
	Track Track `json:"track"`
}

func GetAccessToken(client *http.Client, clientId string, clientSecret string) (string, error) {

	//Création requête pour l'API Spotify
	data := strings.NewReader("grant_type=client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", data)
	if err != nil {
		return "", err
	}

	//Ajout informations d'identification de l'application à l'en-tête de la demande.
	req.Header.Set("Authorization", "Basic "+EncodeCredentials(clientId, clientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//Ajout de paramètres de demande à la requ^te POST
	q := req.URL.Query()
	q.Add("grant_type", "client_credentials")
	//q.Add("scope" )
	req.URL.RawQuery = q.Encode()

	//Envoie de la requête au serveur
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//Lire la réponse de l'API Spotify
	var tokenResponse AccesTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	//Retourne le jeton d'accès (token)
	return tokenResponse.AccesToken, nil
}

func EncodeCredentials(clientId string, clientSecret string) string {
	return base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))
}

func GetSpotifyArtist(artist string) SpotifyArtist {
	var space string
	for i := 0; i < len(artist); i++ {
		if artist[i] == ' ' {
			space += "%20"
		} else {
			space += string(artist[i])
		}
	}

	//Configuration du client HTTP
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	//Obtention d'un jeton d'accès auprès de l'API Spotify
	token, err := GetAccessToken(client, "971b13b2aaf14061a2dd134d0cc18134", "e958e7edc5af4e17b6f55235f4ee62cc")
	if err != nil {
		panic(err)
	}

	//Configuration de la requête à l'API Spotify
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/search?q="+space+"&type=artist&limit=1&offset=0", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	//Envoie de la requête à l'API Spotify
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	// //Analyse de la réponse json de l'API Spotify

	var myDataList SpotifySearchResult

	json.Unmarshal(reponsebytes, &myDataList)

	Artist := myDataList.Artists.Items[0]
	Artist.Albums = GetSpotifyArtistsAlbum(Artist.Id)

	//Affichage données renvoyées par API Spotify
	return Artist
}

func GetSpotifyArtistsAlbum(id string) []SpotifyAlbum {

	// Configuration du client HTTP
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Obtention d'un jeton d'accès auprès de l'API Spotify
	token, err := GetAccessToken(client, "971b13b2aaf14061a2dd134d0cc18134", "e958e7edc5af4e17b6f55235f4ee62cc")
	if err != nil {
		panic(err)
	}

	//Configuration de la requête à l'API Spotify
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/artists/"+id+"/albums?limit=50", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	//Envoie de la requête à l'API Spotify
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	// //Analyse de la réponse json de l'API Spotify

	var myDataList AlbumsResponse
	var Albums []SpotifyAlbum

	json.Unmarshal(reponsebytes, &myDataList)

	for _, j := range myDataList.Items {
		temp := j.Images[2]
		j.Images = []struct {
			Url string "json:\"url\""
		}{temp}
		j.Tracks = GetAlbumTracks(j)
		Albums = append(Albums, j)
	}

	//Affichage données renvoyées par API Spotify
	return Albums
}

func GetSpotifyUser(id string) (SpotifyUser, error) {

	// Configuration du client HTTP
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	var erreur error

	// Obtention d'un jeton d'accès auprès de l'API Spotify
	token, err := GetAccessToken(client, "971b13b2aaf14061a2dd134d0cc18134", "e958e7edc5af4e17b6f55235f4ee62cc")
	if err != nil {
		erreur = err
	}

	//Configuration de la requête à l'API Spotify
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/users/"+id, nil)
	if err != nil {
		erreur = err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	//Envoie de la requête à l'API Spotify
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		erreur = err
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		erreur = err
	}

	// //Analyse de la réponse json de l'API Spotify

	var myUser SpotifyUser

	json.Unmarshal(reponsebytes, &myUser)

	myUser.Playlists = GetUsersPlaylist(myUser)

	for _, i := range myUser.Playlists {
		if len(i.Images) > 2 {
			i.Images = []struct {
				Url string "json:\"url\""
			}{i.Tracks[0].Images[0]}
		}
	}

	//Affichage données renvoyées par API Spotify
	return myUser, erreur
}

func GetUsersPlaylist(user SpotifyUser) []UserPlaylist {

	// Configuration du client HTTP
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Obtention d'un jeton d'accès auprès de l'API Spotify
	token, err := GetAccessToken(client, "971b13b2aaf14061a2dd134d0cc18134", "e958e7edc5af4e17b6f55235f4ee62cc")
	if err != nil {
		panic(err)
	}

	//Configuration de la requête à l'API Spotify
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/users/"+user.Id+"/playlists", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	//Envoie de la requête à l'API Spotify
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	// //Analyse de la réponse json de l'API Spotify

	var myDataList PlaylistsResponse
	var Playlists []UserPlaylist

	json.Unmarshal(reponsebytes, &myDataList)

	for _, j := range myDataList.Items {
		j.Tracks = GetPlaylistTracks(j)
		Playlists = append(Playlists, j)
	}

	//Affichage données renvoyées par API Spotify
	return Playlists
}

func GetAlbumTracks(album SpotifyAlbum) []Track {

	// Configuration du client HTTP
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Obtention d'un jeton d'accès auprès de l'API Spotify
	token, err := GetAccessToken(client, "971b13b2aaf14061a2dd134d0cc18134", "e958e7edc5af4e17b6f55235f4ee62cc")
	if err != nil {
		panic(err)
	}

	//Configuration de la requête à l'API Spotify
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/albums/"+album.Id+"/tracks?&limit=50", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	//Envoie de la requête à l'API Spotify
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	// //Analyse de la réponse json de l'API Spotify

	var myDataList TracksResponse
	var Tracks []Track

	json.Unmarshal(reponsebytes, &myDataList)

	for _, j := range myDataList.Items {
		duration := time.Duration(j.DurationMs) * time.Millisecond
		j.DurationMinute = strconv.Itoa(int(duration.Minutes())) + ":" + strconv.Itoa(int(duration.Seconds())%60)
		j.Images = album.Images
		Tracks = append(Tracks, j)
	}

	//Affichage données renvoyées par API Spotify
	return Tracks
}

func GetPlaylistTracks(playlist UserPlaylist) []Track {

	// Configuration du client HTTP
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Obtention d'un jeton d'accès auprès de l'API Spotify
	token, err := GetAccessToken(client, "971b13b2aaf14061a2dd134d0cc18134", "e958e7edc5af4e17b6f55235f4ee62cc")
	if err != nil {
		panic(err)
	}

	//Configuration de la requête à l'API Spotify
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+playlist.Id+"/tracks?&limit=50", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	//Envoie de la requête à l'API Spotify
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	reponsebytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error creating HTTP request: %v", err)
	}

	// //Analyse de la réponse json de l'API Spotify

	var myDataList PlaylistTracksResponse
	var Tracks []Track

	json.Unmarshal(reponsebytes, &myDataList)

	for _, j := range myDataList.Items {
		j.Track.Images = append(j.Track.Images, struct {
			Url string "json:\"url\""
		}{j.Track.Album.Images[0].URL})
		duration := time.Duration(j.Track.DurationMs) * time.Millisecond
		j.Track.DurationMinute = strconv.Itoa(int(duration.Minutes())) + ":" + strconv.Itoa(int(duration.Seconds())%60)
		Tracks = append(Tracks, j.Track)
	}

	//Affichage données renvoyées par API Spotify
	return Tracks
}
