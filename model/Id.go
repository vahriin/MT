package model

import (
	"fmt"
	"strconv"
)

type Id int

func (id Id) Scan(src interface{}) error {
	_, err := fmt.Sscanf(strconv.FormatInt(src.(int64), 10), "%d", &id)
	return err
}
