package jailbotque

import "testing"

func wrongStackSize(t *testing.T, got, wants uint) {
	t.Errorf("Alert stack had incorrect size. Got %d, wants %d", got, wants)
}

func TestAlertStackAdd(t *testing.T) {
	stack := NewAlertStack()

	if stack.Size() != 0 {
		wrongStackSize(t, stack.Size(), 0)
	}

	alert := &Alert{}
	stack.Add(alert)

	if stack.Size() != 1 {
		wrongStackSize(t, stack.Size(), 1)
	}

	for i := 0; i < 10; i++ {
		stack.Add(alert)
	}
	if stack.Size() != 11 {
		wrongStackSize(t, stack.Size(), 11)
	}
}

func TestAlertStackPop(t *testing.T) {
	stack := NewAlertStack()
	alert := &Alert{}

	for i := 0; i < 10; i++ {
		stack.Add(alert)
	}

	if stack.Size() != 10 {
		wrongStackSize(t, stack.Size(), 10)
	}

	stack.Pop()
	if stack.Size() != 9 {
		wrongStackSize(t, stack.Size(), 9)
	}

	for i := 0; i < 6; i++ {
		stack.Pop()
	}
	if stack.Size() != 3 {
		wrongStackSize(t, stack.Size(), 3)
	}

}
