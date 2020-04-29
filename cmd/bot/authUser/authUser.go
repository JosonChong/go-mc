package authUser

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Tnze/go-mc/yggdrasil"
)

// var Config = struct {
// 	APPName string `default:"app name"`

// 	Account struct {
// 		Username string `required:"true"`
// 		Password string `required:"true"`
// 	}
// }{}

type UserAuth struct {
	Username string           `json:"username"`
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Tokens   yggdrasil.Tokens `json:"tokens"`
}

// func main() {
// 	err := configor.Load(&Config, "config.json")

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	GetToken(Config.Account.Username, Config.Account.Password)
// }

// Get token with username and password
func GetToken(username, password string) *UserAuth {
	isValid, userAuth := verifyLocalToken(username)
	if !isValid {
		var err error
		userAuth, err = authenticate(username, password)
		if err != nil {
			fmt.Println(err)
		}
	}

	return userAuth
}

func authenticate(username string, password string) (*UserAuth, error) {
	resp, err := yggdrasil.Authenticate(username, password)
	if err != nil {
		// fmt.Println("Failed to authenticate", err)
		return nil, err
	}

	fmt.Println("Got a new token for ", username)

	var userAuth UserAuth
	userAuth.Username = username
	userAuth.ID, userAuth.Name = resp.SelectedProfile()
	userAuth.Tokens = resp.GetTokens()

	file, _ := json.MarshalIndent(userAuth, "", " ")
	ioutil.WriteFile("userAuth.json", file, 0644)

	return &userAuth, nil
}

func verifyLocalToken(username string) (bool, *UserAuth) {
	var userAuth UserAuth

	file, err := ioutil.ReadFile("userAuth.json")

	err = json.Unmarshal([]byte(file), &userAuth)

	if userAuth.Tokens.AccessToken != "" && err == nil {
		if userAuth.Username == username {
			isValid, err := userAuth.Tokens.Validate()
			if err == nil && isValid == true {
				fmt.Println("Saved token is valid!")
				return true, &userAuth
			}
		}

		// remove the file if the saved user auth is not valid
		os.Remove("userAuth.json")
	}

	return false, nil
}
