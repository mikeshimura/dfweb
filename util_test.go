package dfweb

import (
  "testing"
  "fmt"
)

func TestUtil(t *testing.T) {
	res:=ConvFromWebDataInd("300","int64")
	fmt.Printf("res %v %T \n",res,res)
}