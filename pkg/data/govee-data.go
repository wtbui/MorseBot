package data

import (
	"os"
	"strings"
	"io/ioutil"
	"go.uber.org/zap"
	"strconv"
)
type GRegistration struct {
	User string
	GKey string
}

func RegisterUser(regStr string) error {
	regList, err := ParseRegistrations(regStr)
	if err != nil {
		return err
	}

	err = UpdateRegistrations(regList)
	if err != nil {
		return err
	}

	return nil	
}

func ParseRegistrations(regStr string) (map[string]*GRegistration, error) {
	regStrList := strings.Split(strings.ReplaceAll(regStr, " ", ""), ",")
	regList := map[string]*GRegistration{}

	for _, reg := range regStrList {
		zap.S().Debug(reg)
		gRegOrig := GRegistration{}
		gReg := &gRegOrig

		singleReg := strings.Split(reg, ":")

		// Error Checking
		if len(singleReg) < 2 {
			continue
		}

		if len(singleReg[0]) == 0 {
			continue
		}

		if len(singleReg[1]) == 0 {
			continue
		}

		gReg.User = singleReg[0]
		gReg.GKey = singleReg[1]

		regList[gReg.User] = gReg
		zap.S().Debug(len(regList))
		zap.S().Debug(gReg.User)
	}

	zap.S().Debug("Found " + strconv.Itoa(len(regList)) + " users in GoveeDB")

	return regList, nil
}

func UpdateRegistrations(regList map[string]*GRegistration) error { 
	oldRegList, err := RetrieveCurrentGDB()
	if err != nil {
		return err
	}

	for user, reg := range oldRegList {
		_, exists := regList[user]
		
		if !exists {
			regList[user] = &GRegistration{reg.User, reg.GKey} 
		} else {
			zap.S().Info("Updated user \"" + user + "\" with govee key \"" + reg.GKey + "\"")
		}
	}

	err = writeRegistrations(regList)

	return err
}

func writeRegistrations(regList map[string]*GRegistration) error {
	dbFile := os.Getenv("GOVEEDB")
	file, err := os.OpenFile(dbFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	for user, reg := range regList {
		_, err = file.WriteString(user + ":" + reg.GKey + ",")
		if err != nil {
			return err
		}
	}
	file.Close()

	return nil
}

func RetrieveCurrentGDB() (map[string]*GRegistration, error) {
	dbFile := os.Getenv("GOVEEDB")
	content, err := ioutil.ReadFile(dbFile)
	if err != nil {
		return nil, err
	}

	zap.S().Debug("Raw Govee Database: " + string(content))
	return ParseRegistrations(string(content))
} 

func UserExist(username string) (bool, error) {
	currentGDB, err := RetrieveCurrentGDB()
	if err != nil {
		return false, err
	}

	_, exists := currentGDB[username]

	return exists, nil
}

func DeleteUser(user string) error {
	currentGDB, err := RetrieveCurrentGDB()
	if err != nil {
		return err
	}

	_, exists := currentGDB[user]
	if exists {
		zap.S().Info("Deleted user \"" + user + "\"")
		delete(currentGDB, user)
	} else {
		zap.S().Info("User \"" + user + "\" not found")
	}

	for user, _ := range currentGDB {
		zap.S().Info(user)
	}
	
	err = writeRegistrations(currentGDB)
	if err != nil {
		return err
	}

	return nil
}
