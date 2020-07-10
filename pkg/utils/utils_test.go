package utils

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestExists(t *testing.T) {
	type existStruct struct {
		name string
		file string
		want bool
	}

	tests := []existStruct{
		{
			name: "Empty",
			file: "",
			want: false,
		},
		{
			name: "Random file",
			file: "/tmp/sdfasdfasgasdfs897986yasdfsdfwerwer2efsdfsdfsdfsd",
			want: false,
		},
	}

	parentDir := os.TempDir()
	testDir, err := ioutil.TempDir(parentDir, "gopicam-test-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(testDir) // clean up

	log.Println("created temp directory:", testDir)

	tests = append(tests, existStruct{name: "Temp Directory", file: testDir, want: true})

	testFile, err := ioutil.TempFile("", "gopicam-test.*.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(testFile.Name()) // clean up

	log.Println("created temp file:", testFile.Name())

	tests = append(tests, existStruct{name: "Temp File", file: testFile.Name(), want: true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := Exists(tt.file)

			if exists != tt.want {
				t.Errorf("want %t; got %t", tt.want, exists)
			}
		})
	}
}

func TestIsDirectory(t *testing.T) {
	type isDirStruct struct {
		name string
		file string
		want bool
	}

	tests := []isDirStruct{
		{
			name: "Empty",
			file: "",
			want: false,
		},
		{
			name: "Random file",
			file: "/tmp/sdfasdfasgasdfs897986yasdfsdfwerwer2efsdfsdfsdfsd",
			want: false,
		},
	}

	testFile, err := ioutil.TempFile("", "gopicam-test.*.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(testFile.Name()) // clean up

	log.Println("created temp file:", testFile.Name())

	tests = append(tests, isDirStruct{name: "Temp File", file: testFile.Name(), want: false})

	parentDir := os.TempDir()
	testDir, err := ioutil.TempDir(parentDir, "gopicam-test-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(testDir) // clean up

	log.Println("created temp directory:", testDir)

	tests = append(tests, isDirStruct{name: "Temp Directory", file: testDir, want: true})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := IsDirectory(tt.file)

			if exists != tt.want {
				t.Errorf("want %t; got %t", tt.want, exists)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name        string
		fileContent []byte
		want        string
	}{
		{
			name:        "Empty",
			fileContent: []byte{},
			want:        "data:text/plain; charset=utf-8;base64,",
		},
		{
			name:        "PNG 1x1",
			fileContent: []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 8, 2, 0, 0, 0, 144, 119, 83, 222, 0, 0, 1, 133, 105, 67, 67, 80, 73, 67, 67, 32, 112, 114, 111, 102, 105, 108, 101, 0, 0, 40, 145, 125, 145, 61, 72, 195, 64, 28, 197, 95, 91, 165, 42, 21, 7, 139, 72, 113, 200, 80, 29, 196, 130, 168, 136, 163, 86, 161, 8, 21, 66, 173, 208, 170, 131, 201, 165, 95, 208, 164, 33, 73, 113, 113, 20, 92, 11, 14, 126, 44, 86, 29, 92, 156, 117, 117, 112, 21, 4, 193, 15, 16, 39, 71, 39, 69, 23, 41, 241, 127, 73, 161, 69, 140, 7, 199, 253, 120, 119, 239, 113, 247, 14, 240, 215, 203, 76, 53, 59, 198, 1, 85, 179, 140, 84, 34, 46, 100, 178, 171, 66, 240, 21, 65, 116, 99, 0, 17, 140, 74, 204, 212, 231, 68, 49, 9, 207, 241, 117, 15, 31, 95, 239, 98, 60, 203, 251, 220, 159, 163, 87, 201, 153, 12, 240, 9, 196, 179, 76, 55, 44, 226, 13, 226, 233, 77, 75, 231, 188, 79, 28, 102, 69, 73, 33, 62, 39, 30, 51, 232, 130, 196, 143, 92, 151, 93, 126, 227, 92, 112, 216, 207, 51, 195, 70, 58, 53, 79, 28, 38, 22, 10, 109, 44, 183, 49, 43, 26, 42, 241, 20, 113, 84, 81, 53, 202, 247, 103, 92, 86, 56, 111, 113, 86, 203, 85, 214, 188, 39, 127, 97, 40, 167, 173, 44, 115, 157, 230, 16, 18, 88, 196, 18, 68, 8, 144, 81, 69, 9, 101, 88, 136, 209, 170, 145, 98, 34, 69, 251, 113, 15, 127, 196, 241, 139, 228, 146, 201, 85, 2, 35, 199, 2, 42, 80, 33, 57, 126, 240, 63, 248, 221, 173, 153, 159, 156, 112, 147, 66, 113, 160, 243, 197, 182, 63, 134, 129, 224, 46, 208, 168, 217, 246, 247, 177, 109, 55, 78, 128, 192, 51, 112, 165, 181, 252, 149, 58, 48, 243, 73, 122, 173, 165, 69, 143, 128, 190, 109, 224, 226, 186, 165, 201, 123, 192, 229, 14, 48, 248, 164, 75, 134, 228, 72, 1, 154, 254, 124, 30, 120, 63, 163, 111, 202, 2, 253, 183, 64, 207, 154, 219, 91, 115, 31, 167, 15, 64, 154, 186, 74, 222, 0, 7, 135, 192, 72, 129, 178, 215, 61, 222, 221, 213, 222, 219, 191, 103, 154, 253, 253, 0, 114, 219, 114, 167, 141, 105, 10, 186, 0, 0, 0, 9, 112, 72, 89, 115, 0, 0, 11, 19, 0, 0, 11, 19, 1, 0, 154, 156, 24, 0, 0, 0, 7, 116, 73, 77, 69, 7, 228, 7, 10, 0, 35, 53, 116, 134, 124, 2, 0, 0, 0, 12, 73, 68, 65, 84, 8, 215, 99, 248, 255, 255, 63, 0, 5, 254, 2, 254, 220, 204, 89, 231, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130},
			want:        "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAIAAACQd1PeAAABhWlDQ1BJQ0MgcHJvZmlsZQAAKJF9kT1Iw0AcxV9bpSoVB4tIcchQHcSCqIijVqEIFUKt0KqDyaVf0KQhSXFxFFwLDn4sVh1cnHV1cBUEwQ8QJ0cnRRcp8X9JoUWMB8f9eHfvcfcO8NfLTDU7xgFVs4xUIi5ksqtC8BVBdGMAEYxKzNTnRDEJz/F1Dx9f72I8y/vcn6NXyZkM8AnEs0w3LOIN4ulNS+e8TxxmRUkhPiceM+iCxI9cl11+41xw2M8zw0Y6NU8cJhYKbSy3MSsaKvEUcVRRNcr3Z1xWOG9xVstV1rwnf2Eop60sc53mEBJYxBJECJBRRQllWIjRqpFiIkX7cQ9/xPGL5JLJVQIjxwIqUCE5fvA/+N2tmZ+ccJNCcaDzxbY/hoHgLtCo2fb3sW03ToDAM3CltfyVOjDzSXqtpUWPgL5t4OK6pcl7wOUOMPikS4bkSAGa/nweeD+jb8oC/bdAz5rbW3Mfpw9AmrpK3gAHh8BIgbLXPd7d1d7bv2ea/f0Acttyp41pCroAAAAJcEhZcwAACxMAAAsTAQCanBgAAAAHdElNRQfkBwoAIzV0hnwCAAAADElEQVQI12P4//8/AAX+Av7czFnnAAAAAElFTkSuQmCC",
		},
		{
			name:        "JPG 1x1",
			fileContent: []byte{255, 216, 255, 224, 0, 16, 74, 70, 73, 70, 0, 1, 1, 1, 0, 72, 0, 72, 0, 0, 255, 226, 2, 176, 73, 67, 67, 95, 80, 82, 79, 70, 73, 76, 69, 0, 1, 1, 0, 0, 2, 160, 108, 99, 109, 115, 4, 48, 0, 0, 109, 110, 116, 114, 82, 71, 66, 32, 88, 89, 90, 32, 7, 228, 0, 7, 0, 9, 0, 22, 0, 26, 0, 42, 97, 99, 115, 112, 65, 80, 80, 76, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 246, 214, 0, 1, 0, 0, 0, 0, 211, 45, 108, 99, 109, 115, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 13, 100, 101, 115, 99, 0, 0, 1, 32, 0, 0, 0, 64, 99, 112, 114, 116, 0, 0, 1, 96, 0, 0, 0, 54, 119, 116, 112, 116, 0, 0, 1, 152, 0, 0, 0, 20, 99, 104, 97, 100, 0, 0, 1, 172, 0, 0, 0, 44, 114, 88, 89, 90, 0, 0, 1, 216, 0, 0, 0, 20, 98, 88, 89, 90, 0, 0, 1, 236, 0, 0, 0, 20, 103, 88, 89, 90, 0, 0, 2, 0, 0, 0, 0, 20, 114, 84, 82, 67, 0, 0, 2, 20, 0, 0, 0, 32, 103, 84, 82, 67, 0, 0, 2, 20, 0, 0, 0, 32, 98, 84, 82, 67, 0, 0, 2, 20, 0, 0, 0, 32, 99, 104, 114, 109, 0, 0, 2, 52, 0, 0, 0, 36, 100, 109, 110, 100, 0, 0, 2, 88, 0, 0, 0, 36, 100, 109, 100, 100, 0, 0, 2, 124, 0, 0, 0, 36, 109, 108, 117, 99, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 12, 101, 110, 85, 83, 0, 0, 0, 36, 0, 0, 0, 28, 0, 71, 0, 73, 0, 77, 0, 80, 0, 32, 0, 98, 0, 117, 0, 105, 0, 108, 0, 116, 0, 45, 0, 105, 0, 110, 0, 32, 0, 115, 0, 82, 0, 71, 0, 66, 109, 108, 117, 99, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 12, 101, 110, 85, 83, 0, 0, 0, 26, 0, 0, 0, 28, 0, 80, 0, 117, 0, 98, 0, 108, 0, 105, 0, 99, 0, 32, 0, 68, 0, 111, 0, 109, 0, 97, 0, 105, 0, 110, 0, 0, 88, 89, 90, 32, 0, 0, 0, 0, 0, 0, 246, 214, 0, 1, 0, 0, 0, 0, 211, 45, 115, 102, 51, 50, 0, 0, 0, 0, 0, 1, 12, 66, 0, 0, 5, 222, 255, 255, 243, 37, 0, 0, 7, 147, 0, 0, 253, 144, 255, 255, 251, 161, 255, 255, 253, 162, 0, 0, 3, 220, 0, 0, 192, 110, 88, 89, 90, 32, 0, 0, 0, 0, 0, 0, 111, 160, 0, 0, 56, 245, 0, 0, 3, 144, 88, 89, 90, 32, 0, 0, 0, 0, 0, 0, 36, 159, 0, 0, 15, 132, 0, 0, 182, 196, 88, 89, 90, 32, 0, 0, 0, 0, 0, 0, 98, 151, 0, 0, 183, 135, 0, 0, 24, 217, 112, 97, 114, 97, 0, 0, 0, 0, 0, 3, 0, 0, 0, 2, 102, 102, 0, 0, 242, 167, 0, 0, 13, 89, 0, 0, 19, 208, 0, 0, 10, 91, 99, 104, 114, 109, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 163, 215, 0, 0, 84, 124, 0, 0, 76, 205, 0, 0, 153, 154, 0, 0, 38, 103, 0, 0, 15, 92, 109, 108, 117, 99, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 12, 101, 110, 85, 83, 0, 0, 0, 8, 0, 0, 0, 28, 0, 71, 0, 73, 0, 77, 0, 80, 109, 108, 117, 99, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 12, 101, 110, 85, 83, 0, 0, 0, 8, 0, 0, 0, 28, 0, 115, 0, 82, 0, 71, 0, 66, 255, 219, 0, 67, 0, 3, 2, 2, 3, 2, 2, 3, 3, 3, 3, 4, 3, 3, 4, 5, 8, 5, 5, 4, 4, 5, 10, 7, 7, 6, 8, 12, 10, 12, 12, 11, 10, 11, 11, 13, 14, 18, 16, 13, 14, 17, 14, 11, 11, 16, 22, 16, 17, 19, 20, 21, 21, 21, 12, 15, 23, 24, 22, 20, 24, 18, 20, 21, 20, 255, 219, 0, 67, 1, 3, 4, 4, 5, 4, 5, 9, 5, 5, 9, 20, 13, 11, 13, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 255, 194, 0, 17, 8, 0, 1, 0, 1, 3, 1, 17, 0, 2, 17, 1, 3, 17, 1, 255, 196, 0, 20, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 255, 196, 0, 20, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 12, 3, 1, 0, 2, 16, 3, 16, 0, 0, 1, 84, 159, 255, 196, 0, 20, 16, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 1, 0, 1, 5, 2, 127, 255, 196, 0, 20, 17, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 3, 1, 1, 63, 1, 127, 255, 196, 0, 20, 17, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 2, 1, 1, 63, 1, 127, 255, 196, 0, 20, 16, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 1, 0, 6, 63, 2, 127, 255, 196, 0, 20, 16, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 1, 0, 1, 63, 33, 127, 255, 218, 0, 12, 3, 1, 0, 2, 0, 3, 0, 0, 0, 16, 159, 255, 196, 0, 20, 17, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 3, 1, 1, 63, 16, 127, 255, 196, 0, 20, 17, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 2, 1, 1, 63, 16, 127, 255, 196, 0, 20, 16, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 218, 0, 8, 1, 1, 0, 1, 63, 16, 127, 255, 217},
			want:        "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/4gKwSUNDX1BST0ZJTEUAAQEAAAKgbGNtcwQwAABtbnRyUkdCIFhZWiAH5AAHAAkAFgAaACphY3NwQVBQTAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA9tYAAQAAAADTLWxjbXMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA1kZXNjAAABIAAAAEBjcHJ0AAABYAAAADZ3dHB0AAABmAAAABRjaGFkAAABrAAAACxyWFlaAAAB2AAAABRiWFlaAAAB7AAAABRnWFlaAAACAAAAABRyVFJDAAACFAAAACBnVFJDAAACFAAAACBiVFJDAAACFAAAACBjaHJtAAACNAAAACRkbW5kAAACWAAAACRkbWRkAAACfAAAACRtbHVjAAAAAAAAAAEAAAAMZW5VUwAAACQAAAAcAEcASQBNAFAAIABiAHUAaQBsAHQALQBpAG4AIABzAFIARwBCbWx1YwAAAAAAAAABAAAADGVuVVMAAAAaAAAAHABQAHUAYgBsAGkAYwAgAEQAbwBtAGEAaQBuAABYWVogAAAAAAAA9tYAAQAAAADTLXNmMzIAAAAAAAEMQgAABd7///MlAAAHkwAA/ZD///uh///9ogAAA9wAAMBuWFlaIAAAAAAAAG+gAAA49QAAA5BYWVogAAAAAAAAJJ8AAA+EAAC2xFhZWiAAAAAAAABilwAAt4cAABjZcGFyYQAAAAAAAwAAAAJmZgAA8qcAAA1ZAAAT0AAACltjaHJtAAAAAAADAAAAAKPXAABUfAAATM0AAJmaAAAmZwAAD1xtbHVjAAAAAAAAAAEAAAAMZW5VUwAAAAgAAAAcAEcASQBNAFBtbHVjAAAAAAAAAAEAAAAMZW5VUwAAAAgAAAAcAHMAUgBHAEL/2wBDAAMCAgMCAgMDAwMEAwMEBQgFBQQEBQoHBwYIDAoMDAsKCwsNDhIQDQ4RDgsLEBYQERMUFRUVDA8XGBYUGBIUFRT/2wBDAQMEBAUEBQkFBQkUDQsNFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBQUFBT/wgARCAABAAEDAREAAhEBAxEB/8QAFAABAAAAAAAAAAAAAAAAAAAACP/EABQBAQAAAAAAAAAAAAAAAAAAAAD/2gAMAwEAAhADEAAAAVSf/8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABBQJ//8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAgBAwEBPwF//8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAgBAgEBPwF//8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQAGPwJ//8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABPyF//9oADAMBAAIAAwAAABCf/8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAgBAwEBPxB//8QAFBEBAAAAAAAAAAAAAAAAAAAAAP/aAAgBAgEBPxB//8QAFBABAAAAAAAAAAAAAAAAAAAAAP/aAAgBAQABPxB//9k=",
		},
		{
			name:        "GIF 1x1",
			fileContent: []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 255, 255, 255, 33, 254, 17, 67, 114, 101, 97, 116, 101, 100, 32, 119, 105, 116, 104, 32, 71, 73, 77, 80, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59},
			want:        "data:image/gif;base64,R0lGODlhAQABAIAAAP///////yH+EUNyZWF0ZWQgd2l0aCBHSU1QACwAAAAAAQABAAACAkQBADs=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base64String := Base64Encode(tt.fileContent)

			if base64String != tt.want {
				t.Errorf("want \"%s\"; got \"%s\"", tt.want, base64String)
			}
		})
	}

}
