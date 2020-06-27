package db

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/boltdb/bolt"

	"github.com/jempe/panki/validator"
)

type Request struct {
	ID         string    `json:"id"`
	FromDevice string    `json:"fromdevice"`
	Data       string    `json:"data"`
	IP         string    `json:"ip"`
	Created    time.Time `json:"created"`
}

type Requests []Request

func (boltdb *DB) GetRequest(requestID string) (request Request, err error) {
	validID, err := validator.UUID(requestID)
	if !validID {
		return request, err
	}

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("requests"))
		v := b.Get([]byte(requestID))

		if v == nil {
			return errors.New("request not found")
		}

		err := json.Unmarshal(v, &request)

		return err
	})

	return request, err
}

func (boltdb *DB) InsertRequest(request Request, fields []string) (requestID string, err error) {

	validationErrorPrefix := "insert_request_error:"

	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(validationErrorPrefix, err)
		return
	}

	var requestData Request

	if request.ID == "" {
		requestID = id.String()

		request.ID = requestID
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := request.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		requestData.ID = request.ID

	}
	if emptyOrContains(fields, "FromDevice") {
		validFromDevice, validFromDeviceErr := request.ValidFromDeviceDefault()
		if !validFromDevice {
			err = validFromDeviceErr
			return
		}

		requestData.FromDevice = request.FromDevice

	}
	if emptyOrContains(fields, "Data") {
		validData, validDataErr := request.ValidDataDefault()
		if !validData {
			err = validDataErr
			return
		}

		requestData.Data = request.Data

	}
	if emptyOrContains(fields, "IP") {
		validIP, validIPErr := request.ValidIPDefault()
		if !validIP {
			err = validIPErr
			return
		}

		requestData.IP = request.IP

	}

	existRequestData, _ := boltdb.GetRequest(request.ID)
	if existRequestData.ID != "" {
		err = errors.New(validationErrorPrefix + " request with ID " + request.ID + " already exists")
		return
	}
	requestData.Created = time.Now().UTC()

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("requests"))

		requestJson, err := json.Marshal(requestData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(request.ID), requestJson)
		return err
	})

	return
}

func (boltdb *DB) DeleteRequest(requestID string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "delete_request_error:"

	validID, err := validator.UUID(requestID)
	if !validID {
		return
	}

	requestData, err := boltdb.GetRequest(requestID)
	if err != nil {
		return
	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("requests"))
		err = b.Delete([]byte(requestData.ID))

		if err == nil {
			rowsAffected = 1
		}
		return err
	})

	if err == nil {
		rowsAffected = 1
	}

	return
}

func (boltdb *DB) UpdateRequest(request Request, fields []string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "update_request_error:"

	validID, err := validator.UUID(request.ID)
	if !validID {
		return
	}

	requestData, err := boltdb.GetRequest(request.ID)
	if err != nil {
		return
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := request.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		requestData.ID = request.ID

	}
	if emptyOrContains(fields, "FromDevice") {
		validFromDevice, validFromDeviceErr := request.ValidFromDeviceDefault()
		if !validFromDevice {
			err = validFromDeviceErr
			return
		}

		requestData.FromDevice = request.FromDevice

	}
	if emptyOrContains(fields, "Data") {
		validData, validDataErr := request.ValidDataDefault()
		if !validData {
			err = validDataErr
			return
		}

		requestData.Data = request.Data

	}
	if emptyOrContains(fields, "IP") {
		validIP, validIPErr := request.ValidIPDefault()
		if !validIP {
			err = validIPErr
			return
		}

		requestData.IP = request.IP

	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("requests"))

		requestJson, err := json.Marshal(requestData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(requestData.ID), requestJson)

		if err == nil {
			rowsAffected = 1
		}

		return err
	})

	return
}

func (boltdb *DB) GetRequestList(offset int, limit int, filters Filters, returnFields []string, sortBy SortBy) (results []Request, totalResults int64, err error) {
	validationErrorPrefix := "get_user_error:"

	if !(filters.Operator == "AND" || filters.Operator == "OR") {
		err = errors.New(validationErrorPrefix + " filter operator error")
	}

	var requestList Requests

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("requests"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var request Request
			err := json.Unmarshal(v, &request)

			includeThis, err := includeThisRequest(filters, request)

			if err != nil {
				return err
			}

			if includeThis {
				resultRequest := Request{ID: request.ID}
				if emptyOrContains(returnFields, "ID") {
					resultRequest.ID = request.ID
				}
				if emptyOrContains(returnFields, "FromDevice") {
					resultRequest.FromDevice = request.FromDevice
				}
				if emptyOrContains(returnFields, "Data") {
					resultRequest.Data = request.Data
				}
				if emptyOrContains(returnFields, "IP") {
					resultRequest.IP = request.IP
				}
				if emptyOrContains(returnFields, "Created") {
					resultRequest.Created = request.Created
				}

				requestList = append(requestList, resultRequest)
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	if sortBy.Direction == "ASC" || sortBy.Direction == "DESC" {
		if sortBy.Field == "ID" && sortBy.Direction == "ASC" {
			sort.Sort(sortByRequestID{requestList})
		} else if sortBy.Field == "ID" && sortBy.Direction == "DESC" {
			sort.Sort(sortByRequestIDDesc{requestList})
		}
		if sortBy.Field == "FromDevice" && sortBy.Direction == "ASC" {
			sort.Sort(sortByRequestFromDevice{requestList})
		} else if sortBy.Field == "FromDevice" && sortBy.Direction == "DESC" {
			sort.Sort(sortByRequestFromDeviceDesc{requestList})
		}
		if sortBy.Field == "Data" && sortBy.Direction == "ASC" {
			sort.Sort(sortByRequestData{requestList})
		} else if sortBy.Field == "Data" && sortBy.Direction == "DESC" {
			sort.Sort(sortByRequestDataDesc{requestList})
		}
		if sortBy.Field == "IP" && sortBy.Direction == "ASC" {
			sort.Sort(sortByRequestIP{requestList})
		} else if sortBy.Field == "IP" && sortBy.Direction == "DESC" {
			sort.Sort(sortByRequestIPDesc{requestList})
		}
		if sortBy.Field == "Created" && sortBy.Direction == "ASC" {
			sort.Sort(sortByRequestCreated{requestList})
		} else if sortBy.Field == "Created" && sortBy.Direction == "DESC" {
			sort.Sort(sortByRequestCreatedDesc{requestList})
		}

	} else {
		err = errors.New(validationErrorPrefix + " sort Direction error")
	}

	totalResults = int64(len(requestList))

	for indexRequest, resultRequest := range requestList {
		if indexRequest >= offset && indexRequest < (offset+limit) {
			results = append(results, resultRequest)
		}
	}

	return
}
func includeThisRequest(filters Filters, request Request) (include bool, err error) {
	validationErrorPrefix := "get_request_error:"

	if len(filters.Conditions) == 0 {
		return true, nil
	}

	if filters.Operator == "AND" {
		include = true
	}

	for _, condition := range filters.Conditions {
		if !(condition.Comparison == "LIKE" || condition.Comparison == "=" || condition.Comparison == ">" || condition.Comparison == "<") {
			err = errors.New(validationErrorPrefix + " condition operator error")
			return false, err
		}

		meetConditionID := false

		if condition.Field == "ID" {
			conditionValueID := condition.Value.(string)

			if condition.Comparison == "=" && request.ID == conditionValueID {
				meetConditionID = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueID, "%") && strings.HasSuffix(conditionValueID, "%") {
					if strings.Contains(request.ID, strings.TrimSuffix(strings.TrimPrefix(conditionValueID, "%"), "%")) {
						meetConditionID = true
					}
				} else if strings.HasPrefix(conditionValueID, "%") {
					if strings.HasSuffix(request.ID, strings.TrimPrefix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if strings.HasSuffix(conditionValueID, "%") {
					if strings.HasPrefix(request.ID, strings.TrimSuffix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if request.ID == conditionValueID {
					meetConditionID = true
				}
			}

			if meetConditionID {
				if filters.Operator == "OR" {
					include = true
					return include, err
				}
			} else {
				if filters.Operator == "AND" {
					include = false
					return include, err
				}
			}
		}

		meetConditionFromDevice := false

		if condition.Field == "FromDevice" {
			conditionValueFromDevice := condition.Value.(string)

			if condition.Comparison == "=" && request.FromDevice == conditionValueFromDevice {
				meetConditionFromDevice = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueFromDevice, "%") && strings.HasSuffix(conditionValueFromDevice, "%") {
					if strings.Contains(request.FromDevice, strings.TrimSuffix(strings.TrimPrefix(conditionValueFromDevice, "%"), "%")) {
						meetConditionFromDevice = true
					}
				} else if strings.HasPrefix(conditionValueFromDevice, "%") {
					if strings.HasSuffix(request.FromDevice, strings.TrimPrefix(conditionValueFromDevice, "%")) {
						meetConditionFromDevice = true
					}
				} else if strings.HasSuffix(conditionValueFromDevice, "%") {
					if strings.HasPrefix(request.FromDevice, strings.TrimSuffix(conditionValueFromDevice, "%")) {
						meetConditionFromDevice = true
					}
				} else if request.FromDevice == conditionValueFromDevice {
					meetConditionFromDevice = true
				}
			}

			if meetConditionFromDevice {
				if filters.Operator == "OR" {
					include = true
					return include, err
				}
			} else {
				if filters.Operator == "AND" {
					include = false
					return include, err
				}
			}
		}

		meetConditionData := false

		if condition.Field == "Data" {
			conditionValueData := condition.Value.(string)

			if condition.Comparison == "=" && request.Data == conditionValueData {
				meetConditionData = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueData, "%") && strings.HasSuffix(conditionValueData, "%") {
					if strings.Contains(request.Data, strings.TrimSuffix(strings.TrimPrefix(conditionValueData, "%"), "%")) {
						meetConditionData = true
					}
				} else if strings.HasPrefix(conditionValueData, "%") {
					if strings.HasSuffix(request.Data, strings.TrimPrefix(conditionValueData, "%")) {
						meetConditionData = true
					}
				} else if strings.HasSuffix(conditionValueData, "%") {
					if strings.HasPrefix(request.Data, strings.TrimSuffix(conditionValueData, "%")) {
						meetConditionData = true
					}
				} else if request.Data == conditionValueData {
					meetConditionData = true
				}
			}

			if meetConditionData {
				if filters.Operator == "OR" {
					include = true
					return include, err
				}
			} else {
				if filters.Operator == "AND" {
					include = false
					return include, err
				}
			}
		}

		meetConditionIP := false

		if condition.Field == "IP" {
			conditionValueIP := condition.Value.(string)

			if condition.Comparison == "=" && request.IP == conditionValueIP {
				meetConditionIP = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueIP, "%") && strings.HasSuffix(conditionValueIP, "%") {
					if strings.Contains(request.IP, strings.TrimSuffix(strings.TrimPrefix(conditionValueIP, "%"), "%")) {
						meetConditionIP = true
					}
				} else if strings.HasPrefix(conditionValueIP, "%") {
					if strings.HasSuffix(request.IP, strings.TrimPrefix(conditionValueIP, "%")) {
						meetConditionIP = true
					}
				} else if strings.HasSuffix(conditionValueIP, "%") {
					if strings.HasPrefix(request.IP, strings.TrimSuffix(conditionValueIP, "%")) {
						meetConditionIP = true
					}
				} else if request.IP == conditionValueIP {
					meetConditionIP = true
				}
			}

			if meetConditionIP {
				if filters.Operator == "OR" {
					include = true
					return include, err
				}
			} else {
				if filters.Operator == "AND" {
					include = false
					return include, err
				}
			}
		}

		meetConditionCreated := false

		if condition.Field == "Created" {
			conditionValueCreated := condition.Value.(time.Time)
			diffCreated := request.Created.Sub(conditionValueCreated)

			if condition.Comparison == "=" && request.Created == conditionValueCreated {
				meetConditionCreated = true
			} else if condition.Comparison == ">" && diffCreated > 0 {
				meetConditionCreated = true
			} else if condition.Comparison == "<" && diffCreated < 0 {
				meetConditionCreated = true
			}

			if meetConditionCreated {
				if filters.Operator == "OR" {
					include = true
					return include, err
				}
			} else {
				if filters.Operator == "AND" {
					include = false
					return include, err
				}
			}
		}
	}

	return include, err
}

func (s Requests) Len() int {
	return len(s)
}
func (s Requests) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortByRequestID struct {
	Requests
}

func (s sortByRequestID) Less(i, j int) bool {
	return s.Requests[i].ID < s.Requests[j].ID
}

type sortByRequestIDDesc struct {
	Requests
}

func (s sortByRequestIDDesc) Less(i, j int) bool {
	return s.Requests[i].ID > s.Requests[j].ID

}

type sortByRequestFromDevice struct {
	Requests
}

func (s sortByRequestFromDevice) Less(i, j int) bool {
	return s.Requests[i].FromDevice < s.Requests[j].FromDevice
}

type sortByRequestFromDeviceDesc struct {
	Requests
}

func (s sortByRequestFromDeviceDesc) Less(i, j int) bool {
	return s.Requests[i].FromDevice > s.Requests[j].FromDevice

}

type sortByRequestData struct {
	Requests
}

func (s sortByRequestData) Less(i, j int) bool {
	return s.Requests[i].Data < s.Requests[j].Data
}

type sortByRequestDataDesc struct {
	Requests
}

func (s sortByRequestDataDesc) Less(i, j int) bool {
	return s.Requests[i].Data > s.Requests[j].Data

}

type sortByRequestIP struct {
	Requests
}

func (s sortByRequestIP) Less(i, j int) bool {
	return s.Requests[i].IP < s.Requests[j].IP
}

type sortByRequestIPDesc struct {
	Requests
}

func (s sortByRequestIPDesc) Less(i, j int) bool {
	return s.Requests[i].IP > s.Requests[j].IP

}

type sortByRequestCreated struct {
	Requests
}

func (s sortByRequestCreated) Less(i, j int) bool {
	diffLastModification := s.Requests[i].Created.Sub(s.Requests[j].Created)
	return diffLastModification < 0
}

type sortByRequestCreatedDesc struct {
	Requests
}

func (s sortByRequestCreatedDesc) Less(i, j int) bool {
	diffLastModification := s.Requests[i].Created.Sub(s.Requests[j].Created)
	return diffLastModification > 0

}

func (request Request) ValidIDDefault() (validField bool, err error) {
	validField, _ = validator.UUID(request.ID)
	if !validField {
		err = errors.New("error_uuid__request___ID")
		return
	}

	return
}
func (request Request) ValidFromDeviceDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(request.FromDevice, 200)
	if !validField {
		err = errors.New("error_maxlength__request___FromDevice")
		return
	}

	return
}
func (request Request) ValidDataDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(request.Data, 2083)
	if !validField {
		err = errors.New("error_maxlength__request___Data")
		return
	}

	return
}
func (request Request) ValidIPDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(request.IP, 200)
	if !validField {
		err = errors.New("error_maxlength__request___IP")
		return
	}

	return
}
func (request Request) ValidCreatedDefault() (validField bool, err error) {
	validField = true

	return
}
