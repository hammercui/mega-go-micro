/**
 * Description
 * version 1.0.0
 * Created by GoLand.
 * Company sdbean
 * Author: hammercui
 * Date: 2021/1/7
 * Time: 14:24
 * Mail: hammercui@163.com
 *
 */
package watch

import (
	"encoding/json"
	simple "github.com/bitly/go-simplejson"
	"testing"
)

var jsonStr = "[{\"host\":\"192.168.2.20\",\"port\":9092}]"

func TestArraryParseTest(t *testing.T) {
	sj := simple.New()
	var bytes = []byte(jsonStr)
	sj.UnmarshalJSON(bytes)
	item,err := sj.Array()
	if err != nil || len(item) == 0{
		t.Fatal(err)
	}else{
		t.Logf("item 0:%+v",toJsonStr(item[0]))
	}
}

func toJsonStr(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}