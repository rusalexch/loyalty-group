package reward

// запрос создания таблицы начислений
const sqlCreateRewards = `
CREATE TABLE IF NOT EXISTS rewards (
	id varchar PRIMARY KEY,
	type varchar NOT NULL,
	reward double precision NOT NULL
);
`

// запрос добавления нового значения начисления
// $1 - идентификатор начисления
// $2 - тип начисления
// $3 - значение начисления
const sqlAddReward = `
INSERT INTO rewards VALUES($1, $2, $3) ON CONFLICT (id) DO UPDATE SET type = $2, reward = $3;
`

// запрос поиска начисления по названию товара
// $1 - название товара
const sqlFindRewards = `
SELECT * FROM rewards r WHERE LOWER($1) LIKE '%' || LOWER(r.id) || '%' limit 1;
`

// запрос поиска начисления по названию товара
// $1 - идентификатор начисления
const sqlFindByID = `
SELECT * FROM rewards r WHERE r.id = $1;
`
