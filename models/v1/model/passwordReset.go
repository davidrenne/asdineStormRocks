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

var PasswordResets modelPasswordResets

type modelPasswordResets struct{}

var collectionPasswordResetsMutex *sync.RWMutex

type PasswordResetJoinItems struct {
	Count int              `json:"Count"`
	Items *[]PasswordReset `json:"Items"`
}

var GoCorePasswordResetsHasBootStrapped bool

func init() {
	collectionPasswordResetsMutex = &sync.RWMutex{}

	//PasswordResets.Index()
	go func() {
		time.Sleep(time.Second * 5)
		PasswordResets.Bootstrap()
	}()
	store.RegisterStore(PasswordResets)
}

func (self *PasswordReset) GetId() string {
	return self.Id.Hex()
}

type PasswordReset struct {
	Id            bson.ObjectId  `json:"Id" storm:"id"`
	UserId        string         `json:"UserId" validate:"true,,,,,,"`
	Complete      bool           `json:"Complete"`
	Url           string         `json:"Url"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id       string `json:"Id"`
		UserId   string `json:"UserId"`
		Complete string `json:"Complete"`
		Url      string `json:"Url"`
	} `json:"Errors" bson:"-"`
}

func (self modelPasswordResets) Single(field string, value interface{}) (retObj PasswordReset, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelPasswordResets) Search(field string, value interface{}) (retObj []PasswordReset, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []PasswordReset{}
	}
	return
}

func (obj modelPasswordResets) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []PasswordReset, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	return
}

func (obj modelPasswordResets) All() (retObj []PasswordReset, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []PasswordReset{}
	}
	return
}

func (obj modelPasswordResets) AllAdvanced(limit int, skip int) (retObj []PasswordReset, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	return
}

func (obj modelPasswordResets) AllByIndex(index string) (retObj []PasswordReset, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []PasswordReset{}
	}
	return
}

func (obj modelPasswordResets) AllByIndexAdvanced(index string, limit int, skip int) (retObj []PasswordReset, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	return
}

func (obj modelPasswordResets) Range(min, max, field string) (retObj []PasswordReset, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []PasswordReset{}
	}
	return
}

func (obj modelPasswordResets) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []PasswordReset, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []PasswordReset{}
		}
		return
	}
	return
}

func (obj modelPasswordResets) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj PasswordReset
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelPasswordResets) NewByReflection() (value reflect.Value) {
	retObj := PasswordReset{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelPasswordResets) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []PasswordReset
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

func (obj modelPasswordResets) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionPasswordResetsMutex.RLock()
		bootstrapped := GoCorePasswordResetsHasBootStrapped
		collectionPasswordResetsMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("PasswordResets has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "PasswordResets"
	query.entityName = "PasswordReset"
	return query
}
func (obj modelPasswordResets) Index() error {
	return dbServices.BoltDB.Init(&PasswordReset{})
}

func (obj modelPasswordResets) BootStrapComplete() {
	collectionPasswordResetsMutex.Lock()
	GoCorePasswordResetsHasBootStrapped = true
	collectionPasswordResetsMutex.Unlock()
}
func (obj modelPasswordResets) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of PasswordResets Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []PasswordReset
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := ""

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "PasswordResets")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for PasswordResets due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("passwordResets", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for PasswordResets: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for PasswordResets: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []PasswordReset
	for _, file := range files {
		var fileBootstrap []PasswordReset
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"PasswordResets", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for PasswordResets: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for PasswordResets: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"PasswordResets", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "PasswordResets")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting PasswordResets", len(v))

	for _, doc := range v {
		var original PasswordReset
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "PasswordResets" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for PasswordResets:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for PasswordResets:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("PasswordResets skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap PasswordResets")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("passwordResets", "PasswordResets")
		}
		if err == nil {
			log.Println("Successfully bootstrapped PasswordResets")
			if actualCount != originalCount {
				logger.Message("PasswordResets counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelPasswordResets) New() *PasswordReset {
	return &PasswordReset{}
}

func (obj *PasswordReset) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *PasswordReset) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("PasswordResets", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *PasswordReset) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *PasswordReset) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *PasswordReset) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("PasswordResets", self.Id.Hex())
	return self.Save()
}

func (self *PasswordReset) ValidateAndClean() error {

	return validateFields(PasswordReset{}, self, reflect.ValueOf(self).Elem())
}

func (self *PasswordReset) Reflect() []Field {

	return Reflect(PasswordReset{})
}

func (self *PasswordReset) Delete() error {
	dbServices.CollectionCache{}.Remove("PasswordResets", self.Id.Hex())
	return dbServices.BoltDB.Delete("PasswordReset", self.Id.Hex())
}

func (self *PasswordReset) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("PasswordResets", self.Id.Hex())
	return dbServices.BoltDB.Delete("PasswordResets", self.Id.Hex())
}

func (self *PasswordReset) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *PasswordReset) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *PasswordReset) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *PasswordReset) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *PasswordReset) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *PasswordReset) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *PasswordReset) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelPasswordResets) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Id":
		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "UserId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Complete":
		obj, ok := x.(bool)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Url":
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

func (obj modelPasswordResets) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Complete":
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
	case "Url":
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
	case "UserId":
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
