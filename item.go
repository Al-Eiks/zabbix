package zabbix

import (
	"fmt"
	"github.com/AlekSi/reflector"
)

type (
	ItemType  int
	ValueType int
	DataType  int
	DeltaType int
)

const (
	ZabbixAgent       ItemType = 0
	SNMPv1Agent       ItemType = 1
	ZabbixTrapper     ItemType = 2
	SimpleCheck       ItemType = 3
	SNMPv2Agent       ItemType = 4
	ZabbixInternal    ItemType = 5
	SNMPv3Agent       ItemType = 6
	ZabbixAgentActive ItemType = 7
	ZabbixAggregate   ItemType = 8
	WebItem           ItemType = 9
	ExternalCheck     ItemType = 10
	DatabaseMonitor   ItemType = 11
	IPMIAgent         ItemType = 12
	SSHAgent          ItemType = 13
	TELNETAgent       ItemType = 14
	Calculated        ItemType = 15
	JMXAgent          ItemType = 16

	Float     ValueType = 0
	Character ValueType = 1
	Log       ValueType = 2
	Unsigned  ValueType = 3
	Text      ValueType = 4

	Decimal     DataType = 0
	Octal       DataType = 1
	Hexadecimal DataType = 2
	Boolean     DataType = 3

	AsIs  DeltaType = 0
	Speed DeltaType = 1
	Delta DeltaType = 2
)

// https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/definitions
type Item struct {
	ItemId      string    `json:"itemid,omitempty"`
	Delay       int       `json:"delay"`
	HostId      string    `json:"hostid"`
	InterfaceId string    `json:"interfaceid,omitempty"`
	Key         string    `json:"key_"`
	Name        string    `json:"name"`
	Type        ItemType  `json:"type"`
	ValueType   ValueType `json:"value_type"`
	DataType    DataType  `json:"data_type"`
	Delta       DeltaType `json:"delta"`
	Description string    `json:"description"`
	Error       string    `json:"error"`
	History     int       `json:"history,omitempty"`
	Trends      int       `json:"trends,omitempty"`

	// Fields below used only when creating applications
	ApplicationIds []string `json:"applications,omitempty"`
}

type Items []Item

// Converts slice to map by key. Panics if there are duplicate keys.
func (items Items) ByKey() (res map[string]Item) {
	res = make(map[string]Item, len(items))
	for _, i := range items {
		_, present := res[i.Key]
		if present {
			panic(fmt.Errorf("Duplicate key %s", i.Key))
		}
		res[i.Key] = i
	}
	return
}

// Wrapper for item.get https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/get
func (api *API) ItemsGet(params Params) (res Items, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	response, err := api.CallWithError("item.get", params)
	if err != nil {
		return
	}

	reflector.MapsToStructs2(response.Result.([]interface{}), &res, reflector.Strconv, "json")
	return
}

// Gets items by application Id.
func (api *API) ItemsGetByApplicationId(id string) (res Items, err error) {
	return api.ItemsGet(Params{"applicationids": id})
}

// Wrapper for item.create: https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/create
func (api *API) ItemsCreate(items Items) (err error) {
	response, err := api.CallWithError("item.create", items)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids := result["itemids"].([]interface{})
	for i, id := range itemids {
		items[i].ItemId = id.(string)
	}
	return
}

// Wrapper for item.delete: https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/delete
// Cleans ItemId in all items elements if call succeed.
func (api *API) ItemsDelete(items Items) (err error) {
	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ItemId
	}

	err = api.ItemsDeleteByIds(ids)
	if err == nil {
		for i := range items {
			items[i].ItemId = ""
		}
	}
	return
}

// Wrapper for item.delete: https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/delete
func (api *API) ItemsDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("item.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	itemids1, ok := result["itemids"].([]interface{})
	l := len(itemids1)
	if !ok {
		// some versions actually return map there
		itemids2 := result["itemids"].(map[string]interface{})
		l = len(itemids2)
	}
	if len(ids) != l {
		err = &ExpectedMore{len(ids), l}
	}
	return
}

// Wrapper for item.get https://www.zabbix.com/documentation/2.0/manual/appendix/api/item/get
func (api *API) GetInterfaceItemProd (nameVoisin string, params Params) (items []string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("item.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if strings.Contains(tmp["key_"].(string), "alias") {
                        parser := strings.Contains(tmp["prevvalue"].(string), nameVoisin)
                        p1 := strings.Contains(nameVoisin, "PRDNETRHP")
                        p2 := strings.Contains(tmp["prevvalue"].(string), "PRDNETRHP")
                        if ((parser) || ((p1) && (p2))) {
                                testAlias := strings.Contains(tmp["key_"].(string), "alias_admin")
                                testAlias2 := strings.Contains(tmp["key_"].(string), "alias_prod")
                                if ((testAlias) || (testAlias2)) {
                                        continue
                                } else {
                                        itemKey := tmp["key_"].(string)
                                        itemKey = strings.TrimPrefix(itemKey, "alias[")
                                        itemKey = strings.TrimPrefix(itemKey, "alias_admin[")
                                        itemKey = strings.TrimPrefix(itemKey, "alias_prod[")
                                        itemKey = strings.TrimSuffix(itemKey, "]")
                                        items = append(items, itemKey)
                                }
                        }
                }
        }
        n := len(items)
        if n == 1 {
                return
        } else {
                var items2 []string
                for i, _ := range items {
                        test2 := strings.Contains(items[i], "GigabitEthernet")
                        if test2  {
                                test4 := StringInSlice(items2, "Aggregation")
                                if test4 == false {
                                        items2 = append(items2, items[i])
                                }
                        } else {
                                test3 := StringInSlice(items2, "GigabitEthernet")
                                if test3 {
                                        items2 = nil
                                }
                        items2 = append(items2, items[i])
                        }
                }
             items = items2
        }
        return
}

//permet de tester le contenu d'une slice
func SliceContains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }
    _, ok := set[item]
    return ok
}

func StringInSlice(list []string, a string) bool {
    for _, b := range list {
        if strings.Contains(b, a) {
            return true
        }
    }
    return false
}


func (api *API) GetItemId(key string, params Params) (itemId string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("item.get", params)
        if err != nil {
                fmt.Println(err.Error())
                return
        }
        result := response.Result.([]interface{})
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if tmp["key_"].(string) == key {
                        itemId = tmp["itemid"].(string)
                }
        }
        return
}

func (api *API) GetInterfaces(nameVoisin string, params Params) (items []string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("item.get", params)
        if err != nil {
                fmt.Println(err.Error())
                return
        }
        result := response.Result.([]interface{})
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if strings.Contains(tmp["key_"].(string), "alias") {
                        fmt.Println(tmp["key_"].(string), tmp["prevvalue"].(string))
                        testAlias := strings.Contains(tmp["key_"].(string), "alias_admin")
                        testAlias2 := strings.Contains(tmp["key_"].(string), "alias_prod")
                        if ((testAlias) || (testAlias2)) {
                                continue
                        }
                        item := tmp["key_"].(string)
                        item = strings.TrimPrefix(item, "alias[")
                        item = strings.TrimPrefix(item, "alias_admin[")
                        item = strings.TrimPrefix(item, "alias_prod[")
                        item = strings.TrimSuffix(item, "]")
                        if strings.Contains(tmp["prevvalue"].(string), nameVoisin) {
                                items = append(items, item)
                        } else if strings.Contains(nameVoisin, "520") {
                                test520 := strings.Contains(tmp["prevvalue"].(string), "520")
                                test521 := strings.Contains(tmp["prevvalue"].(string), "521")
                                test522 := strings.Contains(tmp["prevvalue"].(string), "522")
                                if test520 || test521 || test522 {
                                        items = append(items, item)
                                }
                        } else if strings.Contains(nameVoisin, "510") {
                                test510 := strings.Contains(tmp["prevvalue"].(string), "510")
                                test511 := strings.Contains(tmp["prevvalue"].(string), "511")
                                test512 := strings.Contains(tmp["prevvalue"].(string), "512")
                                if test510 || test511 || test512 {
                                        items = append(items, item)
                                }
                        } else if strings.Contains(nameVoisin, "520") {
                                test500 := strings.Contains(tmp["prevvalue"].(string), "500")
                                test501 := strings.Contains(tmp["prevvalue"].(string), "501")
                                test502 := strings.Contains(tmp["prevvalue"].(string), "502")
                                if test500 || test501 || test502 {
                                        items = append(items, item)
                                }
                        } else if strings.Contains(nameVoisin, "PRDNETRHP") {
                                if strings.Contains(tmp["prevvalue"].(string), "PRDNETRHP") {
                                        items = append(items, item)
                                }
                        }
                }
        }
        fmt.Println(items)
        return
}

func (api *API) GetInterfaceFromItem(nameVoisin string, params Params) (items []string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("item.get", params)
        if err != nil {
                fmt.Println(err.Error())
                return
        }
        result := response.Result.([]interface{})
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if strings.Contains(tmp["key_"].(string), "alias") {
                       parser := strings.Contains(tmp["prevvalue"].(string), nameVoisin)
                       p1 := strings.Contains(nameVoisin, "PRDNETRHP")
                       p2 := strings.Contains(tmp["prevvalue"].(string), "PRDNETRHP")
                       if ((parser) || ((p1) && (p2))) {
                                testAlias := strings.Contains(tmp["key_"].(string), "alias_admin")
                                testAlias2 := strings.Contains(tmp["key_"].(string), "alias_prod")
                                if ((testAlias) || (testAlias2)) {
                                        continue
                                } else {
                                        itemKey := tmp["key_"].(string)
                                        itemKey = strings.TrimPrefix(itemKey, "alias[")
                                        itemKey = strings.TrimPrefix(itemKey, "alias_admin[")
                                        itemKey = strings.TrimPrefix(itemKey, "alias_prod[")
                                        itemKey = strings.TrimSuffix(itemKey, "]")
                                        items = append(items, itemKey)
                                }
                        }
                }
        }
        n := len(items)
        if n == 1 {
                return
        } else {
                var items2 []string
                for i, _ := range items {
                        test2 := strings.Contains(items[i], "GigabitEthernet")
                        if test2  {
                                test4 := StringInSlice(items2, "Aggregation")
                                if test4 == false {
                                        items2 = append(items2, items[i])
                                }
                        } else {
                                test3 := StringInSlice(items2, "GigabitEthernet")
                                if test3 {
                                        items2 = nil
                                }
                        items2 = append(items2, items[i])
                        }
                }
             items = items2
        }
        return
}

func (api *API) GetNeighbors(params Params) (items3 []string, err error) {
        if _, present := params["output"]; !present {
                params["output"] = "extend"
        }
        response, err := api.CallWithError("item.get", params)
        if err != nil {
                return
        }
        result := response.Result.([]interface{})
        var items2 []string
        for _, i := range result {
                tmp := i.(map[string]interface{})
                if strings.Contains(tmp["key_"].(string), "alias") {
                        testAlias := strings.Contains(tmp["key_"].(string), "alias_admin")
                        testAlias2 := strings.Contains(tmp["key_"].(string), "alias_prod")
                        if ((testAlias) || (testAlias2)) {
                                continue
                        } else {
                                itemKey := tmp["key_"].(string)
                                itemKey = strings.TrimPrefix(itemKey, "alias[")
                                itemKey = strings.TrimPrefix(itemKey, "alias_admin[")
                                itemKey = strings.TrimPrefix(itemKey, "alias_prod[")
                                itemKey = strings.TrimSuffix(itemKey, "]")
                                items2 = append(items2, tmp["prevvalue"].(string))
                        }
                }
        }
        sort.Strings(items2)
        for i, _ := range items2 {
                if strings.Contains(items2[i], "Vers") || strings.Contains(items2[i], "LS") || strings.Contains(items2[i], "Portable") || items2[i] == "0" {
                        continue
                }
                element := string(items2[i][0:12])
                if strings.Contains(element, "PRDNETRHP") {
                        element = "PRDNETRHP500"
                }
                if !StringInSlice(items3, element) {
                        items3 = append(items3, element)
                }
        }
        return
}
