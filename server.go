package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	groupie "groupie/API"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var Profil groupie.User

func main() {
	//Demarrage du Serveur
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/artists", ArtistsPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/registration", RegisterPage)
	http.HandleFunc("/profil", ProfilPage)
	http.HandleFunc("/artist", OneArtistPage)
	http.HandleFunc("/map", MapPage)
	http.HandleFunc("/informations", AboutUs)

	var db *sql.DB
	groupie.Initdb(db)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HomePage(w http.ResponseWriter, r *http.Request) {

	var db *sql.DB

	data := groupie.GetArtistsData()
	var mydata groupie.AllArtists
	mydata.User = Profil
	var liked struct {
		ID string `json:"id"`
	}

	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		nb := rand.Intn(52)
		var bool bool
		var tmp groupie.Artists
		for _, j := range mydata.Artists {
			if j.Id == data[nb].Id {
				bool = true
				i--
				break
			}
		}
		if !bool {
			tmp = data[nb]
			mydata.Artists = append(mydata.Artists, tmp)
		}
	}

	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&liked)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if liked.ID != "" {
			var id string
			var bool bool
			for i := 0; i < len(Profil.LikedArtists); i++ {
				if Profil.LikedArtists[i] != ',' {
					id += string(Profil.LikedArtists[i])
				} else {
					if id == liked.ID {
						bool = true
					}
					id = ""
				}
			}
			if !bool {
				Profil.LikedArtists += liked.ID + ","
				groupie.UpdateUser(db, "likedArtists", Profil.LikedArtists, Profil.Email)
			}
		}
	}

	tmpl, _ := template.ParseFiles("./html/pages/index.html", "./html/templates/header.html", "./html/templates/searchbar.html", "./html/templates/footer.html", "./html/templates/vinyle.html")

	tmpl.Execute(w, mydata)
}

func ArtistsPage(w http.ResponseWriter, r *http.Request) {

	var db *sql.DB

	r.ParseForm()

	members := r.Form["nb"]
	name := r.FormValue("search")
	creation := r.FormValue("yearartists")
	album := r.FormValue("yearalbums")
	pays := r.FormValue("pays-select")
	pages := r.FormValue("pages")
	var liked struct {
		ID string `json:"id"`
	}

	data := groupie.AllArtists{Artists: groupie.GetArtistsData()[:12], Pays: groupie.GetPays()}

	data.User = Profil

	if r.Method == "POST" {
		data.Artists = []groupie.Artists{}
		if pages != "" {
			pageint, _ := strconv.Atoi(pages)
			limit := 12 * pageint
			if limit > 52 {
				limit = 52
			}
			data.Artists = append(data.Artists, groupie.GetArtistsData()[12*(pageint-1):limit]...)
		} else if len(members) > 0 {
			for i := 0; i < len(members); i++ {
				data.Artists = append(data.Artists, groupie.GetArtistsByMembers(members[i])...)
			}
		} else if name != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByName(name)...)
		} else if pays != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByLocation(pays)...)
		} else if album != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByFirstAlbumDate(album)...)
		} else if creation != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByCreationDate(creation)...)
		} else {
			err := json.NewDecoder(r.Body).Decode(&liked)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
		}
		if liked.ID != "" {
			var id string
			var bool bool
			for i := 0; i < len(Profil.LikedArtists); i++ {
				if Profil.LikedArtists[i] != ',' {
					id += string(Profil.LikedArtists[i])
				} else {
					if id == liked.ID {
						bool = true
					}
					id = ""
				}
			}
			if !bool {
				Profil.LikedArtists += liked.ID + ","
				groupie.UpdateUser(db, "likedArtists", Profil.LikedArtists, Profil.Email)
			}
		}
	}

	tmpl, _ := template.ParseFiles("./html/pages/allartistspage.html", "./html/templates/header.html", "./html/templates/footer.html", "./html/templates/vinyle.html", "./html/templates/searchbar.html", "./html/templates/searchbar-filter.html")

	tmpl.Execute(w, data)
}

func OneArtistPage(w http.ResponseWriter, r *http.Request) {

	data := groupie.GetOneArtistById(r.FormValue("id"))
	data.SpotifyData = groupie.GetSpotifyArtist(data.Name)

	tmpl, _ := template.ParseFiles("./html/pages/artistpage.html", "./html/templates/header.html", "./html/templates/album-spotify.html", "./html/templates/track.html", "./html/templates/footer.html")

	tmpl.Execute(w, data)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {

	if Profil.Pseudo != "" {
		http.Redirect(w, r, "/profil", http.StatusFound)
	}

	var db *sql.DB
	var erreur int

	pseudo := r.FormValue("pseudo")
	password := r.FormValue("password")

	if r.Method == "POST" {
		erreur = 0
		if !groupie.VerifPseudo(db, pseudo) {
			if !groupie.VerifPassword(db, pseudo, password) {
				Profil = groupie.GetUser(db, pseudo)
				http.Redirect(w, r, "/profil", http.StatusFound)
			} else {
				erreur = 2
			}
		} else {
			erreur = 1
		}
	}

	tmpl, _ := template.ParseFiles("./html/pages/loginpage.html", "./html/templates/header2.html", "./html/templates/footer2.html")

	tmpl.Execute(w, erreur)

}

func RegisterPage(w http.ResponseWriter, r *http.Request) {

	var db *sql.DB
	var erreur int

	pseudo := r.FormValue("pseudo")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmpassword := r.FormValue("confirmpassword")

	if r.Method == "POST" {
		erreur = 0
		if groupie.VerifPseudo(db, pseudo) {
			if groupie.VerifEmail(db, email) {
				if password == confirmpassword {
					groupie.CreateUser(db, groupie.User{Pseudo: pseudo, Email: email, Password: password})
					Profil = groupie.User{Pseudo: pseudo, Email: email, Password: password}
					http.Redirect(w, r, "/profil", http.StatusFound)
				} else {
					erreur = 3
				}
			} else {
				erreur = 2
			}
		} else {
			erreur = 1
		}
	}

	tmpl, _ := template.ParseFiles("./html/pages/registerpage.html", "./html/templates/header2.html", "./html/templates/footer2.html")

	tmpl.Execute(w, erreur)
}

func ProfilPage(w http.ResponseWriter, r *http.Request) {

	var db *sql.DB
	var erreur int

	password := r.FormValue("password")
	chpassword := r.FormValue("chpassword")
	logout := r.FormValue("logout")
	spotify := r.FormValue("spotify")

	spotifydata, _ := groupie.GetSpotifyUser(Profil.SpotifyUser)

	data := groupie.Erreur{Profil: Profil, Err: erreur, SpotifyProfil: spotifydata}

	var liked struct {
		ID string `json:"id"`
	}

	if r.Method == "POST" {
		erreur = 0
		if logout != "" {
			Profil = groupie.User{}
			data.Profil = Profil
			http.Redirect(w, r, "/", http.StatusFound)
		} else if spotify != "" {
			var err error
			data.SpotifyProfil, err = groupie.GetSpotifyUser(spotify)
			if err != nil {
				erreur = 2
			} else {
				Profil.SpotifyUser = spotify
				groupie.UpdateUser(db, "spotifyUser", spotify, Profil.Email)
			}
		} else if password == chpassword && password != "" {
			groupie.UpdateUser(db, "password", password, Profil.Email)
		} else if password != chpassword {
			erreur = 1
		} else {
			err := json.NewDecoder(r.Body).Decode(&liked)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
		}
		if liked.ID != "" {
			Profil.LikedArtists = strings.Replace(Profil.LikedArtists, liked.ID+",", "", 1)
			groupie.UpdateUser(db, "likedArtists", Profil.LikedArtists, Profil.Email)
		}
	}

	data.Like = groupie.GetArtistsLiked(Profil)

	tmpl, _ := template.ParseFiles("./html/pages/profilpage.html", "./html/templates/header2.html", "./html/templates/footer2.html", "./html/templates/vinyle.html", "./html/templates/playlist-spotify.html", "./html/templates/track.html")

	tmpl.Execute(w, data)

}

func MapPage(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	name := r.FormValue("search")
	pays := r.FormValue("pays-select")
	pages := r.FormValue("pages")
	concertDate := r.FormValue("yearartists")
	bubble := r.FormValue("bubble")

	data := groupie.MapData{Artists: groupie.GetArtistsData()[:12], Pays: groupie.GetPays()}
	for _, i := range data.Pays {
		data.Coordonnees += groupie.GetLocationWithAdress(i.Pays)
		data.CooName += i.Pays + ","
	}

	if r.Method == "POST" {
		data.Artists = []groupie.Artists{}
		data.Coordonnees = ""
		if pages != "" {
			pageint, _ := strconv.Atoi(pages)
			limit := 12 * pageint
			if limit > 52 {
				limit = 52
			}
			data.Artists = append(data.Artists, groupie.GetArtistsData()[12*(pageint-1):limit]...)
			data.Coordonnees, data.CooName = groupie.LocationByParams(data)
		} else if name != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByName(name)...)
			data.Coordonnees, data.CooName = groupie.LocationByParams(data)
		} else if pays != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByConcertLocation(pays)...)
			data.Coordonnees, data.CooName = groupie.LocationByParams(data)
		} else if concertDate != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByConcertDate(concertDate)...)
			data.Coordonnees, data.CooName = groupie.LocationByParams(data)
		} else if bubble != "" {
			data.Artists = append(data.Artists, groupie.GetArtistsByConcertLocation(bubble)...)
			data.Coordonnees, data.CooName = groupie.LocationByParams(data)
		}
		tab := strings.Split(data.CooName[:len(data.CooName)-1], ",")
		for _, j := range tab {
			data.CooId += "["
			data.CooArtist += "["
			data.CooDates += "["
			for _, i := range data.Artists {
				for k, v := range i.Relations.DatesLocations {
					if k == j {
						data.CooId += string('"') + strconv.Itoa(i.Id) + string('"') + ","
						data.CooArtist += string('"') + strings.Replace(i.Name, "'", " ", -1) + string('"') + ","
						data.CooDates += string('"') + v[0] + string('"') + ","
					}
				}
			}
			data.CooId = data.CooId[:len(data.CooId)-1] + "],"
			data.CooArtist = data.CooArtist[:len(data.CooArtist)-1] + "],"
			data.CooDates = data.CooDates[:len(data.CooDates)-1] + "],"
		}
		data.CooId = "[" + data.CooId[:len(data.CooId)-1] + "]"
		data.CooArtist = "[" + data.CooArtist[:len(data.CooArtist)-1] + "]"
		data.CooDates = "[" + data.CooDates[:len(data.CooDates)-1] + "]"
	}

	data.Coordonnees = "[" + data.Coordonnees[:len(data.Coordonnees)-1] + "]"
	data.CooName = data.CooName[:len(data.CooName)-1]

	tmpl, _ := template.ParseFiles("./html/pages/mappage.html", "./html/templates/header.html", "./html/templates/searchbar-filter-for-map.html", "./html/templates/result-card.html", "./html/templates/footer.html")

	tmpl.Execute(w, data)

}

func AboutUs(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("./html/pages/about-us.html", "./html/templates/header2.html", "./html/templates/footer2.html")

	tmpl.Execute(w, "")
}
