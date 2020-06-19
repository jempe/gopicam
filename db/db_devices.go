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

type Device struct {
	ID      string    `json:"id"`
	Key     string    `json:"key"`
	Name    string    `json:"name"`
	Secret  string    `json:"secret"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Devices []Device

func (boltdb *DB) GetDevice(deviceID string) (device Device, err error) {
	validID, err := validator.UUID(deviceID)
	if !validID {
		return device, err
	}

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))
		v := b.Get([]byte(deviceID))

		if v == nil {
			return errors.New("device not found")
		}

		err := json.Unmarshal(v, &device)

		return err
	})

	return device, err
}

func (boltdb *DB) InsertDevice(device Device, fields []string) (deviceID string, err error) {

	validationErrorPrefix := "insert_device_error:"

	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(validationErrorPrefix, err)
		return
	}

	var deviceData Device

	if device.ID == "" {
		deviceID = id.String()

		device.ID = deviceID
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := device.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		deviceData.ID = device.ID

	}
	if emptyOrContains(fields, "Key") {
		validKey, validKeyErr := device.ValidKeyDefault()
		if !validKey {
			err = validKeyErr
			return
		}

		deviceData.Key = device.Key

	}
	if emptyOrContains(fields, "Name") {
		validName, validNameErr := device.ValidNameDefault()
		if !validName {
			err = validNameErr
			return
		}

		deviceData.Name = device.Name

	}
	if emptyOrContains(fields, "Secret") {
		validSecret, validSecretErr := device.ValidSecretDefault()
		if !validSecret {
			err = validSecretErr
			return
		}

		deviceData.Secret = device.Secret

	}

	existDeviceData, _ := boltdb.GetDevice(device.ID)
	if existDeviceData.ID != "" {
		err = errors.New(validationErrorPrefix + " device with ID " + device.ID + " already exists")
		return
	}
	deviceData.Created = time.Now().UTC()
	deviceData.Updated = time.Now().UTC()

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))

		deviceJson, err := json.Marshal(deviceData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(device.ID), deviceJson)
		return err
	})

	return
}

func (boltdb *DB) DeleteDevice(deviceID string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "delete_device_error:"

	validID, err := validator.UUID(deviceID)
	if !validID {
		return
	}

	deviceData, err := boltdb.GetDevice(deviceID)
	if err != nil {
		return
	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))
		err = b.Delete([]byte(deviceData.ID))

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

func (boltdb *DB) UpdateDevice(device Device, fields []string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "update_device_error:"

	validID, err := validator.UUID(device.ID)
	if !validID {
		return
	}

	deviceData, err := boltdb.GetDevice(device.ID)
	if err != nil {
		return
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := device.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		deviceData.ID = device.ID

	}
	if emptyOrContains(fields, "Key") {
		validKey, validKeyErr := device.ValidKeyDefault()
		if !validKey {
			err = validKeyErr
			return
		}

		deviceData.Key = device.Key

	}
	if emptyOrContains(fields, "Name") {
		validName, validNameErr := device.ValidNameDefault()
		if !validName {
			err = validNameErr
			return
		}

		deviceData.Name = device.Name

	}
	if emptyOrContains(fields, "Secret") {
		validSecret, validSecretErr := device.ValidSecretDefault()
		if !validSecret {
			err = validSecretErr
			return
		}

		deviceData.Secret = device.Secret

	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))

		deviceJson, err := json.Marshal(deviceData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(deviceData.ID), deviceJson)

		if err == nil {
			rowsAffected = 1
		}

		return err
	})

	return
}

func (boltdb *DB) GetDeviceList(offset int, limit int, filters Filters, returnFields []string, sortBy SortBy) (results []Device, totalResults int64, err error) {
	validationErrorPrefix := "get_user_error:"

	if !(filters.Operator == "AND" || filters.Operator == "OR") {
		err = errors.New(validationErrorPrefix + " filter operator error")
	}

	var deviceList Devices

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var device Device
			err := json.Unmarshal(v, &device)

			includeThis, err := includeThisDevice(filters, device)

			if err != nil {
				return err
			}

			if includeThis {
				resultDevice := Device{ID: device.ID}
				if emptyOrContains(returnFields, "ID") {
					resultDevice.ID = device.ID
				}
				if emptyOrContains(returnFields, "Key") {
					resultDevice.Key = device.Key
				}
				if emptyOrContains(returnFields, "Name") {
					resultDevice.Name = device.Name
				}
				if emptyOrContains(returnFields, "Secret") {
					resultDevice.Secret = device.Secret
				}
				if emptyOrContains(returnFields, "Created") {
					resultDevice.Created = device.Created
				}
				if emptyOrContains(returnFields, "Updated") {
					resultDevice.Updated = device.Updated
				}

				deviceList = append(deviceList, resultDevice)
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	if sortBy.Direction == "ASC" || sortBy.Direction == "DESC" {
		if sortBy.Field == "ID" && sortBy.Direction == "ASC" {
			sort.Sort(sortByDeviceID{deviceList})
		} else if sortBy.Field == "ID" && sortBy.Direction == "DESC" {
			sort.Sort(sortByDeviceIDDesc{deviceList})
		}
		if sortBy.Field == "Key" && sortBy.Direction == "ASC" {
			sort.Sort(sortByDeviceKey{deviceList})
		} else if sortBy.Field == "Key" && sortBy.Direction == "DESC" {
			sort.Sort(sortByDeviceKeyDesc{deviceList})
		}
		if sortBy.Field == "Name" && sortBy.Direction == "ASC" {
			sort.Sort(sortByDeviceName{deviceList})
		} else if sortBy.Field == "Name" && sortBy.Direction == "DESC" {
			sort.Sort(sortByDeviceNameDesc{deviceList})
		}
		if sortBy.Field == "Secret" && sortBy.Direction == "ASC" {
			sort.Sort(sortByDeviceSecret{deviceList})
		} else if sortBy.Field == "Secret" && sortBy.Direction == "DESC" {
			sort.Sort(sortByDeviceSecretDesc{deviceList})
		}
		if sortBy.Field == "Created" && sortBy.Direction == "ASC" {
			sort.Sort(sortByDeviceCreated{deviceList})
		} else if sortBy.Field == "Created" && sortBy.Direction == "DESC" {
			sort.Sort(sortByDeviceCreatedDesc{deviceList})
		}
		if sortBy.Field == "Updated" && sortBy.Direction == "ASC" {
			sort.Sort(sortByDeviceUpdated{deviceList})
		} else if sortBy.Field == "Updated" && sortBy.Direction == "DESC" {
			sort.Sort(sortByDeviceUpdatedDesc{deviceList})
		}

	} else {
		err = errors.New(validationErrorPrefix + " sort Direction error")
	}

	totalResults = int64(len(deviceList))

	for indexDevice, resultDevice := range deviceList {
		if indexDevice >= offset && indexDevice < (offset+limit) {
			results = append(results, resultDevice)
		}
	}

	return
}
func includeThisDevice(filters Filters, device Device) (include bool, err error) {
	validationErrorPrefix := "get_device_error:"

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

			if condition.Comparison == "=" && device.ID == conditionValueID {
				meetConditionID = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueID, "%") && strings.HasSuffix(conditionValueID, "%") {
					if strings.Contains(device.ID, strings.TrimSuffix(strings.TrimPrefix(conditionValueID, "%"), "%")) {
						meetConditionID = true
					}
				} else if strings.HasPrefix(conditionValueID, "%") {
					if strings.HasSuffix(device.ID, strings.TrimPrefix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if strings.HasSuffix(conditionValueID, "%") {
					if strings.HasPrefix(device.ID, strings.TrimSuffix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if device.ID == conditionValueID {
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

		meetConditionKey := false

		if condition.Field == "Key" {
			conditionValueKey := condition.Value.(string)

			if condition.Comparison == "=" && device.Key == conditionValueKey {
				meetConditionKey = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueKey, "%") && strings.HasSuffix(conditionValueKey, "%") {
					if strings.Contains(device.Key, strings.TrimSuffix(strings.TrimPrefix(conditionValueKey, "%"), "%")) {
						meetConditionKey = true
					}
				} else if strings.HasPrefix(conditionValueKey, "%") {
					if strings.HasSuffix(device.Key, strings.TrimPrefix(conditionValueKey, "%")) {
						meetConditionKey = true
					}
				} else if strings.HasSuffix(conditionValueKey, "%") {
					if strings.HasPrefix(device.Key, strings.TrimSuffix(conditionValueKey, "%")) {
						meetConditionKey = true
					}
				} else if device.Key == conditionValueKey {
					meetConditionKey = true
				}
			}

			if meetConditionKey {
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

		meetConditionName := false

		if condition.Field == "Name" {
			conditionValueName := condition.Value.(string)

			if condition.Comparison == "=" && device.Name == conditionValueName {
				meetConditionName = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueName, "%") && strings.HasSuffix(conditionValueName, "%") {
					if strings.Contains(device.Name, strings.TrimSuffix(strings.TrimPrefix(conditionValueName, "%"), "%")) {
						meetConditionName = true
					}
				} else if strings.HasPrefix(conditionValueName, "%") {
					if strings.HasSuffix(device.Name, strings.TrimPrefix(conditionValueName, "%")) {
						meetConditionName = true
					}
				} else if strings.HasSuffix(conditionValueName, "%") {
					if strings.HasPrefix(device.Name, strings.TrimSuffix(conditionValueName, "%")) {
						meetConditionName = true
					}
				} else if device.Name == conditionValueName {
					meetConditionName = true
				}
			}

			if meetConditionName {
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

		meetConditionSecret := false

		if condition.Field == "Secret" {
			conditionValueSecret := condition.Value.(string)

			if condition.Comparison == "=" && device.Secret == conditionValueSecret {
				meetConditionSecret = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueSecret, "%") && strings.HasSuffix(conditionValueSecret, "%") {
					if strings.Contains(device.Secret, strings.TrimSuffix(strings.TrimPrefix(conditionValueSecret, "%"), "%")) {
						meetConditionSecret = true
					}
				} else if strings.HasPrefix(conditionValueSecret, "%") {
					if strings.HasSuffix(device.Secret, strings.TrimPrefix(conditionValueSecret, "%")) {
						meetConditionSecret = true
					}
				} else if strings.HasSuffix(conditionValueSecret, "%") {
					if strings.HasPrefix(device.Secret, strings.TrimSuffix(conditionValueSecret, "%")) {
						meetConditionSecret = true
					}
				} else if device.Secret == conditionValueSecret {
					meetConditionSecret = true
				}
			}

			if meetConditionSecret {
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
			diffCreated := device.Created.Sub(conditionValueCreated)

			if condition.Comparison == "=" && device.Created == conditionValueCreated {
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

		meetConditionUpdated := false

		if condition.Field == "Updated" {
			conditionValueUpdated := condition.Value.(time.Time)
			diffUpdated := device.Updated.Sub(conditionValueUpdated)

			if condition.Comparison == "=" && device.Updated == conditionValueUpdated {
				meetConditionUpdated = true
			} else if condition.Comparison == ">" && diffUpdated > 0 {
				meetConditionUpdated = true
			} else if condition.Comparison == "<" && diffUpdated < 0 {
				meetConditionUpdated = true
			}

			if meetConditionUpdated {
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

func (s Devices) Len() int {
	return len(s)
}
func (s Devices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortByDeviceID struct {
	Devices
}

func (s sortByDeviceID) Less(i, j int) bool {
	return s.Devices[i].ID < s.Devices[j].ID
}

type sortByDeviceIDDesc struct {
	Devices
}

func (s sortByDeviceIDDesc) Less(i, j int) bool {
	return s.Devices[i].ID > s.Devices[j].ID

}

type sortByDeviceKey struct {
	Devices
}

func (s sortByDeviceKey) Less(i, j int) bool {
	return s.Devices[i].Key < s.Devices[j].Key
}

type sortByDeviceKeyDesc struct {
	Devices
}

func (s sortByDeviceKeyDesc) Less(i, j int) bool {
	return s.Devices[i].Key > s.Devices[j].Key

}

type sortByDeviceName struct {
	Devices
}

func (s sortByDeviceName) Less(i, j int) bool {
	return s.Devices[i].Name < s.Devices[j].Name
}

type sortByDeviceNameDesc struct {
	Devices
}

func (s sortByDeviceNameDesc) Less(i, j int) bool {
	return s.Devices[i].Name > s.Devices[j].Name

}

type sortByDeviceSecret struct {
	Devices
}

func (s sortByDeviceSecret) Less(i, j int) bool {
	return s.Devices[i].Secret < s.Devices[j].Secret
}

type sortByDeviceSecretDesc struct {
	Devices
}

func (s sortByDeviceSecretDesc) Less(i, j int) bool {
	return s.Devices[i].Secret > s.Devices[j].Secret

}

type sortByDeviceCreated struct {
	Devices
}

func (s sortByDeviceCreated) Less(i, j int) bool {
	diffLastModification := s.Devices[i].Created.Sub(s.Devices[j].Created)
	return diffLastModification < 0
}

type sortByDeviceCreatedDesc struct {
	Devices
}

func (s sortByDeviceCreatedDesc) Less(i, j int) bool {
	diffLastModification := s.Devices[i].Created.Sub(s.Devices[j].Created)
	return diffLastModification > 0

}

type sortByDeviceUpdated struct {
	Devices
}

func (s sortByDeviceUpdated) Less(i, j int) bool {
	diffLastModification := s.Devices[i].Updated.Sub(s.Devices[j].Updated)
	return diffLastModification < 0
}

type sortByDeviceUpdatedDesc struct {
	Devices
}

func (s sortByDeviceUpdatedDesc) Less(i, j int) bool {
	diffLastModification := s.Devices[i].Updated.Sub(s.Devices[j].Updated)
	return diffLastModification > 0

}

func (device Device) ValidIDDefault() (validField bool, err error) {
	validField, _ = validator.UUID(device.ID)
	if !validField {
		err = errors.New("error_uuid__device___ID")
		return
	}

	return
}
func (device Device) ValidKeyDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(device.Key, 100)
	if !validField {
		err = errors.New("error_maxlength__device___Key")
		return
	}

	return
}
func (device Device) ValidNameDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(device.Name, 100)
	if !validField {
		err = errors.New("error_maxlength__device___Name")
		return
	}

	return
}
func (device Device) ValidSecretDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(device.Secret, 100)
	if !validField {
		err = errors.New("error_maxlength__device___Secret")
		return
	}

	return
}
func (device Device) ValidCreatedDefault() (validField bool, err error) {
	validField = true

	return
}
func (device Device) ValidUpdatedDefault() (validField bool, err error) {
	validField = true

	return
}
