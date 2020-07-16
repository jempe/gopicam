package db

import (
	"testing"
	"time"

	"github.com/jempe/gopicam/pkg/validator"
)

func TestAudioDb(t *testing.T) {
	database, teardown := newTestDB(t)
	defer teardown()

	// Test if DB can be initialized
	t.Run("Init DB", func(t *testing.T) {

		err := database.InitDb()

		if err != nil {
			t.Errorf("error creating database")
		}
	})

	type audioTest struct {
		name    string
		input   Audio
		fields  []string
		wantErr bool
	}

	var tests []audioTest

	sampleAudio := Audio{FileType: "wav", Length: 200, Size: 213124, DeviceTime: 234234234}

	sampleAudioUpdate := Audio{FileType: "aif", Length: 300, Size: 234343124, DeviceTime: 7457534234234}

	tests = append(tests, audioTest{
		name:    "Blank",
		input:   Audio{},
		fields:  []string{},
		wantErr: false,
	})

	tests = append(tests, audioTest{
		name:    "Only File Type",
		input:   sampleAudio,
		fields:  []string{"FileType"},
		wantErr: false,
	})

	tests = append(tests, audioTest{
		name:    "Only Length",
		input:   sampleAudio,
		fields:  []string{"Length"},
		wantErr: false,
	})

	tests = append(tests, audioTest{
		name:    "Only Size",
		input:   sampleAudio,
		fields:  []string{"Size"},
		wantErr: false,
	})

	tests = append(tests, audioTest{
		name:    "Only DeviceTime",
		input:   sampleAudio,
		fields:  []string{"DeviceTime"},
		wantErr: false,
	})

	tests = append(tests, audioTest{
		name: "Long File Type",
		input: Audio{
			FileType: "jhfkghfkghfkgfluygflyugjknkjnkjjlhyvljhb/;kjn;kjb;jkbnḱjhjfhgvhj;hvghvhgcvgfcgfcfgcfchvjhlvhlghvhgvchjvlkhgvkhgfcfjgcfg",
		},
		fields:  []string{},
		wantErr: true,
	})

	// TODO tests with every type of field error

	for _, tt := range tests {
		// Test if DB file has been created
		t.Run(tt.name, func(t *testing.T) {
			startTime := time.Now()

			audioID, err := database.InsertAudio(tt.input, tt.fields)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("Error saving item")
				}
			} else {

				savedAudio, err := database.GetAudio(audioID)

				if err != nil {
					t.Errorf("Error getting item")
				}

				validID, idError := validator.UUID(savedAudio.ID)

				if !validID || idError != nil {
					t.Errorf("ID was not saved correctly")
				}

				if emptyOrContains(tt.fields, "FileType") {
					if savedAudio.FileType != tt.input.FileType {
						t.Errorf("FileType was not saved correctly")
					}
				} else {
					if savedAudio.FileType != "" {
						t.Errorf("FileType default value was not saved correctly")
					}
				}

				if emptyOrContains(tt.fields, "Length") {
					if savedAudio.Length != tt.input.Length {
						t.Errorf("Length was not saved correctly")
					}
				} else {
					if savedAudio.Length != 0 {
						t.Errorf("Length default value was not saved correctly")
					}
				}

				if emptyOrContains(tt.fields, "Size") {
					if savedAudio.Size != tt.input.Size {
						t.Errorf("Size was not saved correctly")
					}
				} else {
					if savedAudio.Size != 0 {
						t.Errorf("Size default value was not saved correctly")
					}
				}

				if emptyOrContains(tt.fields, "DeviceTime") {
					if savedAudio.DeviceTime != tt.input.DeviceTime {
						t.Errorf("DeviceTime was not saved correctly")
					}
				} else {
					if savedAudio.DeviceTime != 0 {
						t.Errorf("DeviceTime default value was not saved correctly")
					}
				}

				currentTime := time.Now()

				if savedAudio.Created != currentTime {
					if savedAudio.Created.After(currentTime) || savedAudio.Created.Before(startTime) {
						t.Errorf("Created was not saved correctly")
					}
				}

				if savedAudio.Updated != currentTime {
					if savedAudio.Updated.After(currentTime) || savedAudio.Updated.Before(startTime) {
						t.Errorf("Updated was not saved correctly")
					}
				}
			}

			sampleAudioUpdate.ID = audioID

			rowsAffected, err := database.UpdateAudio(sampleAudioUpdate, tt.fields)

			if err != nil || rowsAffected != 1 {
				if !tt.wantErr {
					t.Errorf("Error saving item")
				}
			} else {
				savedAudio, err := database.GetAudio(audioID)

				if err != nil {
					t.Errorf("Error getting item")
				}

				validID, idError := validator.UUID(savedAudio.ID)

				if !validID || idError != nil {
					t.Errorf("ID was not saved correctly")
				}

				if emptyOrContains(tt.fields, "FileType") {
					if savedAudio.FileType != sampleAudioUpdate.FileType {
						t.Errorf("FileType was not saved correctly")
					}
				} else {
					if savedAudio.FileType != "" {
						t.Errorf("FileType default value was not saved correctly")
					}
				}

				if emptyOrContains(tt.fields, "Length") {
					if savedAudio.Length != sampleAudioUpdate.Length {
						t.Errorf("Length was not saved correctly")
					}
				} else {
					if savedAudio.Length != 0 {
						t.Errorf("Length default value was not saved correctly")
					}
				}

				if emptyOrContains(tt.fields, "Size") {
					if savedAudio.Size != sampleAudioUpdate.Size {
						t.Errorf("Size was not saved correctly")
					}
				} else {
					if savedAudio.Size != 0 {
						t.Errorf("Size default value was not saved correctly")
					}
				}

				if emptyOrContains(tt.fields, "DeviceTime") {
					if savedAudio.DeviceTime != sampleAudioUpdate.DeviceTime {
						t.Errorf("DeviceTime was not saved correctly")
					}
				} else {
					if savedAudio.DeviceTime != 0 {
						t.Errorf("DeviceTime default value was not saved correctly")
					}
				}

				currentTime := time.Now()

				if !savedAudio.Created.Before(savedAudio.Updated) {
					t.Errorf("Created was not saved correctly")
				}

				if savedAudio.Updated != currentTime {
					if savedAudio.Updated.After(currentTime) || savedAudio.Updated.Before(startTime) {
						t.Errorf("Updated was not saved correctly")
					}
				}

			}
		})
	}
}
