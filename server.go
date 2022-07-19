package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	_ "github.com/lib/pq"
)

const (
	//on docker use docker ip
	host     = "172.26.0.1"
	port     = 5435
	user     = "root"
	password = "root"
	dbname   = "postgres"
)

func main() {

	e := echo.New()

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	e.GET("/", hello)

	e.GET("/actor", getActor)

	e.GET("/actor/:id", getActorById)

	e.DELETE("/actor/:id", deleteActorById)

	e.PUT("/actor/:id", updateActorById)

	//CALL procedure
	e.POST("/actor", insertActor)

	e.Logger.Fatal(e.Start(":3000"))

}

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Main Page")
}

type Actor struct {
	Actor_id    string `json:"actor_id"`
	First_name  string `json:"first_name"`
	Last_name   string `json:"last_name"`
	Last_update string `json:"-"`
}

func getActor(c echo.Context) error {
	db := OpenConnection()
	rows, err := db.Query("SELECT actor_id, first_name, last_name, last_update FROM actor;")

	if err != nil {
		log.Fatal(err)
	}

	var actors []Actor

	for rows.Next() {
		var actor Actor
		rows.Scan(&actor.Actor_id, &actor.First_name, &actor.Last_name, &actor.Last_update)
		actors = append(actors, actor)
	}

	// log.Println(actors[0].First_name)

	fmt.Printf("%+v \n", actors[0])
	defer rows.Close()
	defer db.Close()

	// json.Marshal(people) for encoding the structure
	j, err := json.Marshal(actors)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		// json.Unmarshal(people) for decoding json
		fmt.Println(string(j))
	}
	return c.JSON(http.StatusOK, actors)

}

func getActorById(c echo.Context) error {
	db := OpenConnection()
	id := string(c.Param("id"))
	fmt.Println(id)
	rows, err := db.Query("SELECT actor_id, first_name, last_name, last_update FROM actor where actor_id =" + id)

	if err != nil {
		log.Fatal(err)
		defer rows.Close()
		defer db.Close()
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	var actors []Actor

	for rows.Next() {
		var actor Actor
		rows.Scan(&actor.Actor_id, &actor.First_name, &actor.Last_name, &actor.Last_update)
		actors = append(actors, actor)
	}

	// log.Println(actors[0].First_name)
	fmt.Printf("%+v \n", actors[0])
	defer rows.Close()
	defer db.Close()

	// json.Marshal(people) for encoding the structure
	j, err := json.Marshal(actors)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		// json.Unmarshal(people) for decoding json
		fmt.Println(string(j))
	}
	return c.JSON(http.StatusOK, actors)

}

func deleteActorById(c echo.Context) error {
	db := OpenConnection()
	id := string(c.Param("id"))
	fmt.Println(id)
	_, err := db.Query("DELETE FROM actor where actor_id =" + id)

	if err != nil {
		log.Fatal(err)
		defer db.Close()
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.String(http.StatusOK, "THE ACTOR ID "+id+" HAS BEEN DELETED")

}

func updateActorById(c echo.Context) error {
	db := OpenConnection()
	id := string(c.Param("id"))
	fname := string(c.FormValue("fname"))
	lname := string(c.FormValue("lname"))
	fmt.Println(id)
	_, err := db.Query("UPDATE actor SET first_name= '" + fname + "',last_name='" + lname + "' where actor_id =" + id)

	if err != nil {
		log.Fatal(err)
		defer db.Close()
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	return c.String(http.StatusOK, "THE ACTOR ID "+id+" HAS BEEN UPDATED")

}

type Addactor struct {
	Id_actor string `json:"a_id"`
}

func insertActor(c echo.Context) error {
	db := OpenConnection()
	fname := string(c.FormValue("fname"))
	lname := string(c.FormValue("lname"))
	fmt.Print("CALL addactor('fname', 'lname',null)")
	rows, err := db.Query("CALL addactor('" + fname + "', '" + lname + "',null)")

	if err != nil {
		log.Fatal(err)
	}

	var id_actor []Addactor

	for rows.Next() {
		var actorID Addactor
		rows.Scan(&actorID.Id_actor)
		id_actor = append(id_actor, actorID)
	}
	fmt.Printf(id_actor[0].Id_actor)
	defer rows.Close()
	defer db.Close()

	// j, err := json.Marshal(id_actor)
	// if err != nil {
	// 	fmt.Printf("Error: %s", err.Error())
	// } else {
	// 	fmt.Println(string(j))
	// }
	return c.JSON(http.StatusCreated, id_actor[0].Id_actor)

}
