package pkg

import (
	"sort"
	"strings"
)

// 将string+int类型的map排序后返回，可指定升序还是降序，默认为降序,by soapffz 2022-11-17

type kv struct {
	Key   string
	Value int
}

func SortAMapOfStringAndIntByValue(m map[string]int, Order string) map[string]int {
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		if strings.Contains(Order, "asc") { // 升序
			return ss[i].Value < ss[j].Value
		}
		if strings.Contains(Order, "desc") { // 降序
			return ss[i].Value > ss[j].Value
		}
		return ss[i].Value > ss[j].Value // 降序
	})
	return nil
}
