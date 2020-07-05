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

	"github.com/jempe/gopicam/pkg/validator"
)

type Location struct {
	ID          string    `json:"id"`
	DeviceIndex int       `json:"device_index"`
	Device      string    `json:"device"`
	Latitude    int       `json:"latitude"`
	Longitude   int       `json:"longitude"`
	Accuracy    int       `json:"accuracy"`
	Altitude    int       `json:"altitude"`
	Speed       int       `json:"speed"`
	Battery     int       `json:"battery"`
	DeviceTime  int64     `json:"device_time"`
	BearingTo   int       `json:"bearing_to"`
	Wifi        string    `json:"wifi"`
	Created     time.Time `json:"created"`
}

type Locations []Location

func (boltdb *DB) GetLocation(locationID string) (location Location, err error) {
	validID, err := validator.UUID(locationID)
	if !validID {
		return location, err
	}

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))
		v := b.Get([]byte(locationID))

		if v == nil {
			return errors.New("location not found")
		}

		err := json.Unmarshal(v, &location)

		return err
	})

	return location, err
}

func (boltdb *DB) InsertLocation(location Location, fields []string) (locationID string, err error) {

	validationErrorPrefix := "insert_location_error:"

	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(validationErrorPrefix, err)
		return
	}

	var locationData Location

	if location.ID == "" {
		locationID = id.String()

		location.ID = locationID
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := location.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		locationData.ID = location.ID

	}
	if emptyOrContains(fields, "DeviceIndex") {
		validDeviceIndex, validDeviceIndexErr := location.ValidDeviceIndexDefault()
		if !validDeviceIndex {
			err = validDeviceIndexErr
			return
		}

		locationData.DeviceIndex = location.DeviceIndex

	}
	if emptyOrContains(fields, "Device") {
		validDevice, validDeviceErr := location.ValidDeviceDefault()
		if !validDevice {
			err = validDeviceErr
			return
		}

		locationData.Device = location.Device

	}
	if emptyOrContains(fields, "Latitude") {
		validLatitude, validLatitudeErr := location.ValidLatitudeDefault()
		if !validLatitude {
			err = validLatitudeErr
			return
		}

		locationData.Latitude = location.Latitude

	}
	if emptyOrContains(fields, "Longitude") {
		validLongitude, validLongitudeErr := location.ValidLongitudeDefault()
		if !validLongitude {
			err = validLongitudeErr
			return
		}

		locationData.Longitude = location.Longitude

	}
	if emptyOrContains(fields, "Accuracy") {
		validAccuracy, validAccuracyErr := location.ValidAccuracyDefault()
		if !validAccuracy {
			err = validAccuracyErr
			return
		}

		locationData.Accuracy = location.Accuracy

	}
	if emptyOrContains(fields, "Altitude") {
		validAltitude, validAltitudeErr := location.ValidAltitudeDefault()
		if !validAltitude {
			err = validAltitudeErr
			return
		}

		locationData.Altitude = location.Altitude

	}
	if emptyOrContains(fields, "Speed") {
		validSpeed, validSpeedErr := location.ValidSpeedDefault()
		if !validSpeed {
			err = validSpeedErr
			return
		}

		locationData.Speed = location.Speed

	}
	if emptyOrContains(fields, "Battery") {
		validBattery, validBatteryErr := location.ValidBatteryDefault()
		if !validBattery {
			err = validBatteryErr
			return
		}

		locationData.Battery = location.Battery

	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := location.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		locationData.DeviceTime = location.DeviceTime

	}
	if emptyOrContains(fields, "BearingTo") {
		validBearingTo, validBearingToErr := location.ValidBearingToDefault()
		if !validBearingTo {
			err = validBearingToErr
			return
		}

		locationData.BearingTo = location.BearingTo

	}
	if emptyOrContains(fields, "Wifi") {
		validWifi, validWifiErr := location.ValidWifiDefault()
		if !validWifi {
			err = validWifiErr
			return
		}

		locationData.Wifi = location.Wifi

	}

	existLocationData, _ := boltdb.GetLocation(location.ID)
	if existLocationData.ID != "" {
		err = errors.New(validationErrorPrefix + " location with ID " + location.ID + " already exists")
		return
	}
	locationData.Created = time.Now().UTC()

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))

		locationJson, err := json.Marshal(locationData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(location.ID), locationJson)
		return err
	})

	return
}

func (boltdb *DB) DeleteLocation(locationID string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "delete_location_error:"

	validID, err := validator.UUID(locationID)
	if !validID {
		return
	}

	locationData, err := boltdb.GetLocation(locationID)
	if err != nil {
		return
	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))
		err = b.Delete([]byte(locationData.ID))

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

func (boltdb *DB) UpdateLocation(location Location, fields []string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "update_location_error:"

	validID, err := validator.UUID(location.ID)
	if !validID {
		return
	}

	locationData, err := boltdb.GetLocation(location.ID)
	if err != nil {
		return
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := location.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		locationData.ID = location.ID

	}
	if emptyOrContains(fields, "DeviceIndex") {
		validDeviceIndex, validDeviceIndexErr := location.ValidDeviceIndexDefault()
		if !validDeviceIndex {
			err = validDeviceIndexErr
			return
		}

		locationData.DeviceIndex = location.DeviceIndex

	}
	if emptyOrContains(fields, "Device") {
		validDevice, validDeviceErr := location.ValidDeviceDefault()
		if !validDevice {
			err = validDeviceErr
			return
		}

		locationData.Device = location.Device

	}
	if emptyOrContains(fields, "Latitude") {
		validLatitude, validLatitudeErr := location.ValidLatitudeDefault()
		if !validLatitude {
			err = validLatitudeErr
			return
		}

		locationData.Latitude = location.Latitude

	}
	if emptyOrContains(fields, "Longitude") {
		validLongitude, validLongitudeErr := location.ValidLongitudeDefault()
		if !validLongitude {
			err = validLongitudeErr
			return
		}

		locationData.Longitude = location.Longitude

	}
	if emptyOrContains(fields, "Accuracy") {
		validAccuracy, validAccuracyErr := location.ValidAccuracyDefault()
		if !validAccuracy {
			err = validAccuracyErr
			return
		}

		locationData.Accuracy = location.Accuracy

	}
	if emptyOrContains(fields, "Altitude") {
		validAltitude, validAltitudeErr := location.ValidAltitudeDefault()
		if !validAltitude {
			err = validAltitudeErr
			return
		}

		locationData.Altitude = location.Altitude

	}
	if emptyOrContains(fields, "Speed") {
		validSpeed, validSpeedErr := location.ValidSpeedDefault()
		if !validSpeed {
			err = validSpeedErr
			return
		}

		locationData.Speed = location.Speed

	}
	if emptyOrContains(fields, "Battery") {
		validBattery, validBatteryErr := location.ValidBatteryDefault()
		if !validBattery {
			err = validBatteryErr
			return
		}

		locationData.Battery = location.Battery

	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := location.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		locationData.DeviceTime = location.DeviceTime

	}
	if emptyOrContains(fields, "BearingTo") {
		validBearingTo, validBearingToErr := location.ValidBearingToDefault()
		if !validBearingTo {
			err = validBearingToErr
			return
		}

		locationData.BearingTo = location.BearingTo

	}
	if emptyOrContains(fields, "Wifi") {
		validWifi, validWifiErr := location.ValidWifiDefault()
		if !validWifi {
			err = validWifiErr
			return
		}

		locationData.Wifi = location.Wifi

	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))

		locationJson, err := json.Marshal(locationData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(locationData.ID), locationJson)

		if err == nil {
			rowsAffected = 1
		}

		return err
	})

	return
}

func (boltdb *DB) GetLocationList(offset int, limit int, filters Filters, returnFields []string, sortBy SortBy) (results []Location, totalResults int64, err error) {
	validationErrorPrefix := "get_user_error:"

	if !(filters.Operator == "AND" || filters.Operator == "OR") {
		err = errors.New(validationErrorPrefix + " filter operator error")
	}

	var locationList Locations

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("locations"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var location Location
			err := json.Unmarshal(v, &location)

			includeThis, err := includeThisLocation(filters, location)

			if err != nil {
				return err
			}

			if includeThis {
				resultLocation := Location{ID: location.ID}
				if emptyOrContains(returnFields, "ID") {
					resultLocation.ID = location.ID
				}
				if emptyOrContains(returnFields, "DeviceIndex") {
					resultLocation.DeviceIndex = location.DeviceIndex
				}
				if emptyOrContains(returnFields, "Device") {
					resultLocation.Device = location.Device
				}
				if emptyOrContains(returnFields, "Latitude") {
					resultLocation.Latitude = location.Latitude
				}
				if emptyOrContains(returnFields, "Longitude") {
					resultLocation.Longitude = location.Longitude
				}
				if emptyOrContains(returnFields, "Accuracy") {
					resultLocation.Accuracy = location.Accuracy
				}
				if emptyOrContains(returnFields, "Altitude") {
					resultLocation.Altitude = location.Altitude
				}
				if emptyOrContains(returnFields, "Speed") {
					resultLocation.Speed = location.Speed
				}
				if emptyOrContains(returnFields, "Battery") {
					resultLocation.Battery = location.Battery
				}
				if emptyOrContains(returnFields, "DeviceTime") {
					resultLocation.DeviceTime = location.DeviceTime
				}
				if emptyOrContains(returnFields, "BearingTo") {
					resultLocation.BearingTo = location.BearingTo
				}
				if emptyOrContains(returnFields, "Wifi") {
					resultLocation.Wifi = location.Wifi
				}
				if emptyOrContains(returnFields, "Created") {
					resultLocation.Created = location.Created
				}

				locationList = append(locationList, resultLocation)
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	if sortBy.Direction == "ASC" || sortBy.Direction == "DESC" {
		if sortBy.Field == "ID" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationID{locationList})
		} else if sortBy.Field == "ID" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationIDDesc{locationList})
		}
		if sortBy.Field == "DeviceIndex" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationDeviceIndex{locationList})
		} else if sortBy.Field == "DeviceIndex" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationDeviceIndexDesc{locationList})
		}
		if sortBy.Field == "Device" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationDevice{locationList})
		} else if sortBy.Field == "Device" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationDeviceDesc{locationList})
		}
		if sortBy.Field == "Latitude" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationLatitude{locationList})
		} else if sortBy.Field == "Latitude" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationLatitudeDesc{locationList})
		}
		if sortBy.Field == "Longitude" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationLongitude{locationList})
		} else if sortBy.Field == "Longitude" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationLongitudeDesc{locationList})
		}
		if sortBy.Field == "Accuracy" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationAccuracy{locationList})
		} else if sortBy.Field == "Accuracy" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationAccuracyDesc{locationList})
		}
		if sortBy.Field == "Altitude" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationAltitude{locationList})
		} else if sortBy.Field == "Altitude" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationAltitudeDesc{locationList})
		}
		if sortBy.Field == "Speed" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationSpeed{locationList})
		} else if sortBy.Field == "Speed" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationSpeedDesc{locationList})
		}
		if sortBy.Field == "Battery" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationBattery{locationList})
		} else if sortBy.Field == "Battery" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationBatteryDesc{locationList})
		}
		if sortBy.Field == "DeviceTime" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationDeviceTime{locationList})
		} else if sortBy.Field == "DeviceTime" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationDeviceTimeDesc{locationList})
		}
		if sortBy.Field == "BearingTo" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationBearingTo{locationList})
		} else if sortBy.Field == "BearingTo" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationBearingToDesc{locationList})
		}
		if sortBy.Field == "Wifi" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationWifi{locationList})
		} else if sortBy.Field == "Wifi" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationWifiDesc{locationList})
		}
		if sortBy.Field == "Created" && sortBy.Direction == "ASC" {
			sort.Sort(sortByLocationCreated{locationList})
		} else if sortBy.Field == "Created" && sortBy.Direction == "DESC" {
			sort.Sort(sortByLocationCreatedDesc{locationList})
		}

	} else {
		err = errors.New(validationErrorPrefix + " sort Direction error")
	}

	totalResults = int64(len(locationList))

	for indexLocation, resultLocation := range locationList {
		if indexLocation >= offset && indexLocation < (offset+limit) {
			results = append(results, resultLocation)
		}
	}

	return
}
func includeThisLocation(filters Filters, location Location) (include bool, err error) {
	validationErrorPrefix := "get_location_error:"

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

			if condition.Comparison == "=" && location.ID == conditionValueID {
				meetConditionID = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueID, "%") && strings.HasSuffix(conditionValueID, "%") {
					if strings.Contains(location.ID, strings.TrimSuffix(strings.TrimPrefix(conditionValueID, "%"), "%")) {
						meetConditionID = true
					}
				} else if strings.HasPrefix(conditionValueID, "%") {
					if strings.HasSuffix(location.ID, strings.TrimPrefix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if strings.HasSuffix(conditionValueID, "%") {
					if strings.HasPrefix(location.ID, strings.TrimSuffix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if location.ID == conditionValueID {
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

		meetConditionDeviceIndex := false

		if condition.Field == "DeviceIndex" {
			conditionValueDeviceIndex := condition.Value.(int)

			if condition.Comparison == "=" && location.DeviceIndex == conditionValueDeviceIndex {
				meetConditionDeviceIndex = true
			} else if condition.Comparison == ">" && location.DeviceIndex > conditionValueDeviceIndex {
				meetConditionDeviceIndex = true
			} else if condition.Comparison == "<" && location.DeviceIndex < conditionValueDeviceIndex {
				meetConditionDeviceIndex = true
			}

			if meetConditionDeviceIndex {
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

		meetConditionDevice := false

		if condition.Field == "Device" {
			conditionValueDevice := condition.Value.(string)

			if condition.Comparison == "=" && location.Device == conditionValueDevice {
				meetConditionDevice = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueDevice, "%") && strings.HasSuffix(conditionValueDevice, "%") {
					if strings.Contains(location.Device, strings.TrimSuffix(strings.TrimPrefix(conditionValueDevice, "%"), "%")) {
						meetConditionDevice = true
					}
				} else if strings.HasPrefix(conditionValueDevice, "%") {
					if strings.HasSuffix(location.Device, strings.TrimPrefix(conditionValueDevice, "%")) {
						meetConditionDevice = true
					}
				} else if strings.HasSuffix(conditionValueDevice, "%") {
					if strings.HasPrefix(location.Device, strings.TrimSuffix(conditionValueDevice, "%")) {
						meetConditionDevice = true
					}
				} else if location.Device == conditionValueDevice {
					meetConditionDevice = true
				}
			}

			if meetConditionDevice {
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

		meetConditionLatitude := false

		if condition.Field == "Latitude" {
			conditionValueLatitude := condition.Value.(int)

			if condition.Comparison == "=" && location.Latitude == conditionValueLatitude {
				meetConditionLatitude = true
			} else if condition.Comparison == ">" && location.Latitude > conditionValueLatitude {
				meetConditionLatitude = true
			} else if condition.Comparison == "<" && location.Latitude < conditionValueLatitude {
				meetConditionLatitude = true
			}

			if meetConditionLatitude {
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

		meetConditionLongitude := false

		if condition.Field == "Longitude" {
			conditionValueLongitude := condition.Value.(int)

			if condition.Comparison == "=" && location.Longitude == conditionValueLongitude {
				meetConditionLongitude = true
			} else if condition.Comparison == ">" && location.Longitude > conditionValueLongitude {
				meetConditionLongitude = true
			} else if condition.Comparison == "<" && location.Longitude < conditionValueLongitude {
				meetConditionLongitude = true
			}

			if meetConditionLongitude {
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

		meetConditionAccuracy := false

		if condition.Field == "Accuracy" {
			conditionValueAccuracy := condition.Value.(int)

			if condition.Comparison == "=" && location.Accuracy == conditionValueAccuracy {
				meetConditionAccuracy = true
			} else if condition.Comparison == ">" && location.Accuracy > conditionValueAccuracy {
				meetConditionAccuracy = true
			} else if condition.Comparison == "<" && location.Accuracy < conditionValueAccuracy {
				meetConditionAccuracy = true
			}

			if meetConditionAccuracy {
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

		meetConditionAltitude := false

		if condition.Field == "Altitude" {
			conditionValueAltitude := condition.Value.(int)

			if condition.Comparison == "=" && location.Altitude == conditionValueAltitude {
				meetConditionAltitude = true
			} else if condition.Comparison == ">" && location.Altitude > conditionValueAltitude {
				meetConditionAltitude = true
			} else if condition.Comparison == "<" && location.Altitude < conditionValueAltitude {
				meetConditionAltitude = true
			}

			if meetConditionAltitude {
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

		meetConditionSpeed := false

		if condition.Field == "Speed" {
			conditionValueSpeed := condition.Value.(int)

			if condition.Comparison == "=" && location.Speed == conditionValueSpeed {
				meetConditionSpeed = true
			} else if condition.Comparison == ">" && location.Speed > conditionValueSpeed {
				meetConditionSpeed = true
			} else if condition.Comparison == "<" && location.Speed < conditionValueSpeed {
				meetConditionSpeed = true
			}

			if meetConditionSpeed {
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

		meetConditionBattery := false

		if condition.Field == "Battery" {
			conditionValueBattery := condition.Value.(int)

			if condition.Comparison == "=" && location.Battery == conditionValueBattery {
				meetConditionBattery = true
			} else if condition.Comparison == ">" && location.Battery > conditionValueBattery {
				meetConditionBattery = true
			} else if condition.Comparison == "<" && location.Battery < conditionValueBattery {
				meetConditionBattery = true
			}

			if meetConditionBattery {
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

		meetConditionDeviceTime := false

		if condition.Field == "DeviceTime" {
			conditionValueDeviceTime := condition.Value.(int64)

			if condition.Comparison == "=" && location.DeviceTime == conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == ">" && location.DeviceTime > conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == "<" && location.DeviceTime < conditionValueDeviceTime {
				meetConditionDeviceTime = true
			}

			if meetConditionDeviceTime {
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

		meetConditionBearingTo := false

		if condition.Field == "BearingTo" {
			conditionValueBearingTo := condition.Value.(int)

			if condition.Comparison == "=" && location.BearingTo == conditionValueBearingTo {
				meetConditionBearingTo = true
			} else if condition.Comparison == ">" && location.BearingTo > conditionValueBearingTo {
				meetConditionBearingTo = true
			} else if condition.Comparison == "<" && location.BearingTo < conditionValueBearingTo {
				meetConditionBearingTo = true
			}

			if meetConditionBearingTo {
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

		meetConditionWifi := false

		if condition.Field == "Wifi" {
			conditionValueWifi := condition.Value.(string)

			if condition.Comparison == "=" && location.Wifi == conditionValueWifi {
				meetConditionWifi = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueWifi, "%") && strings.HasSuffix(conditionValueWifi, "%") {
					if strings.Contains(location.Wifi, strings.TrimSuffix(strings.TrimPrefix(conditionValueWifi, "%"), "%")) {
						meetConditionWifi = true
					}
				} else if strings.HasPrefix(conditionValueWifi, "%") {
					if strings.HasSuffix(location.Wifi, strings.TrimPrefix(conditionValueWifi, "%")) {
						meetConditionWifi = true
					}
				} else if strings.HasSuffix(conditionValueWifi, "%") {
					if strings.HasPrefix(location.Wifi, strings.TrimSuffix(conditionValueWifi, "%")) {
						meetConditionWifi = true
					}
				} else if location.Wifi == conditionValueWifi {
					meetConditionWifi = true
				}
			}

			if meetConditionWifi {
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
			diffCreated := location.Created.Sub(conditionValueCreated)

			if condition.Comparison == "=" && location.Created == conditionValueCreated {
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

func (s Locations) Len() int {
	return len(s)
}
func (s Locations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortByLocationID struct {
	Locations
}

func (s sortByLocationID) Less(i, j int) bool {
	return s.Locations[i].ID < s.Locations[j].ID
}

type sortByLocationIDDesc struct {
	Locations
}

func (s sortByLocationIDDesc) Less(i, j int) bool {
	return s.Locations[i].ID > s.Locations[j].ID

}

type sortByLocationDeviceIndex struct {
	Locations
}

func (s sortByLocationDeviceIndex) Less(i, j int) bool {
	return s.Locations[i].DeviceIndex < s.Locations[j].DeviceIndex
}

type sortByLocationDeviceIndexDesc struct {
	Locations
}

func (s sortByLocationDeviceIndexDesc) Less(i, j int) bool {
	return s.Locations[i].DeviceIndex > s.Locations[j].DeviceIndex

}

type sortByLocationDevice struct {
	Locations
}

func (s sortByLocationDevice) Less(i, j int) bool {
	return s.Locations[i].Device < s.Locations[j].Device
}

type sortByLocationDeviceDesc struct {
	Locations
}

func (s sortByLocationDeviceDesc) Less(i, j int) bool {
	return s.Locations[i].Device > s.Locations[j].Device

}

type sortByLocationLatitude struct {
	Locations
}

func (s sortByLocationLatitude) Less(i, j int) bool {
	return s.Locations[i].Latitude < s.Locations[j].Latitude
}

type sortByLocationLatitudeDesc struct {
	Locations
}

func (s sortByLocationLatitudeDesc) Less(i, j int) bool {
	return s.Locations[i].Latitude > s.Locations[j].Latitude

}

type sortByLocationLongitude struct {
	Locations
}

func (s sortByLocationLongitude) Less(i, j int) bool {
	return s.Locations[i].Longitude < s.Locations[j].Longitude
}

type sortByLocationLongitudeDesc struct {
	Locations
}

func (s sortByLocationLongitudeDesc) Less(i, j int) bool {
	return s.Locations[i].Longitude > s.Locations[j].Longitude

}

type sortByLocationAccuracy struct {
	Locations
}

func (s sortByLocationAccuracy) Less(i, j int) bool {
	return s.Locations[i].Accuracy < s.Locations[j].Accuracy
}

type sortByLocationAccuracyDesc struct {
	Locations
}

func (s sortByLocationAccuracyDesc) Less(i, j int) bool {
	return s.Locations[i].Accuracy > s.Locations[j].Accuracy

}

type sortByLocationAltitude struct {
	Locations
}

func (s sortByLocationAltitude) Less(i, j int) bool {
	return s.Locations[i].Altitude < s.Locations[j].Altitude
}

type sortByLocationAltitudeDesc struct {
	Locations
}

func (s sortByLocationAltitudeDesc) Less(i, j int) bool {
	return s.Locations[i].Altitude > s.Locations[j].Altitude

}

type sortByLocationSpeed struct {
	Locations
}

func (s sortByLocationSpeed) Less(i, j int) bool {
	return s.Locations[i].Speed < s.Locations[j].Speed
}

type sortByLocationSpeedDesc struct {
	Locations
}

func (s sortByLocationSpeedDesc) Less(i, j int) bool {
	return s.Locations[i].Speed > s.Locations[j].Speed

}

type sortByLocationBattery struct {
	Locations
}

func (s sortByLocationBattery) Less(i, j int) bool {
	return s.Locations[i].Battery < s.Locations[j].Battery
}

type sortByLocationBatteryDesc struct {
	Locations
}

func (s sortByLocationBatteryDesc) Less(i, j int) bool {
	return s.Locations[i].Battery > s.Locations[j].Battery

}

type sortByLocationDeviceTime struct {
	Locations
}

func (s sortByLocationDeviceTime) Less(i, j int) bool {
	return s.Locations[i].DeviceTime < s.Locations[j].DeviceTime
}

type sortByLocationDeviceTimeDesc struct {
	Locations
}

func (s sortByLocationDeviceTimeDesc) Less(i, j int) bool {
	return s.Locations[i].DeviceTime > s.Locations[j].DeviceTime

}

type sortByLocationBearingTo struct {
	Locations
}

func (s sortByLocationBearingTo) Less(i, j int) bool {
	return s.Locations[i].BearingTo < s.Locations[j].BearingTo
}

type sortByLocationBearingToDesc struct {
	Locations
}

func (s sortByLocationBearingToDesc) Less(i, j int) bool {
	return s.Locations[i].BearingTo > s.Locations[j].BearingTo

}

type sortByLocationWifi struct {
	Locations
}

func (s sortByLocationWifi) Less(i, j int) bool {
	return s.Locations[i].Wifi < s.Locations[j].Wifi
}

type sortByLocationWifiDesc struct {
	Locations
}

func (s sortByLocationWifiDesc) Less(i, j int) bool {
	return s.Locations[i].Wifi > s.Locations[j].Wifi

}

type sortByLocationCreated struct {
	Locations
}

func (s sortByLocationCreated) Less(i, j int) bool {
	diffLastModification := s.Locations[i].Created.Sub(s.Locations[j].Created)
	return diffLastModification < 0
}

type sortByLocationCreatedDesc struct {
	Locations
}

func (s sortByLocationCreatedDesc) Less(i, j int) bool {
	diffLastModification := s.Locations[i].Created.Sub(s.Locations[j].Created)
	return diffLastModification > 0

}

func (location Location) ValidIDDefault() (validField bool, err error) {
	validField, _ = validator.UUID(location.ID)
	if !validField {
		err = errors.New("error_uuid__location___ID")
		return
	}

	return
}
func (location Location) ValidDeviceIndexDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidDeviceDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidLatitudeDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidLongitudeDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidAccuracyDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidAltitudeDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidSpeedDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidBatteryDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidDeviceTimeDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidBearingToDefault() (validField bool, err error) {
	validField = true

	return
}
func (location Location) ValidWifiDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(location.Wifi, 2083)
	if !validField {
		err = errors.New("error_maxlength__location___Wifi")
		return
	}

	return
}
func (location Location) ValidCreatedDefault() (validField bool, err error) {
	validField = true

	return
}
