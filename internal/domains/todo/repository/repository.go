package repository

import (
	"context"

	"github.com/payfazz/fazzlearning-api/internal/domains/todo/model"
	"github.com/payfazz/fazzlearning-api/lib/tx"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// NewTodoRepository is a constructor to construct repo
func NewTodoRepository(q *fazzdb.Query) TodoRepositoryInterface {
	return &todoRepository{
		Q:    q,
		Todo: model.TodoModel(),
	}
}

func (r *todoRepository) Create(ctx context.Context, m *model.Todo) (*int64, error) {
	result, errTrans := tx.RunDefault(r.Q.Db, func(q *fazzdb.Query) (interface{}, error) {
		r, err := r.Q.Use(m).InsertCtx(ctx, false)

		if nil != err {
			return nil, err
		}

		return r, nil
	})

	if nil != errTrans {
		return nil, errTrans
	}

	id := result.(int64)

	return &id, nil
}

// Find a function that used to find the data by id
func (r *todoRepository) Find(ctx context.Context, id int64) (*model.Todo, error) {
	rows, errTrans := tx.RunDefault(r.Q.Db, func(q *fazzdb.Query) (interface{}, error) {
		r, err := r.Q.Use(r.Todo).
			Where("id", id).
			WithLimit(1).
			AllCtx(ctx)

		if nil != err {
			return nil, err
		}

		return r, nil
	})

	if nil != errTrans {
		return nil, errTrans
	}

	results := rows.([]*model.Todo)
	if len(results) == 0 {
		return nil, nil
	}

	result := results[0]

	return result, nil
}

func (r *todoRepository) All(ctx context.Context, conditions []fazzdb.SliceCondition, orders []fazzdb.Order, limit int, offset int) ([]*model.Todo, error) {
	rows, errTrans := tx.RunDefault(r.Q.Db, func(q *fazzdb.Query) (interface{}, error) {
		current := r.Q.Use(r.Todo).
			WhereMany(conditions...).
			OrderByMany(orders...)

		if limit > 0 {
			current.WithLimit(limit)
		}

		if offset > 0 {
			current.WithOffset(offset)
		}

		r, err := current.AllCtx(ctx)

		if nil != err {
			return nil, err
		}

		return r, nil
	})

	if nil != errTrans {
		return nil, errTrans
	}

	result := rows.([]*model.Todo)

	return result, nil
}

func (r *todoRepository) Count(ctx context.Context) (*float64, error) {
	result, errTrans := tx.RunDefault(r.Q.Db, func(q *fazzdb.Query) (interface{}, error) {

		count, err := r.Q.Use(r.Todo).
			CountCtx(ctx)

		if nil != err {
			return nil, err
		}

		return count, nil
	})

	if nil != errTrans {
		return nil, errTrans
	}

	return result.(*float64), nil
}
