package testframes

import (
	"sync"
	"testing"
)

type Hook func(t *testing.T)

type onceHook struct {
	hook Hook
	once *sync.Once
}

func (hook onceHook) Do(t *testing.T) {
	hook.once.Do(func() { hook.hook(t) })
}

type Frame struct {
	t               *testing.T
	beforeHooks     []onceHook
	beforeEachHooks []Hook
	afterEachHooks  []Hook
}

func (f *Frame) Before(hook Hook) {
	f.beforeHooks = append(f.beforeHooks, onceHook{
		hook: hook,
		once: new(sync.Once),
	})
}

func (f *Frame) BeforeEach(hook Hook) {
	f.beforeEachHooks = append(f.beforeEachHooks, hook)
}

func (f *Frame) AfterEach(hook Hook) {
	f.afterEachHooks = append(f.afterEachHooks, hook)
}

func (f *Frame) Run(caseName string, testFunc func(t *testing.T)) {
	f.t.Run(caseName, func(t *testing.T) {
		for _, hook := range f.beforeHooks {
			hook.Do(t)
		}

		for _, hook := range f.beforeEachHooks {
			hook(t)
		}

		testFunc(t)

		for _, hook := range f.afterEachHooks {
			hook(t)
		}
	})
}

func (f *Frame) NextFrame(name string, frameFunc func(*Frame)) {
	f.t.Run(name, func(t *testing.T) {
		frameFunc(&Frame{
			t:               t,
			beforeHooks:     append([]onceHook{}, f.beforeHooks...),
			beforeEachHooks: append([]Hook{}, f.beforeEachHooks...),
			afterEachHooks:  append([]Hook{}, f.afterEachHooks...),
		})
	})
}

func New(t *testing.T) *Frame {
	return &Frame{t: t}
}
