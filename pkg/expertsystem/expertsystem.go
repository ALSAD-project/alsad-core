package expertsystem

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ExpertData for read_data
// ExpertData is the data to be exchanged through HTTP between client and daemon
//
// TODO: Enhance the below problem:
// (Too Lazy to handle different struct, so all data is included and every request)
//
// Pagination is for further rule-based / web-based development
// ---------------------------
type ExpertData struct {
	LinePerPage int
	Page        int

	Filename  string   `json:"filename"`
	Labels    []string `json:"labels"`
	LineCount int      `json:"lineCount"`
	Lines     []string `json:"lines"`
}

// ExpertData ends
// ---------------------------

// Config wrap the env var for daemon
// ---------------------------
type Config struct {
	Port    int
	SrcDir  string
	DestDir string
}

// Configure expertsystem configuration from env var
func Configure() (config Config, err error) {
	_port, err := strconv.Atoi(os.Getenv("REQUEST_PORT"))
	_srcDir := os.Getenv("SOURCE_DIRECTORY")
	_destDir := os.Getenv("DESTINATION_DIRECTORY")

	if (err != nil) || (_srcDir == "") || (_destDir == "") {
		// Error handling for invalid env var
		// ----
		// return config, errors.New("Fatal error invalid enviornment variables")
		// ----
		// Or update to default value by fakeConfigure()
		return fakeConfigure()
	}
	config.Port = _port
	config.SrcDir = _srcDir
	config.DestDir = _destDir
	return config, nil
}

// For Dev Purpose
func fakeConfigure() (config Config, err error) {
	config.Port = 4000
	config.SrcDir = "tmp/input/"
	config.DestDir = "tmp/output/"
	return config, nil
}

// Config ends
// ---------------------------

// Utils
// ----------------------------

// CheckFatal error handler
func CheckFatal(e error) {
	if e != nil {
		panic(e)
	}
}

// NewUUID for filenames
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// Utils end
// ----------------------------
