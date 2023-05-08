package db

const sqlCreateRewards = `
CREATE TABLE IF NOT EXIST rewards (
	id varchar PRIMARY KEY,
	type integer,
	reward double precision
);
`

const sqlAddReward = `
INSERT INTO rewards AS r VALUES($1, $2, $3) ON CONFLICT (id) DO UPDATE SET type = $2 reward = $3;
`

const sqlFindRewards = `
SELECT * FORM rewards r WHERE $1 = '%' || r.id || '%';
`
