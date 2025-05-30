package std

import "testing"

type iFace interface {
	Get() bool
}

type iFaceImp struct {
	Body string
}

func (i iFaceImp) Get() bool {
	return false
}

func TestIsInterfaceNil(t *testing.T) {
	var (
		i1  interface{}
		ifc iFace
	)

	// check initial value
	t.Log(IsInterfaceNil(i1))
	t.Log(IsInterfaceNil(ifc))
	ifc = iFaceImp{}
	t.Log(IsInterfaceNil(ifc))
	ifc = nil
	t.Log(IsInterfaceNil(ifc))
}

func TestAnyVal(t *testing.T) {
	val := AnyVal[int]("1")
	if val == 0 {
		t.Fail()
	}
	t.Log(val)
}
