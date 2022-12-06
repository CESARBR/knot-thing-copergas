package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/CESARBR/knot-thing-copergas/pkg/logging"
)

// Represents the token received from Copergas.
var tokenObject *entities.Token

func convertStringToTime(datetimeInString string) (time.Time, error) {
	return time.Parse(time.RFC3339, datetimeInString)
}

/*
Subtracts three hours from the expiry date contained in the received token, in order to adjust the time zone.
Subtracts an additional 50 minutes to ensure the execution of subsequent requests.
*/
func correctTimeDasedOnTimeZone(originalTime time.Time) time.Time {
	return originalTime.Add(-3 * time.Hour).Add(-50 * time.Minute)
}

func isTokenExperied(expirationDatetime time.Time) bool {
	nowInUTFC := time.Now().UTC()
	if nowInUTFC.After(expirationDatetime) {
		return true
	} else {
		return false
	}
}

func requestToken(client *http.Client, setup entities.CopergasConfig, logger logging.Logger) error {
	url := setup.Endpoints.AuthToken
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	header := MakesTokenRequestHeader(setup)
	request.Header = header

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer closeFile(response, logger)
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &tokenObject)
	if err != nil {
		return err
	}

	return nil
}

func closeFile(response *http.Response, log logging.Logger) {

	err := response.Body.Close()
	if err != nil {
		log.Info("error: ", err)
		os.Exit(1)
	}
}

/*
Wait for the signal that a new token is needed.
When this signal is received, a new token is requested from the Copergas service.
Finally, this token is distributed to the routines that request it.
*/
func TokenHandler(client *http.Client, applicationSetup entities.CopergasConfig, newTokenIsNeeded chan struct{}, obtainedTokenSync chan string, mutex *sync.Mutex, logger logging.Logger) {
	for {
		// Blocks the loop until a signal is received.
		<-newTokenIsNeeded
		mutex.Lock()
		if tokenObject == nil {
			tokenObject = &entities.Token{}
			tokenObject.ExpiresUtc = "2021-01-01T12:27:19-03:00"

		}
		expiresInUTC, err := convertStringToTime(tokenObject.ExpiresUtc)
		if err != nil {
			logger.Infof(err.Error())
		} else {

			expiresInUTCCorrected := correctTimeDasedOnTimeZone(expiresInUTC)

			if isTokenExperied(expiresInUTCCorrected) {
				err := requestToken(client, applicationSetup, logger)
				if err != nil {
					logger.Infof(err.Error())
				}
			}
		}
		obtainedTokenSync <- tokenObject.AccessToken
		mutex.Unlock()

	}

}
