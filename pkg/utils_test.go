package pkg

import (
	"fmt"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	_, err := os.Stat("s")
	fmt.Println(os.IsExist(err))
	fmt.Println(os.IsNotExist(err))
}
