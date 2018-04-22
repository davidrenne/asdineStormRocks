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

var Accounts modelAccounts

type modelAccounts struct{}

var collectionAccountsMutex *sync.RWMutex

type AccountJoinItems struct {
	Count int        `json:"Count"`
	Items *[]Account `json:"Items"`
}

var GoCoreAccountsHasBootStrapped bool

func init() {
	collectionAccountsMutex = &sync.RWMutex{}

	//Accounts.Index()
	go func() {
		time.Sleep(time.Second * 5)
		Accounts.Bootstrap()
	}()
	store.RegisterStore(Accounts)
}

func (self *Account) GetId() string {
	return self.Id.Hex()
}

type Account struct {
	Id               bson.ObjectId              `json:"Id" storm:"id"`
	AccountName      string                     `json:"AccountName" validate:"true,,,,,,"`
	Address1         string                     `json:"Address1" validate:"true,,,,,,"`
	Address2         string                     `json:"Address2"`
	Region           string                     `json:"Region" validate:"false,,,,,,"`
	City             string                     `json:"City" validate:"true,,,,,,"`
	PostCode         string                     `json:"PostCode" validate:"true,,,,,,"`
	CountryId        string                     `json:"CountryId" validate:"true,,,,,,"`
	StateName        string                     `json:"StateName" validate:"false,,,,,,"`
	StateId          string                     `json:"StateId" validate:"false,,,,,,"`
	PrimaryPhone     AccountsPhoneInfo          `json:"PrimaryPhone"`
	SecondaryPhone   AccountsSecondaryPhoneInfo `json:"SecondaryPhone"`
	Email            string                     `json:"Email" validate:"true,email,,,,,"`
	AccountTypeShort string                     `json:"AccountTypeShort"`
	AccountTypeLong  string                     `json:"AccountTypeLong"`
	RelatedAcctId    string                     `json:"RelatedAcctId" storm:"index"`
	IsSystemAccount  bool                       `json:"IsSystemAccount"`
	CreateDate       time.Time                  `json:"CreateDate" bson:"CreateDate"`
	UpdateDate       time.Time                  `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId     string                     `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta    *BootstrapMeta             `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id           string `json:"Id"`
		AccountName  string `json:"AccountName"`
		Address1     string `json:"Address1"`
		Address2     string `json:"Address2"`
		Region       string `json:"Region"`
		City         string `json:"City"`
		PostCode     string `json:"PostCode"`
		CountryId    string `json:"CountryId"`
		StateName    string `json:"StateName"`
		StateId      string `json:"StateId"`
		PrimaryPhone struct {
			Value      string `json:"Value"`
			Numeric    string `json:"Numeric"`
			DialCode   string `json:"DialCode"`
			CountryISO string `json:"CountryISO"`
		} `json:"PrimaryPhone"`
		SecondaryPhone struct {
			Value      string `json:"Value"`
			Numeric    string `json:"Numeric"`
			DialCode   string `json:"DialCode"`
			CountryISO string `json:"CountryISO"`
		} `json:"SecondaryPhone"`
		Email            string `json:"Email"`
		AccountTypeShort string `json:"AccountTypeShort"`
		AccountTypeLong  string `json:"AccountTypeLong"`
		RelatedAcctId    string `json:"RelatedAcctId"`
		IsSystemAccount  string `json:"IsSystemAccount"`
	} `json:"Errors" bson:"-"`

	Views struct {
		UpdateDate    string `json:"UpdateDate" ref:"UpdateDate~DateTime"`
		UpdateFromNow string `json:"UpdateFromNow" ref:"UpdateDate~TimeFromNow"`
	} `json:"Views" bson:"-"`

	Joins struct {
		Country        *Country `json:"Country,omitempty" join:"Countries,Country,CountryId,false,"`
		State          *State   `json:"State,omitempty" join:"States,State,StateId,false,"`
		RelatedAccount *Account `json:"RelatedAccount,omitempty" join:"Accounts,Account,RelatedAcctId,false,"`
		LastUpdateUser *User    `json:"LastUpdateUser,omitempty" join:"Users,User,LastUpdateId,false,"`
	} `json:"Joins" bson:"-"`
}

type AccountsPhoneInfo struct {
	Value      string `json:"Value" validate:"true,,,,,,"`
	Numeric    string `json:"Numeric"`
	DialCode   string `json:"DialCode"`
	CountryISO string `json:"CountryISO"`
}

type AccountsSecondaryPhoneInfo struct {
	Value      string `json:"Value"`
	Numeric    string `json:"Numeric"`
	DialCode   string `json:"DialCode"`
	CountryISO string `json:"CountryISO"`
}

func (self modelAccounts) Single(field string, value interface{}) (retObj Account, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelAccounts) Search(field string, value interface{}) (retObj []Account, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []Account{}
	}
	return
}

func (obj modelAccounts) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []Account, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	return
}

func (obj modelAccounts) All() (retObj []Account, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []Account{}
	}
	return
}

func (obj modelAccounts) AllAdvanced(limit int, skip int) (retObj []Account, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	return
}

func (obj modelAccounts) AllByIndex(index string) (retObj []Account, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []Account{}
	}
	return
}

func (obj modelAccounts) AllByIndexAdvanced(index string, limit int, skip int) (retObj []Account, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	return
}

func (obj modelAccounts) Range(min, max, field string) (retObj []Account, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []Account{}
	}
	return
}

func (obj modelAccounts) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []Account, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []Account{}
		}
		return
	}
	return
}

func (obj modelAccounts) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj Account
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelAccounts) NewByReflection() (value reflect.Value) {
	retObj := Account{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelAccounts) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []Account
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

func (obj modelAccounts) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionAccountsMutex.RLock()
		bootstrapped := GoCoreAccountsHasBootStrapped
		collectionAccountsMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("Accounts has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "Accounts"
	query.entityName = "Account"
	return query
}
func (obj modelAccounts) Index() error {
	return dbServices.BoltDB.Init(&Account{})
}

func (obj modelAccounts) BootStrapComplete() {
	collectionAccountsMutex.Lock()
	GoCoreAccountsHasBootStrapped = true
	collectionAccountsMutex.Unlock()
}
func (obj modelAccounts) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of Accounts Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []Account
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := "IFsKCXsKCQkiSWQiOiAiNTg0MDU3MThmOTRjNjcxYjA1MzUwODU3IiwKCQkiQWNjb3VudE5hbWUiOiAiTXkgQ29tcGFueSIsCgkJIkFkZHJlc3MxIjogIk15IENvbXBhbnkgRHJpdmUiLAoJCSJBZGRyZXNzMiIgOiAiIiwKCQkiUmVnaW9uIjogIlVua25vd24iLAoJCSJJc1N5c3RlbUFjY291bnQiOiB0cnVlLAoJCSJJc0RlZmF1bHRBY2NvdW50IjogdHJ1ZSwKCQkiQ2l0eSI6ICJVbmtub3duIiwKCQkiUG9zdENvZGUiOiIwMDAwMCIsCgkJIkNvdW50cnlJZCI6IjU3ZmU1OTM2ZWQwNzI3ZDg5NDA2ZDA1OCIsCgkJIlByaW1hcnlQaG9uZSI6IHsKCQkJIlZhbHVlIjogIisxIDAwMC0wMDAtMDAwIiwKCQkJIk51bWVyaWMiOiAiMTAwMDAwMDAwMCIsCgkJCSJEaWFsQ29kZSI6ICIxIiwKCQkJIkNvdW50cnlJU08iOiAidXMiCgkJfSwKCQkiU2Vjb25kYXJ5UGhvbmUiOiB7CgkJCSJWYWx1ZSI6ICIrMSAwMDAtMDAwLTAwMCIsCgkJCSJOdW1lcmljIjogIjEwMDAwMDAwMDAiLAoJCQkiRGlhbENvZGUiOiAiMSIsCgkJCSJDb3VudHJ5SVNPIjogInVzIgoJCX0sCgkJIkVtYWlsIjoicm9vdEByb290LWNvbXBhbnkuY29tIiwKCQkiQmlsbGluZ0luZm8iOgoJCXsKCQkJIlZhdWx0UmVmSWQiOiIiCgkJfSwKCQkiQWNjb3VudFR5cGVTaG9ydCI6IiIsCgkJIkFjY291bnRUeXBlTG9uZyI6IiIsCgkJIkRpc2FibGVkRmVhdHVyZXMiOltdLAoJCSJSZWxhdGVkQWNjdElkIjoiIiwKCQkiQ3JlYXRlRGF0ZSI6IjIwMTYtMDgtMjZUMTA6NDk6MDQuNjMwNTM2NDQ2LTA0OjAwIiwKCQkiVXBkYXRlRGF0ZSI6IjIwMTYtMDgtMjZUMTA6NDk6MDQuNjMwNTM2NDQ2LTA0OjAwIiwKCQkiQm9vdHN0cmFwTWV0YSI6IHsKCQl9Cgl9Cl0K"

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "Accounts")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Accounts due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("accounts", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for Accounts: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for Accounts: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []Account
	for _, file := range files {
		var fileBootstrap []Account
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"Accounts", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for Accounts: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for Accounts: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"Accounts", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "Accounts")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting Accounts", len(v))

	for _, doc := range v {
		var original Account
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "Accounts" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for Accounts:  " + doc.Id.Hex() + "  " + err.Error())
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
						log.Println("Failed to bootstrap data for Accounts:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("Accounts skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap Accounts")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("accounts", "Accounts")
		}
		if err == nil {
			log.Println("Successfully bootstrapped Accounts")
			if actualCount != originalCount {
				logger.Message("Accounts counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelAccounts) New() *Account {
	return &Account{}
}

func (obj *Account) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *Account) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("Accounts", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *Account) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *Account) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *Account) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("Accounts", self.Id.Hex())
	return self.Save()
}

func (self *Account) ValidateAndClean() error {

	return validateFields(Account{}, self, reflect.ValueOf(self).Elem())
}

func (self *Account) Reflect() []Field {

	return Reflect(Account{})
}

func (self *Account) Delete() error {
	dbServices.CollectionCache{}.Remove("Accounts", self.Id.Hex())
	return dbServices.BoltDB.Delete("Account", self.Id.Hex())
}

func (self *Account) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("Accounts", self.Id.Hex())
	return dbServices.BoltDB.Delete("Accounts", self.Id.Hex())
}

func (self *Account) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

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

func (self *Account) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *Account) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *Account) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *Account) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *Account) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *Account) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelAccounts) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "CountryISO":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Address2":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "StateId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "PrimaryPhone":
		data, _ := json.Marshal(x)
		var obj AccountsPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "SecondaryPhone":
		data, _ := json.Marshal(x)
		var obj AccountsSecondaryPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "RelatedAcctId":
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
	case "Id":
		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "CountryId":
		obj, ok := x.(string)
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
	case "AccountTypeLong":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "IsSystemAccount":
		obj, ok := x.(bool)
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
	case "AccountName":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "AccountTypeShort":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Address1":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Region":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "City":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "PostCode":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "StateName":
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
	}
	return
}

func (obj modelAccounts) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "StateId":
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
	case "PrimaryPhone":
		if x == nil {
			obj := AccountsPhoneInfo{}
			value = reflect.ValueOf(&obj)
			return
		}

		data, _ := json.Marshal(x)
		var obj AccountsPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "SecondaryPhone":
		if x == nil {
			obj := AccountsSecondaryPhoneInfo{}
			value = reflect.ValueOf(&obj)
			return
		}

		data, _ := json.Marshal(x)
		var obj AccountsSecondaryPhoneInfo
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "RelatedAcctId":
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
	case "Address2":
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
	case "CountryId":
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
	case "AccountTypeLong":
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
	case "IsSystemAccount":
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
	case "AccountTypeShort":
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
	case "AccountName":
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
	case "Region":
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
	case "City":
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
	case "PostCode":
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
	case "StateName":
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
	case "Address1":
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
