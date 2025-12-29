package busi

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	obs "github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type ListBucketResult struct {
	Name      string `xml:"Name"`
	Prefix    string `xml:"Prefix"`
	Marker    string `xml:"Marker"`
	MaxKeys   int    `xml:"MaxKeys"`
	IsTruncated bool   `xml:"IsTruncated"`
	Contents  []struct {
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string `xml:"ETag"`
		Size         int64  `xml:"Size"`
		StorageClass string `xml:"StorageClass"`
	} `xml:"Contents"`
}

type HuaweiOSSClient struct {
	AccessKeyID string
	SecretAccessKey string
	EndPoint string
	BucketName string
	OutputPath string
}

func (hoc *HuaweiOSSClient)ListObjectsSDK(){
	obsClient, err := obs.New(hoc.AccessKeyID, hoc.SecretAccessKey, hoc.EndPoint)
	if err != nil {
		log.Println("create obs client failed:", err)
		return
	}
	defer obsClient.Close()

	input := &obs.ListObjectsInput{}
    // 指定存储桶名称
    input.Bucket = hoc.BucketName
    // 列举桶内对象
    output, err := obsClient.ListObjects(input)
	if(err != nil){
		log.Printf("List objects under the bucket(%s) fail!\n", input.Bucket)
		if obsError, ok := err.(obs.ObsError); ok {
			log.Println("An ObsError was found, which means your request sent to OBS was rejected with an error response.")
			log.Println(obsError.Error())
		} else {
			log.Println("An Exception was found, which means the client encountered an internal problem when attempting to communicate with OBS, for example, the client was unable to access the network.")
			log.Println(err)
		}
		return
	}

    log.Printf("List objects under the bucket(%s) successful!\n", input.Bucket)
    log.Printf("RequestId:%s\n", output.RequestId)
    for index, val := range output.Contents {
        log.Printf("Content[%d]-OwnerId:%s, ETag:%s, Key:%s, LastModified:%s, Size:%d\n",index, val.Owner.ID, val.ETag, val.Key, val.LastModified, val.Size)
    }
}

func (hoc *HuaweiOSSClient)GetObjectSDK(key string){
	obsClient, err := obs.New(hoc.AccessKeyID, hoc.SecretAccessKey, hoc.EndPoint)
	if err != nil {
		log.Println("create obs client failed:", err)
		return
	}
	defer obsClient.Close()

	input := &obs.GetObjectInput{}
    // 指定存储桶名称
    input.Bucket = hoc.BucketName
    // 指定下载对象，此处以 example/objectname 为例。
    input.Key = key
    // 流式下载对象
    output, err := obsClient.GetObject(input)

	if err != nil {
		log.Printf("List objects under the bucket(%s) fail!\n", input.Bucket)
		if obsError, ok := err.(obs.ObsError); ok {
			log.Println("An ObsError was found, which means your request sent to OBS was rejected with an error response.")
			log.Println(obsError.Error())
		} else {
			log.Println("An Exception was found, which means the client encountered an internal problem when attempting to communicate with OBS, for example, the client was unable to access the network.")
			log.Println(err)
		}
	}

	defer output.Body.Close()
	log.Printf("Get object(%s) under the bucket(%s) successful!\n", input.Key, input.Bucket)
    log.Printf("StorageClass:%s, ETag:%s, ContentType:%s, ContentLength:%d, LastModified:%s\n",output.StorageClass, output.ETag, output.ContentType, output.ContentLength, output.LastModified)
	
	fileName := strings.Split(key, "/")[len(strings.Split(key, "/"))-1]
	filePath := hoc.OutputPath +"/"+ fileName
	file, err := os.Create(filePath)
	if err != nil {
		log.Println("create file failed:", err)
		return
	}
	defer file.Close()

	//body, err := ioutil.ReadAll(output.Body)
	// 读取对象内容
	p := make([]byte, 1024)
	var readErr error
	var readCount int
	for {
		readCount, readErr = output.Body.Read(p)
		if readCount > 0 {
			//fmt.Printf("%s", p[:readCount])
			file.Write(p[:readCount])
		}

		if readErr != nil {
			log.Println("read object failed:", readErr)
			break
		}
	}
	log.Printf("Save object(%s) to local file(%s) successful!\n", input.Key, filePath)
}

func (hoc *HuaweiOSSClient)ListObjects(){
	// 调用华为云OBS的List Objects REST API查询桶内对象列表
	// 这里只做简单演示，未做复杂签名校验等，生产应用建议使用官方SDK

	// 构建请求URL
	url := fmt.Sprintf("https://%s.%s", hoc.BucketName, hoc.EndPoint)

	fmt.Println("url:", url)
	
	// 简单版本：未签名（若桶是公开读权限可用，否则需实现签名）
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("create request failed:", err)
		return
	}

	fmt.Println("req:", req)

	// 此处如需签名认证，请实现Authorization头
	req.Header.Set("Authorization", "OBS HPUAZSP3LUEP3S86CUCP:5a5D/j8UFeMIZbbgBGeeOT7CPwY=")
	

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("request obs bucket failed:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("resp:", resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read obs response failed:", err)
		return
	}

	// 解析返回的XML
	var result ListBucketResult
	if err := xml.Unmarshal(body, &result); err != nil {
		log.Println("parse obs xml failed:", err)
		return
	}

	// 打印对象列表
	for _, obj := range result.Contents {
		fmt.Printf("Object Key: %s, Size: %d, LastModified: %s\n", obj.Key, obj.Size, obj.LastModified)
	}
}



