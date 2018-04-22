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

var Users modelUsers

type modelUsers struct{}

var collectionUsersMutex *sync.RWMutex

type UserJoinItems struct {
	Count int     `json:"Count"`
	Items *[]User `json:"Items"`
}

var GoCoreUsersHasBootStrapped bool

func init() {
	collectionUsersMutex = &sync.RWMutex{}

	//Users.Index()
	go func() {
		time.Sleep(time.Second * 5)
		Users.Bootstrap()
	}()
	store.RegisterStore(Users)
}

func (self *User) GetId() string {
	return self.Id.Hex()
}

type User struct {
	Id                    bson.ObjectId     `json:"Id" storm:"id"`
	First                 string            `json:"First" validate:"true,,,,,,"`
	Last                  string            `json:"Last" validate:"true,,,,,,"`
	Email                 string            `json:"Email" storm:"unique" validate:"true,email,,,,,"`
	CompanyName           string            `json:"CompanyName"`
	OfficeName            string            `json:"OfficeName"`
	SkypeId               string            `json:"SkypeId"`
	DefaultAccountId      string            `json:"DefaultAccountId"`
	Phone                 UsersPhoneInfo    `json:"Phone"`
	Ext                   string            `json:"Ext"`
	Mobile                UsersPhoneInfo    `json:"Mobile"`
	Preferences           []UsersPreference `json:"Preferences"`
	JobTitle              string            `json:"JobTitle"`
	Dept                  string            `json:"Dept"`
	Bio                   string            `json:"Bio"`
	PhotoIcon             string            `json:"PhotoIcon"`
	PasswordId            string            `json:"PasswordId"`
	Language              string            `json:"Language"`
	TimeZone              string            `json:"TimeZone"`
	DateFormat            string            `json:"DateFormat"`
	LastLoginDate         time.Time         `json:"LastLoginDate"`
	LastLoginIP           string            `json:"LastLoginIP"`
	LoginAttempts         int               `json:"LoginAttempts"`
	Locked                bool              `json:"Locked"`
	EnforcePasswordChange bool              `json:"EnforcePasswordChange"`
	CreateDate            time.Time         `json:"CreateDate" bson:"CreateDate"`
	UpdateDate            time.Time         `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId          string            `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta         *BootstrapMeta    `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id               string `json:"Id"`
		First            string `json:"First"`
		Last             string `json:"Last"`
		Email            string `json:"Email"`
		CompanyName      string `json:"CompanyName"`
		OfficeName       string `json:"OfficeName"`
		SkypeId          string `json:"SkypeId"`
		DefaultAccountId string `json:"DefaultAccountId"`
		Phone            struct {
			Value      string `json:"Value"`
			Numeric    string `json:"Numeric"`
			DialCode   string `json:"DialCode"`
			CountryISO string `json:"CountryISO"`
		} `json:"Phone"`
		Ext    string `json:"Ext"`
		Mobile struct {
			Value      string `json:"Value"`
			Numeric    string `json:"Numeric"`
			DialCode   string `json:"DialCode"`
			CountryISO string `json:"CountryISO"`
		} `json:"Mobile"`
		Preferences []struct {
			Key   string `json:"Key"`
			Value string `json:"Value"`
		} `json:"Preferences"`
		JobTitle              string `json:"JobTitle"`
		Dept                  string `json:"Dept"`
		Bio                   string `json:"Bio"`
		PhotoIcon             string `json:"PhotoIcon"`
		PasswordId            string `json:"PasswordId"`
		Language              string `json:"Language"`
		TimeZone              string `json:"TimeZone"`
		DateFormat            string `json:"DateFormat"`
		LastLoginDate         string `json:"LastLoginDate"`
		LastLoginIP           string `json:"LastLoginIP"`
		LoginAttempts         string `json:"LoginAttempts"`
		Locked                string `json:"Locked"`
		EnforcePasswordChange string `json:"EnforcePasswordChange"`
	} `json:"Errors" bson:"-"`

	Views struct {
		FullName      string `json:"FullName" ref:"Last~Concatenate:, {first}"`
		UpdateDate    string `json:"UpdateDate" ref:"UpdateDate~DateTime"`
		UpdateFromNow string `json:"UpdateFromNow" ref:"UpdateDate~TimeFromNow"`
		Locked        string `json:"Locked" ref:"Locked~EnabledDisabled"`
	} `json:"Views" bson:"-"`

	Joins struct {
		LastUpdateUser *User     `json:"LastUpdateUser,omitempty" join:"Users,User,LastUpdateId,false,"`
		Password       *Password `json:"Password,omitempty" join:"Passwords,Password,PasswordId,false,"`
	} `json:"Joins" bson:"-"`
}

type UsersPhoneInfo struct {
	Value      string `json:"Value"`
	Numeric    string `json:"Numeric"`
	DialCode   string `json:"DialCode"`
	CountryISO string `json:"CountryISO"`
}

type UsersPreference struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

func (self modelUsers) Single(field string, value interface{}) (retObj User, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelUsers) Search(field string, value interface{}) (retObj []User, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []User{}
	}
	return
}

func (obj modelUsers) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []User, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	return
}

func (obj modelUsers) All() (retObj []User, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []User{}
	}
	return
}

func (obj modelUsers) AllAdvanced(limit int, skip int) (retObj []User, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	return
}

func (obj modelUsers) AllByIndex(index string) (retObj []User, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []User{}
	}
	return
}

func (obj modelUsers) AllByIndexAdvanced(index string, limit int, skip int) (retObj []User, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	return
}

func (obj modelUsers) Range(min, max, field string) (retObj []User, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []User{}
	}
	return
}

func (obj modelUsers) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []User, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []User{}
		}
		return
	}
	return
}

func (obj modelUsers) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj User
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelUsers) NewByReflection() (value reflect.Value) {
	retObj := User{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelUsers) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []User
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

func (obj modelUsers) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionUsersMutex.RLock()
		bootstrapped := GoCoreUsersHasBootStrapped
		collectionUsersMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("Users has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "Users"
	query.entityName = "User"
	return query
}
func (obj modelUsers) Index() error {
	return dbServices.BoltDB.Init(&User{})
}

func (obj modelUsers) BootStrapComplete() {
	collectionUsersMutex.Lock()
	GoCoreUsersHasBootStrapped = true
	collectionUsersMutex.Unlock()
}
func (obj modelUsers) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of Users Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []User
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwoJewoJCSJJZCI6ICI1ODQwNWI3OWY5NGM2NzFiMDUzNTA4NTgiLAoJCSJGaXJzdCI6ICJBZG1pbiIsCgkJIkxhc3QiOiAiQWRtaW4iLAoJCSJFbWFpbCI6ICJhZG1pbiIsCgkJIlBhc3N3b3JkSWQiOiAiNTdjZjJiODFmMzZkMjg2NmVlMzdjMGM1IiwKCQkiRGVmYXVsdEFjY291bnRJZCI6ICI1ODQwNTcxOGY5NGM2NzFiMDUzNTA4NTciLAoJCSJDcmVhdGVEYXRlIjoiMjAxNi0wOC0yNlQxMDo0OTowNC42MzA1MzY0NDYtMDQ6MDAiLAoJCSJVcGRhdGVEYXRlIjoiMjAxNi0wOC0yNlQxMDo0OTowNC42MzA1MzY0NDYtMDQ6MDAiLAoJCSJMYXN0TG9naW5EYXRlIjogIjIwMTYtMDgtMjZUMTA6NDk6MDQuNjMwNTM2NDQ2LTA0OjAwIiwKCQkiTGFzdExvZ2luSVAiOiAiIiwKCQkiTGFuZ3VhZ2UiOiAiZW4iLAoJCSJUaW1lWm9uZSI6ICJVUy9FYXN0ZXJuIiwKCQkiRGF0ZUZvcm1hdCI6ICJtbS9kZC95eXl5IiwKCQkiTGFzdFVwZGF0ZUlkIjoiNTdkOWIzODNkY2JhMGY1MTE3MmYxZjU3IiwKCQkiRW5mb3JjZVBhc3N3b3JkQ2hhbmdlIjogdHJ1ZSwKCQkiQm9vdHN0cmFwTWV0YSI6IHsKCQl9Cgl9LAoJewoJCSJJZCI6ICI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciLAoJCSJGaXJzdCI6ICJTeXN0ZW0iLAoJCSJMYXN0IjogIlN5c3RlbSIsCgkJIkVtYWlsIjogImFub255bW91c0BzeXN0ZW0uY29tIiwKCQkiUGFzc3dvcmRJZCI6ICIiLAoJCSJEZWZhdWx0QWNjb3VudElkIjogIiIsCgkJIkNyZWF0ZURhdGUiOiIyMDE2LTA4LTI2VDEwOjQ5OjA0LjYzMDUzNjQ0Ni0wNDowMCIsCgkJIlVwZGF0ZURhdGUiOiIyMDE2LTA4LTI2VDEwOjQ5OjA0LjYzMDUzNjQ0Ni0wNDowMCIsCgkJIkxhc3RMb2dpbkRhdGUiOiAiMjAxNi0wOC0yNlQxMDo0OTowNC42MzA1MzY0NDYtMDQ6MDAiLAoJCSJMYXN0TG9naW5JUCI6ICIiLAoJCSJMYXN0VXBkYXRlSWQiOiI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciCgl9LAoJewoJCSJJZCI6ICI1ODM1ZWI2MWU5ZjEyODNkNDk1MTE0YzEiLAoJCSJGaXJzdCI6ICJDcm9uIiwKCQkiTGFzdCI6ICJKb2IiLAoJCSJFbWFpbCI6ICJjcm9uam9iQHN5c3RlbS5jb20iLAoJCSJQYXNzd29yZElkIjogIiIsCgkJIkRlZmF1bHRBY2NvdW50SWQiOiAiIiwKCQkiQ3JlYXRlRGF0ZSI6IjIwMTYtMDgtMjZUMTA6NDk6MDQuNjMwNTM2NDQ2LTA0OjAwIiwKCQkiVXBkYXRlRGF0ZSI6IjIwMTYtMDgtMjZUMTA6NDk6MDQuNjMwNTM2NDQ2LTA0OjAwIiwKCQkiTGFzdExvZ2luRGF0ZSI6ICIyMDE2LTA4LTI2VDEwOjQ5OjA0LjYzMDUzNjQ0Ni0wNDowMCIsCgkJIkxhc3RMb2dpbklQIjogIiIsCgkJIkxhc3RVcGRhdGVJZCI6IjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIKCX0KXQo="

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "Users")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Users due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("users", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Users: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for Users: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []User
	for _, file := range files {
		var fileBootstrap []User
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"Users", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for Users: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for Users: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"Users", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "Users")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting Users", len(v))

	for _, doc := range v {
		var original User
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "Users" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for Users:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for Users:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("Users skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap Users")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("users", "Users")
		}
		if err == nil {
			log.Println("Successfully bootstrapped Users")
			if actualCount != originalCount {
				logger.Message("Users counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelUsers) New() *User {
	return &User{}
}

func (obj *User) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *User) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *User) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *User) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *User) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	return self.Save()
}

func (self *User) ValidateAndClean() error {

	return validateFields(User{}, self, reflect.ValueOf(self).Elem())
}

func (self *User) Reflect() []Field {

	return Reflect(User{})
}

func (self *User) Delete() error {
	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	return dbServices.BoltDB.Delete("User", self.Id.Hex())
}

func (self *User) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("Users", self.Id.Hex())
	return dbServices.BoltDB.Delete("Users", self.Id.Hex())
}

func (self *User) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *User) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *User) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *User) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *User) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *User) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *User) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelUsers) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Phone":
		data, _ := json.Marshal(x)
		var obj UsersPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Preferences":
		xArray, ok := x.([]interface{})

		if ok {
			arrayToSet := make([]UsersPreference, len(xArray))
			for i := range xArray {
				inf := xArray[i]
				data, _ := json.Marshal(inf)
				var obj UsersPreference
				err = json.Unmarshal(data, &obj)
				if err != nil {
					return
				}
				arrayToSet[i] = obj
			}

			value = reflect.ValueOf(arrayToSet)
		} else {
			data, _ := json.Marshal(x)
			var obj UsersPreference
			err = json.Unmarshal(data, &obj)
			if err != nil {
				return
			}
			value = reflect.ValueOf(obj)
		}

	case "Bio":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "PasswordId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "First":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "OfficeName":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Mobile":
		data, _ := json.Marshal(x)
		var obj UsersPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "LastLoginIP":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Value":
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
	case "Last":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "SkypeId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Numeric":
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
	case "DefaultAccountId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "JobTitle":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Language":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "EnforcePasswordChange":
		obj, ok := x.(bool)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "CompanyName":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Dept":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "TimeZone":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "DateFormat":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Locked":
		obj, ok := x.(bool)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Email":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Ext":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "PhotoIcon":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "DialCode":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "CountryISO":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "LastLoginDate":
		obj, ok := x.(time.Time)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "LoginAttempts":
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

func (obj modelUsers) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Locked":
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
	case "DateFormat":
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
	case "Ext":
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
	case "PhotoIcon":
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
	case "DialCode":
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
	case "CountryISO":
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
	case "Email":
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
	case "LoginAttempts":
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
	case "LastLoginDate":
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
	case "Preferences":
		if x == nil {
			obj := UsersPreference{}
			value = reflect.ValueOf(&obj)
			return
		}

		data, _ := json.Marshal(x)
		var obj UsersPreference
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Bio":
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
	case "PasswordId":
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
	case "Phone":
		if x == nil {
			obj := UsersPhoneInfo{}
			value = reflect.ValueOf(&obj)
			return
		}

		data, _ := json.Marshal(x)
		var obj UsersPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "OfficeName":
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
	case "Mobile":
		if x == nil {
			obj := UsersPhoneInfo{}
			value = reflect.ValueOf(&obj)
			return
		}

		data, _ := json.Marshal(x)
		var obj UsersPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "LastLoginIP":
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
	case "First":
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
	case "SkypeId":
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
	case "Numeric":
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
	case "Last":
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
	case "DefaultAccountId":
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
	case "JobTitle":
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
	case "Language":
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
	case "EnforcePasswordChange":
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
	case "Dept":
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
	case "TimeZone":
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
	case "CompanyName":
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
