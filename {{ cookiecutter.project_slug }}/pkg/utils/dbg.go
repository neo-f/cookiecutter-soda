package utils

import (
	"encoding/json"
	"fmt"
)

func DBG(v interface{}, msg ...string) {
	bv, err := json.MarshalIndent(v, "", " ")
	if len(msg) == 0 {
		fmt.Println("===========DBG===========")
	} else {
		fmt.Printf("===========%s===========\n", msg[0])
	}
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(bv))
	}
	fmt.Println("=========================")
}

func DbgE(v interface{}, err error) {
	bv, _ := json.Marshal(v)
	fmt.Println(string(bv))
}
