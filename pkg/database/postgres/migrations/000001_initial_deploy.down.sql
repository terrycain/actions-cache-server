BEGIN TRANSACTION;

DROP TABLE cache_part;
DROP TABLE cache;
DROP SEQUENCE seq_cache;
DROP SEQUENCE seq_cache_part;

COMMIT;