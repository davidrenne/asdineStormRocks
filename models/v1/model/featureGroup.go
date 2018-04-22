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

var FeatureGroups modelFeatureGroups

type modelFeatureGroups struct{}

var collectionFeatureGroupsMutex *sync.RWMutex

type FeatureGroupJoinItems struct {
	Count int             `json:"Count"`
	Items *[]FeatureGroup `json:"Items"`
}

var GoCoreFeatureGroupsHasBootStrapped bool

func init() {
	collectionFeatureGroupsMutex = &sync.RWMutex{}

	//FeatureGroups.Index()
	go func() {
		time.Sleep(time.Second * 5)
		FeatureGroups.Bootstrap()
	}()
	store.RegisterStore(FeatureGroups)
}

func (self *FeatureGroup) GetId() string {
	return self.Id.Hex()
}

type FeatureGroup struct {
	Id            bson.ObjectId  `json:"Id" storm:"id"`
	Name          string         `json:"Name" validate:"true,,,,,,"`
	AccountType   string         `json:"AccountType"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id          string `json:"Id"`
		Name        string `json:"Name"`
		AccountType string `json:"AccountType"`
	} `json:"Errors" bson:"-"`

	Views struct {
		UpdateDate    string `json:"UpdateDate" ref:"UpdateDate~DateTime"`
		UpdateFromNow string `json:"UpdateFromNow" ref:"UpdateDate~TimeFromNow"`
	} `json:"Views" bson:"-"`

	Joins struct {
		Features       *FeatureJoinItems `json:"Features,omitempty" join:"Features,Feature,Id,true,FeatureGroupId"`
		LastUpdateUser *User             `json:"LastUpdateUser,omitempty" join:"Users,User,LastUpdateId,false,"`
	} `json:"Joins" bson:"-"`
}

func (self modelFeatureGroups) Single(field string, value interface{}) (retObj FeatureGroup, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelFeatureGroups) Search(field string, value interface{}) (retObj []FeatureGroup, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []FeatureGroup{}
	}
	return
}

func (obj modelFeatureGroups) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []FeatureGroup, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	return
}

func (obj modelFeatureGroups) All() (retObj []FeatureGroup, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []FeatureGroup{}
	}
	return
}

func (obj modelFeatureGroups) AllAdvanced(limit int, skip int) (retObj []FeatureGroup, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	return
}

func (obj modelFeatureGroups) AllByIndex(index string) (retObj []FeatureGroup, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []FeatureGroup{}
	}
	return
}

func (obj modelFeatureGroups) AllByIndexAdvanced(index string, limit int, skip int) (retObj []FeatureGroup, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	return
}

func (obj modelFeatureGroups) Range(min, max, field string) (retObj []FeatureGroup, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []FeatureGroup{}
	}
	return
}

func (obj modelFeatureGroups) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []FeatureGroup, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []FeatureGroup{}
		}
		return
	}
	return
}

func (obj modelFeatureGroups) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj FeatureGroup
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelFeatureGroups) NewByReflection() (value reflect.Value) {
	retObj := FeatureGroup{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelFeatureGroups) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []FeatureGroup
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

func (obj modelFeatureGroups) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionFeatureGroupsMutex.RLock()
		bootstrapped := GoCoreFeatureGroupsHasBootStrapped
		collectionFeatureGroupsMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("FeatureGroups has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "FeatureGroups"
	query.entityName = "FeatureGroup"
	return query
}
func (obj modelFeatureGroups) Index() error {
	return dbServices.BoltDB.Init(&FeatureGroup{})
}

func (obj modelFeatureGroups) BootStrapComplete() {
	collectionFeatureGroupsMutex.Lock()
	GoCoreFeatureGroupsHasBootStrapped = true
	collectionFeatureGroupsMutex.Unlock()
}
func (obj modelFeatureGroups) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of FeatureGroups Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []FeatureGroup
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwoJewoJCSJJZCI6ICI1ODQ1YzJkNTFkNDFjODBhNDA5Y2NlMmMiLAoJCSJOYW1lIjogIkFjY291bnQgUmVsYXRlZCIsCgkJIkFjY291bnRUeXBlIjogIiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wNVQxNDo0MTowOS4xNzQtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDVUMTc6NTU6NTkuMzMzLTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiTmFtZSI6ICIiLAoJCQkiQWNjb3VudFR5cGUiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0ODM4OGMxZDQxYzgyN2I0ZTE0ZGQ2IiwKCQkiTmFtZSI6ICJVc2VyIFJlbGF0ZWQiLAoJCSJBY2NvdW50VHlwZSI6ICIiLAoJCSJDcmVhdGVEYXRlIjogIjIwMTYtMTItMDdUMTE6Mjc6NTYuOTM5LTA1OjAwIiwKCQkiVXBkYXRlRGF0ZSI6ICIyMDE2LTEyLTA3VDExOjQ5OjA4LjU2MS0wNTowMCIsCgkJIkxhc3RVcGRhdGVJZCI6ICI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciLAoJCSJCb290c3RyYXBNZXRhIjogewoJCQkiVmVyc2lvbiI6IDAsCgkJCSJEb21haW4iOiAiIiwKCQkJIlJlbGVhc2VNb2RlIjogIiIsCgkJCSJQcm9kdWN0TmFtZSI6ICIiLAoJCQkiRG9tYWlucyI6IG51bGwsCgkJCSJQcm9kdWN0TmFtZXMiOiBudWxsLAoJCQkiRGVsZXRlUm93IjogZmFsc2UsCgkJCSJBbHdheXNVcGRhdGUiOiB0cnVlCgkJfSwKCQkiRXJyb3JzIjogewoJCQkiSWQiOiAiIiwKCQkJIk5hbWUiOiAiIiwKCQkJIkFjY291bnRUeXBlIjogIiIKCQl9LAoJCSJWaWV3cyI6IHsKCQkJIlVwZGF0ZURhdGUiOiAiIiwKCQkJIlVwZGF0ZUZyb21Ob3ciOiAiIgoJCX0sCgkJIkpvaW5zIjoge30KCX0sCgl7CgkJIklkIjogIjU4NDgzODhjMWQ0MWM4MjdiNGUxNGRkZCIsCgkJIk5hbWUiOiAiU2VydmVyIFNldHRpbmcgUmVsYXRlZCIsCgkJIkFjY291bnRUeXBlIjogIiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wN1QxMToyNzo1Ni45NC0wNTowMCIsCgkJIlVwZGF0ZURhdGUiOiAiMjAxNi0xMi0wN1QxMToyNzo1Ni45NC0wNTowMCIsCgkJIkxhc3RVcGRhdGVJZCI6ICI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciLAoJCSJCb290c3RyYXBNZXRhIjogewoJCQkiVmVyc2lvbiI6IDAsCgkJCSJEb21haW4iOiAiIiwKCQkJIlJlbGVhc2VNb2RlIjogIiIsCgkJCSJQcm9kdWN0TmFtZSI6ICIiLAoJCQkiRG9tYWlucyI6IG51bGwsCgkJCSJQcm9kdWN0TmFtZXMiOiBudWxsLAoJCQkiRGVsZXRlUm93IjogZmFsc2UsCgkJCSJBbHdheXNVcGRhdGUiOiB0cnVlCgkJfSwKCQkiRXJyb3JzIjogewoJCQkiSWQiOiAiIiwKCQkJIk5hbWUiOiAiIiwKCQkJIkFjY291bnRUeXBlIjogIiIKCQl9LAoJCSJWaWV3cyI6IHsKCQkJIlVwZGF0ZURhdGUiOiAiIiwKCQkJIlVwZGF0ZUZyb21Ob3ciOiAiIgoJCX0sCgkJIkpvaW5zIjoge30KCX0sCgl7CgkJIklkIjogIjU4NDlhZTRmMWQ0MWM4NmE4M2Q2ZWZiNiIsCgkJIk5hbWUiOiAiUm9sZSBSZWxhdGVkIiwKCQkiQWNjb3VudFR5cGUiOiAiIiwKCQkiQ3JlYXRlRGF0ZSI6ICIyMDE2LTEyLTA4VDE0OjAyOjM5Ljc0OC0wNTowMCIsCgkJIlVwZGF0ZURhdGUiOiAiMjAxNi0xMi0wOFQxNDowMjozOS43NDgtMDU6MDAiLAoJCSJMYXN0VXBkYXRlSWQiOiAiNTdkOWIzODNkY2JhMGY1MTE3MmYxZjU3IiwKCQkiQm9vdHN0cmFwTWV0YSI6IHsKCQkJIlZlcnNpb24iOiAwLAoJCQkiRG9tYWluIjogIiIsCgkJCSJSZWxlYXNlTW9kZSI6ICIiLAoJCQkiUHJvZHVjdE5hbWUiOiAiIiwKCQkJIkRvbWFpbnMiOiBudWxsLAoJCQkiUHJvZHVjdE5hbWVzIjogbnVsbCwKCQkJIkRlbGV0ZVJvdyI6IGZhbHNlLAoJCQkiQWx3YXlzVXBkYXRlIjogdHJ1ZQoJCX0sCgkJIkVycm9ycyI6IHsKCQkJIklkIjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJBY2NvdW50VHlwZSI6ICIiCgkJfSwKCQkiVmlld3MiOiB7CgkJCSJVcGRhdGVEYXRlIjogIiIsCgkJCSJVcGRhdGVGcm9tTm93IjogIiIKCQl9LAoJCSJKb2lucyI6IHt9Cgl9Cl0="

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "FeatureGroups")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for FeatureGroups due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("featureGroups", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for FeatureGroups: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for FeatureGroups: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []FeatureGroup
	for _, file := range files {
		var fileBootstrap []FeatureGroup
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"FeatureGroups", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for FeatureGroups: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for FeatureGroups: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"FeatureGroups", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "FeatureGroups")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting FeatureGroups", len(v))

	for _, doc := range v {
		var original FeatureGroup
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "FeatureGroups" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for FeatureGroups:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for FeatureGroups:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("FeatureGroups skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap FeatureGroups")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("featureGroups", "FeatureGroups")
		}
		if err == nil {
			log.Println("Successfully bootstrapped FeatureGroups")
			if actualCount != originalCount {
				logger.Message("FeatureGroups counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelFeatureGroups) New() *FeatureGroup {
	return &FeatureGroup{}
}

func (obj *FeatureGroup) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *FeatureGroup) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("FeatureGroups", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *FeatureGroup) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *FeatureGroup) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *FeatureGroup) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("FeatureGroups", self.Id.Hex())
	return self.Save()
}

func (self *FeatureGroup) ValidateAndClean() error {

	return validateFields(FeatureGroup{}, self, reflect.ValueOf(self).Elem())
}

func (self *FeatureGroup) Reflect() []Field {

	return Reflect(FeatureGroup{})
}

func (self *FeatureGroup) Delete() error {
	dbServices.CollectionCache{}.Remove("FeatureGroups", self.Id.Hex())
	return dbServices.BoltDB.Delete("FeatureGroup", self.Id.Hex())
}

func (self *FeatureGroup) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("FeatureGroups", self.Id.Hex())
	return dbServices.BoltDB.Delete("FeatureGroups", self.Id.Hex())
}

func (self *FeatureGroup) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *FeatureGroup) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *FeatureGroup) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *FeatureGroup) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *FeatureGroup) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *FeatureGroup) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *FeatureGroup) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelFeatureGroups) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
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
	case "AccountType":
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

func (obj modelFeatureGroups) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

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
	}
	return
}
