-- Migrate the non-standard "rmb" currency code to the ISO 4217 "cny" code.
-- "RMB" (renminbi) is a common abbreviation but is not a valid ISO 4217 code;
-- the canonical ISO code for the Chinese yuan is "CNY".
UPDATE groups SET currency = 'cny' WHERE currency = 'rmb';
