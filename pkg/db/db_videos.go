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

type Video struct {
	ID         string    `json:"id"`
	FileType   string    `json:"file_type"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Length     int       `json:"length"`
	Size       int       `json:"size"`
	DeviceTime int64     `json:"device_time"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

type Videos []Video

func (boltdb *DB) GetVideo(videoID string) (video Video, err error) {
	validID, err := validator.UUID(videoID)
	if !validID {
		return video, err
	}

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("videos"))
		v := b.Get([]byte(videoID))

		if v == nil {
			return errors.New("video not found")
		}

		err := json.Unmarshal(v, &video)

		return err
	})

	return video, err
}

func (boltdb *DB) InsertVideo(video Video, fields []string) (videoID string, err error) {

	validationErrorPrefix := "insert_video_error:"

	id, err := uuid.NewRandom()

	if err != nil {
		log.Println(validationErrorPrefix, err)
		return
	}

	var videoData Video

	if video.ID == "" {
		videoID = id.String()

		video.ID = videoID
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := video.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		videoData.ID = video.ID

	}
	if emptyOrContains(fields, "FileType") {
		validFileType, validFileTypeErr := video.ValidFileTypeDefault()
		if !validFileType {
			err = validFileTypeErr
			return
		}

		videoData.FileType = video.FileType

	}
	if emptyOrContains(fields, "Width") {
		validWidth, validWidthErr := video.ValidWidthDefault()
		if !validWidth {
			err = validWidthErr
			return
		}

		videoData.Width = video.Width

	}
	if emptyOrContains(fields, "Height") {
		validHeight, validHeightErr := video.ValidHeightDefault()
		if !validHeight {
			err = validHeightErr
			return
		}

		videoData.Height = video.Height

	}
	if emptyOrContains(fields, "Length") {
		validLength, validLengthErr := video.ValidLengthDefault()
		if !validLength {
			err = validLengthErr
			return
		}

		videoData.Length = video.Length

	}
	if emptyOrContains(fields, "Size") {
		validSize, validSizeErr := video.ValidSizeDefault()
		if !validSize {
			err = validSizeErr
			return
		}

		videoData.Size = video.Size

	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := video.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		videoData.DeviceTime = video.DeviceTime

	}

	existVideoData, _ := boltdb.GetVideo(video.ID)
	if existVideoData.ID != "" {
		err = errors.New(validationErrorPrefix + " video with ID " + video.ID + " already exists")
		return
	}
	videoData.Created = time.Now().UTC()
	videoData.Updated = time.Now().UTC()

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("videos"))

		videoJson, err := json.Marshal(videoData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(video.ID), videoJson)
		return err
	})

	return
}

func (boltdb *DB) DeleteVideo(videoID string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "delete_video_error:"

	validID, err := validator.UUID(videoID)
	if !validID {
		return
	}

	videoData, err := boltdb.GetVideo(videoID)
	if err != nil {
		return
	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("videos"))
		err = b.Delete([]byte(videoData.ID))

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

func (boltdb *DB) UpdateVideo(video Video, fields []string) (rowsAffected int64, err error) {
	//validationErrorPrefix := "update_video_error:"

	validID, err := validator.UUID(video.ID)
	if !validID {
		return
	}

	videoData, err := boltdb.GetVideo(video.ID)
	if err != nil {
		return
	}

	if emptyOrContains(fields, "ID") {
		validID, validIDErr := video.ValidIDDefault()
		if !validID {
			err = validIDErr
			return
		}

		videoData.ID = video.ID

	}
	if emptyOrContains(fields, "FileType") {
		validFileType, validFileTypeErr := video.ValidFileTypeDefault()
		if !validFileType {
			err = validFileTypeErr
			return
		}

		videoData.FileType = video.FileType

	}
	if emptyOrContains(fields, "Width") {
		validWidth, validWidthErr := video.ValidWidthDefault()
		if !validWidth {
			err = validWidthErr
			return
		}

		videoData.Width = video.Width

	}
	if emptyOrContains(fields, "Height") {
		validHeight, validHeightErr := video.ValidHeightDefault()
		if !validHeight {
			err = validHeightErr
			return
		}

		videoData.Height = video.Height

	}
	if emptyOrContains(fields, "Length") {
		validLength, validLengthErr := video.ValidLengthDefault()
		if !validLength {
			err = validLengthErr
			return
		}

		videoData.Length = video.Length

	}
	if emptyOrContains(fields, "Size") {
		validSize, validSizeErr := video.ValidSizeDefault()
		if !validSize {
			err = validSizeErr
			return
		}

		videoData.Size = video.Size

	}
	if emptyOrContains(fields, "DeviceTime") {
		validDeviceTime, validDeviceTimeErr := video.ValidDeviceTimeDefault()
		if !validDeviceTime {
			err = validDeviceTimeErr
			return
		}

		videoData.DeviceTime = video.DeviceTime

	}

	err = boltdb.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("videos"))

		videoJson, err := json.Marshal(videoData)

		if err != nil {
			return err
		}

		err = b.Put([]byte(videoData.ID), videoJson)

		if err == nil {
			rowsAffected = 1
		}

		return err
	})

	return
}

func (boltdb *DB) GetVideoList(offset int, limit int, filters Filters, returnFields []string, sortBy SortBy) (results []Video, totalResults int64, err error) {
	validationErrorPrefix := "get_user_error:"

	if !(filters.Operator == "AND" || filters.Operator == "OR") {
		err = errors.New(validationErrorPrefix + " filter operator error")
	}

	var videoList Videos

	err = boltdb.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("videos"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var video Video
			err := json.Unmarshal(v, &video)

			includeThis, err := includeThisVideo(filters, video)

			if err != nil {
				return err
			}

			if includeThis {
				resultVideo := Video{ID: video.ID}
				if emptyOrContains(returnFields, "ID") {
					resultVideo.ID = video.ID
				}
				if emptyOrContains(returnFields, "FileType") {
					resultVideo.FileType = video.FileType
				}
				if emptyOrContains(returnFields, "Width") {
					resultVideo.Width = video.Width
				}
				if emptyOrContains(returnFields, "Height") {
					resultVideo.Height = video.Height
				}
				if emptyOrContains(returnFields, "Length") {
					resultVideo.Length = video.Length
				}
				if emptyOrContains(returnFields, "Size") {
					resultVideo.Size = video.Size
				}
				if emptyOrContains(returnFields, "DeviceTime") {
					resultVideo.DeviceTime = video.DeviceTime
				}
				if emptyOrContains(returnFields, "Created") {
					resultVideo.Created = video.Created
				}
				if emptyOrContains(returnFields, "Updated") {
					resultVideo.Updated = video.Updated
				}

				videoList = append(videoList, resultVideo)
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	if sortBy.Direction == "ASC" || sortBy.Direction == "DESC" {
		if sortBy.Field == "ID" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoID{videoList})
		} else if sortBy.Field == "ID" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoIDDesc{videoList})
		}
		if sortBy.Field == "FileType" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoFileType{videoList})
		} else if sortBy.Field == "FileType" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoFileTypeDesc{videoList})
		}
		if sortBy.Field == "Width" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoWidth{videoList})
		} else if sortBy.Field == "Width" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoWidthDesc{videoList})
		}
		if sortBy.Field == "Height" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoHeight{videoList})
		} else if sortBy.Field == "Height" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoHeightDesc{videoList})
		}
		if sortBy.Field == "Length" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoLength{videoList})
		} else if sortBy.Field == "Length" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoLengthDesc{videoList})
		}
		if sortBy.Field == "Size" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoSize{videoList})
		} else if sortBy.Field == "Size" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoSizeDesc{videoList})
		}
		if sortBy.Field == "DeviceTime" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoDeviceTime{videoList})
		} else if sortBy.Field == "DeviceTime" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoDeviceTimeDesc{videoList})
		}
		if sortBy.Field == "Created" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoCreated{videoList})
		} else if sortBy.Field == "Created" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoCreatedDesc{videoList})
		}
		if sortBy.Field == "Updated" && sortBy.Direction == "ASC" {
			sort.Sort(sortByVideoUpdated{videoList})
		} else if sortBy.Field == "Updated" && sortBy.Direction == "DESC" {
			sort.Sort(sortByVideoUpdatedDesc{videoList})
		}

	} else {
		err = errors.New(validationErrorPrefix + " sort Direction error")
	}

	totalResults = int64(len(videoList))

	for indexVideo, resultVideo := range videoList {
		if indexVideo >= offset && indexVideo < (offset+limit) {
			results = append(results, resultVideo)
		}
	}

	return
}
func includeThisVideo(filters Filters, video Video) (include bool, err error) {
	validationErrorPrefix := "get_video_error:"

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

			if condition.Comparison == "=" && video.ID == conditionValueID {
				meetConditionID = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueID, "%") && strings.HasSuffix(conditionValueID, "%") {
					if strings.Contains(video.ID, strings.TrimSuffix(strings.TrimPrefix(conditionValueID, "%"), "%")) {
						meetConditionID = true
					}
				} else if strings.HasPrefix(conditionValueID, "%") {
					if strings.HasSuffix(video.ID, strings.TrimPrefix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if strings.HasSuffix(conditionValueID, "%") {
					if strings.HasPrefix(video.ID, strings.TrimSuffix(conditionValueID, "%")) {
						meetConditionID = true
					}
				} else if video.ID == conditionValueID {
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

			if condition.Comparison == "=" && video.FileType == conditionValueFileType {
				meetConditionFileType = true
			} else if condition.Comparison == "LIKE" {
				if strings.HasPrefix(conditionValueFileType, "%") && strings.HasSuffix(conditionValueFileType, "%") {
					if strings.Contains(video.FileType, strings.TrimSuffix(strings.TrimPrefix(conditionValueFileType, "%"), "%")) {
						meetConditionFileType = true
					}
				} else if strings.HasPrefix(conditionValueFileType, "%") {
					if strings.HasSuffix(video.FileType, strings.TrimPrefix(conditionValueFileType, "%")) {
						meetConditionFileType = true
					}
				} else if strings.HasSuffix(conditionValueFileType, "%") {
					if strings.HasPrefix(video.FileType, strings.TrimSuffix(conditionValueFileType, "%")) {
						meetConditionFileType = true
					}
				} else if video.FileType == conditionValueFileType {
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

			if condition.Comparison == "=" && video.Width == conditionValueWidth {
				meetConditionWidth = true
			} else if condition.Comparison == ">" && video.Width > conditionValueWidth {
				meetConditionWidth = true
			} else if condition.Comparison == "<" && video.Width < conditionValueWidth {
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

			if condition.Comparison == "=" && video.Height == conditionValueHeight {
				meetConditionHeight = true
			} else if condition.Comparison == ">" && video.Height > conditionValueHeight {
				meetConditionHeight = true
			} else if condition.Comparison == "<" && video.Height < conditionValueHeight {
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

		meetConditionLength := false

		if condition.Field == "Length" {
			conditionValueLength := condition.Value.(int)

			if condition.Comparison == "=" && video.Length == conditionValueLength {
				meetConditionLength = true
			} else if condition.Comparison == ">" && video.Length > conditionValueLength {
				meetConditionLength = true
			} else if condition.Comparison == "<" && video.Length < conditionValueLength {
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

			if condition.Comparison == "=" && video.Size == conditionValueSize {
				meetConditionSize = true
			} else if condition.Comparison == ">" && video.Size > conditionValueSize {
				meetConditionSize = true
			} else if condition.Comparison == "<" && video.Size < conditionValueSize {
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

			if condition.Comparison == "=" && video.DeviceTime == conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == ">" && video.DeviceTime > conditionValueDeviceTime {
				meetConditionDeviceTime = true
			} else if condition.Comparison == "<" && video.DeviceTime < conditionValueDeviceTime {
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
			diffCreated := video.Created.Sub(conditionValueCreated)

			if condition.Comparison == "=" && video.Created == conditionValueCreated {
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
			diffUpdated := video.Updated.Sub(conditionValueUpdated)

			if condition.Comparison == "=" && video.Updated == conditionValueUpdated {
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

func (s Videos) Len() int {
	return len(s)
}
func (s Videos) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type sortByVideoID struct {
	Videos
}

func (s sortByVideoID) Less(i, j int) bool {
	return s.Videos[i].ID < s.Videos[j].ID
}

type sortByVideoIDDesc struct {
	Videos
}

func (s sortByVideoIDDesc) Less(i, j int) bool {
	return s.Videos[i].ID > s.Videos[j].ID

}

type sortByVideoFileType struct {
	Videos
}

func (s sortByVideoFileType) Less(i, j int) bool {
	return s.Videos[i].FileType < s.Videos[j].FileType
}

type sortByVideoFileTypeDesc struct {
	Videos
}

func (s sortByVideoFileTypeDesc) Less(i, j int) bool {
	return s.Videos[i].FileType > s.Videos[j].FileType

}

type sortByVideoWidth struct {
	Videos
}

func (s sortByVideoWidth) Less(i, j int) bool {
	return s.Videos[i].Width < s.Videos[j].Width
}

type sortByVideoWidthDesc struct {
	Videos
}

func (s sortByVideoWidthDesc) Less(i, j int) bool {
	return s.Videos[i].Width > s.Videos[j].Width

}

type sortByVideoHeight struct {
	Videos
}

func (s sortByVideoHeight) Less(i, j int) bool {
	return s.Videos[i].Height < s.Videos[j].Height
}

type sortByVideoHeightDesc struct {
	Videos
}

func (s sortByVideoHeightDesc) Less(i, j int) bool {
	return s.Videos[i].Height > s.Videos[j].Height

}

type sortByVideoLength struct {
	Videos
}

func (s sortByVideoLength) Less(i, j int) bool {
	return s.Videos[i].Length < s.Videos[j].Length
}

type sortByVideoLengthDesc struct {
	Videos
}

func (s sortByVideoLengthDesc) Less(i, j int) bool {
	return s.Videos[i].Length > s.Videos[j].Length

}

type sortByVideoSize struct {
	Videos
}

func (s sortByVideoSize) Less(i, j int) bool {
	return s.Videos[i].Size < s.Videos[j].Size
}

type sortByVideoSizeDesc struct {
	Videos
}

func (s sortByVideoSizeDesc) Less(i, j int) bool {
	return s.Videos[i].Size > s.Videos[j].Size

}

type sortByVideoDeviceTime struct {
	Videos
}

func (s sortByVideoDeviceTime) Less(i, j int) bool {
	return s.Videos[i].DeviceTime < s.Videos[j].DeviceTime
}

type sortByVideoDeviceTimeDesc struct {
	Videos
}

func (s sortByVideoDeviceTimeDesc) Less(i, j int) bool {
	return s.Videos[i].DeviceTime > s.Videos[j].DeviceTime

}

type sortByVideoCreated struct {
	Videos
}

func (s sortByVideoCreated) Less(i, j int) bool {
	diffLastModification := s.Videos[i].Created.Sub(s.Videos[j].Created)
	return diffLastModification < 0
}

type sortByVideoCreatedDesc struct {
	Videos
}

func (s sortByVideoCreatedDesc) Less(i, j int) bool {
	diffLastModification := s.Videos[i].Created.Sub(s.Videos[j].Created)
	return diffLastModification > 0

}

type sortByVideoUpdated struct {
	Videos
}

func (s sortByVideoUpdated) Less(i, j int) bool {
	diffLastModification := s.Videos[i].Updated.Sub(s.Videos[j].Updated)
	return diffLastModification < 0
}

type sortByVideoUpdatedDesc struct {
	Videos
}

func (s sortByVideoUpdatedDesc) Less(i, j int) bool {
	diffLastModification := s.Videos[i].Updated.Sub(s.Videos[j].Updated)
	return diffLastModification > 0

}

func (video Video) ValidIDDefault() (validField bool, err error) {
	validField, _ = validator.UUID(video.ID)
	if !validField {
		err = errors.New("error_uuid__video___ID")
		return
	}

	return
}
func (video Video) ValidFileTypeDefault() (validField bool, err error) {
	validField, _ = validator.MaxLength(video.FileType, 100)
	if !validField {
		err = errors.New("error_maxlength__video___FileType")
		return
	}

	return
}
func (video Video) ValidWidthDefault() (validField bool, err error) {
	validField = true

	return
}
func (video Video) ValidHeightDefault() (validField bool, err error) {
	validField = true

	return
}
func (video Video) ValidLengthDefault() (validField bool, err error) {
	validField = true

	return
}
func (video Video) ValidSizeDefault() (validField bool, err error) {
	validField = true

	return
}
func (video Video) ValidDeviceTimeDefault() (validField bool, err error) {
	validField = true

	return
}
func (video Video) ValidCreatedDefault() (validField bool, err error) {
	validField = true

	return
}
func (video Video) ValidUpdatedDefault() (validField bool, err error) {
	validField = true

	return
}
