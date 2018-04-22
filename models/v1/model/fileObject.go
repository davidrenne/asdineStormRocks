package model

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/asdine/storm"
	"github.com/globalsign/mgo/bson"
	"log"
	"reflect"
	"sync"
	"time"
)

var FileObjects modelFileObjects

type modelFileObjects struct{}

var collectionFileObjectsMutex *sync.RWMutex

type FileObjectJoinItems struct {
	Count int           `json:"Count"`
	Items *[]FileObject `json:"Items"`
}

var GoCoreFileObjectsHasBootStrapped bool

func init() {
	collectionFileObjectsMutex = &sync.RWMutex{}

	//FileObjects.Index()
	go func() {
		time.Sleep(time.Second * 5)
		FileObjects.Bootstrap()
	}()
	store.RegisterStore(FileObjects)
}

func (self *FileObject) GetId() string {
	return self.Id.Hex()
}

type FileObject struct {
	Id             bson.ObjectId  `json:"Id" storm:"id"`
	Name           string         `json:"Name"`
	Path           string         `json:"Path"`
	SingleDownload bool           `json:"SingleDownload"`
	Content        string         `json:"Content"`
	Size           int            `json:"Size"`
	Type           string         `json:"Type"`
	ModifiedUnix   int            `json:"ModifiedUnix"`
	Modified       time.Time      `json:"Modified"`
	MD5            string         `json:"MD5"`
	AccountId      string         `json:"AccountId" storm:"index"`
	CreateDate     time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate     time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId   string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta  *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id             string `json:"Id"`
		Name           string `json:"Name"`
		Path           string `json:"Path"`
		SingleDownload string `json:"SingleDownload"`
		Content        string `json:"Content"`
		Size           string `json:"Size"`
		Type           string `json:"Type"`
		ModifiedUnix   string `json:"ModifiedUnix"`
		Modified       string `json:"Modified"`
		MD5            string `json:"MD5"`
		AccountId      string `json:"AccountId"`
	} `json:"Errors" bson:"-"`
}

func (self modelFileObjects) Single(field string, value interface{}) (retObj FileObject, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelFileObjects) Search(field string, value interface{}) (retObj []FileObject, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []FileObject{}
	}
	return
}

func (obj modelFileObjects) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []FileObject, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	return
}

func (obj modelFileObjects) All() (retObj []FileObject, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []FileObject{}
	}
	return
}

func (obj modelFileObjects) AllAdvanced(limit int, skip int) (retObj []FileObject, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	return
}

func (obj modelFileObjects) AllByIndex(index string) (retObj []FileObject, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []FileObject{}
	}
	return
}

func (obj modelFileObjects) AllByIndexAdvanced(index string, limit int, skip int) (retObj []FileObject, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	return
}

func (obj modelFileObjects) Range(min, max, field string) (retObj []FileObject, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []FileObject{}
	}
	return
}

func (obj modelFileObjects) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []FileObject, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FileObject{}
		}
		return
	}
	return
}

func (obj modelFileObjects) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj FileObject
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelFileObjects) NewByReflection() (value reflect.Value) {
	retObj := FileObject{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelFileObjects) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []FileObject
	q := obj.Query().Filter(filter)
	if len(inFilter) > 0 {
		q = q.In(inFilter)
	}
	if len(excludeFilter) > 0 {
		q = q.Exclude(excludeFilter)
	}
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.All(&retObj)
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelFileObjects) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionFileObjectsMutex.RLock()
		bootstrapped := GoCoreFileObjectsHasBootStrapped
		collectionFileObjectsMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("FileObjects has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "FileObjects"
	query.entityName = "FileObject"
	return query
}
func (obj modelFileObjects) Index() error {
	return dbServices.BoltDB.Init(&FileObject{})
}

func (obj modelFileObjects) BootStrapComplete() {
	collectionFileObjectsMutex.Lock()
	GoCoreFileObjectsHasBootStrapped = true
	collectionFileObjectsMutex.Unlock()
}
func (obj modelFileObjects) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of FileObjects Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []FileObject
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwpdCg=="

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "FileObjects")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for FileObjects due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("fileObjects", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for FileObjects: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for FileObjects: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []FileObject
	for _, file := range files {
		var fileBootstrap []FileObject
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"FileObjects", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for FileObjects: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for FileObjects: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"FileObjects", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "FileObjects")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting FileObjects", len(v))

	for _, doc := range v {
		var original FileObject
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "FileObjects" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for FileObjects:  " + doc.Id.Hex() + "  " + err.Error())
					isError = true
				}
			} else {
				valid := 0x01
				var reason map[string]bool
				reason = make(map[string]bool, 0)

				if doc.BootstrapMeta != nil && doc.BootstrapMeta.Version > 0 && doc.BootstrapMeta.Version <= serverSettings.WebConfig.Application.VersionNumeric {
					valid &= 0x00
					reason["Version Mismatch"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.Domain != "" && doc.BootstrapMeta.Domain != serverSettings.WebConfig.Application.ServerFQDN {
					valid &= 0x00
					reason["FQDN Mismatch With Domain"] = true
				}
				if doc.BootstrapMeta != nil && len(doc.BootstrapMeta.Domains) > 0 && !utils.InArray(serverSettings.WebConfig.Application.ServerFQDN, doc.BootstrapMeta.Domains) {
					valid &= 0x00
					reason["FQDN Mismatch With Domains"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.ProductName != "" && doc.BootstrapMeta.ProductName != serverSettings.WebConfig.Application.ProductName {
					valid &= 0x00
					reason["ProductName does not Match"] = true
				}
				if doc.BootstrapMeta != nil && len(doc.BootstrapMeta.ProductNames) > 0 && !utils.InArray(serverSettings.WebConfig.Application.ProductName, doc.BootstrapMeta.ProductNames) {
					valid &= 0x00
					reason["ProductNames does not Match Product"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.ReleaseMode != "" && doc.BootstrapMeta.ReleaseMode != serverSettings.WebConfig.Application.ReleaseMode {
					valid &= 0x00
					reason["ReleaseMode does not match"] = true
				}

				if valid == 0x01 {
					actualCount += 1
					err = doc.Save()
					if err != nil {
						log.Println("Failed to bootstrap data for FileObjects:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("FileObjects skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap FileObjects")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("fileObjects", "FileObjects")
		}
		if err == nil {
			log.Println("Successfully bootstrapped FileObjects")
			if actualCount != originalCount {
				logger.Message("FileObjects counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelFileObjects) New() *FileObject {
	return &FileObject{}
}

func (obj *FileObject) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *FileObject) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("FileObjects", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *FileObject) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *FileObject) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *FileObject) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("FileObjects", self.Id.Hex())
	return self.Save()
}

func (self *FileObject) ValidateAndClean() error {

	return validateFields(FileObject{}, self, reflect.ValueOf(self).Elem())
}

func (self *FileObject) Reflect() []Field {

	return Reflect(FileObject{})
}

func (self *FileObject) Delete() error {
	dbServices.CollectionCache{}.Remove("FileObjects", self.Id.Hex())
	return dbServices.BoltDB.Delete("FileObject", self.Id.Hex())
}

func (self *FileObject) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("FileObjects", self.Id.Hex())
	return dbServices.BoltDB.Delete("FileObjects", self.Id.Hex())
}

func (self *FileObject) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

	source := reflect.ValueOf(self).Elem()

	var joins []join
	joins, err = getJoins(source, remainingRecursions)

	if len(joins) == 0 {
		return
	}

	s := source
	for _, j := range joins {
		id := reflect.ValueOf(q.CheckForObjectId(s.FieldByName(j.joinFieldRefName).Interface())).String()
		joinsField := s.FieldByName("Joins")
		setField := joinsField.FieldByName(j.joinFieldName)

		endRecursion := false
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Print("Remaining Recursions")
			fmt.Println(fmt.Sprintf("%+v", remainingRecursions))
			fmt.Println(fmt.Sprintf("%+v", j.collectionName))
		}
		if remainingRecursions == j.joinSpecified {
			endRecursion = true
		}
		err = joinField(j, id, setField, j.joinSpecified, q, endRecursion, recursionCount)
		if err != nil {
			return
		}
	}
	return
}

func (self *FileObject) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FileObject) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *FileObject) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *FileObject) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *FileObject) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *FileObject) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelFileObjects) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Size":
		obj, ok := x.(int)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "ModifiedUnix":
		obj, ok := x.(int)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Modified":
		obj, ok := x.(time.Time)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "AccountId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Id":
		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Name":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Path":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Content":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "SingleDownload":
		obj, ok := x.(bool)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Type":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "MD5":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	}
	return
}

func (obj modelFileObjects) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "SingleDownload":
		if x == nil {
			var obj bool
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(bool)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Type":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "MD5":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "ModifiedUnix":
		if x == nil {
			var obj int
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(int)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Modified":
		if x == nil {
			var obj time.Time
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(time.Time)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "AccountId":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Id":
		if x == nil {
			var obj bson.ObjectId
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Name":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Path":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Content":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Size":
		if x == nil {
			var obj int
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(int)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	}
	return
}
