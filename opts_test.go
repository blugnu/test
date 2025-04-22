package test

// func TestGetOpt(t *testing.T) {
// 	// ARRANGE
// 	testcases := []struct {
// 		scenario string
// 		test     func()
// 	}{
// 		{
// 			scenario: "no value of desired type",
// 			test: func() {
// 				n := "default"
// 				result := getOpt(&n, 1, 2)
// 				Equal(t, result, false, "returns")
// 				Equal(t, "default", n)
// 			},
// 		},
// 		{
// 			scenario: "single value of desired type",
// 			test: func() {
// 				var n string
// 				result := getOpt(&n, "value")
// 				Equal(t, result, true, "returns")
// 				Equal(t, "value", n)
// 			},
// 		},
// 		{
// 			scenario: "first of multiple values of desired type",
// 			test: func() {
// 				var n string
// 				result := getOpt(&n, "first", "extra")
// 				Equal(t, result, true, "returns")
// 				Equal(t, "first", n)
// 			},
// 		},
// 	}
// 	for _, tc := range testcases {
// 		t.Run(tc.scenario, func(t *testing.T) {
// 			tc.test()
// 		})
// 	}
// }
