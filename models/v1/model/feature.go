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

var Features modelFeatures

type modelFeatures struct{}

var collectionFeaturesMutex *sync.RWMutex

type FeatureJoinItems struct {
	Count int        `json:"Count"`
	Items *[]Feature `json:"Items"`
}

var GoCoreFeaturesHasBootStrapped bool

func init() {
	collectionFeaturesMutex = &sync.RWMutex{}

	//Features.Index()
	go func() {
		time.Sleep(time.Second * 5)
		Features.Bootstrap()
	}()
	store.RegisterStore(Features)
}

func (self *Feature) GetId() string {
	return self.Id.Hex()
}

type Feature struct {
	Id             bson.ObjectId  `json:"Id" storm:"id"`
	Key            string         `json:"Key" storm:"unique" validate:"true,,,,,,"`
	Name           string         `json:"Name" validate:"true,,,,,,"`
	Description    string         `json:"Description" validate:"true,,,,,,"`
	FeatureGroupId string         `json:"FeatureGroupId" validate:"true,,,,,,"`
	CreateDate     time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate     time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId   string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta  *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id             string `json:"Id"`
		Key            string `json:"Key"`
		Name           string `json:"Name"`
		Description    string `json:"Description"`
		FeatureGroupId string `json:"FeatureGroupId"`
	} `json:"Errors" bson:"-"`

	Views struct {
		UpdateDate    string `json:"UpdateDate" ref:"UpdateDate~DateTime"`
		UpdateFromNow string `json:"UpdateFromNow" ref:"UpdateDate~TimeFromNow"`
	} `json:"Views" bson:"-"`

	Joins struct {
		FeatureGroup   *FeatureGroup `json:"FeatureGroup,omitempty" join:"FeatureGroups,FeatureGroup,FeatureGroupId,false,"`
		LastUpdateUser *User         `json:"LastUpdateUser,omitempty" join:"Users,User,LastUpdateId,false,"`
	} `json:"Joins" bson:"-"`
}

func (self modelFeatures) Single(field string, value interface{}) (retObj Feature, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelFeatures) Search(field string, value interface{}) (retObj []Feature, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []Feature{}
	}
	return
}

func (obj modelFeatures) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []Feature, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	return
}

func (obj modelFeatures) All() (retObj []Feature, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []Feature{}
	}
	return
}

func (obj modelFeatures) AllAdvanced(limit int, skip int) (retObj []Feature, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	return
}

func (obj modelFeatures) AllByIndex(index string) (retObj []Feature, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []Feature{}
	}
	return
}

func (obj modelFeatures) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Feature, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	return
}

func (obj modelFeatures) Range(min, max, field string) (retObj []Feature, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []Feature{}
	}
	return
}

func (obj modelFeatures) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Feature, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Feature{}
		}
		return
	}
	return
}

func (obj modelFeatures) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj Feature
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelFeatures) NewByReflection() (value reflect.Value) {
	retObj := Feature{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelFeatures) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []Feature
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

func (obj modelFeatures) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionFeaturesMutex.RLock()
		bootstrapped := GoCoreFeaturesHasBootStrapped
		collectionFeaturesMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("Features has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "Features"
	query.entityName = "Feature"
	return query
}
func (obj modelFeatures) Index() error {
	return dbServices.BoltDB.Init(&Feature{})
}

func (obj modelFeatures) BootStrapComplete() {
	collectionFeaturesMutex.Lock()
	GoCoreFeaturesHasBootStrapped = true
	collectionFeaturesMutex.Unlock()
}
func (obj modelFeatures) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of Features Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []Feature
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "WwoJewoJCSJJZCI6ICI1ODAzYzRhM2FiODM5MGY3ZTBhN2JjYzQiLAoJCSJLZXkiOiAiQUNDT1VOVF9BREQiLAoJCSJOYW1lIjogIkFkZCBBY2NvdW50IiwKCQkiRGVzY3JpcHRpb24iOiAiQWJpbGl0eSB0byBhZGQgYSBuZXcgYWNjb3VudCB0byB0aGUgc3lzdGVtLiIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDVjMmQ1MWQ0MWM4MGE0MDljY2UyYyIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0wOC0yNlQxMDo0OTowNC42My0wNDowMCIsCgkJIlVwZGF0ZURhdGUiOiAiMjAxNi0wOC0yNlQxMDo0OTowNC42My0wNDowMCIsCgkJIkxhc3RVcGRhdGVJZCI6ICI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciLAoJCSJCb290c3RyYXBNZXRhIjogewoJCQkiVmVyc2lvbiI6IDAsCgkJCSJEb21haW4iOiAiIiwKCQkJIlJlbGVhc2VNb2RlIjogIiIsCgkJCSJQcm9kdWN0TmFtZSI6ICIiLAoJCQkiRG9tYWlucyI6IG51bGwsCgkJCSJQcm9kdWN0TmFtZXMiOiBudWxsLAoJCQkiRGVsZXRlUm93IjogZmFsc2UsCgkJCSJBbHdheXNVcGRhdGUiOiB0cnVlCgkJfSwKCQkiRXJyb3JzIjogewoJCQkiSWQiOiAiIiwKCQkJIktleSI6ICIiLAoJCQkiTmFtZSI6ICIiLAoJCQkiRGVzY3JpcHRpb24iOiAiIiwKCQkJIkZlYXR1cmVHcm91cElkIjogIiIKCQl9LAoJCSJWaWV3cyI6IHsKCQkJIlVwZGF0ZURhdGUiOiAiIiwKCQkJIlVwZGF0ZUZyb21Ob3ciOiAiIgoJCX0sCgkJIkpvaW5zIjoge30KCX0sCgl7CgkJIklkIjogIjU4MDQyNzAxZDExMzBkZmMzNDAyYzM1YiIsCgkJIktleSI6ICJBQ0NPVU5UX1ZJRVciLAoJCSJOYW1lIjogIlZpZXcgQWNjb3VudHMiLAoJCSJEZXNjcmlwdGlvbiI6ICJBYmlsaXR5IHRvIHZpZXcgYWNjb3VudHMuIiwKCQkiRmVhdHVyZUdyb3VwSWQiOiAiNTg0NWMyZDUxZDQxYzgwYTQwOWNjZTJjIiwKCQkiQ3JlYXRlRGF0ZSI6ICIyMDE2LTA4LTI2VDEwOjQ5OjA0LjYzLTA0OjAwIiwKCQkiVXBkYXRlRGF0ZSI6ICIyMDE2LTEyLTA1VDE3OjE5OjQzLjM0MS0wNTowMCIsCgkJIkxhc3RVcGRhdGVJZCI6ICI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciLAoJCSJCb290c3RyYXBNZXRhIjogewoJCQkiVmVyc2lvbiI6IDAsCgkJCSJEb21haW4iOiAiIiwKCQkJIlJlbGVhc2VNb2RlIjogIiIsCgkJCSJQcm9kdWN0TmFtZSI6ICIiLAoJCQkiRG9tYWlucyI6IG51bGwsCgkJCSJQcm9kdWN0TmFtZXMiOiBudWxsLAoJCQkiRGVsZXRlUm93IjogZmFsc2UsCgkJCSJBbHdheXNVcGRhdGUiOiB0cnVlCgkJfSwKCQkiRXJyb3JzIjogewoJCQkiSWQiOiAiIiwKCQkJIktleSI6ICIiLAoJCQkiTmFtZSI6ICIiLAoJCQkiRGVzY3JpcHRpb24iOiAiIiwKCQkJIkZlYXR1cmVHcm91cElkIjogIiIKCQl9LAoJCSJWaWV3cyI6IHsKCQkJIlVwZGF0ZURhdGUiOiAiIiwKCQkJIlVwZGF0ZUZyb21Ob3ciOiAiIgoJCX0sCgkJIkpvaW5zIjoge30KCX0sCgl7CgkJIklkIjogIjU4NDVkZDhlMWQ0MWM4NjE5MjA4NjA0MSIsCgkJIktleSI6ICJBQ0NPVU5UX01PRElGWSIsCgkJIk5hbWUiOiAiTW9kaWZ5IEFjY291bnRzIiwKCQkiRGVzY3JpcHRpb24iOiAiQWJpbGl0eSB0byBtb2RpZnkgYWNjb3VudHMiLAoJCSJGZWF0dXJlR3JvdXBJZCI6ICI1ODQ1YzJkNTFkNDFjODBhNDA5Y2NlMmMiLAoJCSJDcmVhdGVEYXRlIjogIjIwMTYtMTItMDVUMTY6MzU6MTAuMjE5LTA1OjAwIiwKCQkiVXBkYXRlRGF0ZSI6ICIyMDE2LTEyLTA1VDE2OjQ3OjMyLjY5LTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0NWRkOGUxZDQxYzg2MTkyMDg2MDQyIiwKCQkiS2V5IjogIkFDQ09VTlRfREVMRVRFIiwKCQkiTmFtZSI6ICJEZWxldGUgQWNjb3VudHMiLAoJCSJEZXNjcmlwdGlvbiI6ICJBYmlsaXR5IHRvIGRlbGV0ZSBhY2NvdW50cyIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDVjMmQ1MWQ0MWM4MGE0MDljY2UyYyIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wNVQxNjozNToxMC4yMTktMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDVUMTc6Mjk6NTUuMjg0LTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0NWRkOGUxZDQxYzg2MTkyMDg2MDQ0IiwKCQkiS2V5IjogIkFDQ09VTlRfRVhQT1JUIiwKCQkiTmFtZSI6ICJFeHBvcnQgQWNjb3VudHMiLAoJCSJEZXNjcmlwdGlvbiI6ICJBYmlsaXR5IHRvIGV4cG9ydCBhY2NvdW50cyIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDVjMmQ1MWQ0MWM4MGE0MDljY2UyYyIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wNVQxNjozNToxMC4yMTktMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDVUMTY6NDk6MzguMDk4LTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0NzM3NmQxZDQxYzgzZWE3ZDJlYjQ3IiwKCQkiS2V5IjogIkFDQ09VTlRfSU5WSVRFIiwKCQkiTmFtZSI6ICJBY2NvdW50IEludml0ZSIsCgkJIkRlc2NyaXB0aW9uIjogIkludml0ZSBvdGhlcnMgdG8gdGhpcyBhY2NvdW50IiwKCQkiRmVhdHVyZUdyb3VwSWQiOiAiNTg0NWMyZDUxZDQxYzgwYTQwOWNjZTJjIiwKCQkiQ3JlYXRlRGF0ZSI6ICIyMDE2LTEyLTA2VDE3OjEwOjUzLjU2LTA1OjAwIiwKCQkiVXBkYXRlRGF0ZSI6ICIyMDE2LTEyLTA2VDE3OjEwOjUzLjU2LTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0ODNlOGUxZDQxYzgzNWYzMzY4ZDYzIiwKCQkiS2V5IjogIlVTRVJfQUREIiwKCQkiTmFtZSI6ICJBZGQgdXNlciIsCgkJIkRlc2NyaXB0aW9uIjogIkFiaWxpdHkgdG8gYWRkIGEgbmV3ICB1c2VyLiIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDgzODhjMWQ0MWM4MjdiNGUxNGRkNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wN1QxMTo1MzozNC4yMzQtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDdUMTE6NTM6MzQuMjM0LTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0ODNlOGUxZDQxYzgzNWYzMzY4ZDY0IiwKCQkiS2V5IjogIlVTRVJfVklFVyIsCgkJIk5hbWUiOiAiVmlldyB1c2VycyIsCgkJIkRlc2NyaXB0aW9uIjogIkFiaWxpdHkgdG8gdmlldyAgdXNlcnMuIiwKCQkiRmVhdHVyZUdyb3VwSWQiOiAiNTg0ODM4OGMxZDQxYzgyN2I0ZTE0ZGQ2IiwKCQkiQ3JlYXRlRGF0ZSI6ICIyMDE2LTEyLTA3VDExOjUzOjM0LjIzNC0wNTowMCIsCgkJIlVwZGF0ZURhdGUiOiAiMjAxNi0xMi0wN1QxMTo1MzozNC4yMzQtMDU6MDAiLAoJCSJMYXN0VXBkYXRlSWQiOiAiNTdkOWIzODNkY2JhMGY1MTE3MmYxZjU3IiwKCQkiQm9vdHN0cmFwTWV0YSI6IHsKCQkJIlZlcnNpb24iOiAwLAoJCQkiRG9tYWluIjogIiIsCgkJCSJSZWxlYXNlTW9kZSI6ICIiLAoJCQkiUHJvZHVjdE5hbWUiOiAiIiwKCQkJIkRvbWFpbnMiOiBudWxsLAoJCQkiUHJvZHVjdE5hbWVzIjogbnVsbCwKCQkJIkRlbGV0ZVJvdyI6IGZhbHNlLAoJCQkiQWx3YXlzVXBkYXRlIjogdHJ1ZQoJCX0sCgkJIkVycm9ycyI6IHsKCQkJIklkIjogIiIsCgkJCSJLZXkiOiAiIiwKCQkJIk5hbWUiOiAiIiwKCQkJIkRlc2NyaXB0aW9uIjogIiIsCgkJCSJGZWF0dXJlR3JvdXBJZCI6ICIiCgkJfSwKCQkiVmlld3MiOiB7CgkJCSJVcGRhdGVEYXRlIjogIiIsCgkJCSJVcGRhdGVGcm9tTm93IjogIiIKCQl9LAoJCSJKb2lucyI6IHt9Cgl9LAoJewoJCSJJZCI6ICI1ODQ4M2U4ZTFkNDFjODM1ZjMzNjhkNjUiLAoJCSJLZXkiOiAiVVNFUl9NT0RJRlkiLAoJCSJOYW1lIjogIk1vZGlmeSB1c2VycyIsCgkJIkRlc2NyaXB0aW9uIjogIkFiaWxpdHkgdG8gbW9kaWZ5ICB1c2VycyIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDgzODhjMWQ0MWM4MjdiNGUxNGRkNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wN1QxMTo1MzozNC4yMzQtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDdUMTE6NTM6MzQuMjM0LTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0ODNlOGUxZDQxYzgzNWYzMzY4ZDY2IiwKCQkiS2V5IjogIlNFUlZFUl9TRVRUSU5HX01PRElGWSIsCgkJIk5hbWUiOiAiTW9kaWZ5IHNlcnZlciBzZXR0aW5ncyIsCgkJIkRlc2NyaXB0aW9uIjogIkFiaWxpdHkgdG8gbW9kaWZ5ICBzZXJ2ZXIgc2V0dGluZ3MiLAoJCSJGZWF0dXJlR3JvdXBJZCI6ICI1ODQ4Mzg4YzFkNDFjODI3YjRlMTRkZGQiLAoJCSJDcmVhdGVEYXRlIjogIjIwMTYtMTItMDdUMTE6NTM6MzQuMjM0LTA1OjAwIiwKCQkiVXBkYXRlRGF0ZSI6ICIyMDE2LTEyLTA3VDExOjU4OjQ4LjI5NS0wNTowMCIsCgkJIkxhc3RVcGRhdGVJZCI6ICI1N2Q5YjM4M2RjYmEwZjUxMTcyZjFmNTciLAoJCSJCb290c3RyYXBNZXRhIjogewoJCQkiVmVyc2lvbiI6IDAsCgkJCSJEb21haW4iOiAiIiwKCQkJIlJlbGVhc2VNb2RlIjogIiIsCgkJCSJQcm9kdWN0TmFtZSI6ICIiLAoJCQkiRG9tYWlucyI6IG51bGwsCgkJCSJQcm9kdWN0TmFtZXMiOiBudWxsLAoJCQkiRGVsZXRlUm93IjogZmFsc2UsCgkJCSJBbHdheXNVcGRhdGUiOiB0cnVlCgkJfSwKCQkiRXJyb3JzIjogewoJCQkiSWQiOiAiIiwKCQkJIktleSI6ICIiLAoJCQkiTmFtZSI6ICIiLAoJCQkiRGVzY3JpcHRpb24iOiAiIiwKCQkJIkZlYXR1cmVHcm91cElkIjogIiIKCQl9LAoJCSJWaWV3cyI6IHsKCQkJIlVwZGF0ZURhdGUiOiAiIiwKCQkJIlVwZGF0ZUZyb21Ob3ciOiAiIgoJCX0sCgkJIkpvaW5zIjoge30KCX0sCgl7CgkJIklkIjogIjU4NDg1NDE2MWQ0MWM4NmUzZjdlMWNhMCIsCgkJIktleSI6ICJVU0VSX0NIQU5HRV9ST0xFIiwKCQkiTmFtZSI6ICJDaGFuZ2UgYW55IHVzZXIgcm9sZSIsCgkJIkRlc2NyaXB0aW9uIjogIkFsbG93IHRoZSBhYmlsaXR5IHRvIGNoYW5nZSByb2xlcyBmb3IgYW55b25lIGluIHlvdXIgY29tcGFueSIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDgzODhjMWQ0MWM4MjdiNGUxNGRkNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wN1QxMzoyNToyNi4zMTQtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDdUMTM6MjU6MjYuMzE0LTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0OWFlNGYxZDQxYzg2YTgzZDZlZmI3IiwKCQkiS2V5IjogIlJPTEVfVklFVyIsCgkJIk5hbWUiOiAiVmlldyByb2xlcyIsCgkJIkRlc2NyaXB0aW9uIjogIkFiaWxpdHkgdG8gdmlldyByb2xlcyIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDlhZTRmMWQ0MWM4NmE4M2Q2ZWZiNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wOFQxNDowMjozOS43NTEtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDhUMTQ6MDI6MzkuNzUxLTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0OWFlNGYxZDQxYzg2YTgzZDZlZmI4IiwKCQkiS2V5IjogIlJPTEVfQUREIiwKCQkiTmFtZSI6ICJBZGQgcm9sZXMiLAoJCSJEZXNjcmlwdGlvbiI6ICJBYmlsaXR5IHRvIGFkZCByb2xlcyIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDlhZTRmMWQ0MWM4NmE4M2Q2ZWZiNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wOFQxNDowMjozOS43NTEtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDhUMTQ6MDI6MzkuNzUxLTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0OWFlNGYxZDQxYzg2YTgzZDZlZmI5IiwKCQkiS2V5IjogIlJPTEVfTU9ESUZZIiwKCQkiTmFtZSI6ICJNb2RpZnkgcm9sZXMiLAoJCSJEZXNjcmlwdGlvbiI6ICJBYmlsaXR5IHRvIG1vZGlmeSByb2xlcyIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDlhZTRmMWQ0MWM4NmE4M2Q2ZWZiNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wOFQxNDowMjozOS43NTEtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDhUMTQ6MDI6MzkuNzUxLTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0OWFlNGYxZDQxYzg2YTgzZDZlZmJhIiwKCQkiS2V5IjogIlJPTEVfREVMRVRFIiwKCQkiTmFtZSI6ICJEZWxldGUgcm9sZXMiLAoJCSJEZXNjcmlwdGlvbiI6ICJBYmlsaXR5IHRvIGRlbGV0ZSByb2xlcyIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDlhZTRmMWQ0MWM4NmE4M2Q2ZWZiNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wOFQxNDowMjozOS43NTEtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDhUMTQ6MDI6MzkuNzUxLTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfSwKCXsKCQkiSWQiOiAiNTg0OWFlNGYxZDQxYzg2YTgzZDZlZmJjIiwKCQkiS2V5IjogIlJPTEVfQ09QWSIsCgkJIk5hbWUiOiAiQ29weSByb2xlIiwKCQkiRGVzY3JpcHRpb24iOiAiQWJpbGl0eSB0byBjb3B5IGEgcm9sZSIsCgkJIkZlYXR1cmVHcm91cElkIjogIjU4NDlhZTRmMWQ0MWM4NmE4M2Q2ZWZiNiIsCgkJIkNyZWF0ZURhdGUiOiAiMjAxNi0xMi0wOFQxNDowMjozOS43NTEtMDU6MDAiLAoJCSJVcGRhdGVEYXRlIjogIjIwMTYtMTItMDhUMTQ6MDI6MzkuNzUxLTA1OjAwIiwKCQkiTGFzdFVwZGF0ZUlkIjogIjU3ZDliMzgzZGNiYTBmNTExNzJmMWY1NyIsCgkJIkJvb3RzdHJhcE1ldGEiOiB7CgkJCSJWZXJzaW9uIjogMCwKCQkJIkRvbWFpbiI6ICIiLAoJCQkiUmVsZWFzZU1vZGUiOiAiIiwKCQkJIlByb2R1Y3ROYW1lIjogIiIsCgkJCSJEb21haW5zIjogbnVsbCwKCQkJIlByb2R1Y3ROYW1lcyI6IG51bGwsCgkJCSJEZWxldGVSb3ciOiBmYWxzZSwKCQkJIkFsd2F5c1VwZGF0ZSI6IHRydWUKCQl9LAoJCSJFcnJvcnMiOiB7CgkJCSJJZCI6ICIiLAoJCQkiS2V5IjogIiIsCgkJCSJOYW1lIjogIiIsCgkJCSJEZXNjcmlwdGlvbiI6ICIiLAoJCQkiRmVhdHVyZUdyb3VwSWQiOiAiIgoJCX0sCgkJIlZpZXdzIjogewoJCQkiVXBkYXRlRGF0ZSI6ICIiLAoJCQkiVXBkYXRlRnJvbU5vdyI6ICIiCgkJfSwKCQkiSm9pbnMiOiB7fQoJfQpd"

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "Features")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Features due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("features", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Features: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for Features: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []Feature
	for _, file := range files {
		var fileBootstrap []Feature
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"Features", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for Features: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for Features: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"Features", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "Features")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting Features", len(v))

	for _, doc := range v {
		var original Feature
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "Features" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for Features:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for Features:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("Features skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap Features")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("features", "Features")
		}
		if err == nil {
			log.Println("Successfully bootstrapped Features")
			if actualCount != originalCount {
				logger.Message("Features counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelFeatures) New() *Feature {
	return &Feature{}
}

func (obj *Feature) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *Feature) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("Features", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *Feature) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *Feature) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *Feature) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("Features", self.Id.Hex())
	return self.Save()
}

func (self *Feature) ValidateAndClean() error {

	return validateFields(Feature{}, self, reflect.ValueOf(self).Elem())
}

func (self *Feature) Reflect() []Field {

	return Reflect(Feature{})
}

func (self *Feature) Delete() error {
	dbServices.CollectionCache{}.Remove("Features", self.Id.Hex())
	return dbServices.BoltDB.Delete("Feature", self.Id.Hex())
}

func (self *Feature) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("Features", self.Id.Hex())
	return dbServices.BoltDB.Delete("Features", self.Id.Hex())
}

func (self *Feature) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *Feature) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Feature) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Feature) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *Feature) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *Feature) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *Feature) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelFeatures) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "FeatureGroupId":
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
	case "Key":
		obj, ok := x.(string)
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
	case "Description":
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

func (obj modelFeatures) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

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
	case "Description":
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
	case "FeatureGroupId":
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
