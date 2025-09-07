package graphql

import (
	"bytes"
	"encoding/json"
	"io"
	"jello-mark-backend/internal/adapter/outbound/memory"
	"jello-mark-backend/internal/application"
	"jello-mark-backend/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type fixedClock struct{ t time.Time }

func (f fixedClock) Now() time.Time { return f.t }

func TestGraphQL_HealthAndUserFlow(t *testing.T) {
	repo := memory.NewUserRepo()
	clk := fixedClock{t: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}
	idgen := func() domain.UserID { return domain.UserID("u-1") }
	svc := application.NewUserService(repo, clk, idgen)
	h := Handler{Users: svc}

	do := func(query string, vars map[string]interface{}) *httptest.ResponseRecorder {
		body, _ := json.Marshal(map[string]interface{}{"query": query, "variables": vars})
		r := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(body))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		return w
	}

	w := do("query{health}", nil)
	if w.Code != 200 {
		t.Fatalf("health code=%d", w.Code)
	}
	b, _ := io.ReadAll(w.Body)
	if !bytes.Contains(b, []byte("ok")) {
		t.Fatalf("health body=%s", string(b))
	}

	w = do("mutation($name:String!){createUser(name:$name){id name createdAt}}", map[string]interface{}{"name": "alice"})
	if w.Code != 200 {
		t.Fatalf("create code=%d", w.Code)
	}
	b, _ = io.ReadAll(w.Body)
	if !bytes.Contains(b, []byte("alice")) {
		t.Fatalf("create body=%s", string(b))
	}

	w = do("query{users{ id name createdAt }}", nil)
	if w.Code != 200 {
		t.Fatalf("users code=%d", w.Code)
	}
	b, _ = io.ReadAll(w.Body)
	if !bytes.Contains(b, []byte("alice")) {
		t.Fatalf("users body=%s", string(b))
	}
}
