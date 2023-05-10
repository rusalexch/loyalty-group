package order

// запрос на создание типа статус заказа
const sqlCreateOrderStatus = `
CREATE TYPE order_status AS ENUM ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED');
`

// запрос на создание таблицы заказов
const sqlCreateOrders = `
CREATE TABLE IF NOT EXIST orders (
	id varchar PRIMARY KEY,
	status order_status DEFAULT 'REGISTERED',
	accrual double precision DEFAULT null
);
`

// запрос на добавление нового заказа:
// $1 - идентификатор (id) заказа
const sqlAddNewOrder = `INSERT INTO orders VALUES($1) ON CONFLICT (id) DO NOTHING`

// запрос на поиск заказа:
// $1 - идентификатор (id) заказа
const sqlFindByID = `SELECT * FROM orders WHERE id = $1`

// запрос на обновление статуса заказа:
// $1 - идентификатор (id) заказа
// $2 - статус (status) заказа
const sqlUpdateStatus = `UPDATE orders SET status = $2 WHERE id = $1`

// запрос на обновление заказа:
// $1 - идентификатор (id) заказ
// $2 - статус (status) заказа
// $3 - начисления (accrual) заказа
const sqlUpdate = `UPDATE orders SET status = $2, accrual = $3 WHERE id = $1`
