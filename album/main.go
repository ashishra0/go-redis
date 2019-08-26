package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

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
