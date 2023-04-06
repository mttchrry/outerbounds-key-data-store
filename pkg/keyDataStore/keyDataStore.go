package keyDataStore

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	//"github.com/mttchrry/outerbounds-key-data-store/pkg/keyDataStore/model"
)

type KeyDataStore struct {
	Data map[uuid.UUID]string
}

func New(inputFile string) (*KeyDataStore, error){
	readFile, err := os.Open(inputFile)
	defer readFile.Close()

	if err != nil {
		fmt.Printf("error opening file %v, %v", inputFile, err)
        return nil, err
    }
    fileScanner := bufio.NewScanner(readFile)
 
    fileScanner.Split(bufio.ScanLines)
	kds := &KeyDataStore{
		Data: map[uuid.UUID]string{},
	}
    for fileScanner.Scan() {
        //fmt.Println(fileScanner.Text())
		branchArray := strings.SplitN(fileScanner.Text(), " ", 2)
		if len(branchArray) != 2 {
			fmt.Printf("error opening file %v, %v", inputFile, err)
        	return nil, fmt.Errorf("key-value error in at %v", fileScanner.Text())
		}
		key, err := uuid.Parse(branchArray[0])
		if err != nil {
			fmt.Printf("error with key format %v", branchArray[0])
			return nil, err
		}
		kds.Data[key] = branchArray[1]
    }
	return kds, nil
}

func (kds *KeyDataStore) Get(key string) (string, error) {
	id, err := uuid.Parse(key)
	if err != nil {
		err = fmt.Errorf("invalid key err: %v", err)
		return "", err
	}

	val, ok := kds.Data[id]
	if !ok {
		return "", fmt.Errorf("key doesn't exist: %v", key)
	}

	return val, nil
}	