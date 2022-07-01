package sourcemod

import (
	"io/ioutil"
	"log"
	"net/http"
)

const (
	BaseSMDropURL = "https://sm.alliedmods.net/smdrop/"

	LatestWindowsSMVersion = "sourcemod-latest-windows"
)

// Gets latest major-minor-patch version of SourceMod based on major-minor string passed in.
func GetLatestSourceModVersion(latestSmVersion string) (string, error) {
	res, err := http.Get(BaseSMDropURL + latestSmVersion + "/" + LatestWindowsSMVersion)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	res.Body.Close()
	return string(body), nil
}
