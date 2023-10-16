package send

import (
	"github.com/go-redis/redis/v8"
	"time"
	"log"
	"encoding/json"
)

type SendRecordCache struct {
	client *redis.Client
	expire time.Duration
}

func (repo *SendRecordCache)Init(url string,db int,expire time.Duration,passwd string){
	repo.client=redis.NewClient(&redis.Options{
        Addr:     url,
        Password: passwd, 
        DB:       db,  // use default DB
    })
	repo.expire=expire
}

func (repo *SendRecordCache)SaveSendRecord(deviceID string,rec map[string]interface{})(error){
	// Create JSON from the instance data.
    bytes, err := json.Marshal(rec)
	if err!=nil {
		log.Println("save send record error:",err.Error())
		return err
	}
    // Convert bytes to string.
    jsonStr := string(bytes)
	log.Println("saveSendRecord deviceID:"+deviceID+" row:"+jsonStr)
	return repo.client.Set(repo.client.Context(), deviceID, jsonStr, repo.expire).Err()
}

func (repo *SendRecordCache)GetSendRecord(deviceID string)(map[string]interface{},error){
	jsonStr,err:=repo.client.Get(repo.client.Context(), deviceID).Result()
	if err!=nil {
		log.Println("get send record error:",err.Error())
		return nil,err
	}
	// Get byte slice from string.
    bytes := []byte(jsonStr)
	rec:=map[string]interface{}{}
	err = json.Unmarshal(bytes, &rec)
	if err!=nil {
		log.Println("get send record error:",err.Error())
		return nil,err
	}
	return rec,nil
}

