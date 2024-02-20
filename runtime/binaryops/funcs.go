package binaryops

import (
	"banek/runtime/objs"
	"strings"
)

func addInts(left, right objs.Obj) (objs.Obj, bool) {
	return objs.Obj{Int: left.Int + right.Int, Type: objs.Int}, true
}

func subInts(left, right objs.Obj) (objs.Obj, bool) {
	return objs.Obj{Int: left.Int - right.Int, Type: objs.Int}, true
}

func mulInts(left, right objs.Obj) (objs.Obj, bool) {
	return objs.Obj{Int: left.Int * right.Int, Type: objs.Int}, true
}

func divInts(left, right objs.Obj) (objs.Obj, bool) {
	if right.Int == 0 {
		return objs.Obj{}, false
	}

	return objs.Obj{Int: left.Int / right.Int, Type: objs.Int}, true
}

func modInts(left, right objs.Obj) (objs.Obj, bool) {
	if right.Int == 0 {
		return objs.Obj{}, false
	}

	return objs.Obj{Int: left.Int % right.Int, Type: objs.Int}, true
}

func addStrings(left, right objs.Obj) (objs.Obj, bool) {
	str1 := left.AsString()
	str2 := right.AsString()

	return objs.MakeString(str1 + str2), true
}

func repeatStrings(left, right objs.Obj) (objs.Obj, bool) {
	str := left.AsString()
	count := right.Int

	if count < 0 {
		return objs.Obj{}, false
	}

	return objs.MakeString(strings.Repeat(str, count)), true
}

func concatArrays(left, right objs.Obj) (objs.Obj, bool) {
	arr1 := left.AsArray()
	arr2 := right.AsArray()

	return objs.MakeArray(append(arr1, arr2...)), true
}

func repeatArray(left, right objs.Obj) (objs.Obj, bool) {
	arr := left.AsArray()
	count := right.Int

	if count < 0 {
		return objs.Obj{}, false
	}

	newArr := make([]objs.Obj, len(arr)*count)
	for i := 0; i < count; i++ {
		copy(newArr[i*len(arr):], arr)
	}

	return objs.MakeArray(newArr), true
}
