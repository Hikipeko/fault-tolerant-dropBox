package surfstore

import (
	"io"
	"log"
	"os"
	"strings"
)

func ValidFileName(filename string) bool {
	if filename == "index.db" {
		return false
	}
	return !strings.Contains(filename, ",") && !strings.Contains(filename, "/")
}

func GetBlockHashList(client RPCClient, path string) []string {
	hashList := []string{}
	file, err := os.Open(path)
	if err != nil {
		log.Println(err.Error())
	}

	readBuffer := make([]byte, client.BlockSize)
	for {
		l, err := file.Read(readBuffer)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err.Error())
		}
		hash := GetBlockHashString(readBuffer[:l])
		hashList = append(hashList, hash)
	}
	if len(hashList) == 0 {
		hashList = append(hashList, "-1")
	}
	return hashList
}

func SlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range len(a) {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func UpdateLocalIndexes(client RPCClient, localIndexes *map[string]*FileMetaData) {
	files, err := os.ReadDir(client.BaseDir)
	if err != nil {
		log.Println(err.Error())
	}

	for _, file := range files {
		fn := file.Name()
		if !ValidFileName(fn) || file.IsDir() {
			continue
		}
		newHashList := GetBlockHashList(client, ConcatPath(client.BaseDir, file.Name()))

		// if new file, add it to indexes
		if fileMeta, ok := (*localIndexes)[fn]; !ok {
			(*localIndexes)[fn] = &FileMetaData{Filename: fn, Version: 1, BlockHashList: newHashList}
		} else if !SlicesEqual(fileMeta.BlockHashList, newHashList) {
			// modified
			fileMeta.BlockHashList = newHashList
			fileMeta.Version++
		}
	}

	// deleted file
	for fn, fileMeta := range *localIndexes {
		path := ConcatPath(client.BaseDir, fn)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if len(fileMeta.BlockHashList) == 1 && fileMeta.BlockHashList[0] == "0" {
				// already deleted
				continue
			}
			fileMeta.Version++
			fileMeta.BlockHashList = []string{"0"}
		}
	}
}

func DownloadFromServer(client RPCClient, remoteMeta *FileMetaData) error {
	// if file deleted, do nothing
	path := ConcatPath(client.BaseDir, remoteMeta.Filename)
	if len(remoteMeta.BlockHashList) == 1 && remoteMeta.BlockHashList[0] == "0" {
		return os.Remove(path)
	}
	// receive content unordered
	blockStoreMap := make(map[string][]string)
	err := client.GetBlockStoreMap(remoteMeta.BlockHashList, &blockStoreMap)
	if err != nil {
		log.Println(err.Error())
	}

	// hash2data: blockHash -> data
	hash2data := make(map[string]string)
	for server, hashes := range blockStoreMap {
		for _, h := range hashes {
			// empty file
			if h == "-1" {
				continue
			}
			var b Block
			err := client.GetBlock(h, server, &b)
			if err != nil {
				log.Println(err.Error())
			}
			hash2data[h] = string(b.BlockData[:b.BlockSize])
		}
	}

	// assemble data in correct order
	content := ""
	for _, h := range remoteMeta.BlockHashList {
		// empty file
		if h == "-1" {
			continue
		}
		content += hash2data[h]
	}

	// write to file
	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	file, err := os.Create(path)
	if err != nil {
		log.Println(err.Error())
	}
	file.WriteString(content)
	log.Printf("Download file %v from server", remoteMeta.Filename)
	return nil
}

func UploadToServer(client RPCClient, localMeta *FileMetaData) error {
	var newVersion int32

	// hash2server: blockHash -> serverAddr
	blockStoreMap := make(map[string][]string)
	err := client.GetBlockStoreMap(localMeta.BlockHashList, &blockStoreMap)
	if err != nil {
		log.Println(err.Error())
	}
	hash2server := make(map[string]string)
	for server, hashes := range blockStoreMap {
		for _, h := range hashes {
			hash2server[h] = server
		}
	}
	// upload to blockserver only if file not deleted
	if !(len(localMeta.BlockHashList) == 1 && localMeta.BlockHashList[0] == "0") {
		file, err := os.Open(ConcatPath(client.BaseDir, localMeta.Filename))
		if err != nil {
			log.Println(err.Error())
		}
		readBuffer := make([]byte, client.BlockSize)
		for {
			l, err := file.Read(readBuffer)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println(err.Error())
			}
			block := Block{BlockData: readBuffer[:l], BlockSize: int32(l)}
			success := false
			blockHash := GetBlockHashString(block.BlockData)
			err = client.PutBlock(&block, hash2server[blockHash], &success)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	// update metaserver after blockserver to avoid race condition
	// log.Println(localMeta)
	err = client.UpdateFile(localMeta, &newVersion)
	if err != nil {
		return err
	}
	log.Printf("Upload file %v to server", localMeta.Filename)
	localMeta.Version = newVersion
	return nil
}

// Implement the logic for a client syncing with the server here.
func ClientSync(client RPCClient) {
	// 1 load local indexes
	localIndexes, err := LoadMetaFromMetaFile(client.BaseDir)
	if err != nil {
		log.Println(err.Error())
	}
	// 2 update local indexes
	UpdateLocalIndexes(client, &localIndexes)

	// 3 load remote indexes
	remoteIndexes := make(map[string]*FileMetaData)
	err = client.GetFileInfoMap(&remoteIndexes)
	if err != nil {
		log.Println(err.Error())
	}

	// 4 compare local vs remote
	// 4.1 download new server files that client doesn't have
	for fn, remoteMeta := range remoteIndexes {
		if _, ok := localIndexes[fn]; !ok {
			err = DownloadFromServer(client, remoteMeta)
			if err != nil {
				log.Println(err.Error())
			}
			localIndexes[fn] = remoteMeta
		}
	}

	// 4.2 upload new local files that server doesn't have
	for fn, localMeta := range localIndexes {
		if _, ok := remoteIndexes[fn]; !ok {
			err = UploadToServer(client, localMeta)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	// 4.3 conflict resolution, download/upload
	for fn, localMeta := range localIndexes {
		if remoteMeta, ok := remoteIndexes[fn]; ok {
			if localMeta.Version > remoteMeta.Version {
				// local is newer
				// if localMeta.Version-remoteMeta.Version != 1 {
				// 	log.Println("version error")
				// }
				err = UploadToServer(client, localMeta)
				if err != nil {
					log.Println(err.Error())
				}
			} else if localMeta.Version < remoteMeta.Version {
				// server is newer
				err = DownloadFromServer(client, remoteMeta)
				if err != nil {
					log.Println(err.Error())
				}
				localIndexes[fn] = remoteMeta
			} else if !SlicesEqual(localMeta.BlockHashList, remoteMeta.BlockHashList) {
				// conflict, overwrite
				err = DownloadFromServer(client, remoteMeta)
				if err != nil {
					log.Println(err.Error())
				}
				localIndexes[fn] = remoteMeta
			}
		}
	}

	// 5 store updated indexes
	WriteMetaFile(localIndexes, client.BaseDir)
}
