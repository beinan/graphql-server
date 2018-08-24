package ploader

type FutureValue interface {
	Value() (interface{}, error)
}

type FutureValueImpl struct {
	value interface{}
	err   error
	done  chan struct{} //ignore the value in chan
}

func (fv *FutureValueImpl) Value() (interface{}, error) {
	<-fv.done //waiting for the result
	return fv.value, fv.err
}

func MakeFutureValue(producer func() (interface{}, error)) FutureValue {
	fv := &FutureValueImpl{
		done: make(chan struct{}),
	}
	go func() {
		fv.value, fv.err = producer()
		close(fv.done) //notify all the waiting/running Value() goroutine
	}()
	return fv
}
