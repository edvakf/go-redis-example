package main

import (
	"fmt"
	"time"

	redis "gopkg.in/redis.v3"
)

var rd *redis.Client

func init() {
	rd = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func main() {
	flush()
	uniquekey()
	list()
	hash()
	sortedSet()
	bit()
}

func flush() {
	err := rd.FlushAll().Err()
	if err != nil {
		panic(err)
	}
}

func getset() {
	err := rd.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rd.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rd.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

func uniquekey() {
	id, _ := rd.Incr("pk").Result()
	fmt.Println("id", id)
	id, _ = rd.Incr("pk").Result()
	fmt.Println("id", id)
	id, _ = rd.Incr("pk").Result()
	fmt.Println("id", id)
}

func list() {
	err := rd.LPush("foo", "1").Err()
	if err != nil {
		panic(err)
	}
	err = rd.LPush("foo", "2").Err()
	if err != nil {
		panic(err)
	}
	err = rd.RPush("foo", "2").Err()
	if err != nil {
		panic(err)
	}
	val, err := rd.RPopLPush("foo", "foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)

	val, err = rd.LPop("foo").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foo", val)

	vals, err := rd.LRange("foo", 0, 10).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("foos", vals)
}

func hash() {
	err := rd.HSet("bar", "hoge", "1").Err()
	if err != nil {
		panic(err)
	}
	err = rd.HSet("bar", "piyo", "2").Err()
	if err != nil {
		panic(err)
	}
	hoge, err := rd.HGet("bar", "hoge").Int64()
	if err != nil {
		panic(err)
	}
	fmt.Println("bar.hoge", hoge)

	keys, err := rd.HGetAll("bar").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("bar.keys", keys)

	m, err := rd.HGetAllMap("bar").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("bar", m)
}

func sortedSet() {
	err := rd.ZAdd("baz", redis.Z{float64(time.Now().Unix()), "hoge"}).Err()
	if err != nil {
		panic(err)
	}
	cnt, err := rd.ZAdd("baz", redis.Z{float64(time.Now().Unix()) + 10, "fuga"}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("sortedSet count: ", cnt)
	vals, err := rd.ZRevRangeWithScores("baz", 0, 2).Result()
	if err != nil {
		panic(err)
	}
	for _, val := range vals {
		t := time.Unix(int64(val.Score), 0)
		valString, _ := val.Member.(string)
		fmt.Println("key", t)
		fmt.Println("val", valString)
	}
	vals, err = rd.ZRangeByScoreWithScores("baz", redis.ZRangeByScore{Min: "-inf", Max: "+inf", Offset: 0, Count: 2}).Result()
	if err != nil {
		panic(err)
	}
	for _, val := range vals {
		t := time.Unix(int64(val.Score), 0)
		valString, _ := val.Member.(string)
		fmt.Println("key", t)
		fmt.Println("val", valString)
	}
}

func bit() {
	err := rd.SetBit("puu", 1, 1).Err()
	if err != nil {
		panic(err)
	}
	err = rd.SetBit("puu", 10, 1).Err()
	if err != nil {
		panic(err)
	}
	bit, err := rd.GetBit("puu", 10).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("bit 10: ", bit)

	cnt, err := rd.BitCount("puu", &redis.BitCount{0, 20}).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("bit count: ", cnt)
}
