# go-connect-todo

Use `make boot-new` to deploy the infrastructure and then use `make port-forward` to expose to a local 8080 port to test the api via k8s (to connect via the UI).

Use `make boot` to run the service as a local docker container.

Discussion during LLD design
```
/deposit
    account_id (string)
    amount int64

/credit
    account_id (string)
    amount int64

/transaction_history
    account_id (string)
    from (timestamp)
    to (timestamp)

/balance
    account_id (string)

transaction
    id
    uuid (txn uuid)
    account_id
    amount
    created_at
    is_credit (bool)

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL,
    account_id TEXT NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_credit BOOLEAN NOT NULL
);

```