package api

import (
	"fmt"
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

	value := cache.Get(int_key)
	message.Key = int_key
	if value != -1 {
		message.Message = "Value found"
		message.Value = value
		return c.JSON(message)
	}
	message.Message = "value not found (uninitialised/expired)."
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
	message.Message = "Value has been set to key."
	message.Key = data.Key
	message.Value = data.Value
	return c.JSON(message)
}
