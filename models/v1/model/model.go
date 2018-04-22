package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/extensions"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/asaskevich/govalidator"
	"github.com/fatih/camelcase"
	"github.com/globalsign/mgo/bson"
)

const (
	TRANSACTION_DATATYPE_ORIGINAL = 1
	TRANSACTION_DATATYPE_NEW      = 2

	TRANSACTION_CHANGETYPE_INSERT = 1
	TRANSACTION_CHANGETYPE_UPDATE = 2
	TRANSACTION_CHANGETYPE_DELETE = 3

	MGO_RECORD_NOT_FOUND = "not found"

	VALIDATION_ERROR                   = "ValidationError"
	VALIDATION_ERROR_REQUIRED          = "ValidationErrorRequiredFieldMissing"
	VALIDATION_ERROR_EMAIL             = "ValidationErrorInvalidEmail"
	VALIDATION_ERROR_SPECIFIC_REQUIRED = "ValidationFieldSpecificRequired"
	VALIDATION_ERROR_SPECIFIC_EMAIL    = "ValidationFieldSpecificEmailRequired"
)

type modelEntity interface {
	Save() error
	Delete() error
	SaveWithTran(*Transaction) error
	Reflect() []Field
	JoinFields(string, *Query, int) error
	GetId() string
}

type modelCollection interface {
	Rollback(transactionId string) error
}

type collection interface {
	Query() *Query
}

type BootstrapMeta struct {
	Version      int      `json:"Version" bson:"Version"`
	Domain       string   `json:"Domain" bson:"Domain"`
	ReleaseMode  string   `json:"ReleaseMode" bson:"ReleaseMode"`
	ProductName  string   `json:"ProductName" bson:"ProductName"`
	Domains      []string `json:"Domains" bson:"Domains"`
	ProductNames []string `json:"ProductNames" bson:"ProductNames"`
	DeleteRow    bool     `json:"DeleteRow" bson:"DeleteRow"`
	AlwaysUpdate bool     `json:"AlwaysUpdate" bson:"AlwaysUpdate"`
}

type BootstrapSync struct {
	sync.Mutex
	Items [][]byte
}

type tQueue struct {
	sync.RWMutex
	queue map[string]*transactionsToPersist
	ids   map[string][]string
}

type transactionsToPersist struct {
	t             *Transaction
	newItems      []entityTransaction
	originalItems []entityTransaction
	startTime     time.Time
}

type entityTransaction struct {
	changeType int
	committed  bool
	entity     modelEntity
}

type Field struct {
	Name       string
	Label      string
	DataType   string
	IsView     bool
	Validation *dbServices.FieldValidation
}

var transactionQueue tQueue

func init() {
	transactionQueue.ids = make(map[string][]string)
	transactionQueue.queue = make(map[string]*transactionsToPersist)
	go clearTransactionQueue()
}

func (self *transactionsToPersist) UpdateEntity(e modelEntity) (updated bool) {
	for i, _ := range self.newItems {
		item := &self.newItems[i]
		if item.entity.GetId() == e.GetId() {
			item.entity = e
			updated = true
			return
		}
	}
	return
}

func Q(k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{k: v}
}

func QTs(k string, v time.Time) map[string]time.Time {
	return map[string]time.Time{k: v}
}

func RangeQ(k string, min interface{}, max interface{}) map[string]Range {
	var rge map[string]Range
	rge = make(map[string]Range)
	rge[k] = Range{
		Max: max,
		Min: min,
	}
	return rge
}

func MinQ(k string, min interface{}) map[string]Min {
	var rge map[string]Min
	rge = make(map[string]Min)
	rge[k] = Min{
		Min: min,
	}
	return rge
}

func MaxQ(k string, max interface{}) map[string]Max {
	var rge map[string]Max
	rge = make(map[string]Max)
	rge[k] = Max{
		Max: max,
	}
	return rge
}

//Every 12 hours check the transactionQueue and remove any outstanding stale transactions > 48 hours old
func clearTransactionQueue() {

	transactionQueue.Lock()

	for key, value := range transactionQueue.queue {

		if time.Since(value.startTime).Hours() > 48 {
			delete(transactionQueue.queue, key)
		}
	}

	transactionQueue.Unlock()

	time.Sleep(12 * time.Hour)
	clearTransactionQueue()
}

func getBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func decodeBase64(value string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return string(data[:]), nil
}

func getNow() time.Time {
	return time.Now()
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func IsValidationError(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == VALIDATION_ERROR || err.Error() == VALIDATION_ERROR_EMAIL {
		return true
	}
	return false
}

func validateFields(x interface{}, objectToUpdate interface{}, val reflect.Value) error {

	isError := false
	collection := strings.Replace(reflect.TypeOf(x).String()+"s", "model.", "", -1)

	obj := reflect.ValueOf(objectToUpdate)
	methodGetID := obj.MethodByName("GetId")
	inID := []reflect.Value{}
	objectID := methodGetID.Call(inID)

	for key, value := range dbServices.GetValidationTags(x) {

		fieldValue := dbServices.GetReflectionFieldValue(key, objectToUpdate)
		validations := strings.Split(value, ",")

		if validations[0] != "" {
			if err := validateRequired(fieldValue, validations[0]); err != nil {
				dbServices.SetFieldValue("Errors."+key, val, VALIDATION_ERROR_SPECIFIC_REQUIRED)
				if store.OnChange != nil {
					store.OnChange(collection, objectID[0].String(), "Errors."+key, VALIDATION_ERROR_SPECIFIC_REQUIRED, nil)
				}
				isError = true
			}
		}
		if validations[1] != "" {

			cleanup, err := validateType(fieldValue, validations[1])

			if err != nil {
				if err.Error() == VALIDATION_ERROR_EMAIL {
					dbServices.SetFieldValue("Errors."+key, val, VALIDATION_ERROR_SPECIFIC_EMAIL)
					if store.OnChange != nil {
						store.OnChange(collection, objectID[0].String(), "Errors."+key, VALIDATION_ERROR_SPECIFIC_EMAIL, nil)
					}
				}
				isError = true
			}

			if cleanup != "" {
				dbServices.SetFieldValue(key, val, cleanup)
			}

		}

	}

	if isError {
		return errors.New(VALIDATION_ERROR)
	}

	return nil
}

func validateRequired(value string, tagValue string) error {
	if tagValue == "true" {
		if value == "" {
			return errors.New(VALIDATION_ERROR_REQUIRED)
		}
		return nil
	}
	return nil
}

func validateType(value string, tagValue string) (string, error) {
	switch tagValue {
	case dbServices.VALIDATION_TYPE_EMAIL:
		return "", validateEmail(value)
	}
	return "", nil
}

func validateEmail(value string) error {
	if !govalidator.IsEmail(value) {
		return errors.New(VALIDATION_ERROR_EMAIL)
	}
	return nil
}

func getJoins(x reflect.Value, remainingRecursions string) (joins []join, err error) {
	if remainingRecursions == "" {
		return
	}

	fields := strings.Split(remainingRecursions, ".")
	fieldName := fields[0]

	joinsField := x.FieldByName("Joins")
	if joinsField.Kind() != reflect.Struct {
		return
	}

	if fieldName == JOIN_ALL {
		for i := 0; i < joinsField.NumField(); i++ {

			typeField := joinsField.Type().Field(i)
			name := typeField.Name
			tagValue := typeField.Tag.Get("join")
			splitValue := strings.Split(tagValue, ",")
			var j join
			j.collectionName = splitValue[0]
			j.joinSchemaName = splitValue[1]
			j.joinFieldRefName = splitValue[2]
			j.isMany = extensions.StringToBool(splitValue[3])
			j.joinForeignFieldName = splitValue[4]
			j.joinFieldName = name
			j.joinSpecified = JOIN_ALL
			joins = append(joins, j)
		}
	} else {
		typeField, ok := joinsField.Type().FieldByName(fieldName)
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println(fmt.Sprintf("%+v", fields), joinsField.Kind(), ok)
		}

		if ok == false {
			msg := "Could not resolve a field (getJoins model): " + fmt.Sprintf("%+v", remainingRecursions) + " on " + x.Type().String() + " object"
			if serverSettings.WebConfig.Application.LogJoinQueries {
				fmt.Println(msg)
			}
			//new := errors.New(msg)
			//err = new
			return
		}
		name := typeField.Name
		tagValue := typeField.Tag.Get("join")
		splitValue := strings.Split(tagValue, ",")
		var j join
		j.collectionName = splitValue[0]
		j.joinSchemaName = splitValue[1]
		j.joinFieldRefName = splitValue[2]
		j.isMany = extensions.StringToBool(splitValue[3])
		j.joinForeignFieldName = splitValue[4]
		j.joinFieldName = name
		j.joinSpecified = strings.Replace(remainingRecursions, fieldName+".", "", 1)
		if strings.Contains(j.joinSpecified, "Count") && j.joinSpecified[:5] == "Count" {
			j.joinSpecified = "Count"
		}
		joins = append(joins, j)
	}
	return
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

func Reflect(obj interface{}) []Field {
	var ret []Field
	val := reflect.ValueOf(obj)

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if typeField.Name != "Errors" && typeField.Name != "Joins" && typeField.Name != "BootstrapMeta" {
			if typeField.Name == "Views" {
				for f := 0; f < val.FieldByName("Views").NumField(); f++ {
					field := Field{}
					field.IsView = true
					name := val.FieldByName("Views").Type().Field(f).Name
					namePart := camelcase.Split(name)
					for x := 0; x < len(namePart); x++ {
						if x > 0 {
							namePart[x] = strings.ToLower(namePart[x])
						}
					}
					field.Name = val.FieldByName("Views").Type().Field(f).Name
					field.Label = strings.Join(namePart[:], " ")
					field.DataType = val.FieldByName("Views").Type().Field(f).Type.Name()
					validate := val.FieldByName("Views").Type().Field(f).Tag.Get("validate")
					if validate != "" {
						field.Validation = &dbServices.FieldValidation{}
						parts := strings.Split(validate, ",")
						field.Validation.Required = extensions.StringToBool(parts[0])
						field.Validation.Type = parts[1]
						field.Validation.Min = parts[2]
						field.Validation.Max = parts[3]
						field.Validation.Length = parts[4]
						field.Validation.LengthMax = parts[5]
						field.Validation.LengthMin = parts[6]
					}
					ret = append(ret, field)
				}
			} else {
				field := Field{}
				validate := typeField.Tag.Get("validate")
				if validate != "" {
					field.Validation = &dbServices.FieldValidation{}
					parts := strings.Split(validate, ",")
					field.Validation.Required = extensions.StringToBool(parts[0])
					field.Validation.Type = parts[1]
					field.Validation.Min = parts[2]
					field.Validation.Max = parts[3]
					field.Validation.Length = parts[4]
					field.Validation.LengthMax = parts[5]
					field.Validation.LengthMin = parts[6]
				}
				name := typeField.Name
				namePart := camelcase.Split(name)
				for x := 0; x < len(namePart); x++ {
					if x > 0 {
						namePart[x] = strings.ToLower(namePart[x])
					}
				}
				field.Name = typeField.Name
				field.Label = strings.Join(namePart[:], " ")
				field.DataType = typeField.Type.Name()
				ret = append(ret, field)
			}
		}
	}
	return ret
}

func JoinEntity(collectionQ *Query, y interface{}, j join, id string, fieldToSet reflect.Value, remainingRecursions string, q *Query, endRecursion bool, recursionCount int) (err error) {

	defer func() {
		if r := recover(); r != nil {
			msg := "Panic Recovered at model.JoinEntity:  Failed to join " + j.joinSchemaName + " with id:" + id + "  Error:" + fmt.Sprintf("%+v", r)
			err = errors.New(msg)
			if serverSettings.WebConfig.Application.LogJoinQueries && err != nil {
				fmt.Println("err Recursion Line 399->" + fmt.Sprintf("%+v", err))
			}
			return
		}
	}()

	if serverSettings.WebConfig.Application.LogJoinQueries {
		fmt.Println("!!!!!!!!!!!!ddd!!!!!!!!")
		fmt.Printf("%+v", fieldToSet)
		fmt.Println("!!!!!!!!!!!!id!!!!!!!!")
		fmt.Printf("%+v", id)
	}

	//Add any whitelisting or blacklisting of fields
	if len(j.whiteListedFields) > 0 {
		collectionQ = collectionQ.Whitelist(collectionQ.entityName, j.whiteListedFields)
	} else if len(j.blackListedFields) > 0 {
		collectionQ = collectionQ.Blacklist(collectionQ.entityName, j.blackListedFields)
	}

	if IsZeroOfUnderlyingType(fieldToSet.Interface()) || j.isMany {
		if j.isMany && id != "" {
			if remainingRecursions == "Count" {
				cnt, err := collectionQ.ToggleLogFlag(true).Filter(Q(j.joinForeignFieldName, id)).Count(y)
				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("JoinEntity() Recursion Count Only err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
				}
				if serverSettings.WebConfig.Application.LogJoinQueries && err != nil {
					fmt.Println("err Recursion Line 413->" + fmt.Sprintf("%+v", err))
				}
				if err != nil {
					// err = errCnt
					return err
				}
				countField := fieldToSet.Elem().FieldByName("Count")
				countField.Set(reflect.ValueOf(cnt))
				return err
			}
			err = collectionQ.ToggleLogFlag(true).Filter(Q(j.joinForeignFieldName, id)).All(y)
			if serverSettings.WebConfig.Application.LogJoinQueries {
				collectionQ.LogQuery("JoinEntity({" + j.joinForeignFieldName + ": " + id + "}) Recursion Many err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
			}
		} else if id != "" {
			if j.joinForeignFieldName == "" {
				err = collectionQ.ToggleLogFlag(true).ById(id, y)
				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("JoinEntity() Recursion Single By Id (" + id + ") err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
				}
			} else {
				err = collectionQ.ToggleLogFlag(true).Filter(Q(j.joinForeignFieldName, id)).One(y)
				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("JoinEntity({" + j.joinForeignFieldName + ": " + id + "}) Recursion Single err?->" + fmt.Sprintf("%+v", err) + " j->" + fmt.Sprintf("%+v", j))
				}
			}
		}

		if err == nil && id != "" {
			if endRecursion == false && recursionCount > 0 {
				recursionCount--

				in := []reflect.Value{}
				in = append(in, reflect.ValueOf(remainingRecursions))
				in = append(in, reflect.ValueOf(q))
				in = append(in, reflect.ValueOf(recursionCount))

				if j.isMany {

					myArray := reflect.ValueOf(y).Elem()
					for i := 0; i < myArray.Len(); i++ {
						s := myArray.Index(i)
						method := s.Addr().MethodByName("JoinFields")
						values := method.Call(in)
						if values[0].Interface() != nil {
							err = values[0].Interface().(error)
						}
					}
				} else {
					err = CallMethod(y, "JoinFields", in)
				}
			}
			if err != nil && serverSettings.WebConfig.Application.LogJoinQueries {
				fmt.Println("err Recursion Line 465->" + fmt.Sprintf("%+v", err))
			}
			if err == nil {
				if j.isMany {

					itemsField := fieldToSet.Elem().FieldByName("Items")
					countField := fieldToSet.Elem().FieldByName("Count")
					itemsField.Set(reflect.ValueOf(y))
					countField.Set(reflect.ValueOf(reflect.ValueOf(y).Elem().Len()))
					//if serverSettings.WebConfig.Application.LogJoinQueries {
					//fmt.Println("!!!!!!!!!!!!reflected pointer for many!!!!!!!!")
					//fmt.Printf("%+v", itemsField)
					//fmt.Println("!!!!!!!!!!!!test interface!!!!!!!!")
					//fmt.Printf("%+v", itemsField.Interface())
					//}
				} else {
					fieldToSet.Set(reflect.ValueOf(y))
					//if serverSettings.WebConfig.Application.LogJoinQueries {
					//fmt.Println("!!!!!!!!!!!!reflected pointer for single row!!!!!!!!")
					//fmt.Printf("%+v", fieldToSet)
					//fmt.Println("!!!!!!!!!!!!test interface!!!!!!!!")
					//fmt.Printf("%+v", fieldToSet.Interface())
					//}

				}

				//if serverSettings.WebConfig.Application.LogJoinQueries {
				//	fmt.Println("!!!!!!!!!!!!reflected and set value to!!!!!!!!")
				//	fmt.Printf("%+v", reflect.ValueOf(y))
				//}

				if q.renderViews {
					err = q.processViews(y)
					if err != nil && serverSettings.WebConfig.Application.LogJoinQueries {
						collectionQ.LogQuery("err Recursion Line 479->" + fmt.Sprintf("%+v", err))
					}
					if err != nil {
						return
					}
				}

			}
		} else {
			if serverSettings.WebConfig.Application.LogJoinQueries && err != nil {
				fmt.Println("err Recursion Line 495->" + fmt.Sprintf("%+v", err))
			}
		}
	} else {
		if endRecursion == false && recursionCount > 0 {
			recursionCount--
			method := fieldToSet.MethodByName("JoinFields")
			in := []reflect.Value{}
			in = append(in, reflect.ValueOf(remainingRecursions))
			in = append(in, reflect.ValueOf(q))
			in = append(in, reflect.ValueOf(recursionCount))
			values := method.Call(in)
			if values[0].Interface() == nil {

				if serverSettings.WebConfig.Application.LogJoinQueries {
					collectionQ.LogQuery("Recursion returning due to nil values[0] interface")
				}
				err = nil
				return
			}
			err = values[0].Interface().(error)
			if err != nil && serverSettings.WebConfig.Application.LogJoinQueries {
				fmt.Println("err Recursion Line 503->" + fmt.Sprintf("%+v", err))
			}
		}
	}
	return
}

func CallMethod(i interface{}, methodName string, in []reflect.Value) (err error) {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		values := finalMethod.Call(in)
		if values[0].Interface() == nil {
			err = nil
			return
		}
		err = values[0].Interface().(error)
		return
	}

	// return or panic, method not found of either type
	return nil
}

func NewObjectId() string {
	return bson.NewObjectId().Hex()
}

func BootstrapMongoDump(directoryName string, collectionName string) (err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at model.BootstrapMongoDump(): " + fmt.Sprintf("%+v", r))
			return
		}
	}()

	path := serverSettings.APP_LOCATION + "/db/bootstrap/" + directoryName + "/mongoDump"

	if extensions.DoesFileExist(path) == false {
		return
	}

	path = path + "/" + directoryName + "Dump.json"

	if extensions.DoesFileExist(path) == false {
		return
	}

	dbName := serverSettings.WebConfig.DbConnection.Database

	var commandPath string
	if runtime.GOOS == "linux" {
		commandPath = "/usr/bin/mongoimport"
	} else if runtime.GOOS == "darwin" {
		commandPath = "mongoimport"
	}
	err = exec.Command(commandPath, "--db", dbName, "--collection", collectionName, "--file", path, "--upsert").Run()
	return

}

func BootstrapDirectory(directoryName string, collectionCount int) (files [][]byte, err error, directoryFound bool) {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic Recovered at model.BootstrapDirectory(): " + fmt.Sprintf("%+v", r))
			return
		}
	}()

	var syncedItems BootstrapSync
	var wg sync.WaitGroup
	path := serverSettings.APP_LOCATION + "/db/bootstrap/" + directoryName + "/dist"

	if extensions.DoesFileExist(path) == false {
		return
	}

	directoryFound = true

	err = fileCache.LoadCachedManifestFromKeyIntoMemory(directoryName)
	if err != nil {
		return
	}

	err = filepath.Walk(path, func(path string, f os.FileInfo, errWalk error) (err error) {

		if errWalk != nil {
			err = errWalk
			return
		}

		var readFile bool
		if !f.IsDir() && !fileCache.DoesHashExistInManifestCache(directoryName, f.Name()) {
			fileCache.UpdateManifestMemoryCache(directoryName, f.Name(), extensions.Int64ToInt32(f.Size()))
			readFile = true
		}

		if !f.IsDir() && fileCache.DoesHashExistInManifestCache(directoryName, f.Name()) {
			var cachedSize int
			fileCache.ByteManifest.Lock()
			cachedSize = fileCache.ByteManifest.Cache[directoryName][f.Name()]
			fileCache.ByteManifest.Unlock()
			actualRetailPrice := extensions.Int64ToInt32(f.Size())
			if cachedSize != actualRetailPrice {
				fileCache.UpdateManifestMemoryCache(directoryName, f.Name(), extensions.Int64ToInt32(f.Size()))
				readFile = true
				log.Println(f.Name() + " is being read because of a difference in size (cached:" + extensions.IntToString(cachedSize) + " new bytes: " + extensions.IntToString(actualRetailPrice) + ")")
			}
		}

		if f.IsDir() || collectionCount == 0 {
			readFile = true
		}

		if readFile && filepath.Ext(f.Name()) == ".json" {
			wg.Add(1)

			go func() {
				defer wg.Done()
				jsonData, err := ioutil.ReadFile(path)
				if err != nil {
					return
				}
				syncedItems.Lock()
				syncedItems.Items = append(syncedItems.Items, jsonData)
				syncedItems.Unlock()

			}()
		}
		return
	})

	if err == nil {
		fileCache.WriteManifestCacheFile(directoryName)
	}
	wg.Wait()
	files = syncedItems.Items

	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func ResolveEntity(key string) modelEntity {

	switch key {
	case "Account":
		return &Account{}
	case "AccountRole":
		return &AccountRole{}
	case "ActivityLog":
		return &ActivityLog{}
	case "AppError":
		return &AppError{}
	case "Country":
		return &Country{}
	case "Feature":
		return &Feature{}
	case "FeatureGroup":
		return &FeatureGroup{}
	case "FileObject":
		return &FileObject{}
	case "Password":
		return &Password{}
	case "PasswordReset":
		return &PasswordReset{}
	case "Role":
		return &Role{}
	case "RoleFeature":
		return &RoleFeature{}
	case "ServerSetting":
		return &ServerSetting{}
	case "State":
		return &State{}
	case "User":
		return &User{}
	}
	return nil
}

func ResolveField(collectionName string, fieldName string) string {

	switch collectionName + fieldName {
	case "RoleId":
		return "int"
	case "RoleName":
		return "string"
	case "RoleAccountId":
		return "string"
	case "RoleCanDelete":
		return "bool"
	case "RoleAccountType":
		return "string"
	case "RoleShortName":
		return "string"
	case "RoleFeatureId":
		return "int"
	case "RoleFeatureRoleId":
		return "string"
	case "RoleFeatureFeatureId":
		return "string"
	case "AccountId":
		return "int"
	case "AccountAccountName":
		return "string"
	case "AccountAddress1":
		return "string"
	case "AccountAddress2":
		return "string"
	case "AccountRegion":
		return "string"
	case "AccountCity":
		return "string"
	case "AccountPostCode":
		return "string"
	case "AccountCountryId":
		return "string"
	case "AccountStateName":
		return "string"
	case "AccountStateId":
		return "string"
	case "AccountPrimaryPhone":
		return "object"
	case "AccountSecondaryPhone":
		return "object"
	case "AccountEmail":
		return "string"
	case "AccountAccountTypeShort":
		return "string"
	case "AccountAccountTypeLong":
		return "string"
	case "AccountRelatedAcctId":
		return "string"
	case "AccountIsSystemAccount":
		return "bool"
	case "UserId":
		return "int"
	case "UserFirst":
		return "string"
	case "UserLast":
		return "string"
	case "UserEmail":
		return "string"
	case "UserCompanyName":
		return "string"
	case "UserOfficeName":
		return "string"
	case "UserSkypeId":
		return "string"
	case "UserDefaultAccountId":
		return "string"
	case "UserPhone":
		return "object"
	case "UserExt":
		return "string"
	case "UserMobile":
		return "object"
	case "UserPreferences":
		return "objectArray"
	case "UserJobTitle":
		return "string"
	case "UserDept":
		return "string"
	case "UserBio":
		return "string"
	case "UserPhotoIcon":
		return "string"
	case "UserPasswordId":
		return "string"
	case "UserLanguage":
		return "string"
	case "UserTimeZone":
		return "string"
	case "UserDateFormat":
		return "string"
	case "UserLastLoginDate":
		return "dateTime"
	case "UserLastLoginIP":
		return "string"
	case "UserLoginAttempts":
		return "int"
	case "UserLocked":
		return "bool"
	case "UserEnforcePasswordChange":
		return "bool"
	case "SecondaryPhoneInfoValue":
		return "string"
	case "SecondaryPhoneInfoNumeric":
		return "string"
	case "SecondaryPhoneInfoDialCode":
		return "string"
	case "SecondaryPhoneInfoCountryISO":
		return "string"
	case "FeatureGroupId":
		return "int"
	case "FeatureGroupName":
		return "string"
	case "FeatureGroupAccountType":
		return "string"
	case "CountryId":
		return "int"
	case "CountryIso":
		return "string"
	case "CountryName":
		return "string"
	case "PreferenceKey":
		return "string"
	case "PreferenceValue":
		return "string"
	case "AccountRoleId":
		return "int"
	case "AccountRoleAccountId":
		return "string"
	case "AccountRoleUserId":
		return "string"
	case "AccountRoleRoleId":
		return "string"
	case "StateId":
		return "int"
	case "StateShort":
		return "string"
	case "StateName":
		return "string"
	case "StateAlternativeName":
		return "string"
	case "StateCountry":
		return "string"
	case "AppErrorId":
		return "int"
	case "AppErrorAccountId":
		return "string"
	case "AppErrorUserId":
		return "string"
	case "AppErrorClientSide":
		return "bool"
	case "AppErrorUrl":
		return "string"
	case "AppErrorMessage":
		return "string"
	case "AppErrorStackShown":
		return "string"
	case "PasswordResetId":
		return "int"
	case "PasswordResetUserId":
		return "string"
	case "PasswordResetComplete":
		return "bool"
	case "PasswordResetUrl":
		return "string"
	case "FeatureId":
		return "int"
	case "FeatureKey":
		return "string"
	case "FeatureName":
		return "string"
	case "FeatureDescription":
		return "string"
	case "FeatureFeatureGroupId":
		return "string"
	case "PasswordId":
		return "int"
	case "PasswordValue":
		return "string"
	case "ActivityLogId":
		return "int"
	case "ActivityLogAccountId":
		return "string"
	case "ActivityLogUserId":
		return "string"
	case "ActivityLogEntity":
		return "string"
	case "ActivityLogEntityId":
		return "string"
	case "ActivityLogAction":
		return "string"
	case "ActivityLogValue":
		return "string"
	case "ServerSettingId":
		return "int"
	case "ServerSettingName":
		return "string"
	case "ServerSettingKey":
		return "string"
	case "ServerSettingCategory":
		return "string"
	case "ServerSettingValue":
		return "string"
	case "ServerSettingAny":
		return "interface"
	case "FileObjectId":
		return "int"
	case "FileObjectName":
		return "string"
	case "FileObjectPath":
		return "string"
	case "FileObjectSingleDownload":
		return "bool"
	case "FileObjectContent":
		return "string"
	case "FileObjectSize":
		return "int"
	case "FileObjectType":
		return "string"
	case "FileObjectModifiedUnix":
		return "int"
	case "FileObjectModified":
		return "dateTime"
	case "FileObjectMD5":
		return "string"
	case "FileObjectAccountId":
		return "string"
	}
	return ""
}

func ResolveCollection(key string) (collection, error) {

	if serverSettings.WebConfig.Application.LogJoinQueries {
		fmt.Println(key)
	}
	switch key {
	case "Accounts":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Accounts")
		}
		return &modelAccounts{}, nil
	case "AccountRoles":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! AccountRoles")
		}
		return &modelAccountRoles{}, nil
	case "ActivityLogs":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! ActivityLogs")
		}
		return &modelActivityLogs{}, nil
	case "AppErrors":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! AppErrors")
		}
		return &modelAppErrors{}, nil
	case "Countries":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Countries")
		}
		return &modelCountries{}, nil
	case "Features":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Features")
		}
		return &modelFeatures{}, nil
	case "FeatureGroups":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! FeatureGroups")
		}
		return &modelFeatureGroups{}, nil
	case "FileObjects":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! FileObjects")
		}
		return &modelFileObjects{}, nil
	case "Passwords":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Passwords")
		}
		return &modelPasswords{}, nil
	case "PasswordResets":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! PasswordResets")
		}
		return &modelPasswordResets{}, nil
	case "Roles":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Roles")
		}
		return &modelRoles{}, nil
	case "RoleFeatures":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! RoleFeatures")
		}
		return &modelRoleFeatures{}, nil
	case "ServerSettings":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! ServerSettings")
		}
		return &modelServerSettings{}, nil
	case "States":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! States")
		}
		return &modelStates{}, nil
	case "Users":
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Println("in case!! Users")
		}
		return &modelUsers{}, nil
	}
	return nil, errors.New("Failed to resolve collection:  " + key)
}

func ResolveHistoryCollection(key string) modelCollection {

	return nil
}

func joinField(j join, id string, fieldToSet reflect.Value, remainingRecursions string, q *Query, endRecursion bool, recursionCount int) (err error) {

	c, err2 := ResolveCollection(j.collectionName)
	if serverSettings.WebConfig.Application.LogJoinQueries {
		fmt.Println("joinFieldLogging")
		fmt.Println(fmt.Sprintf("%+v", j.collectionName))
		fmt.Println("c")
		fmt.Println(fmt.Sprintf("%+v", c))
		fmt.Println("err2")
		fmt.Println(fmt.Sprintf("%+v", err2))
	}
	if err2 != nil {
		err = errors.New("Failed to resolve collection:  " + j.collectionName)
		return
	}
	switch j.joinSchemaName {
	case "Account":
		var y Account
		if j.isMany {

			obj := fieldToSet.Interface().(*AccountJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Account
			var ji AccountJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "AccountRole":
		var y AccountRole
		if j.isMany {

			obj := fieldToSet.Interface().(*AccountRoleJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []AccountRole
			var ji AccountRoleJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "ActivityLog":
		var y ActivityLog
		if j.isMany {

			obj := fieldToSet.Interface().(*ActivityLogJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []ActivityLog
			var ji ActivityLogJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "AppError":
		var y AppError
		if j.isMany {

			obj := fieldToSet.Interface().(*AppErrorJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []AppError
			var ji AppErrorJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Country":
		var y Country
		if j.isMany {

			obj := fieldToSet.Interface().(*CountryJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Country
			var ji CountryJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Feature":
		var y Feature
		if j.isMany {

			obj := fieldToSet.Interface().(*FeatureJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Feature
			var ji FeatureJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "FeatureGroup":
		var y FeatureGroup
		if j.isMany {

			obj := fieldToSet.Interface().(*FeatureGroupJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []FeatureGroup
			var ji FeatureGroupJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "FileObject":
		var y FileObject
		if j.isMany {

			obj := fieldToSet.Interface().(*FileObjectJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []FileObject
			var ji FileObjectJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Password":
		var y Password
		if j.isMany {

			obj := fieldToSet.Interface().(*PasswordJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Password
			var ji PasswordJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "PasswordReset":
		var y PasswordReset
		if j.isMany {

			obj := fieldToSet.Interface().(*PasswordResetJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []PasswordReset
			var ji PasswordResetJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "Role":
		var y Role
		if j.isMany {

			obj := fieldToSet.Interface().(*RoleJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []Role
			var ji RoleJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "RoleFeature":
		var y RoleFeature
		if j.isMany {

			obj := fieldToSet.Interface().(*RoleFeatureJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []RoleFeature
			var ji RoleFeatureJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "ServerSetting":
		var y ServerSetting
		if j.isMany {

			obj := fieldToSet.Interface().(*ServerSettingJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []ServerSetting
			var ji ServerSettingJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "State":
		var y State
		if j.isMany {

			obj := fieldToSet.Interface().(*StateJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []State
			var ji StateJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	case "User":
		var y User
		if j.isMany {

			obj := fieldToSet.Interface().(*UserJoinItems)
			if obj != nil {
				for i, _ := range *obj.Items {
					item := &(*obj.Items)[i]
					item.JoinFields(remainingRecursions, q, recursionCount)
				}
				return
			}

			var z []User
			var ji UserJoinItems
			fieldToSet.Set(reflect.ValueOf(&ji))
			JoinEntity(c.Query(), &z, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		} else {
			JoinEntity(c.Query(), &y, j, id, fieldToSet, remainingRecursions, q, endRecursion, recursionCount)
		}
		return
	}
	err = errors.New("Failed to resolve schema :  " + j.joinSchemaName)
	return
}

const (
	FIELD_ACCOUNT_ID                    = "Id"
	FIELD_ACCOUNT_ACCOUNTNAME           = "AccountName"
	FIELD_ACCOUNT_ADDRESS1              = "Address1"
	FIELD_ACCOUNT_ADDRESS2              = "Address2"
	FIELD_ACCOUNT_REGION                = "Region"
	FIELD_ACCOUNT_CITY                  = "City"
	FIELD_ACCOUNT_POSTCODE              = "PostCode"
	FIELD_ACCOUNT_COUNTRYID             = "CountryId"
	FIELD_ACCOUNT_STATENAME             = "StateName"
	FIELD_ACCOUNT_STATEID               = "StateId"
	FIELD_ACCOUNT_PRIMARYPHONE          = "PrimaryPhone"
	FIELD_ACCOUNT_SECONDARYPHONE        = "SecondaryPhone"
	FIELD_ACCOUNT_EMAIL                 = "Email"
	FIELD_ACCOUNT_ACCOUNTTYPESHORT      = "AccountTypeShort"
	FIELD_ACCOUNT_ACCOUNTTYPELONG       = "AccountTypeLong"
	FIELD_ACCOUNT_RELATEDACCTID         = "RelatedAcctId"
	FIELD_ACCOUNT_ISSYSTEMACCOUNT       = "IsSystemAccount"
	FIELD_ROLE_ID                       = "Id"
	FIELD_ROLE_NAME                     = "Name"
	FIELD_ROLE_ACCOUNTID                = "AccountId"
	FIELD_ROLE_CANDELETE                = "CanDelete"
	FIELD_ROLE_ACCOUNTTYPE              = "AccountType"
	FIELD_ROLE_SHORTNAME                = "ShortName"
	FIELD_ROLEFEATURE_ID                = "Id"
	FIELD_ROLEFEATURE_ROLEID            = "RoleId"
	FIELD_ROLEFEATURE_FEATUREID         = "FeatureId"
	FIELD_ACCOUNTROLE_ID                = "Id"
	FIELD_ACCOUNTROLE_ACCOUNTID         = "AccountId"
	FIELD_ACCOUNTROLE_USERID            = "UserId"
	FIELD_ACCOUNTROLE_ROLEID            = "RoleId"
	FIELD_STATE_ID                      = "Id"
	FIELD_STATE_SHORT                   = "Short"
	FIELD_STATE_NAME                    = "Name"
	FIELD_STATE_ALTERNATIVENAME         = "AlternativeName"
	FIELD_STATE_COUNTRY                 = "Country"
	FIELD_APPERROR_ID                   = "Id"
	FIELD_APPERROR_ACCOUNTID            = "AccountId"
	FIELD_APPERROR_USERID               = "UserId"
	FIELD_APPERROR_CLIENTSIDE           = "ClientSide"
	FIELD_APPERROR_URL                  = "Url"
	FIELD_APPERROR_MESSAGE              = "Message"
	FIELD_APPERROR_STACKSHOWN           = "StackShown"
	FIELD_USER_ID                       = "Id"
	FIELD_USER_FIRST                    = "First"
	FIELD_USER_LAST                     = "Last"
	FIELD_USER_EMAIL                    = "Email"
	FIELD_USER_COMPANYNAME              = "CompanyName"
	FIELD_USER_OFFICENAME               = "OfficeName"
	FIELD_USER_SKYPEID                  = "SkypeId"
	FIELD_USER_DEFAULTACCOUNTID         = "DefaultAccountId"
	FIELD_USER_PHONE                    = "Phone"
	FIELD_USER_EXT                      = "Ext"
	FIELD_USER_MOBILE                   = "Mobile"
	FIELD_USER_PREFERENCES              = "Preferences"
	FIELD_USER_JOBTITLE                 = "JobTitle"
	FIELD_USER_DEPT                     = "Dept"
	FIELD_USER_BIO                      = "Bio"
	FIELD_USER_PHOTOICON                = "PhotoIcon"
	FIELD_USER_PASSWORDID               = "PasswordId"
	FIELD_USER_LANGUAGE                 = "Language"
	FIELD_USER_TIMEZONE                 = "TimeZone"
	FIELD_USER_DATEFORMAT               = "DateFormat"
	FIELD_USER_LASTLOGINDATE            = "LastLoginDate"
	FIELD_USER_LASTLOGINIP              = "LastLoginIP"
	FIELD_USER_LOGINATTEMPTS            = "LoginAttempts"
	FIELD_USER_LOCKED                   = "Locked"
	FIELD_USER_ENFORCEPASSWORDCHANGE    = "EnforcePasswordChange"
	FIELD_SECONDARYPHONEINFO_VALUE      = "Value"
	FIELD_SECONDARYPHONEINFO_NUMERIC    = "Numeric"
	FIELD_SECONDARYPHONEINFO_DIALCODE   = "DialCode"
	FIELD_SECONDARYPHONEINFO_COUNTRYISO = "CountryISO"
	FIELD_FEATUREGROUP_ID               = "Id"
	FIELD_FEATUREGROUP_NAME             = "Name"
	FIELD_FEATUREGROUP_ACCOUNTTYPE      = "AccountType"
	FIELD_COUNTRY_ID                    = "Id"
	FIELD_COUNTRY_ISO                   = "Iso"
	FIELD_COUNTRY_NAME                  = "Name"
	FIELD_PREFERENCE_KEY                = "Key"
	FIELD_PREFERENCE_VALUE              = "Value"
	FIELD_ACTIVITYLOG_ID                = "Id"
	FIELD_ACTIVITYLOG_ACCOUNTID         = "AccountId"
	FIELD_ACTIVITYLOG_USERID            = "UserId"
	FIELD_ACTIVITYLOG_ENTITY            = "Entity"
	FIELD_ACTIVITYLOG_ENTITYID          = "EntityId"
	FIELD_ACTIVITYLOG_ACTION            = "Action"
	FIELD_ACTIVITYLOG_VALUE             = "Value"
	FIELD_PASSWORDRESET_ID              = "Id"
	FIELD_PASSWORDRESET_USERID          = "UserId"
	FIELD_PASSWORDRESET_COMPLETE        = "Complete"
	FIELD_PASSWORDRESET_URL             = "Url"
	FIELD_FEATURE_ID                    = "Id"
	FIELD_FEATURE_KEY                   = "Key"
	FIELD_FEATURE_NAME                  = "Name"
	FIELD_FEATURE_DESCRIPTION           = "Description"
	FIELD_FEATURE_FEATUREGROUPID        = "FeatureGroupId"
	FIELD_PASSWORD_ID                   = "Id"
	FIELD_PASSWORD_VALUE                = "Value"
	FIELD_SERVERSETTING_ID              = "Id"
	FIELD_SERVERSETTING_NAME            = "Name"
	FIELD_SERVERSETTING_KEY             = "Key"
	FIELD_SERVERSETTING_CATEGORY        = "Category"
	FIELD_SERVERSETTING_VALUE           = "Value"
	FIELD_SERVERSETTING_ANY             = "Any"
	FIELD_FILEOBJECT_ID                 = "Id"
	FIELD_FILEOBJECT_NAME               = "Name"
	FIELD_FILEOBJECT_PATH               = "Path"
	FIELD_FILEOBJECT_SINGLEDOWNLOAD     = "SingleDownload"
	FIELD_FILEOBJECT_CONTENT            = "Content"
	FIELD_FILEOBJECT_SIZE               = "Size"
	FIELD_FILEOBJECT_TYPE               = "Type"
	FIELD_FILEOBJECT_MODIFIEDUNIX       = "ModifiedUnix"
	FIELD_FILEOBJECT_MODIFIED           = "Modified"
	FIELD_FILEOBJECT_MD5                = "MD5"
	FIELD_FILEOBJECT_ACCOUNTID          = "AccountId"
)
