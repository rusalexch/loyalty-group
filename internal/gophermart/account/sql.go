package account

// sqlCreateTable - запрос на создание таблицы если не существует
const sqlCreateTable = `
CREATE TABLE IF NOT EXISTS transactions (
	id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
	type varchar,
	order_id bigint,
	amount double precision NOT NULL,
	processed_at timestamp,
	CONSTRAINT order_transaction_key FOREIGN KEY (order_id) REFERENCES user_orders (id)
);
`

// sqlAddDebit - запрос на добавления начисления
// $1 - идентификатор заказа
// $2 - сумма начисления
const sqlAddDebit = `
INSERT INTO transactions VALUES (DEFAULT, 'DEBIT', $1, $2, now())
`

// sqlAddCredit - запрос на добавление списания
// $1 - идентификатор заказа на который тратиться списание
// $2 - сумма списания
const sqlAddCredit = `
INSERT INTO transactions VALUES (DEFAULT, 'CREDIT', $1, $2, now())
`

// sqlGetCurrentAmount - запрос на получение сумм по поступлениям и списаниям
// $1 - идентификатор пользователя
const sqlGetUserCurrentBalance = `
WITH tr AS (
	SELECT t.* FROM transactions t
	JOIN user_orders o ON o.id = t.order_id
	WHERE o.user_id = $1
)
SELECT 
COALESCE((SELECT sum(amount) FROM tr WHERE type = 'DEBIT'), 0) as debit, 
COALESCE((SELECT sum(amount) FROM tr WHERE type = 'CREDIT'), 0) as credit;
`

const sqlGetUserCredit = `
SELECT * FROM transactions t
JOIN user_orders o ON o.id = t.order_id
WHERE o.user_id = $1 AND t."type" = 'CREDIT';
`
