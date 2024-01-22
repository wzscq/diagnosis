package idm

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type IdmSyncLock struct {
	client *redis.Client
	expire time.Duration
}

func (isl *IdmSyncLock)Init(url string,db int,expire time.Duration,passwd string){
	isl.client=redis.NewClient(&redis.Options{
        Addr:     url,
        Password: passwd, 
        DB:       db,  // use default DB
    })
	isl.expire=expire
}

func (isl *IdmSyncLock)Lock()(bool){
	return isl.client.SetNX(isl.client.Context(),"IdmSyncLock", 1, isl.expire).Val()
}

func (isl *IdmSyncLock)Unlock(){
	isl.client.Del(isl.client.Context(),"IdmSyncLock")
}

