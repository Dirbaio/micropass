package main // import "dirba.io/micropass"
import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Hello World!")
	db := &Database{}
	err := databaseSave("db", "omgwtflol", db)
	if err != nil {
		log.Fatal(err)
	}
	db, err = databaseLoad("db", "omgwtflol")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("yay!")

}
