package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GenerateDatabase(dbConfig string) (string, error) {
	cfgDir := filepath.Dir(dbConfig)
	cfgFile, err := os.Open(dbConfig)
	if err != nil {
		return "", err
	}
	cfgBytes, err := ioutil.ReadAll(cfgFile)
	if err != nil {
		return "", err
	}

	var entries []packageRecord
	err = json.Unmarshal(cfgBytes, &entries)
	if err != nil {
		return "", err
	}
	for i, entry := range entries {
		log.Println(entry)
		entries[i].Manifest = filepathToString(filepath.Join(cfgDir, entry.Manifest))
		entries[i].Recipe = filepathToString(filepath.Join(cfgDir, entry.Recipe))
	}
	for _, entry := range entries {
		log.Println(entry)
	}

	cfgFilename := filepath.Base(dbConfig)
	dbFilename := strings.TrimSuffix(cfgFilename, filepath.Ext(cfgFilename))
	dbPath := filepath.Join(cfgDir, dbFilename) + ".sqlite"
	err = createDatabase(dbPath, entries)
	return dbPath, err
}

func filepathToString(fpath string) string {
	f, err := os.Open(fpath)
	if err != nil {
		log.Fatalln(err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	return string(bytes)
}
