// Code generated by counterfeiter. DO NOT EDIT.
package checkfakes

import (
	"sync"

	"github.com/at-silva/ddapi/check"
)

type FakeParamsChecker struct {
	CheckStub        func(map[string]interface{}, string) error
	checkMutex       sync.RWMutex
	checkArgsForCall []struct {
		arg1 map[string]interface{}
		arg2 string
	}
	checkReturns struct {
		result1 error
	}
	checkReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeParamsChecker) Check(arg1 map[string]interface{}, arg2 string) error {
	fake.checkMutex.Lock()
	ret, specificReturn := fake.checkReturnsOnCall[len(fake.checkArgsForCall)]
	fake.checkArgsForCall = append(fake.checkArgsForCall, struct {
		arg1 map[string]interface{}
		arg2 string
	}{arg1, arg2})
	stub := fake.CheckStub
	fakeReturns := fake.checkReturns
	fake.recordInvocation("Check", []interface{}{arg1, arg2})
	fake.checkMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeParamsChecker) CheckCallCount() int {
	fake.checkMutex.RLock()
	defer fake.checkMutex.RUnlock()
	return len(fake.checkArgsForCall)
}

func (fake *FakeParamsChecker) CheckCalls(stub func(map[string]interface{}, string) error) {
	fake.checkMutex.Lock()
	defer fake.checkMutex.Unlock()
	fake.CheckStub = stub
}

func (fake *FakeParamsChecker) CheckArgsForCall(i int) (map[string]interface{}, string) {
	fake.checkMutex.RLock()
	defer fake.checkMutex.RUnlock()
	argsForCall := fake.checkArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeParamsChecker) CheckReturns(result1 error) {
	fake.checkMutex.Lock()
	defer fake.checkMutex.Unlock()
	fake.CheckStub = nil
	fake.checkReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeParamsChecker) CheckReturnsOnCall(i int, result1 error) {
	fake.checkMutex.Lock()
	defer fake.checkMutex.Unlock()
	fake.CheckStub = nil
	if fake.checkReturnsOnCall == nil {
		fake.checkReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeParamsChecker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkMutex.RLock()
	defer fake.checkMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeParamsChecker) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ check.ParamsChecker = new(FakeParamsChecker)
