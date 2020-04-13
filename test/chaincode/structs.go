package main

type TestStruct struct {
	DocType          string                                  `json:"docType" istate:"TestStruct_docType"`
	ID               string                                  `json:"id" istate:"TestStruct_id" primary:"true"`
	AnArray          []int                                   `json:"anArray" istate:"TestStruct_anArray"`
	AMap             map[int]int                             `json:"aMap" istate:"TestStruct_aMap"`
	AStruct          SomeStruct                              `json:"aStruct" istate:"TestStruct_aStruct"`
	AnInt            int64                                   `json:"anInt" istate:"TestStruct_anInt"`
	A3DArray         [][][]int                               `json:"a3DArray" istate:"TestStruct_a3DArray"`
	A2DArray         [][]int                                 `json:"a2DArray" istate:"TestStruct_a2DArray"`
	AComplexMapSlice []map[string][]map[int]struct{}         `json:"aComplexMapSlice" istate:"TestStruct_aComplexMapSlice"`
	AMapStruct       []map[string]SomeStruct                 `json:"aMapStruct" istate:"TestStruct_aMapStruct"`
	AMultiStruct     MultiStruct                             `json:"aMultiStruct" istate:"TestStruct_aMultiStruct"`
	AMultiMap        map[string]map[string]map[string]string `json:"aMultiMap" istate:"TestStruct_aMultiMap"`
}

// type TestStruct struct {
// 	DocType string         `json:"docType" istate:"0"`
// 	ID      string         `json:"id" istate:"1" primary:"true"`
// 	AnArray []int          `json:"anArray" istate:"2"`
// 	AMap    map[string]int `json:"aMap" istate:"3"`
// 	AStruct SomeStruct     `json:"aStruct" istate:"4"`
// }

type SomeStruct struct {
	Val string `json:"val"`
}

type MultiStruct struct {
	MultiVal SomeStruct `json:"multiVal"`
}

type ReadStateInput struct {
	ID string `json:"id"`
}

type DeleteStateInput struct {
	ID string `json:"id"`
}

type QueryInput struct {
	QueryString string `json:""queryString`
}

type QueryOut struct {
	Result []TestStruct `json:"result"`
	Count  int          `json:"count"`
}
