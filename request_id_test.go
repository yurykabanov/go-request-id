package requestid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFromContext(t *testing.T) {
	tests := []struct {
		name string
		ctx context.Context
		value string
		ok bool
	}{
		{
			name: "nil context",
			ctx: nil,
			value: "",
			ok: false,
		},
		{
			name: "context without value",
			ctx: context.Background(),
			value: "",
			ok: false,
		},
		{
			name: "context with value",
			ctx: context.WithValue(context.Background(), ContextKey, "value"),
			value: "value",
			ok: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, ok := GetFromContext(test.ctx)

			assert.Equal(t, test.value, value)
			assert.Equal(t, test.ok, ok)
		})
	}
}

func TestWithRequestId(t *testing.T) {
	ctx := WithRequestId(context.Background(), "value")
	assert.Equal(t, "value", ctx.Value(ContextKey))
}

type testHandler struct {
}

func (*testHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name string
		requestHeader string
		responseHeader string
	}{
		{
			name: "without header",
			requestHeader: "",
			responseHeader: "some_generated_value",
		},
		{
			name: "with header",
			requestHeader: "some_value",
			responseHeader: "some_value",
		},
	}

	ts := httptest.NewServer(Middleware(&testHandler{}, func() string {
		return "some_generated_value"
	}))
	client := http.Client{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL, nil)
			req.Header.Set(Header, test.requestHeader)
			assert.NoError(t, err)

			res, err := client.Do(req)
			assert.NoError(t, err)
			if res != nil {
				defer res.Body.Close()
			}

			assert.Equal(t, test.responseHeader, res.Header.Get(Header))
		})
	}
}

func TestDefaultRequestIdProvider(t *testing.T) {
	v := DefaultRequestIdProvider()
	assert.NotEmpty(t, v)
	assert.Equal(t, 32, len(v))
}
