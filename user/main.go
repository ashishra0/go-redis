package main

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// User struct defines the user object
type User struct {
	Username  string
	MobileID  int
	Email     string
	FirstName string
	LastName  string
}

func main() {
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()

	err := ping(conn)
	if err != nil {
		fmt.Println(err)
	}
	err = set(conn)
	if err != nil {
		fmt.Println(err)
	}
	err = get(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = setStruct(conn)
	if err != nil {
		fmt.Println(err)
	}

	err = getStruct(conn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("success!")
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func ping(c redis.Conn) error {
	// Send PING command to redis
	// PING command returns a redis "Simple String"
	// use redis.String to convert interface type to string
	s, err := redis.String(c.Do("PING"))
	if err != nil {
		return err
	}

	fmt.Printf("PING response = %s\n", s)
	return nil
}

// set executes the redis SET command
func set(c redis.Conn) error {
	_, err := c.Do("SET", "Favorite Movie", "Avengers Endgame")
	if err != nil {
		return err
	}
	_, err = c.Do("SET", "Release Year", 2019)
	if err != nil {
		return err
	}
	return nil
}

// get executes the redis GET command
func get(c redis.Conn) error {
	// Simple get command with string helper
	key := "Favorite Movie"
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %s\n", key, s)

	// Simple get command with Int helper
	key = "Release Year"
	i, err := redis.Int(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %d\n", key, i)
	return nil
}

func setStruct(c redis.Conn) error {
	const objectPrefix string = "user:"

	usr := User{
		Username:  "otto",
		MobileID:  123455666,
		Email:     "otto@micro.com",
		FirstName: "Otto",
		LastName:  "Maddox",
	}
	// serialize User object to json
	json, err := json.Marshal(usr)
	if err != nil {
		return err
	}
	// Set object
	_, err = c.Do("SET", objectPrefix+usr.Username, json)
	if err != nil {
		return err
	}
	return nil
}

func getStruct(c redis.Conn) error {
	const objectPrefix string = "user:"
	username := "otto"
	s, err := redis.String(c.Do("GET", objectPrefix+username))
	if err == redis.ErrNil {
		fmt.Println("User does not exist")
	} else if err != nil {
		return err
	}
	usr := User{}
	err = json.Unmarshal([]byte(s), &usr)
	fmt.Printf("%+v\n", usr)
	return nil
}
