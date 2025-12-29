package fulldata

import (
	"testing"
	"digimatrix.com/diagnosis/crv"
	"log"
	"digimatrix.com/diagnosis/common"
	"os"
	"path/filepath"
)

func _TestGetFullDataFileName(t *testing.T) {
	CRVClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:"carapiv2",
	}
	
	filename,errcode:=GetFullDataFileName("32","carapiv2",CRVClient)
	if errcode!=common.ResultSuccess {
		t.Errorf("GetFullDataFileName failed")
	}
	log.Println(filename)	
}

func _TestGetFullDataFileNames(t *testing.T) {
	CRVClient := &crv.CRVClient{
		Server: "http://127.0.0.1:8200",
		Token:"carapiv2",
	}
	
	ids:=[]string{
		"76",
		"77",
		"78",
		"79",
	}
	filename,errcode:=GetFullDataFileNames(&ids,"carapiv2",CRVClient)
	if errcode!=common.ResultSuccess {
		t.Errorf("GetFullDataFileName failed")
	}
	log.Println(filename)	
}

func TestCreateZipFile(t *testing.T) {
	// Create temporary test files
	tempDir := ""//t.TempDir()
	log.Println(tempDir)
	
	// Create first test file
	file1Path := filepath.Join(tempDir, "test1.txt")
	err := os.WriteFile(file1Path, []byte("Test content 1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	
	// Create second test file
	file2Path := filepath.Join(tempDir, "test2.txt")
	err = os.WriteFile(file2Path, []byte("Test content 2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	
	// Files to be zipped
	filesToZip := []string{file1Path, file2Path}
	
	// Create zip file
	zipFilePath := filepath.Join(tempDir, "test.zip")
	err = CreateZipFile(&filesToZip, zipFilePath)
	if err != nil {
		t.Errorf("CreateZipFile failed: %v", err)
	}
}
