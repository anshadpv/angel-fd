-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banks (
    fsi varchar(20) NOT NULL,
    name varchar(50) NULL,
    image_url varchar(100) NULL,
    insurance_description varchar(50) NULL,
    min_investment_amount int4 NULL,
    insured_amount int4 NULL,
    about json NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by varchar(50) NOT NULL,
    updated_by varchar(50) NOT NULL,
    CONSTRAINT banks_pkey PRIMARY KEY (fsi)
);

INSERT INTO public.banks
(fsi, "name", image_url, insurance_description, min_investment_amount, insured_amount, about, created_at, updated_at, created_by, updated_by)
VALUES
('STFCIN', 'Shriram Finance', 'https://d3usff6y6s0r8b.cloudfront.net/fd_logo/STFCIN.png', 'Insured upto 5L', 5000, 0, '{"aboutInfo":[{"Trusted":"2 Cr+ transacting customers"},{"Profitable":"Annual profit of ₹6,000 Cr+"},{"Highly Ranked":"Rated AA+/ Stable by ICRA"}],"earlyWithdrawal":"No penalty","lockinPeriod":"No lock in"}'::json, '2024-06-06 12:03:33.433', '2024-06-06 12:03:33.433', 'test', 'test'),
('BJFLIN', 'Bajaj Finance', 'https://d3usff6y6s0r8b.cloudfront.net/fd_logo/BJFLIN.png', 'Insured upto 5L', 15000, 0, '{"aboutInfo":[{"Trusted":"5 Cr+ transacting customers"},{"Profitable":"Annual profit of ₹10,000 Cr+"},{"Highly Ranked":"CRISIL: AAA and ICR: AAA Rated"}],"earlyWithdrawal":"No penalty","lockinPeriod":"No lock in"}'::json, '2024-06-06 12:00:59.853', '2024-06-06 12:00:59.853', 'test', 'test'),
('UTKSIN', 'Utkarsh Small Finance Bank', 'https://d3usff6y6s0r8b.cloudfront.net/fd_logo/UTKSIN.png', 'Insured upto 5L', 1000, 500000, '{"aboutInfo":[{"DICGC Insured":"100% insured up to ₹5 Lakh"},{"Trusted":"30 Lakh+ customers"},{"Highly Ranked":"2nd Rank among SFBs in India"}],"earlyWithdrawal":"No penalty","lockinPeriod":"No lock in"}'::json, '2024-06-06 12:00:31.508', '2024-06-06 12:00:31.508', 'test', 'test'),
('SMCBIN', 'Shivalik Small Finance Bank', 'https://d3usff6y6s0r8b.cloudfront.net/fd_logo/SMCBIN.png', 'Insured upto 5L', 1000, 500000, '{"aboutInfo":[{"DICGC Insured":"100% insured up to ₹5 Lakh"},{"Credible":"Founded in 1998"},{"Trusted":"6 Lakh+ customers"}],"earlyWithdrawal":"No penalty","lockinPeriod":"7 Days"}'::json, '2024-06-06 11:59:51.952', '2024-06-06 11:59:51.952', 'test', 'test');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE banks DROP CONSTRAINT banks_pkey;
drop table banks;
-- +goose StatementEnd