-- public.balance definition

-- Drop table

-- DROP TABLE public.balance;

CREATE TABLE public.balance (
	account int4 NOT NULL,
	balance int4 NOT NULL,
	CONSTRAINT balance_un UNIQUE (account)
);

INSERT INTO public.balance (account,balance) VALUES
(1,100),
(2,100);