package httpexpect

// Array provides methods to inspect attached []interface{} object
// (Go representation of JSON array).
type Array struct {
	checker Checker
	value   []interface{}
}

// NewArray returns a new Array given a checker used to report failures
// and value to be inspected.
//
// Both checker and value should not be nil. If value is nil, failure is reported.
//
// Example:
//  array := NewArray(NewAssertChecker(t), []interface{}{"foo", 123})
func NewArray(checker Checker, value []interface{}) *Array {
	if value == nil {
		checker.Fail("expected non-nil array value")
	} else {
		value, _ = canonArray(checker, value)
	}
	return &Array{checker, value}
}

// Raw returns underlying value attached to Array.
// This is the value originally passed to NewArray, converted to canonical form.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  assert.Equal(t, []interface{}{"foo", 123.0}, array.Raw())
func (a *Array) Raw() []interface{} {
	return a.value
}

// Length returns a new Number object that may be used to inspect array length.
//
// Example:
//  array := NewArray(checker, []interface{}{1, 2, 3})
//  array.Length().Equal(3)
func (a *Array) Length() *Number {
	return NewNumber(a.checker.Clone(), float64(len(a.value)))
}

// Element returns a new Value object that may be used to inspect array element
// for given index.
//
// If index is out of array bounds, Element reports failure and returns empty
// (but non-nil) value.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.Element(0).String().Equal("foo")
//  array.Element(1).Number().Equal(123)
func (a *Array) Element(index int) *Value {
	if len(a.value) <= index {
		a.checker.Fail("expected array with length > %d, got %v", index, a.value)
		return NewValue(a.checker.Clone(), nil)
	}
	return NewValue(a.checker.Clone(), a.value[index])
}

// Empty succeedes if array is empty.
//
// Example:
//  array := NewArray(checker, []interface{}{})
//  array.Empty()
func (a *Array) Empty() *Array {
	expected := []interface{}{}
	a.checker.Equal(expected, a.value)
	return a
}

// NotEmpty succeedes if array is non-empty.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.NotEmpty()
func (a *Array) NotEmpty() *Array {
	expected := []interface{}{}
	a.checker.NotEqual(expected, a.value)
	return a
}

// Equal succeedes if array is equal to another array.
// Before comparison, both arrays are converted to canonical form.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.Equal([]interface{}{"foo", 123})
func (a *Array) Equal(v []interface{}) *Array {
	expected, ok := canonArray(a.checker, v)
	if !ok {
		return a
	}
	a.checker.Equal(expected, a.value)
	return a
}

// NotEqual succeedes if array is not equal to another array.
// Before comparison, both arrays are converted to canonical form.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.NotEqual([]interface{}{123, "foo"})
func (a *Array) NotEqual(v []interface{}) *Array {
	expected, ok := canonArray(a.checker, v)
	if !ok {
		return a
	}
	a.checker.NotEqual(expected, a.value)
	return a
}

// Contains succeedes if array contains all given elements (in any order).
// Before comparison, array and all elements are converted to canonical form.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.Contains(123, "foo")
func (a *Array) Contains(v ...interface{}) *Array {
	elements, ok := canonArray(a.checker, v)
	if !ok {
		return a
	}
	for _, e := range elements {
		if !a.containsElement(e) {
			a.checker.Fail("expected array containing %v, got %v", e, a.value)
		}
	}
	return a
}

// NotContains succeedes if array contains none of given elements.
// Before comparison, array and all elements are converted to canonical form.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.NotContains("bar")         // success
//  array.NotContains("bar", "foo")  // failure (array contains "foo")
func (a *Array) NotContains(v ...interface{}) *Array {
	elements, ok := canonArray(a.checker, v)
	if !ok {
		return a
	}
	for _, e := range elements {
		if a.containsElement(e) {
			a.checker.Fail("expected array NOT containing %v, got %v", e, a.value)
		}
	}
	return a
}

// Elements succeedes if array contains all given elements, in given order, and only them.
// Before comparison, array and all elements are converted to canonical form.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.Elements("foo", 123)
//
// This calls are equivalent:
//  array.Elelems("a", "b")
//  array.Equal([]interface{}{"a", "b"})
func (a *Array) Elements(v ...interface{}) *Array {
	return a.Equal(v)
}

// ElementsAnyOrder succeedes if array contains all given elements, in any order, and only
// them. Before comparison, array and all elements are converted to canonical form.
//
// Example:
//  array := NewArray(checker, []interface{}{"foo", 123})
//  array.ElementsAnyOrder(123, "foo")
//
// This calls are equivalent:
//  array.ElementsAnyOrder("a", "b")
//  array.ElementsAnyOrder("b", "a")
func (a *Array) ElementsAnyOrder(v ...interface{}) *Array {
	elements, ok := canonArray(a.checker, v)
	if !ok {
		return a
	}
	if len(elements) != len(a.value) {
		a.checker.Fail("expected array len == %d, got %d", len(elements), len(a.value))
		return a
	}
	for _, e := range elements {
		if !a.containsElement(e) {
			a.checker.Fail("expected array containing %v, got %v", e, a.value)
		}
	}
	return a
}

func (a *Array) containsElement(expected interface{}) bool {
	for _, e := range a.value {
		if a.checker.Compare(expected, e) {
			return true
		}
	}
	return false
}
