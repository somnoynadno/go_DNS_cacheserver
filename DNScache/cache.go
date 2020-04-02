package DNScache

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

var redisOptions = redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
}

func RedisHealthCheck() error {
	client := redis.NewClient(&redisOptions)
	defer client.Close()

	_, err := client.Ping().Result()

	return err
}

func CheckAnswerInCache(query DNSQuery) ([]DNSAnswer, error) {
	client := redis.NewClient(&redisOptions)
	defer client.Close()

	key := base64.StdEncoding.EncodeToString(query.Inflate())

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	result, err := client.Get(key).Result()
	ttl, err    := client.TTL(key).Result()

	if err != nil || ttl <= 0 {
		return nil, nil
	}

	var answers []DNSAnswer

	for _, a := range strings.Split(result, " ") {
		ans, err := base64.StdEncoding.DecodeString(a)
		if err != nil {
			return nil, err
		}

		if len(ans) > 0 {
			dnsAns := *newDNSAnswer(ans)

			bs := make([]byte, 4)
			var arr [4]byte

			binary.BigEndian.PutUint32(bs, uint32(ttl.Seconds()))
			copy(arr[:], bs[:4])
			dnsAns.TTL = arr

			answers = append(answers, dnsAns)
		}
	}

	return answers, nil
}

func PutAnswersInCache(query DNSQuery, answers []DNSAnswer) error {
	client := redis.NewClient(&redisOptions)
	defer client.Close()

	value := ""
	key := base64.StdEncoding.EncodeToString(query.Inflate())

	ttlMin := uint32(7200)

	for _, a := range answers {
		ttl := binary.BigEndian.Uint32(a.TTL[:])
		if  ttl < ttlMin {
			ttlMin = ttl
		}
		value += base64.StdEncoding.EncodeToString(a.Inflate())
		value += " "
	}

	duration := time.Duration(ttlMin) * time.Second
	_, err := client.Set(key, value, duration).Result()

	return err
}
