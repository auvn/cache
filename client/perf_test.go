package client

import (
	"testing"
)

var (
	testString   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	testValue    = []byte(testString)
	testValueLen = len(testValue)
)

var (
	Host     = "localhost"
	Addrs    = []string{"localhost:1234"}
	Pass     = ""
	PoolSize = 50
)

func check(b *testing.B, i interface{}, e error) {
	if e != nil {
		b.Fatal(e)
	}
}

func runBenchmark(b *testing.B, fn func(Cache) error) {
	c := New(&Options{
		Addrs:    Addrs,
		PoolSize: PoolSize,
		Auth:     Pass,
	})
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := fn(c)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.ReportAllocs()
}

func BenchmarkClientSet(b *testing.B) {
	var key = "key_set"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.Set(key, testValue).Bool()
			return err
		},
	)
}

func BenchmarkClientGet(b *testing.B) {
	var key = "key_set"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.Get(key).Bytes()
			return err
		},
	)
}

func BenchmarkClientSetGet(b *testing.B) {
	var key = "key_set_get"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.Set(key, testValue).Bool()
			if err != nil {
				return err
			}
			_, err = c.Get(key).Bytes()
			return err
		},
	)
}

func BenchmarkClientSetDel(b *testing.B) {
	var key = "key_set_del"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.Set(key, testValue).Bool()
			if err != nil {
				return err
			}
			_, err = c.Del(key).Int()
			return err
		},
	)
}

func benchmarkClientLPush(b *testing.B, llen int) {
	var key = "lpush"
	var values = make([][]byte, llen)
	for i, _ := range values {
		values[i] = []byte(testString)
	}
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.LPush(key, values...).Int()
			return err
		},
	)
}

func BenchmarkClientLPush10(b *testing.B) {
	benchmarkClientLPush(b, 10)
}

func BenchmarkClientLPush50(b *testing.B) {
	benchmarkClientLPush(b, 50)
}

func BenchmarkClientLPush100(b *testing.B) {
	benchmarkClientLPush(b, 100)
}

func benchmarkClientRPush(b *testing.B, llen int) {
	var key = "rpush"
	var values = make([][]byte, llen)
	for i, _ := range values {
		values[i] = []byte(testString)
	}
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.RPush(key, values...).Int()
			return err
		},
	)
}

func BenchmarkClientRPush10(b *testing.B) {
	benchmarkClientRPush(b, 10)
}

func BenchmarkClientRPush50(b *testing.B) {
	benchmarkClientRPush(b, 50)
}

func BenchmarkClientRPush100(b *testing.B) {
	benchmarkClientRPush(b, 100)
}

func BenchmarkClientRPop(b *testing.B) {
	var key = "rpush"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.RPop(key).Bytes()
			return err
		},
	)
}
func BenchmarkClientLPop(b *testing.B) {
	var key = "rpush"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.LPop(key).Bytes()
			return err
		},
	)
}

func BenchmarkClientLPushLPop(b *testing.B) {
	var key = "lpush_lpop"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.LPush(key, testValue).Int()
			if err != nil {
				return err
			}
			_, err = c.LPop(key).Bytes()
			return err
		},
	)
}

func BenchmarkClientLPushRPop(b *testing.B) {
	var key = "lpush_rpop"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.LPush(key, testValue).Int()
			if err != nil {
				return err
			}
			_, err = c.RPop(key).Bytes()
			return err
		},
	)
}

func BenchmarkClientRPushRPop(b *testing.B) {
	var key = "rpush_rpop"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.RPush(key, testValue).Int()
			if err != nil {
				return err
			}
			_, err = c.RPop(key).Bytes()
			return err
		},
	)
}

func BenchmarkClientRPushLPop(b *testing.B) {
	var key = "rpush_lpop"
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.RPush(key, testValue).Int()
			if err != nil {
				return err
			}
			_, err = c.LPop(key).Bytes()
			return err
		},
	)
}

func BenchmarkClientHSet(b *testing.B) {
	var key = "hset"
	var hashKey = []byte(testString)
	var hashValue = []byte(testString + testString)
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.HSet(key, hashKey, hashValue).Bool()
			return err
		},
	)
}

func BenchmarkClientHGet(b *testing.B) {
	var key = "hset"
	var hashKey = []byte(testString)
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.HGet(key, hashKey).Bytes()
			return err
		},
	)
}

func BenchmarkClientHSetHGet(b *testing.B) {
	var key = "hset_hget"
	var hashKey = []byte(testString)
	var hashValue = []byte(testString + testString)
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.HSet(key, hashKey, hashValue).Bool()
			if err != nil {
				return err
			}
			_, err = c.HGet(key, hashKey).Bytes()
			return err
		},
	)
}

func BenchmarkClientHSetHDel(b *testing.B) {
	var key = "hset_hdel"
	var hashKey = []byte(testString)
	var hashValue = []byte(testString + testString)
	runBenchmark(
		b,
		func(c Cache) error {
			_, err := c.HSet(key, hashKey, hashValue).Bool()
			if err != nil {
				return err
			}
			_, err = c.HDel(key, hashKey).Int()
			return err
		},
	)
}
