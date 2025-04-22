package test

// func TestEqual(t *testing.T) {
// 	// ARRANGE
// 	testcases := []struct {
// 		scenario string
// 		act      func(*testing.T)
// 		assert   func(HelperTest)
// 	}{
// 		// these tests should pass
// 		{scenario: "Equal(1,1)",
// 			act: func(t *testing.T) { Equal(t, 1, 1) },
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "NotEqual(1,2)",
// 			act: func(t *testing.T) { NotEqual(t, 1, 2) },
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},

// 		// these tests should fail
// 		{scenario: "Equal(1,2)",
// 			act: func(t *testing.T) { Equal(t, 1, 2) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("Equal(1,2)/equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: 2",
// 					"got   : 1",
// 				})
// 			}},
// 		{scenario: "Equal(0,255,FormatHex)",
// 			act: func(t *testing.T) { Equal(t, 0, 255, FormatHex) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("Equal(0,255,FormatHex)/equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: ff",
// 					"got   : 00",
// 				})
// 			}},
// 		{scenario: "NotEqual(1,1)",
// 			act: func(t *testing.T) { NotEqual(t, 1, 1) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("NotEqual(1,1)/not_equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not: 1",
// 					"got   : 1",
// 				})
// 			}},

// 		// shallow vs deep equality
// 		{scenario: "Equal(*a,*a)",
// 			act: func(t *testing.T) {
// 				a := &struct{ int }{1}
// 				Equal(t, a, a)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "Equal(&{1},&{1})",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{1}
// 				Equal(t, g, w)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("Equal(&{1},&{1})/equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: &{1}",
// 					"got   : &{1}",
// 				})
// 			}},
// 		{scenario: "Equal(&{1},&{1},test.ShallowEquality)",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{1}
// 				Equal(t, g, w, ShallowEquality)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("Equal(&{1},&{1},test.ShallowEquality)/equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: &{1}",
// 					"got   : &{1}",
// 					"method: test.ShallowEquality",
// 				})
// 			}},
// 		{scenario: "Equal(&{1},&{1},test.DeepEquality)",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{1}
// 				Equal(t, g, w, DeepEquality)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "Equal(&{1},&{2},test.DeepEquality)",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{2}
// 				Equal(t, g, w, DeepEquality)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("Equal(&{1},&{2},test.DeepEquality)/equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: &{2}",
// 					"got   : &{1}",
// 					"method: test.DeepEquality",
// 				})
// 			}},
// 		{scenario: "NotEqual(&{1},&{1},test.DeepEquality)",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{1}
// 				NotEqual(t, g, w, DeepEquality)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("NotEqual(&{1},&{1},test.DeepEquality)/not_equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not: &{1}",
// 					"got   : &{1}",
// 					"method: test.DeepEquality",
// 				})
// 			}},

// 		// custom comparison
// 		{scenario: "Equal(1,2,custom{w=2*g})",
// 			act: func(t *testing.T) { Equal(t, 1, 2, func(g, w int) bool { return w == 2*g }) },
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "Equal(1,1,custom{w=2*g})",
// 			act: func(t *testing.T) { Equal(t, 1, 1, func(g, w int) bool { return w == 2*g }) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("Equal(1,1,custom{w=2*g})/equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: 1",
// 					"got   : 1",
// 					"method: func(got, wanted) bool {...}",
// 				})
// 			}},
// 		{scenario: "NotEqual(1,1,custom{w=2*g})",
// 			act: func(t *testing.T) { NotEqual(t, 1, 1, func(g, w int) bool { return w == 2*g }) },
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "NotEqual(1,2,custom{w=2*g})",
// 			act: func(t *testing.T) { NotEqual(t, 1, 2, func(g, w int) bool { return w == 2*g }) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("NotEqual(1,2,custom{w=2*g})/not_equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not: 2",
// 					"got   : 1",
// 					"method: func(got, wanted) bool {...}",
// 				})
// 			}},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			tc.assert(TestHelper(t, tc.act))
// 		})
// 	}
// }

// func TestDeepEqual(t *testing.T) {
// 	// ARRANGE
// 	testcases := []struct {
// 		scenario string
// 		act      func(*testing.T)
// 		assert   func(HelperTest)
// 	}{
// 		// these tests should pass
// 		{scenario: "DeepEqual(1,1)",
// 			act: func(t *testing.T) { DeepEqual(t, 1, 1) },
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "NotDeepEqual(1,2)",
// 			act: func(t *testing.T) { NotDeepEqual(t, 1, 2) },
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},

// 		// these tests should fail
// 		{scenario: "DeepEqual(1,2)",
// 			act: func(t *testing.T) { DeepEqual(t, 1, 2) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("DeepEqual(1,2)/deep_equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: 2",
// 					"got   : 1",
// 				})
// 			}},
// 		{scenario: "DeepEqual(0,255, FormatHex)",
// 			act: func(t *testing.T) { DeepEqual(t, 0, 255, FormatHex) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: ff",
// 					"got   : 00",
// 				})
// 			}},
// 		{scenario: "NotDeepEqual(1,1)",
// 			act: func(t *testing.T) { NotDeepEqual(t, 1, 1) },
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("NotDeepEqual(1,1)/not_deep_equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not: 1",
// 					"got   : 1",
// 				})
// 			}},

// 		{scenario: "DeepEqual(*a,*a)",
// 			act: func(t *testing.T) {
// 				a := &struct{ int }{1}
// 				DeepEqual(t, a, a)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "DeepEqual(&{1},&{1})",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{1}
// 				DeepEqual(t, g, w)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},

// 		// these tests should fail
// 		{scenario: "DeepEqual(&{1},&{1},test.ShallowEquality)",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{1}
// 				DeepEqual(t, g, w, ShallowEquality)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains([]string{
// 					"invalid test:",
// 					"option #1 is an unsupported type: test.Equality",
// 				})
// 			}},
// 		{scenario: "DeepEqual(&{1},&{2})",
// 			act: func(t *testing.T) {
// 				g := &struct{ int }{1}
// 				w := &struct{ int }{2}
// 				DeepEqual(t, g, w)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("DeepEqual(&{1},&{2})/deep_equal")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: &{2}",
// 					"got   : &{1}",
// 				})
// 			}},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			tc.assert(TestHelper(t, tc.act))
// 		})
// 	}
// }

// func TestIsNil(t *testing.T) {
// 	// ARRANGE
// 	testcases := []struct {
// 		scenario string
// 		act      func(*testing.T)
// 		assert   func(HelperTest)
// 	}{
// 		{scenario: "IsNil(zero-value any)",
// 			act: func(t *testing.T) {
// 				var a any
// 				IsNil(t, a)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil any)",
// 			act: func(t *testing.T) {
// 				var a any = nil
// 				IsNil(t, a)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil error)",
// 			act: func(t *testing.T) {
// 				IsNil(t, error(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil any)",
// 			act: func(t *testing.T) {
// 				IsNil(t, nil)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil slice)",
// 			act: func(t *testing.T) {
// 				IsNil(t, []byte(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil map)",
// 			act: func(t *testing.T) {
// 				IsNil(t, map[int]int(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil chan)",
// 			act: func(t *testing.T) {
// 				IsNil(t, chan struct{}(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil ptr)",
// 			act: func(t *testing.T) {
// 				IsNil(t, (*int)(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},
// 		{scenario: "IsNil(nil func)",
// 			act: func(t *testing.T) {
// 				IsNil(t, (func())(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			},
// 		},

// 		{scenario: "IsNil(non-nil error)",
// 			act: func(t *testing.T) {
// 				IsNil(t, errors.New("err"))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNil(non-nil_error)/is_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"unexpected error: *errors.errorString: err",
// 				})
// 			},
// 		},
// 		{scenario: "IsNil(non-nil slice)",
// 			act: func(t *testing.T) {
// 				IsNil(t, []byte{})
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNil(non-nil_slice)/is_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: nil ([]uint8)",
// 					"got   : not nil",
// 				})
// 			}},
// 		{scenario: "IsNil(non-nil map)",
// 			act: func(t *testing.T) {
// 				IsNil(t, map[int]int{})
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNil(non-nil_map)/is_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: nil (map[int]int)",
// 					"got   : not nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNil(non-nil chan)",
// 			act: func(t *testing.T) {
// 				ch := make(chan struct{})
// 				IsNil(t, ch)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNil(non-nil_chan)/is_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: nil (chan struct {})",
// 					"got   : not nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNil(non-nil ptr)",
// 			act: func(t *testing.T) {
// 				v := 42
// 				IsNil(t, &v)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNil(non-nil_ptr)/is_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: nil (*int)",
// 					"got   : not nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNil(non-nil func)",
// 			act: func(t *testing.T) {
// 				fn := func() {}
// 				IsNil(t, fn)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNil(non-nil_func)/is_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: nil (func())",
// 					"got   : not nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNil(non-nilable)",
// 			act: func(t *testing.T) {
// 				IsNil(t, 1)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNil(non-nilable)/is_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"test.IsNil: invalid test: values of type 'int' are not nilable",
// 				})
// 			},
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			tc.assert(TestHelper(t, tc.act))
// 		})
// 	}
// }

// func TestIsNotNil(t *testing.T) {
// 	// ARRANGE
// 	testcases := []struct {
// 		scenario string
// 		act      func(*testing.T)
// 		assert   func(HelperTest)
// 	}{
// 		{scenario: "IsNotNil(error)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, errors.New("err"))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "IsNotNil(slice)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, []byte{})
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "IsNotNil(map)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, map[int]int{})
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "IsNotNil(chan)",
// 			act: func(t *testing.T) {
// 				ch := make(chan struct{})
// 				IsNotNil(t, ch)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "IsNotNil(&int)",
// 			act: func(t *testing.T) {
// 				v := 42
// 				IsNotNil(t, &v)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "IsNotNil(func)",
// 			act: func(t *testing.T) {
// 				fn := func() {}
// 				IsNotNil(t, fn)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidPass()
// 				test.Report.IsEmpty()
// 			}},
// 		{scenario: "IsNotNil(nil error)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, error(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(nil_error)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not nil",
// 					"got   : nil",
// 				})
// 			}},
// 		{scenario: "IsNotNil(nil)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, nil)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(nil)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not nil",
// 					"got   : nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNotNil(nil slice)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, []byte(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(nil_slice)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not nil ([]uint8)",
// 					"got   : nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNotNil(nil map)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, map[int]int(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(nil_map)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not nil (map[int]int)",
// 					"got   : nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNotNil(nil chan)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, chan struct{}(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(nil_chan)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not nil (chan struct {})",
// 					"got   : nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNotNil(nil *int)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, (*int)(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(nil_*int)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not nil (*int)",
// 					"got   : nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNotNil(nil func)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, (func())(nil))
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(nil_func)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"wanted: not nil (func())",
// 					"got   : nil",
// 				})
// 			},
// 		},
// 		{scenario: "IsNotNil(int)",
// 			act: func(t *testing.T) {
// 				IsNotNil(t, 1)
// 			},
// 			assert: func(test HelperTest) {
// 				test.DidFail()
// 				test.Report.Contains("IsNotNil(int)/is_not_nil")
// 				test.Report.Contains([]string{
// 					currentFilename(),
// 					"invalid test: values of type 'int' are not nilable",
// 				})
// 			},
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			tc.assert(TestHelper(t, tc.act))
// 		})
// 	}
// }
