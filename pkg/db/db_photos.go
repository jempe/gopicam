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

type Photo struct {
	ID         string    `json:"id"`
	FileType   string    `json:"file_type"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Size       int       `json:"size"`
	DeviceTime int64     `json:"device_time"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type Photos []Photo

func (boltdb *DB) GetPhoto(photoID string) (photo Photo, err error) {
	validID, err := validator.UUID(photoID)
	if !validID {
		return photo, err
	}

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("photos"))
		v := b.Get([]byte(photoID))

		if v == nil {
			return errors.New("photo not found")
		}

		err := json.Unmarshal(v, &photo)

		return err
	})

	return photo, err
}

func (boltdb *DB) InsertPhoto(photo Photo, fields []string) (photoID string, err error) {

	validationErrorPrefix := "insert_photo_error:"

	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(validationErrorPrefix, err)
		return
	}

	var photoData Photo

	if photo.ID == "" {
		photoID = id.String()

		photo.ID = photoID
	}

	validID, validIDErr := photo.ValidIDDefault()
	if !validID {
		err = validIDErr
		return
	}

	photoData.ID = photo.ID
	if emptyOrContains(fields, "FileType") {
		validFileType, validFileTypeErr := photo.ValidFileTypeDefault()
		if !validFileType {
			err = validFileTypeErr
			return
		}

		photoData.FileType = photo.FileType
	}
	if emptyOrContains(fields, "Width") {
		validWidth, validWidthErr := photo.ValidWidthDefault()
		if !validWidth {
			err = validWidthErr
			return
		}

		photoData.Width = photo.Width
	}
	if emptyOrContains(fields, "Height") {
		validHeight, validHeightErr := photo.ValidHeightDefault()
		if !validHeight {
			err = validHeightErr
			return
		}

		photoData.Height = photo.Height
	}
	if emptyOrContains(fields, "Size") {
		validSize, validSizeErr := photo.ValidSizeDefault()
		if !validSize {
			err = validSizeErr
			return
		}

		photoData.Size = photo.Size
	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := photo.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		photoData.DeviceTime = photo.DeviceTime
	}

	existPhotoData, _ := boltdb.GetPhoto(photo.ID)
	if existPhotoData.ID != "" {
		err = errors.New(validationErrorPrefix + " photo with ID " + photo.ID + " already exists")
		return
	}
	photoData.Created = time.Now().UTC()
	photoData.Updated = time.Now().UTC()

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("photos"))

		photoJson, err := json.Marshal(photoData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(photo.ID), photoJson)
		return err
	})

	return
}

func (boltdb *DB) DeletePhoto(photoID string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "delete_photo_error:"

	validID, err := validator.UUID(photoID)
	if !validID {
		return
	}

	photoData, err := boltdb.GetPhoto(photoID)
	if err != nil {
		return
	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("photos"))
		err = b.Delete([]byte(photoData.ID))

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

func (boltdb *DB) UpdatePhoto(photo Photo, fields []string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "update_photo_error:"

	validID, err := validator.UUID(photo.ID)
	if !validID {
		return
	}

	photoData, err := boltdb.GetPhoto(photo.ID)
	if err != nil {
		return
	}

	validID, validIDErr := photo.ValidIDDefault()
	if !validID {
		err = validIDErr
		return
	}

	photoData.ID = photo.ID
	if emptyOrContains(fields, "FileType") {
		validFileType, validFileTypeErr := photo.ValidFileTypeDefault()
		if !validFileType {
			err = validFileTypeErr
			return
		}

		photoData.FileType = photo.FileType
	}
	if emptyOrContains(fields, "Width") {
		validWidth, validWidthErr := photo.ValidWidthDefault()
		if !validWidth {
			err = validWidthErr
			return
		}

		photoData.Width = photo.Width
	}
	if emptyOrContains(fields, "Height") {
		validHeight, validHeightErr := photo.ValidHeightDefault()
		if !validHeight {
			err = validHeightErr
			return
		}

		photoData.Height = photo.Height
	}
	if emptyOrContains(fields, "Size") {
		validSize, validSizeErr := photo.ValidSizeDefault()
		if !validSize {
			err = validSizeErr
			return
		}

		photoData.Size = photo.Size
	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := photo.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		photoData.DeviceTime = photo.DeviceTime
	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("photos"))

		photoJson, err := json.Marshal(photoData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(photoData.ID), photoJson)

		if err == nil {
			rowsAffected = 1
		}

		return err
	})

	return
}

func (boltdb *DB) GetPhotoList(offset int, limit int, filters Filters, returnFields []string, sortBy SortBy) (results []Photo, totalResults int64, err error) {
	validationErrorPrefix := "get_user_error:"

	if !(filters.Operator == "AND" || filters.Operator == "OR") {
		err = errors.New(validationErrorPrefix + " filter operator error")
	}

	var photoList Photos

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("photos"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var photo Photo
			err := json.Unmarshal(v, &photo)

			includeThis, err := includeThisPhoto(filters, photo)

			if err != nil {
				return err
			}

			if includeThis {
				resultPhoto := Photo{ID: photo.ID}
				if emptyOrContains(returnFields, "ID") {
					resultPhoto.ID = photo.ID
				}
				if emptyOrContains(returnFields, "FileType") {
					resultPhoto.FileType = photo.FileType
				}
				if emptyOrContains(returnFields, "Width") {
					resultPhoto.Width = photo.Width
				}
				if emptyOrContains(returnFields, "Height") {
					resultPhoto.Height = photo.Height
				}
				if emptyOrContains(returnFields, "Size") {
					resultPhoto.Size = photo.Size
				}
				if emptyOrContains(returnFields, "DeviceTime") {
					resultPhoto.DeviceTime = photo.DeviceTime
				}
				if emptyOrContains(returnFields, "Created") {
					resultPhoto.Created = photo.Created
				}
				if emptyOrContains(returnFields, "Updated") {
					resultPhoto.Updated = photo.Updated
				}

				photoList = append(photoList, resultPhoto)
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	if sortBy.Direction == "ASC" || sortBy.Direction == "DESC" {
		if sortBy.Field == "ID" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoID{photoList})
		} else if sortBy.Field == "ID" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoIDDesc{photoList})
		}
		if sortBy.Field == "FileType" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoFileType{photoList})
		} else if sortBy.Field == "FileType" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoFileTypeDesc{photoList})
		}
		if sortBy.Field == "Width" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoWidth{photoList})
		} else if sortBy.Field == "Width" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoWidthDesc{photoList})
		}
		if sortBy.Field == "Height" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoHeight{photoList})
		} else if sortBy.Field == "Height" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoHeightDesc{photoList})
		}
		if sortBy.Field == "Size" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoSize{photoList})
		} else if sortBy.Field == "Size" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoSizeDesc{photoList})
		}
		if sortBy.Field == "DeviceTime" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoDeviceTime{photoList})
		} else if sortBy.Field == "DeviceTime" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoDeviceTimeDesc{photoList})
		}
		if sortBy.Field == "Created" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoCreated{photoList})
		} else if sortBy.Field == "Created" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoCreatedDesc{photoList})
		}
		if sortBy.Field == "Updated" && sortBy.Direction == "ASC" {
			sort.Sort(sortByPhotoUpdated{photoList})
		} else if sortBy.Field == "Updated" && sortBy.Direction == "DESC" {
			sort.Sort(sortByPhotoUpdatedDesc{photoList})
		}

	} else {
		err = errors.New(validationErrorPrefix + " sort Direction error")
	}

	totalResults = int64(len(photoList))

	for indexPhoto, resultPhoto := range photoList {
		if indexPhoto >= offset && indexPhoto < (offset+limit) {
			results = append(results, resultPhoto)
		}
	}

	return
}
func includeThisPhoto(filters Filters, photo Photo) (include bool, err error) {
	validationErrorPrefix := "get_photo_error:"

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

			if condition.Comparison == "=" && photo.ID == conditionValueID {
				meetConditionID = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueID, "%") && strings.HasSuffix(conditionValueID, "%") {
					if strings.Contains(photo.ID, strings.TrimSuffix(strings.TrimPrefix(conditionValueID, "%"), "%")) {
						meetConditionID = true
					}
				} else if strings.HasPrefix(conditionValueID, "%") {
					if strings.HasSuffix(photo.ID, strings.TrimPrefix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if strings.HasSuffix(conditionValueID, "%") {
					if strings.HasPrefix(photo.ID, strings.TrimSuffix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if photo.ID == conditionValueID {
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

		meetConditionFileType := false

		if condition.Field == "FileType" {
			conditionValueFileType := condition.Value.(string)

			if condition.Comparison == "=" && photo.FileType == conditionValueFileType {
				meetConditionFileType = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueFileType, "%") && strings.HasSuffix(conditionValueFileType, "%") {
					if strings.Contains(photo.FileType, strings.TrimSuffix(strings.TrimPrefix(conditionValueFileType, "%"), "%")) {
						meetConditionFileType = true
					}
				} else if strings.HasPrefix(conditionValueFileType, "%") {
					if strings.HasSuffix(photo.FileType, strings.TrimPrefix(conditionValueFileType, "%")) {
						meetConditionFileType = true
					}
				} else if strings.HasSuffix(conditionValueFileType, "%") {
					if strings.HasPrefix(photo.FileType, strings.TrimSuffix(conditionValueFileType, "%")) {
						meetConditionFileType = true
					}
				} else if photo.FileType == conditionValueFileType {
					meetConditionFileType = true
				}
			}

			if meetConditionFileType {
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

		meetConditionWidth := false

		if condition.Field == "Width" {
			conditionValueWidth := condition.Value.(int)

			if condition.Comparison == "=" && photo.Width == conditionValueWidth {
				meetConditionWidth = true
			} else if condition.Comparison == ">" && photo.Width > conditionValueWidth {
				meetConditionWidth = true
			} else if condition.Comparison == "<" && photo.Width < conditionValueWidth {
				meetConditionWidth = true
			}

			if meetConditionWidth {
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

		meetConditionHeight := false

		if condition.Field == "Height" {
			conditionValueHeight := condition.Value.(int)

			if condition.Comparison == "=" && photo.Height == conditionValueHeight {
				meetConditionHeight = true
			} else if condition.Comparison == ">" && photo.Height > conditionValueHeight {
				meetConditionHeight = true
			} else if condition.Comparison == "<" && photo.Height < conditionValueHeight {
				meetConditionHeight = true
			}

			if meetConditionHeight {
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

		meetConditionSize := false

		if condition.Field == "Size" {
			conditionValueSize := condition.Value.(int)

			if condition.Comparison == "=" && photo.Size == conditionValueSize {
				meetConditionSize = true
			} else if condition.Comparison == ">" && photo.Size > conditionValueSize {
				meetConditionSize = true
			} else if condition.Comparison == "<" && photo.Size < conditionValueSize {
				meetConditionSize = true
			}

			if meetConditionSize {
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

			if condition.Comparison == "=" && photo.DeviceTime == conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == ">" && photo.DeviceTime > conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == "<" && photo.DeviceTime < conditionValueDeviceTime {
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

		meetConditionCreated := false

		if condition.Field == "Created" {
			conditionValueCreated := condition.Value.(time.Time)
			diffCreated := photo.Created.Sub(conditionValueCreated)

			if condition.Comparison == "=" && photo.Created == conditionValueCreated {
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
			diffUpdated := photo.Updated.Sub(conditionValueUpdated)

			if condition.Comparison == "=" && photo.Updated == conditionValueUpdated {
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

func (s Photos) Len() int {
	return len(s)
}
func (s Photos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortByPhotoID struct {
	Photos
}

func (s sortByPhotoID) Less(i, j int) bool {
	return s.Photos[i].ID < s.Photos[j].ID
}

type sortByPhotoIDDesc struct {
	Photos
}

func (s sortByPhotoIDDesc) Less(i, j int) bool {
	return s.Photos[i].ID > s.Photos[j].ID

}

type sortByPhotoFileType struct {
	Photos
}

func (s sortByPhotoFileType) Less(i, j int) bool {
	return s.Photos[i].FileType < s.Photos[j].FileType
}

type sortByPhotoFileTypeDesc struct {
	Photos
}

func (s sortByPhotoFileTypeDesc) Less(i, j int) bool {
	return s.Photos[i].FileType > s.Photos[j].FileType

}

type sortByPhotoWidth struct {
	Photos
}

func (s sortByPhotoWidth) Less(i, j int) bool {
	return s.Photos[i].Width < s.Photos[j].Width
}

type sortByPhotoWidthDesc struct {
	Photos
}

func (s sortByPhotoWidthDesc) Less(i, j int) bool {
	return s.Photos[i].Width > s.Photos[j].Width

}

type sortByPhotoHeight struct {
	Photos
}

func (s sortByPhotoHeight) Less(i, j int) bool {
	return s.Photos[i].Height < s.Photos[j].Height
}

type sortByPhotoHeightDesc struct {
	Photos
}

func (s sortByPhotoHeightDesc) Less(i, j int) bool {
	return s.Photos[i].Height > s.Photos[j].Height

}

type sortByPhotoSize struct {
	Photos
}

func (s sortByPhotoSize) Less(i, j int) bool {
	return s.Photos[i].Size < s.Photos[j].Size
}

type sortByPhotoSizeDesc struct {
	Photos
}

func (s sortByPhotoSizeDesc) Less(i, j int) bool {
	return s.Photos[i].Size > s.Photos[j].Size

}

type sortByPhotoDeviceTime struct {
	Photos
}

func (s sortByPhotoDeviceTime) Less(i, j int) bool {
	return s.Photos[i].DeviceTime < s.Photos[j].DeviceTime
}

type sortByPhotoDeviceTimeDesc struct {
	Photos
}

func (s sortByPhotoDeviceTimeDesc) Less(i, j int) bool {
	return s.Photos[i].DeviceTime > s.Photos[j].DeviceTime

}

type sortByPhotoCreated struct {
	Photos
}

func (s sortByPhotoCreated) Less(i, j int) bool {
	diffLastModification := s.Photos[i].Created.Sub(s.Photos[j].Created)
	return diffLastModification < 0
}

type sortByPhotoCreatedDesc struct {
	Photos
}

func (s sortByPhotoCreatedDesc) Less(i, j int) bool {
	diffLastModification := s.Photos[i].Created.Sub(s.Photos[j].Created)
	return diffLastModification > 0

}

type sortByPhotoUpdated struct {
	Photos
}

func (s sortByPhotoUpdated) Less(i, j int) bool {
	diffLastModification := s.Photos[i].Updated.Sub(s.Photos[j].Updated)
	return diffLastModification < 0
}

type sortByPhotoUpdatedDesc struct {
	Photos
}

func (s sortByPhotoUpdatedDesc) Less(i, j int) bool {
	diffLastModification := s.Photos[i].Updated.Sub(s.Photos[j].Updated)
	return diffLastModification > 0

}

func (photo Photo) ValidIDDefault() (validField bool, err error) {
	validField, _ = validator.UUID(photo.ID)
	if !validField {
		err = errors.New("error_uuid__photo___ID")
		return
	}

	return
}
func (photo Photo) ValidFileTypeDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(photo.FileType, 100)
	if !validField {
		err = errors.New("error_maxlength__photo___FileType")
		return
	}

	return
}
func (photo Photo) ValidWidthDefault() (validField bool, err error) {
	validField = true

	return
}
func (photo Photo) ValidHeightDefault() (validField bool, err error) {
	validField = true

	return
}
func (photo Photo) ValidSizeDefault() (validField bool, err error) {
	validField = true

	return
}
func (photo Photo) ValidDeviceTimeDefault() (validField bool, err error) {
	validField = true

	return
}
func (photo Photo) ValidCreatedDefault() (validField bool, err error) {
	validField = true

	return
}
func (photo Photo) ValidUpdatedDefault() (validField bool, err error) {
	validField = true

	return
}
