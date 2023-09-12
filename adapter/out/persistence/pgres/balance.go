package pgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sqr "github.com/Masterminds/squirrel"
	"github.com/dmitriyGarden/test-transaction-service/model"
	"github.com/google/uuid"
)

const balanceTable = "transactionschema.balance"

func (d *DB) GetBalance(ctx context.Context, uid uuid.UUID) (int64, error) {
	q, args, err := sqr.Select("balance").
		From(balanceTable).
		Where(sqr.Eq{"uid": uid}).
		PlaceholderFormat(sqr.Dollar).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("ToSql: %w", err)
	}
	res := int64(0)
	err = d.db.GetContext(ctx, &res, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, model.ErrNotFound
		}
		return 0, fmt.Errorf("db.GetContext: %w", err)
	}
	return res, nil
}

func (d *DB) AddBalance(ctx context.Context, uid uuid.UUID, sum int64) error {
	q := `INSERT INTO transactionschema.balance (uid, balance)
	VALUES ($1, $2)
	ON CONFLICT (uid)
	DO UPDATE
	SET balance = balance.balance + $2
	WHERE balance.uid = $1
`
	_, err := d.db.ExecContext(ctx, q, uid, sum)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}
	return nil
}
