package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	testCases := []struct {
		name string
		args []string
		env  map[string]string
		want config
	}{
		{
			"default",
			[]string{
				"palindrome",
			},
			nil,
			config{
				defaultHTTPAddr,
				defaultStrictPalindrome,
				defaultMongoURI,
			},
		},
		{
			"args only",
			[]string{
				"palindrome",
				"-http-addr=:8081",
				"-strict-palindrome=false",
				"-mongo-uri=mongodb://localhost:27017",
			},
			nil,
			config{
				":8081",
				false,
				"mongodb://localhost:27017",
			},
		},
		{
			"envs only",
			[]string{
				"palindrome",
			},
			map[string]string{
				"HTTP_ADDR":         ":8081",
				"STRICT_PALINDROME": "false",
				"MONGO_URI":         "mongodb://localhost:27017",
			},
			config{
				":8081",
				false,
				"mongodb://localhost:27017",
			},
		},
		{
			"args and envs",
			[]string{
				"palindrome",
				"-http-addr=:8081",
				"-strict-palindrome=false",
				"-mongo-uri=mongodb://user1:pass1@localhost:27017",
			},
			map[string]string{
				"HTTP_ADDR":         ":8082",
				"STRICT_PALINDROME": "true",
				"MONGO_URI":         "mongodb://user2:pass2@localhost:27017",
			},
			config{
				":8081",
				false,
				"mongodb://user1:pass1@localhost:27017",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer setEnv(getEnv(tc.env))
			setEnv(tc.env)
			cfg := parseConfig(tc.args)
			require.Equal(t, tc.want, cfg)
		})
	}
}

func getEnv(env map[string]string) map[string]string {
	retEnv := make(map[string]string, len(env))
	for k := range env {
		v := os.Getenv(k)
		retEnv[k] = v
	}
	return retEnv
}

func setEnv(env map[string]string) {
	for k, v := range env {
		os.Setenv(k, v)
	}
}

func TestHealthz(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/healthz", nil)
	healthz(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "ok", w.Body.String())
}
