package zabbix


import (
//       "fmt"
        "strings"
)


type Graph struct {
//      Graphid    string       `json:"graphid,omitempty"`
        Name       string       `json:"name"`
        Gitems     GraphItems   `json:"gitems,omitempty"`
        Height     int          `json:"height"`
        Width      int          `json:"width"`
}
type Graphs []Graph


type GraphItem struct {
//      Gitemid    string      `json:"gitemid,omitempty"`
        Color      string      `json:"color"`
        Itemid     string      `json:"itemid"`
}
type GraphItems []GraphItem



func (api *API) GraphGet(intName string, params Params) (graphIds []string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("graph.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _,i := range result {
                tmp := i.(map[string]interface{})
                graphId := tmp["graphid"].(string)
                if strings.Contains(tmp["name"].(string), intName) {
                            graphIds = append(graphIds, graphId)
                }
        }
        return
}

func (api *API) GetGraphName(params Params) (graphName string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("graph.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _,i := range result {
                tmp := i.(map[string]interface{})
                graphName = tmp["name"].(string)
        }
        return
}

func (api *API) GetItemKey(params Params) (itemKey string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("graphitem.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if tmp["key_"].(string) != "" {
                        itemKey = tmp["key_"].(string)
                }
        }
        return
}


func (api *API) GetGraphDetails(params Params) (graphDetails []interface{}, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("graphitem.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        graphDetails = result
        return
}

func (api *API) CheckHostPresence(hostId string, params Params) (res bool, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("graphitem.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _,i := range result {
                tmp := i.(map[string]interface{})
                if tmp["hostid"].(string) == hostId {
                        res = true
                        break
                } else {
                        res = false
                }
        }
        return
}


func (api *API) GetGraphItems(graphId string, params Params) (graphItems []string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("graphitem.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _,i := range result {
                tmp := i.(map[string]interface{})
                if tmp["itemid"].(string) != "" {
                        graphItem := tmp["itemid"].(string)
                        graphItems = append(graphItems, graphItem)
                }
        }

        return
}


func (api *API) GetGraphItemColor(graphItemId string, params Params) (graphItemColor string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("graphitem.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _,i := range result {
                tmp := i.(map[string]interface{})
                if tmp["itemid"].(string) == graphItemId {
                        graphItemColor = tmp["color"].(string)
                }
        }
        return
}
