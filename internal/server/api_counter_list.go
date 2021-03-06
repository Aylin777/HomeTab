package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/systemz/hometab/internal/model"
)

type CounterApi struct {
	Id         uint              `json:"id"`
	Name       string            `json:"name"`
	Tags       []string          `json:"tags"`
	Seconds    uint              `json:"seconds"`
	SecondsF   string            `json:"secondsF"`
	InProgress bool              `json:"inProgress"`
	Stats      model.CounterList `json:"stats"`
	Sessions   []CounterLogApi   `json:"sessions"`
}

type CounterLogApi struct {
	Duration          uint      `json:"durationS"`
	DurationFormatted string    `json:"durationSF"`
	Start             time.Time `json:"start"`
	End               time.Time `json:"end"`
}

func ApiCounterList(w http.ResponseWriter, r *http.Request) {
	authDeviceOk, deviceInfo := DeviceApiCheckAuth(w, r)
	authUserOk, userInfo := CheckApiAuth(w, r)
	// deny access if neither auth method works
	if !authUserOk && !authDeviceOk {
		w.Write([]byte{})
		return
	}

	var userId uint
	if authDeviceOk {
		userId = deviceInfo.UserId
	}
	if authUserOk {
		userId = userInfo.Id
	}

	// gather data, convert from DB model to API model
	var counters []CounterApi
	err, dbCounters := model.CountersLatestListAndroid(userId)
	if err != nil {
		logrus.Error(err)
	}
	for _, counter := range dbCounters {
		counters = append(counters, CounterApi{
			Id:         counter.Id,
			Name:       counter.Name,
			Tags:       []string{counter.Tags},
			Seconds:    counter.SecondsAll,
			InProgress: counter.Running == 1,
		})
	}
	logrus.Debugf("%+v", counters)

	// prepare JSON
	counterList, err := json.MarshalIndent(counters, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(counterList)

}

type CounterApiPagination struct {
	Pagination struct {
		AllRecords int `json:"allRecords"`
	} `json:"pagination"`
	Counters []CounterApi `json:"counters"`
}

type PaginateQueryRequest struct {
	Query string `json:"q"`
}

func ApiCounterListPagination(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := userInfo.Id

	// get limitStr on one page
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	// get nextIdStr on one page
	nextIdStr := r.URL.Query().Get("nextId")
	nextId, err := strconv.Atoi(nextIdStr)
	if err != nil || nextId < 1 {
		nextId = 0
	}
	// get prevIdStr on one page
	prevIdStr := r.URL.Query().Get("prevId")
	prevId, err := strconv.Atoi(prevIdStr)
	if err != nil || prevId < 1 {
		prevId = 0
	}

	// get search term
	decoder := json.NewDecoder(r.Body)
	var counterQuery PaginateQueryRequest
	decoder.Decode(&counterQuery)
	searchTerm := counterQuery.Query

	// gather data, convert from DB model to API model
	var rawRes CounterApiPagination
	var counters []CounterApi
	dbCounters, allRecords := model.CountersLongListPaginate(userId, limit, nextId, prevId, searchTerm)
	for _, counter := range dbCounters {
		counters = append(counters, CounterApi{
			Id:         counter.Id,
			Name:       counter.Name,
			Tags:       []string{counter.Tags},
			Seconds:    counter.SecondsAll,
			SecondsF:   counter.SecondsAllFormatted,
			InProgress: counter.Running == 1,
		})
	}
	// prevent null result in JSON, make empty array instead
	//rawRes.Counters = make([]CounterApi, 0)
	rawRes.Counters = counters
	rawRes.Pagination.AllRecords = allRecords

	if len(rawRes.Counters) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// prepare JSON result
	counterList, err := json.MarshalIndent(rawRes, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(counterList)
}
