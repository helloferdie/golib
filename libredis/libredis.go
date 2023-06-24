package libredis

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/helloferdie/golib/liblogger"

	"github.com/redis/go-redis/v9"
)

var initialize = false
var enable = false
var opt = redis.Options{}

// Client - Redis client instance
type Client struct {
	HasInitialize bool
	Redis         *redis.Client
	Duration      time.Duration
	Enable        bool
	ctx           context.Context
}

// loadConfig -
func loadConfig() {
	if !initialize {
		db, err := strconv.Atoi(os.Getenv("redis_db"))
		if err != nil {
			db = 0
		}

		if os.Getenv("redis") == "1" {
			enable = true
		}

		opt = redis.Options{
			Addr:     os.Getenv("redis_address"),
			Username: os.Getenv("redis_username"),
			Password: os.Getenv("redis_password"),
			DB:       db,
		}
		initialize = true
	}
}

// Initialize - Initialize client
func (cl *Client) Initialize() {
	loadConfig()

	cl.ctx = context.Background()
	cl.Redis = redis.NewClient(&opt)
	cl.Duration = time.Minute * time.Duration(5)
	cl.Enable = enable
	cl.HasInitialize = true
}

// Set - Set key & value
func (cl *Client) Set(key string, value interface{}, force bool) error {
	if !cl.HasInitialize {
		cl.Initialize()
	}

	err := errors.New("unknown")
	if force {
		err = cl.Redis.Set(cl.ctx, key, value, cl.Duration).Err()
	} else {
		// Set only if not exist
		err = cl.Redis.SetNX(cl.ctx, key, value, cl.Duration).Err()
	}

	if err != nil {
		tmpErr := errors.New("Failed to set redis key")
		liblogger.Log(nil, true).Errorf("Error: %v %s", tmpErr, key)
		return tmpErr
	}
	return nil
}

// SetCustomDuration - Set key & value with custom duration
func (cl *Client) SetCustomDuration(key string, value interface{}, force bool, duration time.Duration) error {
	if !cl.HasInitialize {
		cl.Initialize()
	}

	err := errors.New("unknown")
	if force {
		err = cl.Redis.Set(cl.ctx, key, value, duration).Err()
	} else {
		err = cl.Redis.SetNX(cl.ctx, key, value, duration).Err()
	}

	if err != nil {
		tmpErr := errors.New("Failed to set redis key")
		liblogger.Log(nil, true).Errorf("Error: %v %s", tmpErr, key)
		return tmpErr
	}
	return nil
}

// Get - Get value based on provided key
func (cl *Client) Get(key string) (interface{}, bool, error) {
	if !cl.HasInitialize {
		cl.Initialize()
	}

	if !enable {
		return "", false, nil
	}

	val, err := cl.Redis.Get(cl.ctx, key).Result()
	if err == redis.Nil {
		return val, false, nil
	} else if err != nil {
		tmpErr := errors.New("Failed to get redis key")
		liblogger.Log(nil, true).Errorf("Error: %v %s", tmpErr, key)
		return val, false, tmpErr
	}
	return val, true, nil
}

// GetUnmarshal - Get unmarshal value based on provided key
func (cl *Client) GetUnmarshal(key string, dt interface{}) (bool, error) {
	if !cl.HasInitialize {
		cl.Initialize()
	}

	if !enable {
		return false, nil
	}

	rdVal, err := cl.Redis.Get(cl.ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		tmpErr := errors.New("Failed to get redis key")
		liblogger.Log(nil, true).Errorf("Error: %v %s", tmpErr, key)
		return false, tmpErr
	}

	errJSON := json.Unmarshal([]byte(rdVal), &dt)
	if errJSON != nil {
		tmpErr := errors.New("Failed to get unmarshal redis key")
		liblogger.Log(nil, true).Errorf("Error: %v %s", tmpErr, key)
		return false, nil
	}
	return true, nil
}

// Delete - Delete value based on provided key
func (cl *Client) Delete(key string) (int64, bool, error) {
	if !cl.HasInitialize {
		cl.Initialize()
	}

	if !enable {
		return 0, false, nil
	}

	val, err := cl.Redis.Del(cl.ctx, key).Result()
	if err == redis.Nil {
		return val, false, nil
	} else if err != nil {
		tmpErr := errors.New("Failed to delete redis key")
		liblogger.Log(nil, true).Errorf("Error: %v %s", tmpErr, key)
		return val, false, tmpErr
	}
	return val, true, nil
}
