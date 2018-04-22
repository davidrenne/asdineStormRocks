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

var AccountRoles modelAccountRoles

type modelAccountRoles struct{}

var collectionAccountRolesMutex *sync.RWMutex

type AccountRoleJoinItems struct {
	Count int            `json:"Count"`
	Items *[]AccountRole `json:"Items"`
}

var GoCoreAccountRolesHasBootStrapped bool

func init() {
	collectionAccountRolesMutex = &sync.RWMutex{}

	//AccountRoles.Index()
	go func() {
		time.Sleep(time.Second * 5)
		AccountRoles.Bootstrap()
	}()
	store.RegisterStore(AccountRoles)
}

func (self *AccountRole) GetId() string {
	return self.Id.Hex()
}

type AccountRole struct {
	Id            bson.ObjectId  `json:"Id" storm:"id"`
	AccountId     string         `json:"AccountId"`
	UserId        string         `json:"UserId"`
	RoleId        string         `json:"RoleId"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id        string `json:"Id"`
		AccountId string `json:"AccountId"`
		UserId    string `json:"UserId"`
		RoleId    string `json:"RoleId"`
	} `json:"Errors" bson:"-"`

	Joins struct {
		User    *User    `json:"User,omitempty" join:"Users,User,UserId,false,"`
		Account *Account `json:"Account,omitempty" join:"Accounts,Account,AccountId,false,"`
		Role    *Role    `json:"Role,omitempty" join:"Roles,Role,RoleId,false,"`
	} `json:"Joins" bson:"-"`
}

func (self modelAccountRoles) Single(field string, value interface{}) (retObj AccountRole, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelAccountRoles) Search(field string, value interface{}) (retObj []AccountRole, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []AccountRole{}
	}
	return
}

func (obj modelAccountRoles) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []AccountRole, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	return
}

func (obj modelAccountRoles) All() (retObj []AccountRole, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []AccountRole{}
	}
	return
}

func (obj modelAccountRoles) AllAdvanced(limit int, skip int) (retObj []AccountRole, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	return
}

func (obj modelAccountRoles) AllByIndex(index string) (retObj []AccountRole, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []AccountRole{}
	}
	return
}

func (obj modelAccountRoles) AllByIndexAdvanced(index string, limit int, skip int) (retObj []AccountRole, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	return
}

func (obj modelAccountRoles) Range(min, max, field string) (retObj []AccountRole, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []AccountRole{}
	}
	return
}

func (obj modelAccountRoles) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []AccountRole, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []AccountRole{}
		}
		return
	}
	return
}

func (obj modelAccountRoles) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj AccountRole
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelAccountRoles) NewByReflection() (value reflect.Value) {
	retObj := AccountRole{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelAccountRoles) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []AccountRole
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

func (obj modelAccountRoles) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionAccountRolesMutex.RLock()
		bootstrapped := GoCoreAccountRolesHasBootStrapped
		collectionAccountRolesMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("AccountRoles has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "AccountRoles"
	query.entityName = "AccountRole"
	return query
}
func (obj modelAccountRoles) Index() error {
	return dbServices.BoltDB.Init(&AccountRole{})
}

func (obj modelAccountRoles) BootStrapComplete() {
	collectionAccountRolesMutex.Lock()
	GoCoreAccountRolesHasBootStrapped = true
	collectionAccountRolesMutex.Unlock()
}
func (obj modelAccountRoles) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of AccountRoles Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []AccountRole
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwoJewoJCSJJZCI6IjU4NDE5YWU2YWI1YjZmM2RkNTM3NzI3MiIsCgkJIkFjY291bnRJZCI6IjU4NDA1NzE4Zjk0YzY3MWIwNTM1MDg1NyIsCgkJIlVzZXJJZCI6IjU4NDA1Yjc5Zjk0YzY3MWIwNTM1MDg1OCIsCgkJIlJvbGVJZCI6IjU3YzA3ZWYzZGNiYTBmN2EwYmUzMzhiOCIsCgkJIkxhc3RVcGRhdGVJZCI6IjU4NDA1Yjc5Zjk0YzY3MWIwNTM1MDg1OCIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJBbHdheXNVcGRhdGUiOnRydWUKCQl9Cgl9Cl0K"

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "AccountRoles")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for AccountRoles due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("accountRoles", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for AccountRoles: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for AccountRoles: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []AccountRole
	for _, file := range files {
		var fileBootstrap []AccountRole
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"AccountRoles", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for AccountRoles: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for AccountRoles: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"AccountRoles", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "AccountRoles")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting AccountRoles", len(v))

	for _, doc := range v {
		var original AccountRole
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "AccountRoles" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for AccountRoles:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for AccountRoles:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("AccountRoles skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap AccountRoles")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("accountRoles", "AccountRoles")
		}
		if err == nil {
			log.Println("Successfully bootstrapped AccountRoles")
			if actualCount != originalCount {
				logger.Message("AccountRoles counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelAccountRoles) New() *AccountRole {
	return &AccountRole{}
}

func (obj *AccountRole) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *AccountRole) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("AccountRoles", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *AccountRole) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *AccountRole) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *AccountRole) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("AccountRoles", self.Id.Hex())
	return self.Save()
}

func (self *AccountRole) ValidateAndClean() error {

	return validateFields(AccountRole{}, self, reflect.ValueOf(self).Elem())
}

func (self *AccountRole) Reflect() []Field {

	return Reflect(AccountRole{})
}

func (self *AccountRole) Delete() error {
	dbServices.CollectionCache{}.Remove("AccountRoles", self.Id.Hex())
	return dbServices.BoltDB.Delete("AccountRole", self.Id.Hex())
}

func (self *AccountRole) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("AccountRoles", self.Id.Hex())
	return dbServices.BoltDB.Delete("AccountRoles", self.Id.Hex())
}

func (self *AccountRole) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *AccountRole) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *AccountRole) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *AccountRole) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *AccountRole) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *AccountRole) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *AccountRole) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelAccountRoles) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "AccountId":
		obj, ok := x.(string)
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
	case "RoleId":
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

func (obj modelAccountRoles) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "RoleId":
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
