package goapitosdk

import "testing"

func TestDeriveRestBaseURL(t *testing.T) {
	t.Parallel()
	cases := []struct {
		graphql string
		want    string
	}{
		{"http://localhost:5050/system/graphql", "http://localhost:5050/secured"},
		{"http://localhost:5050/secured/graphql", "http://localhost:5050/secured"},
		{"http://localhost:5050/graphql", "http://localhost:5050"},
	}
	for _, tc := range cases {
		if got := deriveRestBaseURL(tc.graphql); got != tc.want {
			t.Fatalf("deriveRestBaseURL(%q) = %q, want %q", tc.graphql, got, tc.want)
		}
	}
}
