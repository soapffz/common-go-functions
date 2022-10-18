package pkg

// 移除数组中重复元素并返回数组 by soapffz 2022-10-18

func RemoveRepeatedElementInSlice(arr []string) []string {
	newArr := make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
