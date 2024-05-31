package device

import (
	"fmt"
	"testing"
)

func TestGetBaseBoardID(t *testing.T) {
	_, err := GetBaseBoardID()
	fmt.Println(err)
}
