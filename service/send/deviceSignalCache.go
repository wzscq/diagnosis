package send

import (
	"github.com/go-redis/redis/v8"
	"time"
	"log"
	"encoding/json"
)

type DeviceSignalCache struct {
	client *redis.Client
	expire time.Duration
}

func (repo *DeviceSignalCache)Init(url string,db int,expire time.Duration){
	repo.client=redis.NewClient(&redis.Options{
        Addr:     url,
        Password: "", // no password set
        DB:       db,  // use default DB
    })
	repo.expire=expire
}

func (repo *DeviceSignalCache)SaveSignalList(deviceID string,busiType string,signalList *map[string]interface{})(error){
	// Create JSON from the instance data.
    bytes, err := json.Marshal(*signalList)
	if err!=nil {
		log.Println("SaveDiagSignal error:",err.Error())
		return err
	}
    // Convert bytes to string.
    jsonStr := string(bytes)
	log.Println("SaveDiagSignal deviceID:"+deviceID+" content:"+jsonStr)
	return repo.client.Set(repo.client.Context(), deviceID+":"+busiType, jsonStr, repo.expire).Err()
}

func (repo *DeviceSignalCache)GetSignalList(deviceID string,busiType string)(*map[string]interface{},error){
	jsonStr,err:=repo.client.Get(repo.client.Context(), deviceID+":"+busiType).Result()
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
	return &rec,nil
}

