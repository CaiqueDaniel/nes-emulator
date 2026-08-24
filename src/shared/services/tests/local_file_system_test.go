package shared_services_local_test

import (
	shared_services "nes-emu/src/shared/services"
	"os"
	"testing"
)

func Test_it_should_be_able_to_save_a_file(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()
	err := sut.SaveFile("./tmp/test.txt", []byte("Hello, World!"))

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	data, err := os.ReadFile("./tmp/test.txt")

	if err != nil {
		t.Errorf("Not able to read file")
	}

	if string(data) != "Hello, World!" {
		t.Errorf("File content does not match, got: %s", string(data))
	}

	os.RemoveAll("./tmp")
}

func Test_it_should_be_able_to_read_a_file(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()

	sut.SaveFile("./tmp/test.txt", []byte("Hello, World!"))
	data, err := sut.ReadFile("./tmp/test.txt")

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if string(data) != "Hello, World!" {
		t.Errorf("File content does not match, got: %s", string(data))
	}

	os.RemoveAll("./tmp")
}

func Test_it_should_check_if_file_exists(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()
	sut.SaveFile("./tmp/test.txt", []byte("Hello, World!"))

	exists := sut.ExistsFile("./tmp/test.txt")

	if !exists {
		t.Errorf("Expected file to exist, but it doesn't")
	}

	os.RemoveAll("./tmp")
}

func Test_it_should_copy_a_file(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()
	sut.SaveFile("./tmp/source.txt", []byte("Hello, World!"))
	err := sut.CopyFile("./tmp/source.txt", "./tmp/destination.txt")

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	data, err := sut.ReadFile("./tmp/destination.txt")

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
	if string(data) != "Hello, World!" {
		t.Errorf("File content does not match, got: %s", string(data))
	}
	os.RemoveAll("./tmp")
}

func Test_it_should_list_files_in_a_directory(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()
	sut.SaveFile("./tmp/test.txt", []byte("Hello, World!"))
	sut.SaveFile("./tmp/test2.txt", []byte("Hello, World!"))

	files, err := sut.ListFiles("./tmp")

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, but got %d", len(files))
	}

	if files[0] != "tmp/test.txt" {
		t.Errorf("Expected tmp/test.txt, but got %s", files[0])
	}

	if files[1] != "tmp/test2.txt" {
		t.Errorf("Expected tmp/test2.txt, but got %s", files[1])
	}

	os.RemoveAll("./tmp")
}

func Test_it_should_get_file_name(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()
	fileName := sut.GetFileName("./tmp/test.txt")

	if fileName != "test" {
		t.Errorf("Expected test, but got %s", fileName)
	}
}

func Test_it_should_get_file_extension(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()
	fileExtension := sut.GetFileExtension("./tmp/test.txt")

	if fileExtension != ".txt" {
		t.Errorf("Expected .txt, but got %s", fileExtension)
	}
}

func Test_it_should_get_file_md5_hash(t *testing.T) {
	sut := shared_services.NewLocalFileSystem()
	content := []byte("Hello, World!")
	// MD5 of "Hello, World!" is 65a8e27d8879283831b664bd8b7f0ad4
	expectedHash := "65a8e27d8879283831b664bd8b7f0ad4"

	sut.SaveFile("./tmp/test.txt", content)
	hash, err := sut.GetFileMD5Hash("./tmp/test.txt")

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if hash != expectedHash {
		t.Errorf("Expected %s, but got %s", expectedHash, hash)
	}

	// Test nonexistent file should return empty string and an error
	fallbackHash, err := sut.GetFileMD5Hash("./tmp/nonexistent.txt")

	if err == nil {
		t.Error("Expected an error for nonexistent file, but got nil")
	}

	if fallbackHash != "" {
		t.Errorf("Expected empty hash for nonexistent file, but got %s", fallbackHash)
	}

	os.RemoveAll("./tmp")
}
