package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

// Album struct describes the album structure
type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func main() {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// calling the hmSet()
	err = hmSet(conn)
	if err != nil {
		fmt.Println(err)
	}
	// calling the hGet()
	err = hGet(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = hgetAll(conn)
	if err != nil {
		fmt.Println(err)
	}
}

func hmSet(c redis.Conn) error {
	_, err := c.Do("HMSET", "album:2", "title", "Electronic Ladyland", "artist", "jimi Hendrix", "price", 4.95, "likes", 8)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Electronic Ladyland added")
	return nil
}

func hGet(c redis.Conn) error {
	title, err := redis.String(c.Do("HGET", "album:2", "title"))
	if err != nil {
		log.Fatal(err)
	}
	artist, err := redis.String(c.Do("HGET", "album:2", "artist"))
	if err != nil {
		log.Fatal(err)
	}
	price, err := redis.Float64(c.Do("HGET", "album:2", "price"))
	if err != nil {
		log.Fatal(err)
	}
	likes, err := redis.Int(c.Do("HGET", "album:2", "likes"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s by %s: £%.2f [%d likes]\n", title, artist, price, likes)
	return nil
}

func hgetAll(c redis.Conn) error {
	reply, err := redis.StringMap(c.Do("HGETALL", "album:2"))
	if err != nil {
		log.Fatal(err)
	}
	// Use the populateAlbum helper function to create a new Album
	// object from the map[string]string.
	album, err := populateAlbum(reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", album)
	return nil
}

// Create, populate and return a pointer to a new Album struct, based
// on data from a map[string]string
func populateAlbum(reply map[string]string) (*Album, error) {
	var err error
	album := new(Album)
	album.Title = reply["title"]
	album.Artist = reply["artist"]

	album.Price, err = strconv.ParseFloat(reply["price"], 64)
	if err != nil {
		return nil, err
	}
	album.Likes, err = strconv.Atoi(reply["likes"])
	if err != nil {
		log.Fatal(err)
	}
	return album, nil
}
