package api

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
	c "github.com/kshwedha/GoLRUcache/src/content"
)

type Data struct {
	Key    int `json:"key"`
	Value  any `json:"value"`
	Expiry int `json:"expiry"`
}

type Message struct {
	Message string
	Key     int
	Value   any
}

var cache c.LRUCache = c.Constructor(1024)

func Root(c *fiber.Ctx) error {
	var message Message
	message.Message = "Ahh!! you caught the root."
	return c.JSON(message)
}

func GetHandler(c *fiber.Ctx) error {
	var message Message
	key := c.Params("key")
	int_key, err := strconv.Atoi(key)
	if err != nil {
		// fmt.Fprintf(c, "Pass a valid integer key")
		message.Message = "Pass a valid integer key"
		return c.JSON(message)
	}

	if reflect.TypeOf(key) != nil {
		value := cache.Get(int_key)
		if value != -1 {
			// fmt.Fprintf(c, "value of key: %s is %v", key, value)
			// return nil
			message.Message = "Value found"
			message.Key = int_key
			message.Value = value
			return c.JSON(message)
		}
		// fmt.Fprintf(c, "value not found (uninitialised/expired) for: %s", key)
		// return nil
		message.Message = "value not found (uninitialised/expired)."
		message.Key = int_key
		return c.JSON(message)
	}
	// fmt.Fprintf(c, "Invalid key %s", key)
	// return nil
	message.Message = "Invalid Key."
	return c.JSON(message)
}

func SetHandler(c *fiber.Ctx) error {
	var data Data
	var message Message

	if err := c.BodyParser(&data); err != nil {
		fmt.Println(err)
		return err
	}
	cache.Set(data.Key, data.Value, data.Expiry)

	// fmt.Fprintf(c, "Value %v has been set to key %d\n", c.Params(str_val), c.Params(data.Key))
	message.Message = "Value has been set to key."
	message.Key = data.Key
	message.Value = data.Value
	return c.JSON(message)
}
