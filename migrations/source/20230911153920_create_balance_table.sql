-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactionschema.balance (
    uid uuid,
    balance integer,
    PRIMARY KEY (uid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactionschema.balance;
-- +goose StatementEnd
