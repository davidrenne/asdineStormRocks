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

var Roles modelRoles

type modelRoles struct{}

var collectionRolesMutex *sync.RWMutex

type RoleJoinItems struct {
	Count int     `json:"Count"`
	Items *[]Role `json:"Items"`
}

var GoCoreRolesHasBootStrapped bool

func init() {
	collectionRolesMutex = &sync.RWMutex{}

	//Roles.Index()
	go func() {
		time.Sleep(time.Second * 5)
		Roles.Bootstrap()
	}()
	store.RegisterStore(Roles)
}

func (self *Role) GetId() string {
	return self.Id.Hex()
}

type Role struct {
	Id            bson.ObjectId  `json:"Id" storm:"id"`
	Name          string         `json:"Name" validate:"true,,,,,,"`
	AccountId     string         `json:"AccountId"`
	CanDelete     bool           `json:"CanDelete"`
	AccountType   string         `json:"AccountType" validate:"true,,,,,,"`
	ShortName     string         `json:"ShortName"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id          string `json:"Id"`
		Name        string `json:"Name"`
		AccountId   string `json:"AccountId"`
		CanDelete   string `json:"CanDelete"`
		AccountType string `json:"AccountType"`
		ShortName   string `json:"ShortName"`
	} `json:"Errors" bson:"-"`

	Views struct {
		UpdateDate    string `json:"UpdateDate" ref:"UpdateDate~DateTime"`
		UpdateFromNow string `json:"UpdateFromNow" ref:"UpdateDate~TimeFromNow"`
	} `json:"Views" bson:"-"`

	Joins struct {
		RoleFeatures   *RoleFeatureJoinItems `json:"RoleFeatures,omitempty" join:"RoleFeatures,RoleFeature,Id,true,RoleId"`
		LastUpdateUser *User                 `json:"LastUpdateUser,omitempty" join:"Users,User,LastUpdateId,false,"`
	} `json:"Joins" bson:"-"`
}

func (self modelRoles) Single(field string, value interface{}) (retObj Role, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelRoles) Search(field string, value interface{}) (retObj []Role, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []Role{}
	}
	return
}

func (obj modelRoles) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []Role, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	return
}

func (obj modelRoles) All() (retObj []Role, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []Role{}
	}
	return
}

func (obj modelRoles) AllAdvanced(limit int, skip int) (retObj []Role, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	return
}

func (obj modelRoles) AllByIndex(index string) (retObj []Role, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []Role{}
	}
	return
}

func (obj modelRoles) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Role, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	return
}

func (obj modelRoles) Range(min, max, field string) (retObj []Role, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []Role{}
	}
	return
}

func (obj modelRoles) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Role, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Role{}
		}
		return
	}
	return
}

func (obj modelRoles) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj Role
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelRoles) NewByReflection() (value reflect.Value) {
	retObj := Role{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelRoles) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []Role
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

func (obj modelRoles) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionRolesMutex.RLock()
		bootstrapped := GoCoreRolesHasBootStrapped
		collectionRolesMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("Roles has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "Roles"
	query.entityName = "Role"
	return query
}
func (obj modelRoles) Index() error {
	return dbServices.BoltDB.Init(&Role{})
}

func (obj modelRoles) BootStrapComplete() {
	collectionRolesMutex.Lock()
	GoCoreRolesHasBootStrapped = true
	collectionRolesMutex.Unlock()
}
func (obj modelRoles) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of Roles Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []Role
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwoJewoJCSJJZCI6IjU3YzA3ZWYzZGNiYTBmN2EwYmUzMzhiOCIsCgkJIk5hbWUiOiJBY2NvdW50IEFkbWluaXN0cmF0b3IiLAoJCSJBY2NvdW50SWQiOiIiLAoJCSJDYW5EZWxldGUiOmZhbHNlLAoJCSJDcmVhdGVEYXRlIjoiMjAxNi0wOC0yNlQxMDo0OTowNC42MzA1MzY0NDYtMDQ6MDAiLAoJCSJVcGRhdGVEYXRlIjoiMjAxNi0wOC0yNlQxMDo0OTowNC42MzA1MzY0NDYtMDQ6MDAiLAoJCSJBY2NvdW50VHlwZSI6ImN1c3QiLAoJCSJCb290c3RyYXBNZXRhIjogewoJCQkiVmVyc2lvbiI6IDAsCgkJCSJEb21haW4iOiAiIiwKCQkJIlJlbGVhc2VNb2RlIjogIiIsCgkJCSJQcm9kdWN0TmFtZSI6ICIiLAoJCQkiRG9tYWlucyI6IG51bGwsCgkJCSJQcm9kdWN0TmFtZXMiOiBudWxsLAoJCQkiRGVsZXRlUm93IjogZmFsc2UsCgkJCSJBbHdheXNVcGRhdGUiOiB0cnVlCgkJfSwKCQkiU2hvcnROYW1lIjoiQURNSU4iLAoJCSJMYXN0VXBkYXRlSWQiOiI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciCgl9Cl0K"

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "Roles")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Roles due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("roles", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Roles: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for Roles: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []Role
	for _, file := range files {
		var fileBootstrap []Role
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"Roles", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for Roles: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for Roles: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"Roles", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "Roles")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting Roles", len(v))

	for _, doc := range v {
		var original Role
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "Roles" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for Roles:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for Roles:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("Roles skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap Roles")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("roles", "Roles")
		}
		if err == nil {
			log.Println("Successfully bootstrapped Roles")
			if actualCount != originalCount {
				logger.Message("Roles counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelRoles) New() *Role {
	return &Role{}
}

func (obj *Role) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *Role) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("Roles", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *Role) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *Role) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *Role) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("Roles", self.Id.Hex())
	return self.Save()
}

func (self *Role) ValidateAndClean() error {

	return validateFields(Role{}, self, reflect.ValueOf(self).Elem())
}

func (self *Role) Reflect() []Field {

	return Reflect(Role{})
}

func (self *Role) Delete() error {
	dbServices.CollectionCache{}.Remove("Roles", self.Id.Hex())
	return dbServices.BoltDB.Delete("Role", self.Id.Hex())
}

func (self *Role) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("Roles", self.Id.Hex())
	return dbServices.BoltDB.Delete("Roles", self.Id.Hex())
}

func (self *Role) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *Role) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Role) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Role) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *Role) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *Role) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *Role) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelRoles) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Name":
		obj, ok := x.(string)
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
	case "CanDelete":
		obj, ok := x.(bool)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "AccountType":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "ShortName":
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
	}
	return
}

func (obj modelRoles) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
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
	case "CanDelete":
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
	case "AccountType":
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
	case "ShortName":
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
