package binaryops

import (
	"banek/runtime/objs"
	"slices"
	"strconv"
	"strings"
)

func addInts(left, right objs.Obj) objs.Obj {
	return objs.MakeInt(left.AsInt() + right.AsInt())
}

func subInts(left, right objs.Obj) objs.Obj {
	return objs.MakeInt(left.AsInt() - right.AsInt())
}

func mulInts(left, right objs.Obj) objs.Obj {
	return objs.MakeInt(left.AsInt() * right.AsInt())
}

func divInts(left, right objs.Obj) objs.Obj {
	return objs.MakeInt(left.AsInt() / right.AsInt())
}

func modInts(left, right objs.Obj) objs.Obj {
	return objs.MakeInt(left.AsInt() % right.AsInt())
}

func addStrings(left, right objs.Obj) objs.Obj {
	str1 := left.AsString()
	str2 := right.AsString()

	return objs.MakeString(str1 + str2)
}

type NegativeRepeatCountError int

func (err NegativeRepeatCountError) Error() string {
	return "negative repeat count: " + strconv.Itoa(int(err))
}

func repeatStrings(left, right objs.Obj) objs.Obj {
	str := left.AsString()
	count := right.AsInt()

	if count < 0 {
		panic(NegativeRepeatCountError(count))
	}

	return objs.MakeString(strings.Repeat(str, count))
}

func concatArrays(left, right objs.Obj) objs.Obj {
	arr1 := left.AsArray()
	arr2 := right.AsArray()

	return objs.MakeArray(slices.Concat(arr1, arr2))
}

func repeatArray(left, right objs.Obj) objs.Obj {
	arr := left.AsArray()
	count := right.AsInt()

	if count < 0 {
		panic(NegativeRepeatCountError(count))
	}

	newArr := make([]objs.Obj, len(arr)*count)
	for i := 0; i < count; i++ {
		copy(newArr[i*len(arr):], arr)
	}

	return objs.MakeArray(newArr)
}
