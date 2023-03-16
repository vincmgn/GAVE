package groupie

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id           int    `json:"id"`
	Pseudo       string `json:"pseudo"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	LikedArtists string `json:"likedArtists"`
	SpotifyUser  string `json:"spotifyUser"`
}

type Erreur struct {
	Profil        User        `json:"profil"`
	SpotifyProfil SpotifyUser `json:"spotifyProfil"`
	Err           int         `json:"err"`
	Like          []Artists   `json:"like"`
}

func Initdb(db *sql.DB) {
	db, err := sql.Open("sqlite3", "API/groupieProfil.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS Users(id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, pseudo TEXT, email TEXT, password TEXT, likedArtists TEXT, spotifyUser TEXT)")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CreateUser(db *sql.DB, user User) {

	db, err := sql.Open("sqlite3", "API/groupieProfil.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	User := "INSERT INTO Users (pseudo, email, password, likedArtists, spotifyUser) VALUES (?,?,?,?,?)"
	stmt, err := db.Prepare(User)
	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := stmt.Exec(user.Pseudo, user.Email, user.Password, user.LikedArtists, user.SpotifyUser)
	if err != nil {
		fmt.Println(err)
		fmt.Println(res)
		return
	}

	user.Password, _ = HashPassword(user.Password)

	return
}

func GetUser(db *sql.DB, Pseudo string) User {
	var user User

	db, err := sql.Open("sqlite3", "API/groupieProfil.db")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	err = db.QueryRow("SELECT * FROM Users WHERE pseudo=?", Pseudo).Scan(&user.Id, &user.Pseudo, &user.Email, &user.Password, &user.LikedArtists, &user.SpotifyUser)
	if err != nil {
		panic(err)
	}

	return user
}

func UpdateUser(db *sql.DB, param string, updater string, email string) {

	db, err := sql.Open("sqlite3", "API/groupieProfil.db")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	stmt, err := db.Prepare("UPDATE Users SET " + param + "=? WHERE email=?")
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := stmt.Exec(updater, email)
	if err != nil {
		fmt.Println(err)
		fmt.Println(res)
		return
	}
}

func VerifPseudo(db *sql.DB, pseudo string) bool {
	var bool bool
	var result string

	db, err := sql.Open("sqlite3", "API/groupieProfil.db")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	err = db.QueryRow("SELECT pseudo FROM Users WHERE pseudo=?", pseudo).Scan(&result)
	if err != nil {
		bool = true
	}

	return bool
}

func VerifEmail(db *sql.DB, email string) bool {
	var bool bool
	var result string

	db, err := sql.Open("sqlite3", "API/groupieProfil.db")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	err = db.QueryRow("SELECT email FROM Users WHERE email=?", email).Scan(&result)
	if err != nil {
		bool = true
	}

	return bool
}

func VerifPassword(db *sql.DB, pseudo string, password string) bool {
	var bool bool
	var result string

	db, err := sql.Open("sqlite3", "API/groupieProfil.db")
	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	err = db.QueryRow("SELECT password FROM Users WHERE pseudo=?", pseudo).Scan(&result)
	if CheckPasswordHash(password, result) {
		bool = true
	}

	return bool
}
