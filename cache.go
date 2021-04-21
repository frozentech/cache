package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-redis/redis"
)

var (
	// Redis defines a Redis client
	Redis *redis.Client
)

// MD5 calculates the MD5 hash of a string
func MD5(raw string) string {
	m := md5.New()
	io.WriteString(m, raw)
	return fmt.Sprintf("%x", m.Sum(nil))
}

// ID genetates the token id
func ID(tablename string, id string) string {
	return MD5(fmt.Sprintf("%s_%s", tablename, id))
}

// Get get record in memory
func Get(model interface{}, tablename string, id string) (err error) {

	if Redis == nil {
		err = fmt.Errorf("redis cache not initialize")
		return
	}

	reply := Redis.Get(ID(tablename, id))

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
func Put(model interface{}, tablename string, id string, expiry int) (err error) {

	if Redis == nil {
		err = fmt.Errorf("redis cache not initialize")
		return
	}

	status := Redis.Set(ID(tablename, id), BYTE(model), time.Hour*time.Duration(expiry))
	err = status.Err()
	return

}

// BYTE converts object to byte
func BYTE(model interface{}) []byte {
	body, _ := json.Marshal(model)
	return body
}
