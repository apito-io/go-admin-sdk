package goapitosdk

import "testing"

func TestDeriveSystemAPIBaseURL(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"http://localhost:5050/system/graphql", "http://localhost:5050"},
		{"http://localhost:5050/secured/graphql/", "http://localhost:5050"},
		{"http://localhost:5050/system/graphql/", "http://localhost:5050"},
	}
	for _, tc := range cases {
		if got := deriveSystemAPIBaseURL(tc.in); got != tc.want {
			t.Fatalf("deriveSystemAPIBaseURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
