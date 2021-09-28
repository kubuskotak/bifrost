package bifrost

import "testing"

func toDelimited(tb testing.TB) {
	cases := [][]string{
		{"testCase", "test@case"},
		{"TestCase", "test@case"},
		{"Test Case", "test@case"},
		{" Test Case", "test@case"},
		{"Test Case ", "test@case"},
		{" Test Case ", "test@case"},
		{"test", "test"},
		{"test_case", "test@case"},
		{"Test", "test"},
		{"", ""},
		{"ManyManyWords", "many@many@words"},
		{"manyManyWords", "many@many@words"},
		{"AnyKind of_string", "any@kind@of@string"},
		{"numbers2and55with000", "numbers@2@and@55@with@000"},
		{"JSONData", "json@data"},
		{"userID", "user@id"},
		{"AAAbbb", "aa@abbb"},
		{"test-case", "test@case"},
	}
	for _, i := range cases {
		in := i[0]
		out := i[1]
		result := ToDelimited(in, '@')
		if result != out {
			tb.Errorf("%q (%q != %q)", in, result, out)
		}
	}
}

func TestToDelimited(t *testing.T) { toDelimited(t) }

func BenchmarkToDelimited(b *testing.B) {
	benchmarkSnakeTest(b, toDelimited)
}

func benchmarkSnakeTest(b *testing.B, fn func(testing.TB)) {
	for n := 0; n < b.N; n++ {
		fn(b)
	}
}
