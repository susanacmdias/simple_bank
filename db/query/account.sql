-- name: createAccount :one 
INSERT INTO accounts(
    owner,
    balance, 
    currency
) VALUES(
    $1, $2, $3 
) RETURNING *;

-- name: getAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: getAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: listAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: updateAccount :one
UPDATE accounts 
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: addAccountBalance :one
UPDATE accounts 
SET balance = balance + sqlc.arg(ammount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: deleteAccount :one
DELETE FROM accounts 
WHERE id = $1
RETURNING *;