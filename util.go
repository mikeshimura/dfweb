package dfweb

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/mikeshimura/dbflute/df"
	"io/ioutil"
	"reflect"
	"strconv"
	"time"
)

var OpMap map[string]string
var OpMap2 map[string]string

func ErrorRecover(c *gin.Context, tx *sql.Tx) {
	errx := recover()
	if errx != nil {
		df.TxRollback(tx)
		errs:=fmt.Sprintf("%v",errx)
		if errs==""{
			errs="System Error"
		}
		rmap := SetErrorMessage(errs)
		c.JSON(200, rmap)
	} else {
		df.TxCommit(tx)
	}
}

func GetBody(c *gin.Context) []byte {
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		panic("Request Body Read  Error")
	}
	return body
}

func GetBodyString(c *gin.Context) string {
	return string(GetBody(c))
}
func GetBodyJson(c *gin.Context) map[string]interface{} {
	var dat map[string]interface{}
	byt := GetBody(c)
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	return dat
}

func SetErrorMessage(errorMessage string) map[string]interface{} {
	return SetFetchResult(errorMessage, -1, 0, 0)
}
func SetNormalFetchResult(ar *df.List) map[string]interface{} {
	return SetFetchResult(ar.GetAsArray(), 0, 0, ar.Size())
}
func SetSingleFetchResult(data interface{}) map[string]interface{} {
	return SetFetchResult(data, 0, 0, 1)
}
func SetFetchResult(data interface{}, status int, startRow int,
	totalRow int) map[string]interface{} {
	response := GetResponse(data, status, startRow, totalRow)
	var result map[string]interface{} = make(map[string]interface{})
	result["response"] = response
	return result
}
func GetResponse(data interface{}, status int, startRow int,
	totalRow int) map[string]interface{} {
	var dat map[string]interface{} = make(map[string]interface{})
	dat["status"] = status
	dat["startRow"] = startRow
	dat["endRow"] = startRow + totalRow - 1
	dat["totalRows"] = totalRow
	dat["data"] = data
	return dat
}

func CreateMd5(data string) string {
	bytes := []byte(data)
	res := fmt.Sprintf("%x", md5.Sum(bytes))
	fmt.Println("hash " + res)
	return res
}

func NewEntityToMap(entity *df.Entity, elems []string) map[string]interface{} {
	rmap := make(map[string]interface{})
	EntityToMap(entity, rmap, elems)
	return rmap
}
func EntityToMap(entity *df.Entity, rmap map[string]interface{}, elems []string) {
	for _, e := range elems {
		var edata interface{} = df.GetEntityValue(entity, e)
		rmap[e] = ConvWebData(edata)
	}
}
func MapToEntity(rmap map[string]interface{}, entity *df.Entity, table string, update bool) {
	meta := df.DBMetaProvider_I.TableDbNameInstanceMap[table]
	for propertyName := range rmap {
		colInfo := (*meta).GetColumnInfoByPropertyName(propertyName)
		if colInfo == nil {
			continue
		}
		value := rmap[propertyName]
		if value == nil{
			continue
		}
		argType := df.GetType(value)
		if argType == "string" {
			svalue := value.(string)
			if svalue == "" {
				EntitySetNull(entity, propertyName, colInfo, update)
				continue
			}
		}
		df.SetEntityValue(entity, propertyName,
			ConvFromWebData(value, colInfo, argType))
	}
}
func EntitySetNull(
	entity *df.Entity, propertyName string, colInfo *df.ColumnInfo,update bool) {
	switch colInfo.GoType {
	case "string":
		df.SetEntityValue(entity, propertyName, "")
	case "sql.NullString":
		df.SetEntityValue(entity, propertyName, new(sql.NullString))
	case "pq.NullTime":
		nt := new(pq.NullTime)
		nt.Valid = false
		df.SetEntityValue(entity, propertyName, nt)
	case "df.NullDate":
		df.SetEntityValue(entity, propertyName, new(df.NullDate))
	case "df.NullTimestamp":
		df.SetEntityValue(entity, propertyName, new(df.NullTimestamp))
	case "df.NullNumeric":
		df.SetEntityValue(entity, propertyName, new(df.NullNumeric))
	case "df.MysqlNullDate":
		df.SetEntityValue(entity, propertyName, new(df.MysqlNullDate))
	case "df.MysqlNullTime":
		df.SetEntityValue(entity, propertyName, new(df.MysqlNullTime))
	case "df.MysqlNullTimestamp":
		df.SetEntityValue(entity, propertyName, new(df.MysqlNullTimestamp))
	case "sql.NullInt64":
		df.SetEntityValue(entity, propertyName, new(sql.NullInt64))
	case "sql.NullFloat64":
		df.SetEntityValue(entity, propertyName, new(sql.NullFloat64))
	case "sql.NullBool":
		df.SetEntityValue(entity, propertyName, new(sql.NullBool))
	default:
		if update {
			panic("必須項目未入力です。 項目名:" + propertyName)
		}
	}
}
func ConvFromWebData(
	arg interface{}, colInfo *df.ColumnInfo, argType string) interface{} {
	goType := colInfo.GoType
	if argType == goType {
		return arg
	}
	switch goType {
	case "sql.NullString":
		switch argType {
		case "string":
			return df.CreateNullString(arg.(string))
		}
	case "int64":
		switch argType {
		case "string":
			cv, err := strconv.ParseInt(arg.(string), 10, 64)
			if err != nil {
				panic("int64に変換出来ません:" + arg.(string))
			}
			return cv
		case "float64":
			return int64(arg.(float64))
		}
	case "sql.NullInt64":
		switch argType {
		case "string":
			cv, err := strconv.ParseInt(arg.(string), 10, 64)
			if err != nil {
				panic("int64に変換出来ません:" + arg.(string))
			}
			return df.CreateNullInt64(cv)
		case "float64":
			return df.CreateNullInt64(int64(arg.(float64)))
		}
	case "float64":
		switch argType {
		case "string":
			cv, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				panic("float64に変換出来ません:" + arg.(string))
			}
			return cv
		case "float64":
			return arg.(float64)
		}
	case "sql.NullFloat64":
		switch argType {
		case "string":
			cv, err := strconv.ParseFloat(arg.(string), 64)
			if err != nil {
				panic("float64に変換出来ません:" + arg.(string))
			}
			return df.CreateNullFloat64(cv)
		case "float64":
			return df.CreateNullFloat64(arg.(float64))
		}
	case "sql.NullBool":
		switch argType {
		case "string":
			cv, err := strconv.ParseBool(arg.(string))
			if err != nil {
				panic("Booleanに変換出来ません:" + arg.(string))
			}
			return df.CreateNullBool(cv)
		case "bool":
			return df.CreateNullBool(arg.(bool))
		}
	case "df.Numeric":
		switch argType {
		case "string":
			cv, err := df.CreateNumeric(arg.(string))
			if err != nil {
				panic("Numericに変換出来ません:" + arg.(string))
			}
			return cv
		}
	case "df.NullNumeric":
		switch argType {
		case "string":
			cv, err := df.CreateNullNumeric(arg.(string))
			if err != nil {
				panic("Numericに変換出来ません:" + arg.(string))
			}
			return cv
		}
	case "time.Time":
		switch argType {
		case "string":
			res, err := time.Parse(df.DISP_SQL_DEFAULT_TIME_FORMAT, arg.(string))
			if err != nil {
				panic("Timeに変換出来ません。：" + arg.(string))
			}
			return res
		}
	case "df.Date":
		switch argType {
		case "string":
			res, err := time.Parse(df.DISP_SQL_DEFAULT_DATE_FORMAT, arg.(string))
			if err != nil {
				panic("Timeに変換出来ません。：" + arg.(string))
			}
			return df.CreateDate(res)
		}
	case "df.Timestamp":
		switch argType {
		case "string":
			args := arg.(string)
			parse := df.DISP_SQL_DEFAULT_TIMESTAMP_FORMAT
			if len(args) == len(df.DISP_SQL_DEFAULT_DATE_FORMAT) {
				parse = df.DISP_SQL_DEFAULT_DATE_FORMAT
			}
			res, err := time.Parse(parse, args)
			if err != nil {
				panic("Timestampに変換出来ません。：" + args)
			}
			return df.CreateTimestamp(res)
		}
	case "pq.NullTime":
		switch argType {
		case "string":
			res, err := time.Parse(df.DISP_SQL_DEFAULT_TIME_FORMAT, arg.(string))
			if err != nil {
				panic("Timeに変換出来ません。：" + arg.(string))
			}
			return df.CreateNullTime(res)
		}
	case "df.NullDate":
		switch argType {
		case "string":
			res, err := time.Parse(df.DISP_SQL_DEFAULT_DATE_FORMAT, arg.(string))
			if err != nil {
				panic("Timeに変換出来ません。：" + arg.(string))
			}
			return df.CreateNullDate(res)
		}
	case "df.NullTimestamp":
		switch argType {
		case "string":
			args := arg.(string)
			parse := df.DISP_SQL_DEFAULT_TIMESTAMP_FORMAT
			if len(args) == len(df.DISP_SQL_DEFAULT_DATE_FORMAT) {
				parse = df.DISP_SQL_DEFAULT_DATE_FORMAT
			}
			res, err := time.Parse(parse, args)
			if err != nil {
				panic("Timestampに変換出来ません。：" + args)
			}
			return df.CreateNullTimestamp(res)
		}
	}

	if argType != goType {
		panic("この組み合わせは規定されていません。 argType:" + argType + " goType:" + goType)
	}
	return arg
}

func ConvWebData(arg interface{}) interface{} {
	switch arg.(type) {
	case time.Time:
		tv := arg.(time.Time)
		return tv.Format(df.DISP_SQL_DEFAULT_TIME_FORMAT)
	case *time.Time:
		tv := arg.(*time.Time)
		return (*tv).Format(df.DISP_SQL_DEFAULT_TIME_FORMAT)
	case df.Date:
		dv := arg.(df.Date)
		return dv.Date.Format(df.DISP_SQL_DEFAULT_DATE_FORMAT)
	case *df.Date:
		dv := arg.(*df.Date)
		return (*dv).Date.Format(df.DISP_SQL_DEFAULT_DATE_FORMAT)
	case df.Timestamp:
		dv := arg.(df.Timestamp)
		return dv.Timestamp.Format(df.DISP_SQL_DEFAULT_TIMESTAMP_FORMAT)
	case *df.Timestamp:
		dv := arg.(*df.Timestamp)
		return (*dv).Timestamp.Format(df.DISP_SQL_DEFAULT_TIMESTAMP_FORMAT)
	case df.Numeric:
		nv := arg.(df.Numeric)
		return nv.String()
	case *df.Numeric:
		nv := arg.(*df.Numeric)
		return (*nv).String()
	case df.MysqlDate:
		nv := arg.(df.MysqlDate)
		return nv.String()
	case *df.MysqlDate:
		nv := arg.(*df.MysqlDate)
		return (*nv).String()
	case df.MysqlTime:
		nv := arg.(df.MysqlTime)
		return nv.String()
	case *df.MysqlTime:
		nv := arg.(*df.MysqlTime)
		return (*nv).String()
	case df.MysqlTimestamp:
		nv := arg.(df.MysqlTimestamp)
		return nv.String()
	case *df.MysqlTimestamp:
		nv := arg.(*df.MysqlTimestamp)
		return (*nv).String()
	case sql.NullString:
		nsv := arg.(sql.NullString)
		if nsv.Valid {
			return nsv.String
		}
		return ""
	case *sql.NullString:
		nsv := arg.(*sql.NullString)
		if nsv.Valid {
			return (*nsv).String
		}
		return ""
	case pq.NullTime:
		ntv := arg.(pq.NullTime)
		if ntv.Valid {
			return ntv.Time.Format(df.DISP_SQL_DEFAULT_TIME_FORMAT)
		}
		return ""
	case *pq.NullTime:
		ntv := arg.(*pq.NullTime)
		if ntv.Valid {
			return (*ntv).Time.Format(df.DISP_SQL_DEFAULT_TIME_FORMAT)
		}
		return ""
	case df.NullDate:
		ndv := arg.(df.NullDate)
		if ndv.Valid {
			return ndv.Date.Format(df.DISP_SQL_DEFAULT_DATE_FORMAT)
		}
		return ""
	case *df.NullDate:
		ndv := arg.(*df.NullDate)
		if ndv.Valid {
			return (*ndv).Date.Format(df.DISP_SQL_DEFAULT_DATE_FORMAT)
		}
		return ""
	case df.NullTimestamp:
		ntsv := arg.(df.NullTimestamp)
		if ntsv.Valid {
			return fmt.Sprint(ntsv.Timestamp.Format(df.DISP_SQL_DEFAULT_TIMESTAMP_FORMAT))
		}
		return ""
	case *df.NullTimestamp:
		ntsv := arg.(*df.NullTimestamp)
		if ntsv.Valid {
			return fmt.Sprint((*ntsv).Timestamp.Format(df.DISP_SQL_DEFAULT_TIMESTAMP_FORMAT))
		}
		return ""
	case df.NullNumeric:
		nnv := arg.(df.NullNumeric)
		if nnv.Valid {
			return nnv.String()
		}
		return ""
	case *df.NullNumeric:
		nnv := arg.(*df.NullNumeric)
		if nnv.Valid {
			return (*nnv).String()
		}
		return ""
	case df.MysqlNullDate:
		nnv := arg.(df.MysqlNullDate)
		if nnv.Valid {
			return nnv.String()
		}
		return ""
	case *df.MysqlNullDate:
		nnv := arg.(*df.MysqlNullDate)
		if nnv.Valid {
			return (*nnv).String()
		}
		return ""
	case df.MysqlNullTime:
		nnv := arg.(df.MysqlNullTime)
		if nnv.Valid {
			return nnv.String()
		}
		return ""
	case *df.MysqlNullTime:
		nnv := arg.(*df.MysqlNullTime)
		if nnv.Valid {
			return (*nnv).String()
		}
		return ""
	case df.MysqlNullTimestamp:
		nnv := arg.(df.MysqlNullTimestamp)
		if nnv.Valid {
			return nnv.String()
		}
		return ""
	case *df.MysqlNullTimestamp:
		nnv := arg.(*df.MysqlNullTimestamp)
		if nnv.Valid {
			return (*nnv).String()
		}
		return ""
	case sql.NullInt64:
		nnv := arg.(sql.NullInt64)
		if nnv.Valid {
			return int(nnv.Int64)
		}
		return ""
	case *sql.NullInt64:
		nnv := arg.(*sql.NullInt64)
		if nnv.Valid {
			return int(nnv.Int64)
		}
		return ""
	case sql.NullFloat64:
		nnv := arg.(sql.NullFloat64)
		if nnv.Valid {
			return nnv.Float64
		}
		return ""
	case *sql.NullFloat64:
		nnv := arg.(*sql.NullFloat64)
		if nnv.Valid {
			return nnv.Float64
		}
		return ""
	case sql.NullBool:
		nnv := arg.(sql.NullBool)
		if nnv.Valid {
			return nnv.Bool
		}
		return ""
	case *sql.NullBool:
		nnv := arg.(*sql.NullBool)
		if nnv.Valid {
			return nnv.Bool
		}
		return ""

	case []byte:
		return "[]byte"
	case *[]byte:
		return "[]byte"
	}
	return arg
}
func Invoke(method reflect.Value, colInfo *df.ColumnInfo,
	svalue string, lso *df.LikeSearchOption) {
	param := ConvFromWebData(svalue, colInfo, "string")
	if lso != nil {
		method.Call([]reflect.Value{reflect.ValueOf(param), reflect.ValueOf(lso)})
	} else {
		method.Call([]reflect.Value{reflect.ValueOf(param)})
	}
}

func GetOpMap() map[string]string {
	if OpMap == nil {
		OpMap = make(map[string]string)
		OpMap["="] = "Equal"
		OpMap["<>"] = "NotEqual"
		OpMap[">"] = "GreaterThan"
		OpMap["<"] = "LessThan"
		OpMap[">="] = "GreaterEqual"
		OpMap["<="] = "LessEqual"
		OpMap["contains"] = "LikeSearch"
		OpMap["starts with"] = "LikeSearch"
		OpMap["ends with"] = "LikeSearch"
		OpMap["like"] = "LikeSearch"
		OpMap["between"] = "GreaterEqual"
		OpMap["exclude"] = "LessThan"
	}
	return OpMap
}
func GetOpMap2() map[string]string {
	if OpMap2 == nil {
		OpMap2 = make(map[string]string)
		OpMap2["between"] = "LessEqual"
		OpMap2["exclude"] = "GreaterThan"
	}
	return OpMap2
}
func SetCriteria(query interface{}, emap map[string]interface{}, table string) {
	meta := df.DBMetaProvider_I.TableDbNameInstanceMap[table]
	opMap := GetOpMap()
	opMap2 := GetOpMap2()
	field := (emap["fieldName"]).(string)
	colInfo := (*meta).GetColumnInfoByPropertyName(field)
	operator := (emap["operator"]).(string)
	op := opMap[operator]
	setter := "Set" + df.InitCap(field) + "_" + op
	op2 := opMap2[operator]
	setter2 := "Set" + df.InitCap(field) + "_" + op2
	start := (emap["start"]).(string)
	end := (emap["end"]).(string)
	fmt.Println("setter:" + setter)
	fmt.Println("sette2:" + setter2)
	fmt.Printf(" start:%v%T\n", start, start)
	fmt.Printf("end:%v%T\n", end, end)
	if operator == "" || start == "" {
		return
	}
	setterMethod := reflect.ValueOf(query).MethodByName(setter)
	setterMethod2 := reflect.ValueOf(query).MethodByName(setter2)
	if op == "LikeSearch" || op == "NotLikeSearch" {
		lso := new(df.LikeSearchOption)
		if operator == "contains" {
			lso.LikeContain()
		}
		if operator == "starts with" {
			lso.LikePrefix()
		}
		if operator == "ends with" {
			lso.LikeSuffix()
		}
		Invoke(setterMethod, colInfo, start, lso)
	} else {
		Invoke(setterMethod, colInfo, start, nil)
	}

	if setterMethod2.IsValid() {
		if end == "" {
			panic("2番目の項目のデータが未入力です。")
		}
		Invoke(setterMethod2, colInfo, end, nil)
	}

}
