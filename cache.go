package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/frozentech/random"
	"github.com/go-redis/redis"
)

var (
	// Redis defines a Redis client
	Redis *redis.Client
)

// ID genetates the token id
func ID(tablename string, id string) string {
	return random.MD5(fmt.Sprintf("%s_%s", tablename, id))
}

// Get get record in memory
func Get(model interface{}, tablename string, id string) (err error) {
	var (
		ctx = context.Background()
	)

	if Redis == nil {
		err = fmt.Errorf("redis cache not initialize")
		return
	}

	reply := Redis.Get(ctx, ID(tablename, id))

	if err = reply.Err(); err != nil {
		return
	}

	if reply.String() == "" {
		err = fmt.Errorf("no result found")
		return
	}

	ec, err := reply.Bytes()
	if err != nil {
		return
	}

	if err = json.Unmarshal(ec, model); err != nil {
		return
	}

	return
}

// Put put record in cache
func Put(model interface{}, tablename string, id string, hour int) (err error) {
	var (
		ctx = context.Background()
	)

	if Redis == nil {
		err = fmt.Errorf("redis cache not initialize")
		return
	}

	status := Redis.Set(ctx, ID(tablename, id), BYTE(model), time.Hour*time.Duration(hour))
	err = status.Err()
	return

}

// BYTE converts object to byte
func BYTE(model interface{}) []byte {
	body, _ := json.Marshal(model)
	return body
}
