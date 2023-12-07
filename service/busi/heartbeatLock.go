package busi

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type HeartbeatLock struct {
	client *redis.Client
	expire time.Duration
}

func (hbl *HeartbeatLock)Init(url string,db int,expire time.Duration,passwd string){
	hbl.client=redis.NewClient(&redis.Options{
        Addr:     url,
        Password: passwd, 
        DB:       db,  // use default DB
    })
	hbl.expire=expire
}

func (hbl *HeartbeatLock)Lock(deviceID string)(bool){
	return hbl.client.SetNX(hbl.client.Context(),"heartbeat_"+deviceID, 1, hbl.expire).Val()
}

func (hbl *HeartbeatLock)Unlock(deviceID string){
	hbl.client.Del(hbl.client.Context(),"heartbeat_"+deviceID)
}

