package testframe

import "testing"

// Cleanup is an interface that allows for cleanup functions to be registered
type Cleanup interface {
	Cleanup(fn func())
}

// Nil is a no-op implementation of the TestingT interface.  It is used to
// provide a sentinel implementation of a TestingT that should be interpreted
// as "no test frame".
//
// It implements the Cleanup interface, delegating to another implementation;
// the delegate is the TestingT in the testframe at the time that the Nil frame
// is created.
//
// The other methods are no-ops, implemented to satisfy the TestingT interface.
type Nil struct {
	T Cleanup
}

func (n Nil) Cleanup(fn func()) { n.T.Cleanup(fn) }

func (Nil) Name() string { return "NilFrame" }

func (Nil) Run(string, func(t *testing.T)) bool { return true }

func (Nil) Error(...any) { /* no-op */ }

func (Nil) Errorf(string, ...any) { /* no-op */ }

func (Nil) Fail() { /* no-op */ }

func (Nil) FailNow() { /* no-op */ }

func (Nil) Failed() bool { return false }

func (Nil) Fatal(...any) { /* no-op */ }

func (Nil) Fatalf(string, ...any) { /* no-op */ }

func (Nil) Helper() { /* no-op */ }

func (Nil) Parallel() { /* no-op */ }

func (Nil) Setenv(string, string) { /* no-op */ }

func (Nil) SkipNow() { /* no-op */ }
