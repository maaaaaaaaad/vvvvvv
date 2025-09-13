package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"jello-mark-backend/internal/application"
	"jello-mark-backend/internal/domain"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type fixedClock struct{ t time.Time }

type testUserRepo struct {
	mu    sync.RWMutex
	items []domain.User
}

func newTestUserRepo() *testUserRepo {
	return &testUserRepo{
		items: make(
			[]domain.User,
			0,
			16,
		),
	}
}

func (r *testUserRepo) Create(
	ctx context.Context,
	u domain.User,
) error {
	if u.ID == "" {
		return errors.New("id required")
	}
	r.mu.Lock()
	r.items = append(
		r.items,
		u,
	)
	r.mu.Unlock()
	return nil
}

func (r *testUserRepo) List(
	ctx context.Context,
	limit int,
) (
	[]domain.User,
	error,
) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if limit > len(r.items) {
		limit = len(r.items)
	}
	res := make(
		[]domain.User,
		limit,
	)
	copy(
		res,
		r.items[:limit],
	)
	return res, nil
}

func (f fixedClock) Now() time.Time { return f.t }

func TestGraphQL_HealthAndUserFlow(t *testing.T) {
	repo := newTestUserRepo()
	clk := fixedClock{
		t: time.Date(
			2024,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		),
	}
	idgen := func() domain.UserID { return domain.UserID("u-1") }
	svc := application.NewUserService(
		repo,
		clk,
		idgen,
	)
	h := NewGraphQLHandler(
		svc,
		false,
	)

	do := func(
		query string,
		vars map[string]interface{},
	) *httptest.ResponseRecorder {
		body, _ := json.Marshal(
			map[string]interface{}{
				"query":     query,
				"variables": vars,
			},
		)
		r := httptest.NewRequest(
			http.MethodPost,
			"/graphql",
			bytes.NewReader(body),
		)
		w := httptest.NewRecorder()
		h.ServeHTTP(
			w,
			r,
		)
		return w
	}

	w := do(
		"query{health}",
		nil,
	)
	if w.Code != 200 {
		t.Fatalf(
			"health code=%d",
			w.Code,
		)
	}
	b, _ := io.ReadAll(w.Body)
	if !bytes.Contains(
		b,
		[]byte("ok"),
	) {
		t.Fatalf(
			"health body=%s",
			string(b),
		)
	}

	w = do(
		"mutation($name:String!){createUser(name:$name){id name createdAt}}",
		map[string]interface{}{"name": "alice"},
	)
	if w.Code != 200 {
		t.Fatalf(
			"create code=%d",
			w.Code,
		)
	}
	b, _ = io.ReadAll(w.Body)
	if !bytes.Contains(
		b,
		[]byte("alice"),
	) {
		t.Fatalf(
			"create body=%s",
			string(b),
		)
	}

	w = do(
		"query{users{ id name createdAt }}",
		nil,
	)
	if w.Code != 200 {
		t.Fatalf(
			"users code=%d",
			w.Code,
		)
	}
	b, _ = io.ReadAll(w.Body)
	if !bytes.Contains(
		b,
		[]byte("alice"),
	) {
		t.Fatalf(
			"users body=%s",
			string(b),
		)
	}
}
