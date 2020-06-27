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

type Audio struct {
	ID         string    `json:"id"`
	FileType   string    `json:"file_type"`
	Length     int       `json:"length"`
	Size       int       `json:"size"`
	DeviceTime int64     `json:"device_time"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type Audios []Audio

func (boltdb *DB) GetAudio(audioID string) (audio Audio, err error) {
	validID, err := validator.UUID(audioID)
	if !validID {
		return audio, err
	}

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("audios"))
		v := b.Get([]byte(audioID))

		if v == nil {
			return errors.New("audio not found")
		}

		err := json.Unmarshal(v, &audio)

		return err
	})

	return audio, err
}

func (boltdb *DB) InsertAudio(audio Audio, fields []string) (audioID string, err error) {

	validationErrorPrefix := "insert_audio_error:"

	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(validationErrorPrefix, err)
		return
	}

	var audioData Audio

	if audio.ID == "" {
		audioID = id.String()

		audio.ID = audioID
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := audio.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		audioData.ID = audio.ID

	}
	if emptyOrContains(fields, "FileType") {
		validFileType, validFileTypeErr := audio.ValidFileTypeDefault()
		if !validFileType {
			err = validFileTypeErr
			return
		}

		audioData.FileType = audio.FileType

	}
	if emptyOrContains(fields, "Length") {
		validLength, validLengthErr := audio.ValidLengthDefault()
		if !validLength {
			err = validLengthErr
			return
		}

		audioData.Length = audio.Length

	}
	if emptyOrContains(fields, "Size") {
		validSize, validSizeErr := audio.ValidSizeDefault()
		if !validSize {
			err = validSizeErr
			return
		}

		audioData.Size = audio.Size

	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := audio.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		audioData.DeviceTime = audio.DeviceTime

	}

	existAudioData, _ := boltdb.GetAudio(audio.ID)
	if existAudioData.ID != "" {
		err = errors.New(validationErrorPrefix + " audio with ID " + audio.ID + " already exists")
		return
	}
	audioData.Created = time.Now().UTC()
	audioData.Updated = time.Now().UTC()

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("audios"))

		audioJson, err := json.Marshal(audioData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(audio.ID), audioJson)
		return err
	})

	return
}

func (boltdb *DB) DeleteAudio(audioID string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "delete_audio_error:"

	validID, err := validator.UUID(audioID)
	if !validID {
		return
	}

	audioData, err := boltdb.GetAudio(audioID)
	if err != nil {
		return
	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("audios"))
		err = b.Delete([]byte(audioData.ID))

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

func (boltdb *DB) UpdateAudio(audio Audio, fields []string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "update_audio_error:"

	validID, err := validator.UUID(audio.ID)
	if !validID {
		return
	}

	audioData, err := boltdb.GetAudio(audio.ID)
	if err != nil {
		return
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := audio.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		audioData.ID = audio.ID

	}
	if emptyOrContains(fields, "FileType") {
		validFileType, validFileTypeErr := audio.ValidFileTypeDefault()
		if !validFileType {
			err = validFileTypeErr
			return
		}

		audioData.FileType = audio.FileType

	}
	if emptyOrContains(fields, "Length") {
		validLength, validLengthErr := audio.ValidLengthDefault()
		if !validLength {
			err = validLengthErr
			return
		}

		audioData.Length = audio.Length

	}
	if emptyOrContains(fields, "Size") {
		validSize, validSizeErr := audio.ValidSizeDefault()
		if !validSize {
			err = validSizeErr
			return
		}

		audioData.Size = audio.Size

	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := audio.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		audioData.DeviceTime = audio.DeviceTime

	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("audios"))

		audioJson, err := json.Marshal(audioData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(audioData.ID), audioJson)

		if err == nil {
			rowsAffected = 1
		}

		return err
	})

	return
}

func (boltdb *DB) GetAudioList(offset int, limit int, filters Filters, returnFields []string, sortBy SortBy) (results []Audio, totalResults int64, err error) {
	validationErrorPrefix := "get_user_error:"

	if !(filters.Operator == "AND" || filters.Operator == "OR") {
		err = errors.New(validationErrorPrefix + " filter operator error")
	}

	var audioList Audios

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("audios"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var audio Audio
			err := json.Unmarshal(v, &audio)

			includeThis, err := includeThisAudio(filters, audio)

			if err != nil {
				return err
			}

			if includeThis {
				resultAudio := Audio{ID: audio.ID}
				if emptyOrContains(returnFields, "ID") {
					resultAudio.ID = audio.ID
				}
				if emptyOrContains(returnFields, "FileType") {
					resultAudio.FileType = audio.FileType
				}
				if emptyOrContains(returnFields, "Length") {
					resultAudio.Length = audio.Length
				}
				if emptyOrContains(returnFields, "Size") {
					resultAudio.Size = audio.Size
				}
				if emptyOrContains(returnFields, "DeviceTime") {
					resultAudio.DeviceTime = audio.DeviceTime
				}
				if emptyOrContains(returnFields, "Created") {
					resultAudio.Created = audio.Created
				}
				if emptyOrContains(returnFields, "Updated") {
					resultAudio.Updated = audio.Updated
				}

				audioList = append(audioList, resultAudio)
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	if sortBy.Direction == "ASC" || sortBy.Direction == "DESC" {
		if sortBy.Field == "ID" && sortBy.Direction == "ASC" {
			sort.Sort(sortByAudioID{audioList})
		} else if sortBy.Field == "ID" && sortBy.Direction == "DESC" {
			sort.Sort(sortByAudioIDDesc{audioList})
		}
		if sortBy.Field == "FileType" && sortBy.Direction == "ASC" {
			sort.Sort(sortByAudioFileType{audioList})
		} else if sortBy.Field == "FileType" && sortBy.Direction == "DESC" {
			sort.Sort(sortByAudioFileTypeDesc{audioList})
		}
		if sortBy.Field == "Length" && sortBy.Direction == "ASC" {
			sort.Sort(sortByAudioLength{audioList})
		} else if sortBy.Field == "Length" && sortBy.Direction == "DESC" {
			sort.Sort(sortByAudioLengthDesc{audioList})
		}
		if sortBy.Field == "Size" && sortBy.Direction == "ASC" {
			sort.Sort(sortByAudioSize{audioList})
		} else if sortBy.Field == "Size" && sortBy.Direction == "DESC" {
			sort.Sort(sortByAudioSizeDesc{audioList})
		}
		if sortBy.Field == "DeviceTime" && sortBy.Direction == "ASC" {
			sort.Sort(sortByAudioDeviceTime{audioList})
		} else if sortBy.Field == "DeviceTime" && sortBy.Direction == "DESC" {
			sort.Sort(sortByAudioDeviceTimeDesc{audioList})
		}
		if sortBy.Field == "Created" && sortBy.Direction == "ASC" {
			sort.Sort(sortByAudioCreated{audioList})
		} else if sortBy.Field == "Created" && sortBy.Direction == "DESC" {
			sort.Sort(sortByAudioCreatedDesc{audioList})
		}
		if sortBy.Field == "Updated" && sortBy.Direction == "ASC" {
			sort.Sort(sortByAudioUpdated{audioList})
		} else if sortBy.Field == "Updated" && sortBy.Direction == "DESC" {
			sort.Sort(sortByAudioUpdatedDesc{audioList})
		}

	} else {
		err = errors.New(validationErrorPrefix + " sort Direction error")
	}

	totalResults = int64(len(audioList))

	for indexAudio, resultAudio := range audioList {
		if indexAudio >= offset && indexAudio < (offset+limit) {
			results = append(results, resultAudio)
		}
	}

	return
}
func includeThisAudio(filters Filters, audio Audio) (include bool, err error) {
	validationErrorPrefix := "get_audio_error:"

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

			if condition.Comparison == "=" && audio.ID == conditionValueID {
				meetConditionID = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueID, "%") && strings.HasSuffix(conditionValueID, "%") {
					if strings.Contains(audio.ID, strings.TrimSuffix(strings.TrimPrefix(conditionValueID, "%"), "%")) {
						meetConditionID = true
					}
				} else if strings.HasPrefix(conditionValueID, "%") {
					if strings.HasSuffix(audio.ID, strings.TrimPrefix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if strings.HasSuffix(conditionValueID, "%") {
					if strings.HasPrefix(audio.ID, strings.TrimSuffix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if audio.ID == conditionValueID {
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

			if condition.Comparison == "=" && audio.FileType == conditionValueFileType {
				meetConditionFileType = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueFileType, "%") && strings.HasSuffix(conditionValueFileType, "%") {
					if strings.Contains(audio.FileType, strings.TrimSuffix(strings.TrimPrefix(conditionValueFileType, "%"), "%")) {
						meetConditionFileType = true
					}
				} else if strings.HasPrefix(conditionValueFileType, "%") {
					if strings.HasSuffix(audio.FileType, strings.TrimPrefix(conditionValueFileType, "%")) {
						meetConditionFileType = true
					}
				} else if strings.HasSuffix(conditionValueFileType, "%") {
					if strings.HasPrefix(audio.FileType, strings.TrimSuffix(conditionValueFileType, "%")) {
						meetConditionFileType = true
					}
				} else if audio.FileType == conditionValueFileType {
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

		meetConditionLength := false

		if condition.Field == "Length" {
			conditionValueLength := condition.Value.(int)

			if condition.Comparison == "=" && audio.Length == conditionValueLength {
				meetConditionLength = true
			} else if condition.Comparison == ">" && audio.Length > conditionValueLength {
				meetConditionLength = true
			} else if condition.Comparison == "<" && audio.Length < conditionValueLength {
				meetConditionLength = true
			}

			if meetConditionLength {
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

			if condition.Comparison == "=" && audio.Size == conditionValueSize {
				meetConditionSize = true
			} else if condition.Comparison == ">" && audio.Size > conditionValueSize {
				meetConditionSize = true
			} else if condition.Comparison == "<" && audio.Size < conditionValueSize {
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

			if condition.Comparison == "=" && audio.DeviceTime == conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == ">" && audio.DeviceTime > conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == "<" && audio.DeviceTime < conditionValueDeviceTime {
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
			diffCreated := audio.Created.Sub(conditionValueCreated)

			if condition.Comparison == "=" && audio.Created == conditionValueCreated {
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
			diffUpdated := audio.Updated.Sub(conditionValueUpdated)

			if condition.Comparison == "=" && audio.Updated == conditionValueUpdated {
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

func (s Audios) Len() int {
	return len(s)
}
func (s Audios) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortByAudioID struct {
	Audios
}

func (s sortByAudioID) Less(i, j int) bool {
	return s.Audios[i].ID < s.Audios[j].ID
}

type sortByAudioIDDesc struct {
	Audios
}

func (s sortByAudioIDDesc) Less(i, j int) bool {
	return s.Audios[i].ID > s.Audios[j].ID

}

type sortByAudioFileType struct {
	Audios
}

func (s sortByAudioFileType) Less(i, j int) bool {
	return s.Audios[i].FileType < s.Audios[j].FileType
}

type sortByAudioFileTypeDesc struct {
	Audios
}

func (s sortByAudioFileTypeDesc) Less(i, j int) bool {
	return s.Audios[i].FileType > s.Audios[j].FileType

}

type sortByAudioLength struct {
	Audios
}

func (s sortByAudioLength) Less(i, j int) bool {
	return s.Audios[i].Length < s.Audios[j].Length
}

type sortByAudioLengthDesc struct {
	Audios
}

func (s sortByAudioLengthDesc) Less(i, j int) bool {
	return s.Audios[i].Length > s.Audios[j].Length

}

type sortByAudioSize struct {
	Audios
}

func (s sortByAudioSize) Less(i, j int) bool {
	return s.Audios[i].Size < s.Audios[j].Size
}

type sortByAudioSizeDesc struct {
	Audios
}

func (s sortByAudioSizeDesc) Less(i, j int) bool {
	return s.Audios[i].Size > s.Audios[j].Size

}

type sortByAudioDeviceTime struct {
	Audios
}

func (s sortByAudioDeviceTime) Less(i, j int) bool {
	return s.Audios[i].DeviceTime < s.Audios[j].DeviceTime
}

type sortByAudioDeviceTimeDesc struct {
	Audios
}

func (s sortByAudioDeviceTimeDesc) Less(i, j int) bool {
	return s.Audios[i].DeviceTime > s.Audios[j].DeviceTime

}

type sortByAudioCreated struct {
	Audios
}

func (s sortByAudioCreated) Less(i, j int) bool {
	diffLastModification := s.Audios[i].Created.Sub(s.Audios[j].Created)
	return diffLastModification < 0
}

type sortByAudioCreatedDesc struct {
	Audios
}

func (s sortByAudioCreatedDesc) Less(i, j int) bool {
	diffLastModification := s.Audios[i].Created.Sub(s.Audios[j].Created)
	return diffLastModification > 0

}

type sortByAudioUpdated struct {
	Audios
}

func (s sortByAudioUpdated) Less(i, j int) bool {
	diffLastModification := s.Audios[i].Updated.Sub(s.Audios[j].Updated)
	return diffLastModification < 0
}

type sortByAudioUpdatedDesc struct {
	Audios
}

func (s sortByAudioUpdatedDesc) Less(i, j int) bool {
	diffLastModification := s.Audios[i].Updated.Sub(s.Audios[j].Updated)
	return diffLastModification > 0

}

func (audio Audio) ValidIDDefault() (validField bool, err error) {
	validField, _ = validator.UUID(audio.ID)
	if !validField {
		err = errors.New("error_uuid__audio___ID")
		return
	}

	return
}
func (audio Audio) ValidFileTypeDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(audio.FileType, 100)
	if !validField {
		err = errors.New("error_maxlength__audio___FileType")
		return
	}

	return
}
func (audio Audio) ValidLengthDefault() (validField bool, err error) {
	validField = true

	return
}
func (audio Audio) ValidSizeDefault() (validField bool, err error) {
	validField = true

	return
}
func (audio Audio) ValidDeviceTimeDefault() (validField bool, err error) {
	validField = true

	return
}
func (audio Audio) ValidCreatedDefault() (validField bool, err error) {
	validField = true

	return
}
func (audio Audio) ValidUpdatedDefault() (validField bool, err error) {
	validField = true

	return
}
