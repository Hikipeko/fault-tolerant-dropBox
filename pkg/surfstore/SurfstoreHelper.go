package surfstore

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

/* Hash Related */
func GetBlockHashBytes(blockData []byte) []byte {
	h := sha256.New()
	h.Write(blockData)
	return h.Sum(nil)
}

func GetBlockHashString(blockData []byte) string {
	blockHash := GetBlockHashBytes(blockData)
	return hex.EncodeToString(blockHash)
}

/* File Path Related */
func ConcatPath(baseDir, fileDir string) string {
	return baseDir + "/" + fileDir
}

/*
	Writing Local Metadata File Related
*/

const createTable string = `create table if not exists indexes (
		fileName TEXT, 
		version INT,
		hashIndex INT,
		hashValue TEXT
	);`

const insertTuple string = `INSERT INTO indexes (fileName, version, hashIndex, hashValue) VALUES (?, ?, ?, ?);`

// WriteMetaFile writes the file meta map back to local metadata file index.db
func WriteMetaFile(fileMetas map[string]*FileMetaData, baseDir string) error {
	// remove index.db file if it exists
	outputMetaPath := ConcatPath(baseDir, DEFAULT_META_FILENAME)
	if _, err := os.Stat(outputMetaPath); err == nil {
		e := os.Remove(outputMetaPath)
		if e != nil {
			log.Println("Error During Meta Write Back")
		}
	}
	db, err := sql.Open("sqlite3", outputMetaPath)
	if err != nil {
		log.Println("Error During Meta Write Back")
	}
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Println("Error During Meta Write Back")
	}
	statement.Exec()
	
	// write to indexes
	statement, err = db.Prepare(insertTuple)
	for _, fileMeta := range fileMetas {
		for i, hashValue := range fileMeta.BlockHashList {
			statement.Exec(fileMeta.Filename, fileMeta.Version, i, hashValue)
		}
	}
	return nil
}

/*
Reading Local Metadata File Related
*/
const getDistinctFileName string = `SELECT DISTINCT filename, version FROM indexes;`

const getTuplesByFileName string = `SELECT hashIndex, hashValue FROM indexes WHERE fileName = ? ORDER BY hashIndex ASC;`

// LoadMetaFromMetaFile loads the local metadata file into a file meta map.
// The key is the file's name and the value is the file's metadata.
// You can use this function to load the index.db file in this project.
func LoadMetaFromMetaFile(baseDir string) (fileMetaMap map[string]*FileMetaData, e error) {
	metaFilePath, _ := filepath.Abs(ConcatPath(baseDir, DEFAULT_META_FILENAME))
	fileMetaMap = make(map[string]*FileMetaData)
	metaFileStats, e := os.Stat(metaFilePath)
	if e != nil || metaFileStats.IsDir() {
		return fileMetaMap, nil
	}
	db, err := sql.Open("sqlite3", metaFilePath)
	if err != nil {
		log.Fatal("Error When Opening Meta")
	}
	
	// read indexes
	rows, err := db.Query(getDistinctFileName)
	if err != nil {
		return fileMetaMap, fmt.Errorf("LoadMetaFromMetaFile error")
	}

	var filename string
	var version int32
	for rows.Next() {
		rows.Scan(&filename, &version)
		fileMetaMap[filename] = &FileMetaData{Filename: filename, Version: version}
	}

	var hashIndex int
	var hashValue string
	for filename, fileMeta := range fileMetaMap {
		newRows, err := db.Query(getTuplesByFileName, filename)
		if err != nil {
			return fileMetaMap, fmt.Errorf("LoadMetaFromMetaFile error")
		}

		i := 0
		for newRows.Next() {
			newRows.Scan(&hashIndex, &hashValue)
			if i != hashIndex {
				return fileMetaMap, fmt.Errorf("hashIndex error")
			}
			i++
			fileMeta.BlockHashList = append(fileMeta.BlockHashList, hashValue)
		}
	}
	return fileMetaMap, nil
}

/*
	Debugging Related
*/

// PrintMetaMap prints the contents of the metadata map.
// You might find this function useful for debugging.
func PrintMetaMap(metaMap map[string]*FileMetaData) {

	fmt.Println("--------BEGIN PRINT MAP--------")

	for _, filemeta := range metaMap {
		fmt.Println("\t", filemeta.Filename, filemeta.Version)
		for _, blockHash := range filemeta.BlockHashList {
			fmt.Println("\t", blockHash)
		}
	}

	fmt.Println("---------END PRINT MAP--------")

}
