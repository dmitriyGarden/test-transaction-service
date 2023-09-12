-- +goose Up
-- +goose StatementBegin
CREATE schema transactionschema;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP schema transactionschema;
-- +goose StatementEnd
