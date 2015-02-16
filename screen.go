package zabbix

import (
//        "encoding/csv"
//        "fmt"
//      "strconv"
//       "log"
//        "strings"
//        "os"
//        "github.com/AlekSi/zabbix"
//      "flag"
)


func (api *API) GetScreenElem(screenName string, params Params) (screenItems []string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("screen.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if tmp["name"].(string) == screenName {
                        tmp2 := tmp["screenitems"].([]interface{})
                        for _,j := range tmp2 {
                                k := j.(map[string]interface{})
                                if k["resourcetype"] == "0" {
                                        item := k["resourceid"].(string)
                                        screenItems = append(screenItems, item)
                                }
                        }
                   }
        }
        return
}


func (api *API) CheckScreen(screenName string, params Params) (screenId string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("screen.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if tmp["name"].(string) == screenName {
                        screenId = tmp["screenid"].(string)
                }
        }
        return
}
