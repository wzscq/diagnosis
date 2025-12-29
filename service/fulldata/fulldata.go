package fulldata

import (
	"digimatrix.com/diagnosis/crv"
	"digimatrix.com/diagnosis/common"
	"fmt"
	"time"
	"os"
	"archive/zip"
	"io"
	"path/filepath"
	"log"
)

var queryFields = []map[string]interface{}{
	{"field": "id"},
	{"field": "file_name"},
}

var gTempFileIndex int = 0

func GetFullDataFileName(dataID string, token string, crvClient *crv.CRVClient) (string, int) {
	commonRep:=crv.CommonReq{
		ModelID:"full_data_rec",
		Filter:&map[string]interface{}{
			"id":dataID,
		},
		Fields:&queryFields,
	}

	rsp,errCode:=crvClient.Query(&commonRep,token)
	if errCode!=common.ResultSuccess {
		return "",errCode
	}
	list,_:=rsp.Result["list"].([]interface{})
	row,_:=list[0].(map[string]interface{})
	fileName,_:=row["file_name"].(string)
	return fileName,common.ResultSuccess
}

func GetTempFileName()(string){
	// Generate a unique temporary filename based on current time and index
	// Format: yyyymmddhhmmss + 3-digit index
	now := time.Now()
	timeStr := now.Format("20060102150405")
	
	// Format the index as a 3-digit string with leading zeros
	indexStr := fmt.Sprintf("%03d", gTempFileIndex)
	
	// Increment the global index for next use
	gTempFileIndex++
	
	// Combine time and index to create unique filename
	return timeStr + indexStr
}

func CreateZipFile(filesInZip *[]string,zipFileName string)(error){
	// Create a new zip file
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add each file to the zip
	for _, filePath := range *filesInZip {
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %v", filePath, err)
		}
		defer file.Close()

		// Get file info for header
		fileInfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to get file info for %s: %v", filePath, err)
		}

		// Create zip file header
		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return fmt.Errorf("failed to create header for %s: %v", filePath, err)
		}

		// Use base name of file in the archive
		header.Name = filepath.Base(filePath)
		header.Method = zip.Deflate

		// Create writer for the file within the zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip entry for %s: %v", filePath, err)
		}

		// Copy file contents to the zip
		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("failed to write %s to zip: %v", filePath, err)
		}
	}

	return nil
}

func GetFullDataFileNames(dataIDs *[]string, token string, crvClient *crv.CRVClient,FullDataConf *common.FullData)(string, int){
	commonRep:=crv.CommonReq{
		ModelID:"full_data_rec",
		Filter:&map[string]interface{}{
			"id":map[string]interface{}{
				"Op.in":dataIDs,
			},
		},
		Fields:&queryFields,
	}

	rsp,errCode:=crvClient.Query(&commonRep,token)
	if errCode!=common.ResultSuccess {
		return "",errCode
	}
	list,_:=rsp.Result["list"].([]interface{})


	if len(list)==1 {
		row,_:=list[0].(map[string]interface{})
		tempFileName,_:=row["file_name"].(string)
		return tempFileName,common.ResultSuccess
	}
	
	tempFileName:=GetTempFileName()
	
	filesInZip:=[]string{}
	
	// Write each file_name to the temporary file, one per line
	for _, item := range list {
		if row, ok := item.(map[string]interface{}); ok {
			if fileName, ok := row["file_name"].(string); ok {
				filesInZip = append(filesInZip, fileName)
			}
		}
	}

	//create zip file
	CreateZipFile(&filesInZip,FullDataConf.TempPath+tempFileName+".zip")

	return FullDataConf.URLPrefix+"/getFullDataFile/"+tempFileName+".zip", common.ResultSuccess
}

func RemoveTempFiles(tempPath string, duration time.Duration){
	// Get all files in the temporary directory
	files, err := os.ReadDir(tempPath)
	if err != nil {
		log.Printf("Error reading directory %s: %v", tempPath, err)
		return
	}

	// Get current time
	now := time.Now()
	// Check each file
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Get file info to check modification time
		filePath := filepath.Join(tempPath, file.Name())
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Printf("Error getting file info for %s: %v", filePath, err)
			continue
		}

		// Check if file is older than 24 hours
		if now.Sub(fileInfo.ModTime()) > duration {
			// Remove the file
			err = os.Remove(filePath)
			if err != nil {
				log.Printf("Error removing file %s: %v", filePath, err)
			} else {
				log.Printf("Removed old temporary file: %s", filePath)
			}
		}
	}
}