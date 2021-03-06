package query

import (
	"context"

	"github.com/payfazz/fazzlearning-api/config"
	"github.com/payfazz/fazzlearning-api/internal/domains/todo/model"
	"github.com/payfazz/fazzlearning-api/internal/domains/todo/repository"
	"github.com/payfazz/fazzlearning-api/internal/values/todo"
	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// NewTodoQuery is a function to create new Query Instance
func NewTodoQuery() TodoQueryInterface {
	db := config.GetDb()
	q := fazzdb.QueryDb(db, config.GetIfQueryConfig(config.I_QUERY_CONFIG))

	return &todoQuery{
		Repository: repository.NewTodoRepository(q),
	}
}

func (q *todoQuery) Find(ctx context.Context, id int64) (*model.Todo, error) {
	return q.Repository.Find(ctx, id)
}

func (q *todoQuery) All(ctx context.Context, queryParams map[string]string) ([]*model.Todo, error) {
	offset := 0
	limit := 20

	if "" != queryParams["limit"] {
		limit = (formatter.StringToInteger(queryParams["limit"]))
	}

	if "" != queryParams["page"] {
		offset = (formatter.StringToInteger(queryParams["page"]) - 1) * limit
	}

	conditions := []fazzdb.SliceCondition{
		{Connector: fazzdb.CO_NONE, Field: "status", Operator: fazzdb.OP_LIKE, Value: todo.ACTIVE},
		{Connector: fazzdb.CO_OR, Field: "status", Operator: fazzdb.OP_LIKE, Value: todo.INPROGRESS},
	}

	return q.Repository.All(ctx, conditions, nil, limit, offset)
}

func (q *todoQuery) Count(ctx context.Context) (*float64, error) {
	return q.Repository.Count(ctx)
}
