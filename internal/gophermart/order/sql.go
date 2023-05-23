package order

// sqlCreateTable запрос на создание таблицы
const sqlCreateTable = `
CREATE TABLE IF NOT EXISTS user_orders (
	id bigint PRIMARY KEY,
	user_id integer,
	status order_status DEFAULT 'REGISTERED',
	accrual double precision DEFAULT null,
	uploaded_at timestamp,
	CONSTRAINT user_order_key FOREIGN KEY (user_id) REFERENCES users (id)
);
`

// sqlAdd - запрос на добавление нового заказ
// $1 - идентификатор заказа
// $2 - идентификатор пользователя
const sqlAdd = `
INSERT INTO user_orders VALUES ($1, $2, 'REGISTERED', DEFAULT, now());
`

// sqlFindByID - запрос на поиск заказа по ID
// $1 - идентификатор заказа
const sqlFindByID = `
SELECT * FROM user_orders WHERE id = $1;
`

// sqlFundByUserID - запрос на поиск заказов по идентификатору пользователя
// $1 - идентификатор пользователя
const sqlFundByUserID = `
SELECT * FROM user_orders WHERE user_id = $1 ORDER BY uploaded_at DESC;
`

// sqlFindREgistered - запрос на поиск заказов со статусом REGISTERED
const sqlFindRegistered = `
SELECT * FROM user_orders WHERE status = 'REGISTERED';
`

// sqlUpdateOrder - запрос на обновление заказа
// $1 - идентификатор заказа
// $2 - новое значение статуса заказа
// $3 - новое значение начисления на заказ
const sqlUpdateOrder = `
UPDATE user_orders SET status = $2, accrual = $3 WHERE id = $1
`
