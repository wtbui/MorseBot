package data

import (
	"os"
	"strings"
	"fmt"
	"io/ioutil"
)
type GRegistration struct {
	User string
	GKey string
}
 
func RegisterGKey(regStr string) error {
	regList, err := ParseRegistrations(regStr)
	if err != nil {
		return err
	}

	err = WriteRegistrations(regList)
	if err != nil {
		return err
	}

	return nil	
}

func ParseRegistrations(regStr string) (map[string]GRegistration, error) {
	regStrList := strings.Split(strings.ReplaceAll(regStr, " ", ""), ",")
	regList := map[string]GRegistration{}

	for _, reg := range regStrList {
		gReg := GRegistration{}

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
	}

	return regList, nil
}

func WriteRegistrations(regList map[string]GRegistration) error {
	dbFile := os.Getenv("GOVEEDB")
	file, err := os.OpenFile(dbFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
  
	oldRegList, err := RetrieveCurrentGDB()
	if err != nil {
		return err
	}

	for user, reg := range oldRegList {
		_, exists := regList[user]
		
		if !exists {
			regList[user] = GRegistration{reg.User, reg.GKey} 
		}
	}

	for user, reg := range regList {
		_, exists := oldRegList[user]

		if !exists {
			fmt.Println("Registered user \"" + user + "\" with govee key \"" + reg.GKey + "\"")
		} else {
			fmt.Println("Updated user \"" + user + "\" with govee key \"" + reg.GKey + "\"")
		}

		_, err = file.WriteString(user + ":" + reg.GKey + ",")
		if err != nil {
			return err
		}
	}
	file.Close()

	return nil
}

func RetrieveCurrentGDB() (map[string]GRegistration, error) {
	dbFile := os.Getenv("GOVEEDB")
	content, err := ioutil.ReadFile(dbFile)
	if err != nil {
		return nil, err
	}
	
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
