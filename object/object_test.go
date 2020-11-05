package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is tatsuya"}
	diff2 := &String{Value: "My name is tatsuya"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content has different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content has different hash keys")
	}
	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content has same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	n1 := &Integer{Value: int64(10)}
	n2 := &Integer{Value: int64(10)}
	diff1 := &Integer{Value: int64(20)}
	diff2 := &Integer{Value: int64(20)}

	if n1.HashKey() != n2.HashKey() {
		t.Errorf("integer with same content has different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("integer with same content has different hash keys")
	}
	if n1.HashKey() == diff1.HashKey() {
		t.Errorf("integer with different content has same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	t1 := &Boolean{Value: true}
	t2 := &Boolean{Value: true}
	f1 := &Boolean{Value: false}
	f2 := &Boolean{Value: false}

	if t1.HashKey() != t2.HashKey() {
		t.Errorf("boolean with same content has different hash keys")
	}
	if f1.HashKey() != f2.HashKey() {
		t.Errorf("boolean with same content has different hash keys")
	}
	if t1.HashKey() == f1.HashKey() {
		t.Errorf("boolean with different content has same hash keys")
	}
}
