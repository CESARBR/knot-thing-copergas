package collector

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/CESARBR/knot-thing-copergas/internal/integration/knot/entities"
	"github.com/CESARBR/knot-thing-copergas/pkg/logging"
)

// Represents the access token.
type bearerToken struct {
	accessToken string
}

func (bt *bearerToken) setAccessToken(accessToken string) {
	bt.accessToken = accessToken
}
func (bt *bearerToken) getAccessToken() string {
	return bt.accessToken
}

var bToken *bearerToken = &bearerToken{}

// Establish the interface of the mechanisms for validating the HTTP response received from the copergas system.
type connectionValidation interface {
	execute()
	getResponse() *http.Response
	getTokenSync() chan struct{}
	getObtainedTokenSync() chan string
	getCodVar() int
	getCommunicationChannel() chan entities.ReceivedData
}

// Represents the response to the request for instantaneous measurements.
type instantaneousVariableResponse struct {
	response                         *http.Response
	communicationChannel             chan entities.ReceivedData
	tokenSync                        chan struct{}
	intervalBetweenRequestsInMinutes float32
	obtainedTokenSync                chan string
	codVar                           int
}

// Captures instantaneous measurements for a specific variable, identified in the Copergás system by an integer code.
func GetInstantaneousMeasurement(
	variableCode int,
	intervalBetweenRequestsInMinutes float32,
	client *http.Client,
	data chan entities.ReceivedData,
	applicationSetup entities.CopergasConfig,
	getTokenSync chan struct{},
	obtainedTokenSync chan string,
	logger logging.Logger,
) {
	// Creates a request for the endpoint of interest.
	instantaneousMeasurementUrl := fmt.Sprintf("%s/%d", applicationSetup.Endpoints.Variable, variableCode)
	request, err := http.NewRequest(http.MethodGet, instantaneousMeasurementUrl, nil)
	if err != nil {
		logger.Errorf(err.Error())
	}

	for {

		// Creates the HTTP request header and specifies the authentication token.
		header := http.Header{
			"Accept":        []string{"application/json"},
			"Authorization": []string{bToken.getAccessToken()},
		}
		request.Header = header

		// Executes the request.
		response, err := client.Do(request)

		if err != nil {
			httpRequestError := errors.New("HTTP request error")
			newReceivedData := entities.ReceivedData{}
			newReceivedData.CodVar = variableCode
			newReceivedData.Error = httpRequestError
			data <- newReceivedData
			continue
		}
		// Information used in validating the REST service response.
		apiResponse := &instantaneousVariableResponse{
			response:                         response,
			communicationChannel:             data,
			tokenSync:                        getTokenSync,
			intervalBetweenRequestsInMinutes: intervalBetweenRequestsInMinutes,
			obtainedTokenSync:                obtainedTokenSync,
			codVar:                           variableCode,
		}

		// Creates the decorator that validates the received response.
		connectionValidationDecorator := connectionDecorator{
			connection: apiResponse,
		}

		connectionValidationDecorator.execute()
	}

}

func (ivr *instantaneousVariableResponse) getResponse() *http.Response {
	return ivr.response
}

func (ivr *instantaneousVariableResponse) getTokenSync() chan struct{} {
	return ivr.tokenSync
}
func (ivr *instantaneousVariableResponse) getObtainedTokenSync() chan string {
	return ivr.obtainedTokenSync
}
func (ivr *instantaneousVariableResponse) getCodVar() int {
	return ivr.codVar
}
func (ivr *instantaneousVariableResponse) getCommunicationChannel() chan entities.ReceivedData {
	return ivr.communicationChannel
}

type connectionDecorator struct {
	connection connectionValidation
}

// Validates the HTTP code of the response received.
func (cd *connectionDecorator) execute() {
	var statusCode int = cd.connection.getResponse().StatusCode
	switch statusCode {
	case http.StatusOK:
		cd.connection.execute()
	// Authentication error. A new token will be requested..
	case http.StatusUnauthorized:
		getTokenSyncChan := cd.connection.getTokenSync()
		// Signals a routine that a new token should be requested.
		getTokenSyncChan <- struct{}{}

		// It waits until a new token is obtained and registers the updated token so that it can be used in the next requests.
		getObtainedTokenSyncChan := cd.connection.getObtainedTokenSync()
		newToken := <-getObtainedTokenSyncChan
		bearerToken := fmt.Sprintf("%s %s", "Bearer", newToken)
		bToken.setAccessToken(bearerToken)

	default:
		httpResponseError := fmt.Errorf("%s: %d", "HTTP connection error with status code", statusCode)
		newReceivedData := entities.ReceivedData{}
		newReceivedData.CodVar = cd.connection.getCodVar()
		newReceivedData.Error = httpResponseError
		commucationChannel := cd.connection.getCommunicationChannel()
		commucationChannel <- newReceivedData
	}
}

// Processes the received data.
func (ivr *instantaneousVariableResponse) execute() {
	newReceivedData := entities.ReceivedData{}
	newReceivedData.CodVar = ivr.getCodVar()
	newReceivedData.Error = nil

	communicationChannel := ivr.getCommunicationChannel()
	responseBody := ivr.getResponse().Body
	defer responseBody.Close()
	bodyBytes, err := ioutil.ReadAll(responseBody)

	if err != nil {
		newReceivedData.Error = err
		communicationChannel <- newReceivedData
		return
	} else {
		var variableObject entities.Variable
		unmarshalErr := json.Unmarshal(bodyBytes, &variableObject)
		if err != nil {
			newReceivedData.Error = unmarshalErr
			communicationChannel <- newReceivedData
			return
		} else {
			newReceivedData.Data = variableObject
			communicationChannel <- newReceivedData
			// waits a specified number of minutes before executing the next request.
			Wait(ivr.intervalBetweenRequestsInMinutes)

		}

	}

}
