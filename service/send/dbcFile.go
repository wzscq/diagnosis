package send

import (
	"digimatrix.com/diagnosis/common"
	"log"
	"encoding/base64"
	"os"
	"strings"
)

const (
	CC_FILECONTENT = "contentBase64"
	CC_FILENAME = "name"
)

func saveDBCFile(path string,fieldValue interface{})(string,int){
	log.Println("start saveDBCFile save ... ")
	mapValue,_:=fieldValue.(map[string]interface{})
	mapList,ok:=mapValue["list"]
	if !ok {
		log.Println("saveDBCFile end without file")
		return "",common.ResultSuccess
	}

	list,ok:=mapList.([]interface{})
	if !ok || len(list)<=0 {
		log.Println("saveDBCFile end without file")
		return "",common.ResultSuccess
	}

	for _,row:=range list {
		mapRow,ok:=row.(map[string]interface{})
		if(!ok){
			continue
		}
		return createFileRow(path,mapRow)
	}

	log.Println("end saveDBCFile")
	return "",common.ResultSuccess
}

func createFileRow(path string,row map[string]interface{})(string,int){
	log.Println("createFileRow ... ")
	nameCol,_:=row[CC_FILENAME]
	name,_:= nameCol.(string)
	
	contentCol,_:=row[CC_FILECONTENT]
	contentBase64,_:= contentCol.(string)
	
	errorCode:=saveFileRow(path,name,contentBase64)
	if errorCode!=common.ResultSuccess {
		return "",errorCode
	}

	log.Println("createFileRow end ")
	return name,common.ResultSuccess
}

func saveFileRow(path,name,contentBase64 string)(int){
	//判断并创建文件路径
	err := os.MkdirAll(path, 0750)
	if err != nil && !os.IsExist(err) {
		log.Println("create dir error:", err)
		return common.ResultCreateDirError
	}

	//Base64转码
	//log.Printf("file content: %s",contentBase64)
	//去掉url头信息
	typeIndex:=strings.Index(contentBase64, "base64,")
	if typeIndex>0 {
		contentBase64=contentBase64[typeIndex+7:]
	}
	fileContent := make([]byte, base64.StdEncoding.DecodedLen(len(contentBase64)))
	n, err := base64.StdEncoding.Decode(fileContent, []byte(contentBase64))
	if err != nil {
		log.Println("decode error:", err)
		return common.ResultBase64DecodeError
	}
	fileContent = fileContent[:n]

	//保存文件
	file,err:=os.Create(path+name)
	if err != nil {
		log.Println("create file error:", err)
		return common.ResultCreateFileError
	}

	if _, err := file.Write(fileContent); err != nil {
		log.Println("write file error:", err)
		return common.ResultCreateFileError
	}

	if err := file.Close(); err != nil {
		log.Println("close file error:", err)
		return common.ResultCreateFileError
	}

	return common.ResultSuccess
}