// Code generated by counterfeiter. DO NOT EDIT.
package powerattackforecastfakes

import (
	"sync"

	"github.com/chadius/terosbattleserver/entity/damagedistribution"
	"github.com/chadius/terosbattleserver/usecase/powerattackforecast"
)

type FakeVersusContextStrategy struct {
	CalculateStub        func(powerattackforecast.AttackerContext, powerattackforecast.DefenderContext)
	calculateMutex       sync.RWMutex
	calculateArgsForCall []struct {
		arg1 powerattackforecast.AttackerContext
		arg2 powerattackforecast.DefenderContext
	}
	CanCriticalStub        func() bool
	canCriticalMutex       sync.RWMutex
	canCriticalArgsForCall []struct {
	}
	canCriticalReturns struct {
		result1 bool
	}
	canCriticalReturnsOnCall map[int]struct {
		result1 bool
	}
	CriticalHitDamageStub        func() *damagedistribution.DamageDistribution
	criticalHitDamageMutex       sync.RWMutex
	criticalHitDamageArgsForCall []struct {
	}
	criticalHitDamageReturns struct {
		result1 *damagedistribution.DamageDistribution
	}
	criticalHitDamageReturnsOnCall map[int]struct {
		result1 *damagedistribution.DamageDistribution
	}
	CriticalHitThresholdStub        func() int
	criticalHitThresholdMutex       sync.RWMutex
	criticalHitThresholdArgsForCall []struct {
	}
	criticalHitThresholdReturns struct {
		result1 int
	}
	criticalHitThresholdReturnsOnCall map[int]struct {
		result1 int
	}
	NormalDamageStub        func() *damagedistribution.DamageDistribution
	normalDamageMutex       sync.RWMutex
	normalDamageArgsForCall []struct {
	}
	normalDamageReturns struct {
		result1 *damagedistribution.DamageDistribution
	}
	normalDamageReturnsOnCall map[int]struct {
		result1 *damagedistribution.DamageDistribution
	}
	ToHitStub        func() *damagedistribution.ToHitComparison
	toHitMutex       sync.RWMutex
	toHitArgsForCall []struct {
	}
	toHitReturns struct {
		result1 *damagedistribution.ToHitComparison
	}
	toHitReturnsOnCall map[int]struct {
		result1 *damagedistribution.ToHitComparison
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVersusContextStrategy) Calculate(arg1 powerattackforecast.AttackerContext, arg2 powerattackforecast.DefenderContext) {
	fake.calculateMutex.Lock()
	fake.calculateArgsForCall = append(fake.calculateArgsForCall, struct {
		arg1 powerattackforecast.AttackerContext
		arg2 powerattackforecast.DefenderContext
	}{arg1, arg2})
	stub := fake.CalculateStub
	fake.recordInvocation("Calculate", []interface{}{arg1, arg2})
	fake.calculateMutex.Unlock()
	if stub != nil {
		fake.CalculateStub(arg1, arg2)
	}
}

func (fake *FakeVersusContextStrategy) CalculateCallCount() int {
	fake.calculateMutex.RLock()
	defer fake.calculateMutex.RUnlock()
	return len(fake.calculateArgsForCall)
}

func (fake *FakeVersusContextStrategy) CalculateCalls(stub func(powerattackforecast.AttackerContext, powerattackforecast.DefenderContext)) {
	fake.calculateMutex.Lock()
	defer fake.calculateMutex.Unlock()
	fake.CalculateStub = stub
}

func (fake *FakeVersusContextStrategy) CalculateArgsForCall(i int) (powerattackforecast.AttackerContext, powerattackforecast.DefenderContext) {
	fake.calculateMutex.RLock()
	defer fake.calculateMutex.RUnlock()
	argsForCall := fake.calculateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeVersusContextStrategy) CanCritical() bool {
	fake.canCriticalMutex.Lock()
	ret, specificReturn := fake.canCriticalReturnsOnCall[len(fake.canCriticalArgsForCall)]
	fake.canCriticalArgsForCall = append(fake.canCriticalArgsForCall, struct {
	}{})
	stub := fake.CanCriticalStub
	fakeReturns := fake.canCriticalReturns
	fake.recordInvocation("CanCritical", []interface{}{})
	fake.canCriticalMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVersusContextStrategy) CanCriticalCallCount() int {
	fake.canCriticalMutex.RLock()
	defer fake.canCriticalMutex.RUnlock()
	return len(fake.canCriticalArgsForCall)
}

func (fake *FakeVersusContextStrategy) CanCriticalCalls(stub func() bool) {
	fake.canCriticalMutex.Lock()
	defer fake.canCriticalMutex.Unlock()
	fake.CanCriticalStub = stub
}

func (fake *FakeVersusContextStrategy) CanCriticalReturns(result1 bool) {
	fake.canCriticalMutex.Lock()
	defer fake.canCriticalMutex.Unlock()
	fake.CanCriticalStub = nil
	fake.canCriticalReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeVersusContextStrategy) CanCriticalReturnsOnCall(i int, result1 bool) {
	fake.canCriticalMutex.Lock()
	defer fake.canCriticalMutex.Unlock()
	fake.CanCriticalStub = nil
	if fake.canCriticalReturnsOnCall == nil {
		fake.canCriticalReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.canCriticalReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeVersusContextStrategy) CriticalHitDamage() *damagedistribution.DamageDistribution {
	fake.criticalHitDamageMutex.Lock()
	ret, specificReturn := fake.criticalHitDamageReturnsOnCall[len(fake.criticalHitDamageArgsForCall)]
	fake.criticalHitDamageArgsForCall = append(fake.criticalHitDamageArgsForCall, struct {
	}{})
	stub := fake.CriticalHitDamageStub
	fakeReturns := fake.criticalHitDamageReturns
	fake.recordInvocation("CriticalHitDamage", []interface{}{})
	fake.criticalHitDamageMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVersusContextStrategy) CriticalHitDamageCallCount() int {
	fake.criticalHitDamageMutex.RLock()
	defer fake.criticalHitDamageMutex.RUnlock()
	return len(fake.criticalHitDamageArgsForCall)
}

func (fake *FakeVersusContextStrategy) CriticalHitDamageCalls(stub func() *damagedistribution.DamageDistribution) {
	fake.criticalHitDamageMutex.Lock()
	defer fake.criticalHitDamageMutex.Unlock()
	fake.CriticalHitDamageStub = stub
}

func (fake *FakeVersusContextStrategy) CriticalHitDamageReturns(result1 *damagedistribution.DamageDistribution) {
	fake.criticalHitDamageMutex.Lock()
	defer fake.criticalHitDamageMutex.Unlock()
	fake.CriticalHitDamageStub = nil
	fake.criticalHitDamageReturns = struct {
		result1 *damagedistribution.DamageDistribution
	}{result1}
}

func (fake *FakeVersusContextStrategy) CriticalHitDamageReturnsOnCall(i int, result1 *damagedistribution.DamageDistribution) {
	fake.criticalHitDamageMutex.Lock()
	defer fake.criticalHitDamageMutex.Unlock()
	fake.CriticalHitDamageStub = nil
	if fake.criticalHitDamageReturnsOnCall == nil {
		fake.criticalHitDamageReturnsOnCall = make(map[int]struct {
			result1 *damagedistribution.DamageDistribution
		})
	}
	fake.criticalHitDamageReturnsOnCall[i] = struct {
		result1 *damagedistribution.DamageDistribution
	}{result1}
}

func (fake *FakeVersusContextStrategy) CriticalHitThreshold() int {
	fake.criticalHitThresholdMutex.Lock()
	ret, specificReturn := fake.criticalHitThresholdReturnsOnCall[len(fake.criticalHitThresholdArgsForCall)]
	fake.criticalHitThresholdArgsForCall = append(fake.criticalHitThresholdArgsForCall, struct {
	}{})
	stub := fake.CriticalHitThresholdStub
	fakeReturns := fake.criticalHitThresholdReturns
	fake.recordInvocation("CriticalHitThreshold", []interface{}{})
	fake.criticalHitThresholdMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVersusContextStrategy) CriticalHitThresholdCallCount() int {
	fake.criticalHitThresholdMutex.RLock()
	defer fake.criticalHitThresholdMutex.RUnlock()
	return len(fake.criticalHitThresholdArgsForCall)
}

func (fake *FakeVersusContextStrategy) CriticalHitThresholdCalls(stub func() int) {
	fake.criticalHitThresholdMutex.Lock()
	defer fake.criticalHitThresholdMutex.Unlock()
	fake.CriticalHitThresholdStub = stub
}

func (fake *FakeVersusContextStrategy) CriticalHitThresholdReturns(result1 int) {
	fake.criticalHitThresholdMutex.Lock()
	defer fake.criticalHitThresholdMutex.Unlock()
	fake.CriticalHitThresholdStub = nil
	fake.criticalHitThresholdReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakeVersusContextStrategy) CriticalHitThresholdReturnsOnCall(i int, result1 int) {
	fake.criticalHitThresholdMutex.Lock()
	defer fake.criticalHitThresholdMutex.Unlock()
	fake.CriticalHitThresholdStub = nil
	if fake.criticalHitThresholdReturnsOnCall == nil {
		fake.criticalHitThresholdReturnsOnCall = make(map[int]struct {
			result1 int
		})
	}
	fake.criticalHitThresholdReturnsOnCall[i] = struct {
		result1 int
	}{result1}
}

func (fake *FakeVersusContextStrategy) NormalDamage() *damagedistribution.DamageDistribution {
	fake.normalDamageMutex.Lock()
	ret, specificReturn := fake.normalDamageReturnsOnCall[len(fake.normalDamageArgsForCall)]
	fake.normalDamageArgsForCall = append(fake.normalDamageArgsForCall, struct {
	}{})
	stub := fake.NormalDamageStub
	fakeReturns := fake.normalDamageReturns
	fake.recordInvocation("NormalDamage", []interface{}{})
	fake.normalDamageMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVersusContextStrategy) NormalDamageCallCount() int {
	fake.normalDamageMutex.RLock()
	defer fake.normalDamageMutex.RUnlock()
	return len(fake.normalDamageArgsForCall)
}

func (fake *FakeVersusContextStrategy) NormalDamageCalls(stub func() *damagedistribution.DamageDistribution) {
	fake.normalDamageMutex.Lock()
	defer fake.normalDamageMutex.Unlock()
	fake.NormalDamageStub = stub
}

func (fake *FakeVersusContextStrategy) NormalDamageReturns(result1 *damagedistribution.DamageDistribution) {
	fake.normalDamageMutex.Lock()
	defer fake.normalDamageMutex.Unlock()
	fake.NormalDamageStub = nil
	fake.normalDamageReturns = struct {
		result1 *damagedistribution.DamageDistribution
	}{result1}
}

func (fake *FakeVersusContextStrategy) NormalDamageReturnsOnCall(i int, result1 *damagedistribution.DamageDistribution) {
	fake.normalDamageMutex.Lock()
	defer fake.normalDamageMutex.Unlock()
	fake.NormalDamageStub = nil
	if fake.normalDamageReturnsOnCall == nil {
		fake.normalDamageReturnsOnCall = make(map[int]struct {
			result1 *damagedistribution.DamageDistribution
		})
	}
	fake.normalDamageReturnsOnCall[i] = struct {
		result1 *damagedistribution.DamageDistribution
	}{result1}
}

func (fake *FakeVersusContextStrategy) ToHit() *damagedistribution.ToHitComparison {
	fake.toHitMutex.Lock()
	ret, specificReturn := fake.toHitReturnsOnCall[len(fake.toHitArgsForCall)]
	fake.toHitArgsForCall = append(fake.toHitArgsForCall, struct {
	}{})
	stub := fake.ToHitStub
	fakeReturns := fake.toHitReturns
	fake.recordInvocation("ToHit", []interface{}{})
	fake.toHitMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVersusContextStrategy) ToHitCallCount() int {
	fake.toHitMutex.RLock()
	defer fake.toHitMutex.RUnlock()
	return len(fake.toHitArgsForCall)
}

func (fake *FakeVersusContextStrategy) ToHitCalls(stub func() *damagedistribution.ToHitComparison) {
	fake.toHitMutex.Lock()
	defer fake.toHitMutex.Unlock()
	fake.ToHitStub = stub
}

func (fake *FakeVersusContextStrategy) ToHitReturns(result1 *damagedistribution.ToHitComparison) {
	fake.toHitMutex.Lock()
	defer fake.toHitMutex.Unlock()
	fake.ToHitStub = nil
	fake.toHitReturns = struct {
		result1 *damagedistribution.ToHitComparison
	}{result1}
}

func (fake *FakeVersusContextStrategy) ToHitReturnsOnCall(i int, result1 *damagedistribution.ToHitComparison) {
	fake.toHitMutex.Lock()
	defer fake.toHitMutex.Unlock()
	fake.ToHitStub = nil
	if fake.toHitReturnsOnCall == nil {
		fake.toHitReturnsOnCall = make(map[int]struct {
			result1 *damagedistribution.ToHitComparison
		})
	}
	fake.toHitReturnsOnCall[i] = struct {
		result1 *damagedistribution.ToHitComparison
	}{result1}
}

func (fake *FakeVersusContextStrategy) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.calculateMutex.RLock()
	defer fake.calculateMutex.RUnlock()
	fake.canCriticalMutex.RLock()
	defer fake.canCriticalMutex.RUnlock()
	fake.criticalHitDamageMutex.RLock()
	defer fake.criticalHitDamageMutex.RUnlock()
	fake.criticalHitThresholdMutex.RLock()
	defer fake.criticalHitThresholdMutex.RUnlock()
	fake.normalDamageMutex.RLock()
	defer fake.normalDamageMutex.RUnlock()
	fake.toHitMutex.RLock()
	defer fake.toHitMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVersusContextStrategy) recordInvocation(key string, args []interface{}) {
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

var _ powerattackforecast.VersusContextStrategy = new(FakeVersusContextStrategy)
