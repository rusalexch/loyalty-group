package order

// запрос на создание типа статус заказа
const sqlCreateOrderStatus = `
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
			CREATE TYPE order_status AS ENUM ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED');
    END IF;
END$$;
`

// sqlCreateOrders запрос на создание таблицы заказов
const sqlCreateOrders = `
CREATE TABLE IF NOT EXISTS orders (
	id bigint PRIMARY KEY,
	status order_status DEFAULT 'REGISTERED',
	accrual double precision DEFAULT null
);
`

// sqlAddNewOrder запрос на добавление нового заказа:
// $1 - идентификатор (id) заказа
const sqlAddNewOrder = `INSERT INTO orders VALUES($1) ON CONFLICT (id) DO NOTHING;`

// sqlFindByID запрос на поиск заказа:
// $1 - идентификатор (id) заказа
const sqlFindByID = `SELECT * FROM orders WHERE id = $1;`

// sqlUpdateStatus запрос на обновление статуса заказа:
// $1 - идентификатор (id) заказа
// $2 - статус (status) заказа
const sqlUpdateStatus = `UPDATE orders SET status = $2 WHERE id = $1;`

// sqlUpdate запрос на обновление заказа:
// $1 - идентификатор (id) заказ
// $2 - статус (status) заказа
// $3 - начисления (accrual) заказа
const sqlUpdate = `UPDATE orders SET status = $2, accrual = $3 WHERE id = $1;`

// sqlDelete запрос на удаление заказа
// $1 - номер заказа
const sqlDelete = `DELETE FROM orders WHERE id = $1;`

const sqlFindRegistered = `SELECT id FROM orders WHERE status = 'REGISTERED';`
