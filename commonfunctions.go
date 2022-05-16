package commonfunctions

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"errors"
	"strconv"

	"github.com/gitstliu/log4go"
)

type TimeSpan struct {
	startNS int64
	endNS   int64
}

type NodeToRecord interface {
	ToRecord() string
}

func PanicHandler() {
	if r := recover(); r != nil {
		log4go.Error("Run time Error %v", r)
		panic(r)
	}
}

func InterfacesToStrings(src []interface{}) []string {
	result := []string{}
	for _, value := range src {
		result = append(result, value.(string))
	}
	return result
}

func (this *TimeSpan) Start() {
	this.startNS = time.Now().UnixNano()
}

func (this *TimeSpan) End() {
	this.endNS = time.Now().UnixNano()
}

func (this *TimeSpan) GetTimeSpanMS() float64 {

	return float64(this.endNS-this.startNS) / 1000000
}

func GetFilesWithFolder(folderName string) ([]string, error) {

	result := make([]string, 0, 1000)

	infos, readDirError := ioutil.ReadDir(folderName)

	if readDirError != nil {
		return nil, readDirError
	}

	for _, info := range infos {

		if !info.IsDir() {
			result = append(result, folderName+"/"+info.Name())
		}
	}

	return result, nil
}

func ObjectToJson(value interface{}) (string, error) {
	meta, err := json.Marshal(value)
	return string(meta), err
}

func ObjectsToJson(values []interface{}) ([]interface{}, error) {
	result := [](interface{}){}

	for _, currValue := range values {
		meta, err := ObjectToJson(currValue)

		if err != nil {
			return nil, err
		} else {
			result = append(result, meta)
		}
	}

	return result, nil
}

func DecodeGzipBytes(meta []byte) ([]byte, error) {
	b := bytes.Buffer{}
	b.Write(meta)
	r, _ := gzip.NewReader(&b)
	defer r.Close()
	datas, readErr := ioutil.ReadAll(r)

	if readErr != nil {
		return nil, readErr
	}

	return datas, nil
}

func EncodeGzipBytes(meta []byte) []byte {
	b := bytes.Buffer{}
	w := gzip.NewWriter(&b)
	defer w.Close()

	w.Write(meta)
	w.Flush()

	return b.Bytes()
}

func JsonToObject(meta string, result interface{}) error {
	return json.Unmarshal(StringToBytes(meta), result)
}

func StringToBytes(value string) []byte {
	return []byte(value)
}

func MetaToJsonContent(value string) string {
	return fmt.Sprintf("{%s}", value)
}

func Int64ToBytes(value int64) []byte {

	result := []byte{}

	buffer := bytes.NewBuffer(result)
	binary.Write(buffer, binary.BigEndian, &value)

	return buffer.Bytes()
}

func IntToBool(value int) bool {
	if value == 1 {
		return true
	} else {
		return false
	}
}

func IsGzipEncode(header http.Header) bool {

	log4go.Debug("header : %v", header)

	value := header["Content-Encoding"]

	return value != nil && strings.EqualFold(value[0], "gzip")
}

func ValsToStrs(values []interface{}) ([]string, error) {
	result := []string{}
	for index, value := range values {
		var key string
		if value == nil {
			return nil, errors.New(fmt.Sprintf("index %d is nil!!", index))
		}

		switch value.(type) {
		case float64:
			ft := value.(float64)
			key = strconv.FormatFloat(ft, 'f', -1, 64)
		case float32:
			ft := value.(float32)
			key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
		case int:
			it := value.(int)
			key = strconv.Itoa(it)
		case uint:
			it := value.(uint)
			key = strconv.Itoa(int(it))
		case int8:
			it := value.(int8)
			key = strconv.Itoa(int(it))
		case uint8:
			it := value.(uint8)
			key = strconv.Itoa(int(it))
		case int16:
			it := value.(int16)
			key = strconv.Itoa(int(it))
		case uint16:
			it := value.(uint16)
			key = strconv.Itoa(int(it))
		case int32:
			it := value.(int32)
			key = strconv.Itoa(int(it))
		case uint32:
			it := value.(uint32)
			key = strconv.Itoa(int(it))
		case int64:
			it := value.(int64)
			key = strconv.FormatInt(it, 10)
		case uint64:
			it := value.(uint64)
			key = strconv.FormatUint(it, 10)
		case string:
			key = value.(string)
		case []byte:
			key = string(value.([]byte))
		default:
			newValue, _ := json.Marshal(value)
			key = string(newValue)
		}

		result = append(result, key)
	}

	return result, nil

}

func StringInStrings(values []string, value string) int {
	result := -1

	for index, currValue := range values {
		if currValue == value {
			result = index
			return result
		}
	}
	return result
}
