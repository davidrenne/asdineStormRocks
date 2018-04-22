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

var ServerSettings modelServerSettings

type modelServerSettings struct{}

var collectionServerSettingsMutex *sync.RWMutex

type ServerSettingJoinItems struct {
	Count int              `json:"Count"`
	Items *[]ServerSetting `json:"Items"`
}

var GoCoreServerSettingsHasBootStrapped bool

func init() {
	collectionServerSettingsMutex = &sync.RWMutex{}

	//ServerSettings.Index()
	go func() {
		time.Sleep(time.Second * 5)
		ServerSettings.Bootstrap()
	}()
	store.RegisterStore(ServerSettings)
}

func (self *ServerSetting) GetId() string {
	return self.Id.Hex()
}

type ServerSetting struct {
	Id            bson.ObjectId  `json:"Id" storm:"id"`
	Name          string         `json:"Name"`
	Key           string         `json:"Key"`
	Category      string         `json:"Category"`
	Value         string         `json:"Value"`
	Any           interface{}    `json:"Any"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id       string `json:"Id"`
		Name     string `json:"Name"`
		Key      string `json:"Key"`
		Category string `json:"Category"`
		Value    string `json:"Value"`
		Any      string `json:"Any"`
	} `json:"Errors" bson:"-"`
}

func (self modelServerSettings) Single(field string, value interface{}) (retObj ServerSetting, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelServerSettings) Search(field string, value interface{}) (retObj []ServerSetting, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []ServerSetting{}
	}
	return
}

func (obj modelServerSettings) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []ServerSetting, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	return
}

func (obj modelServerSettings) All() (retObj []ServerSetting, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []ServerSetting{}
	}
	return
}

func (obj modelServerSettings) AllAdvanced(limit int, skip int) (retObj []ServerSetting, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	return
}

func (obj modelServerSettings) AllByIndex(index string) (retObj []ServerSetting, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []ServerSetting{}
	}
	return
}

func (obj modelServerSettings) AllByIndexAdvanced(index string, limit int, skip int) (retObj []ServerSetting, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	return
}

func (obj modelServerSettings) Range(min, max, field string) (retObj []ServerSetting, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []ServerSetting{}
	}
	return
}

func (obj modelServerSettings) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []ServerSetting, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ServerSetting{}
		}
		return
	}
	return
}

func (obj modelServerSettings) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj ServerSetting
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelServerSettings) NewByReflection() (value reflect.Value) {
	retObj := ServerSetting{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelServerSettings) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []ServerSetting
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

func (obj modelServerSettings) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionServerSettingsMutex.RLock()
		bootstrapped := GoCoreServerSettingsHasBootStrapped
		collectionServerSettingsMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("ServerSettings has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "ServerSettings"
	query.entityName = "ServerSetting"
	return query
}
func (obj modelServerSettings) Index() error {
	return dbServices.BoltDB.Init(&ServerSetting{})
}

func (obj modelServerSettings) BootStrapComplete() {
	collectionServerSettingsMutex.Lock()
	GoCoreServerSettingsHasBootStrapped = true
	collectionServerSettingsMutex.Unlock()
}
func (obj modelServerSettings) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of ServerSettings Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []ServerSetting
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwoJewoJCSJJZCI6ICI1N2ZiYTRiYjA0MTI0NTJkNmQ1MjQ4YzYiLAoJCSJOYW1lIjogIkxvY2tvdXQgQXR0ZW1wdHMiLAoJCSJLZXkiOiAibG9ja291dEF0dGVtcHRzIiwKCQkiQ2F0ZWdvcnkiOiAidXNlcnMiLAoJCSJWYWx1ZSI6ICIwIiwKCQkiQ3JlYXRlRGF0ZSI6IjIwMTYtMDgtMjZUMTA6NDk6MDQuNjMwNTM2NDQ2LTA0OjAwIiwKCQkiVXBkYXRlRGF0ZSI6IjIwMTYtMDgtMjZUMTA6NDk6MDQuNjMwNTM2NDQ2LTA0OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjoiNTdkOWIzODNkY2JhMGY1MTE3MmYxZjU3IiwKCQkiQm9vdHN0cmFwTWV0YSI6IHsKCQkJIkFsd2F5c1VwZGF0ZSI6IGZhbHNlCgkJfQoJfSwKCXsKCQkiSWQiOiAiNTk3ZTMxNTQ2MGU2NTdkOWI3MDU2M2FhIiwKCQkiTmFtZSIgOiAiVGltZSBab25lIiwKCQkiS2V5IiA6ICJUaW1lWm9uZSIsCgkJIkNhdGVnb3J5IiA6ICJ0aW1lU2V0dGluZ3MiLAoJCSJWYWx1ZSIgOiAiMCIsCgkJIkFueSIgOiAiMTkyLjE2OC4wLjEiLAoJCSJDcmVhdGVEYXRlIjoiMjAxNi0wOC0yNlQxMDo0OTowNC42MzA1MzY0NDYtMDQ6MDAiLAoJCSJVcGRhdGVEYXRlIjoiMjAxNi0wOC0yNlQxMDo0OTowNC42MzA1MzY0NDYtMDQ6MDAiLAoJCSJMYXN0VXBkYXRlSWQiOiI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciLAoJCSJCb290c3RyYXBNZXRhIjogewoJCQkiQWx3YXlzVXBkYXRlIjogZmFsc2UKCQl9Cgl9Cl0K"

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "ServerSettings")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for ServerSettings due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("serverSettings", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for ServerSettings: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for ServerSettings: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []ServerSetting
	for _, file := range files {
		var fileBootstrap []ServerSetting
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"ServerSettings", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for ServerSettings: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for ServerSettings: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"ServerSettings", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "ServerSettings")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting ServerSettings", len(v))

	for _, doc := range v {
		var original ServerSetting
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "ServerSettings" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for ServerSettings:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for ServerSettings:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("ServerSettings skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap ServerSettings")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("serverSettings", "ServerSettings")
		}
		if err == nil {
			log.Println("Successfully bootstrapped ServerSettings")
			if actualCount != originalCount {
				logger.Message("ServerSettings counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelServerSettings) New() *ServerSetting {
	return &ServerSetting{}
}

func (obj *ServerSetting) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *ServerSetting) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("ServerSettings", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *ServerSetting) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *ServerSetting) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *ServerSetting) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("ServerSettings", self.Id.Hex())
	return self.Save()
}

func (self *ServerSetting) ValidateAndClean() error {

	return validateFields(ServerSetting{}, self, reflect.ValueOf(self).Elem())
}

func (self *ServerSetting) Reflect() []Field {

	return Reflect(ServerSetting{})
}

func (self *ServerSetting) Delete() error {
	dbServices.CollectionCache{}.Remove("ServerSettings", self.Id.Hex())
	return dbServices.BoltDB.Delete("ServerSetting", self.Id.Hex())
}

func (self *ServerSetting) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("ServerSettings", self.Id.Hex())
	return dbServices.BoltDB.Delete("ServerSettings", self.Id.Hex())
}

func (self *ServerSetting) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *ServerSetting) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *ServerSetting) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *ServerSetting) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *ServerSetting) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *ServerSetting) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *ServerSetting) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelServerSettings) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Value":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Any":
		value = reflect.ValueOf(x)
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
	case "Key":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Category":
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

func (obj modelServerSettings) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Any":
		value = reflect.ValueOf(x)
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
	case "Key":
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
	case "Category":
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
	case "Value":
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
	}
	return
}
