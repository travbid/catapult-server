package lib

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return db, err
	}
	err = db.Ping()
	return db, err
}

func createDatabase(fpath string, records []packageRecord) error {
	db, err := sql.Open("sqlite3", fpath)
	if err != nil {
		return err
	}
	create := `CREATE TABLE "Packages" (
		"hash"	TEXT NOT NULL UNIQUE,
		"pkg_name"	TEXT NOT NULL,
		"version"	TEXT NOT NULL,
		"manifest"	TEXT NOT NULL,
		"recipe"	TEXT NOT NULL,
		"user"	TEXT NOT NULL,
		"channel"	TEXT NOT NULL,
		"datetime_added"	INTEGER NOT NULL CHECK(TYPEOF("datetime_added") = 'integer'),
		PRIMARY KEY("hash")
)`
	_, err = db.Exec(create)
	if err != nil {
		return err
	}

	insert := `INSERT INTO Packages
	VALUES (?, ?, ?, ?, ?, ?, ?, ?); `
	for _, record := range records {
		hash := computeHash(record)
		_, err = db.Exec(insert, hash, record.PkgName, record.Version, record.Manifest, record.Recipe, record.User, record.Channel, record.DatetimeAddedMs)
		if err != nil {
			return err
		}
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

func computeHash(record packageRecord) string {
	pkgNameHash := sha256.Sum256([]byte(record.PkgName))
	versionHash := sha256.Sum256([]byte(record.Version))
	manifestHash := sha256.Sum256([]byte(record.Manifest))
	recipeHash := sha256.Sum256([]byte(record.Recipe))
	userHash := sha256.Sum256([]byte(record.User))
	channelHash := sha256.Sum256([]byte(record.Channel))

	combinedHash := append(pkgNameHash[:], versionHash[:]...)
	combinedHash = append(combinedHash, manifestHash[:]...)
	combinedHash = append(combinedHash, recipeHash[:]...)
	combinedHash = append(combinedHash, userHash[:]...)
	combinedHash = append(combinedHash, channelHash[:]...)
	totalHash := sha256.Sum256(combinedHash)

	return fmt.Sprintf("%x", totalHash)
}
