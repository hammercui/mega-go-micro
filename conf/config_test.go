package conf

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_parseFlag(t *testing.T) {
	parseFlag()
	if flagConf.env == ""{
		t.Fatalf("parse err")
	}else{
		fmt.Printf("parse success: %+v \n",flagConf)
	}
}

func TestInitConfig(t *testing.T)  {
	InitConfig()
	if conf.App == nil{
		t.Fatalf("init config err")
	}else{
		str,_ := json.Marshal(conf)
		fmt.Printf("init config success: %s \n",str)
	}
}