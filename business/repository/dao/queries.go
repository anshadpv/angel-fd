package dao

const (
	InsertWebookEvent = `INSERT INTO webhook_events (client_code, vendor, tracking_id, event_type, institution, type, amount, tenure_months, tenure_days, failure_reason, created_by, updated_by) VALUES($1, NULLIF($2, ''), NULLIF($3, ''), NULLIF($4, ''), NULLIF($5, ''), NULLIF($6, ''), NULLIF($7, 0), NULLIF($8, 0), NULLIF($9, 0), NULLIF($10, ''), $11, $12)`
	FetchAllFDDetails = `WITH RankedPlans AS (
		SELECT
			p.fsi AS "fsi",
			b.name AS "name",
			p.plan_type AS "type",
			p.tenure_years AS "tenureYears",
			p.tenure_months AS "tenureMonths",
			p.tenure_days AS "tenureDays",
			p.interest_rate AS "interestRate",
			p.lockin_months AS "lockinMonths",
			COALESCE(p.women_benefit, 0) AS "womenBenefit",
			COALESCE(p.senior_citizen_benefit, 0) AS "seniorCitizen",
			b.image_url AS "imageUrl",
			p.is_mostbought as "isMostbought",
			'' AS "description",
			CASE
				WHEN p.is_insured = true THEN COALESCE(b.insured_amount, 0)
				ELSE 0
			END AS "insuredAmount",
			ROW_NUMBER() OVER (PARTITION BY p.fsi ORDER BY p.interest_rate DESC) AS "row_num"
		FROM
			plans p
		LEFT JOIN
			banks b ON p.fsi = b.fsi
		WHERE p.is_active = true
	)
	SELECT
		"fsi",
		"name",
		"type",
		"tenureYears",
		"tenureMonths",
		"tenureDays",
		"interestRate",
		"lockinMonths",
		"womenBenefit",
		"seniorCitizen",
		"imageUrl",
		"isMostbought",
		"description",
		"insuredAmount"
	FROM
		RankedPlans
	WHERE
		"row_num" = 1;
	
`

	FetchFDHighIR = `WITH RankedPlans AS (
	SELECT
		p.fsi AS "fsi",
		b.name AS "name",
		p.plan_type AS "type",
		p.tenure_years AS "tenureYears",
		p.tenure_months AS "tenureMonths",
		p.tenure_days AS "tenureDays",
		p.interest_rate AS "interestRate",
		p.lockin_months AS "lockinMonths",
		COALESCE(p.women_benefit, 0) AS "womenBenefit",
		COALESCE(p.senior_citizen_benefit, 0) AS "seniorCitizen",
		b.image_url AS "imageUrl",
		p.is_mostbought as "isMostbought",
		'' AS "description",
		CASE
			WHEN p.is_insured = true THEN COALESCE(b.insured_amount, 0)
			ELSE 0
		END AS "insuredAmount",
		ROW_NUMBER() OVER (PARTITION BY p.fsi ORDER BY p.interest_rate DESC) AS "row_num"
	FROM
		plans p
	LEFT JOIN
		banks b ON p.fsi = b.fsi
	WHERE p.is_active = true
)
SELECT
	"fsi",
	"name",
	"type",
	"tenureYears",
	"tenureMonths",
	"tenureDays",
	"interestRate",
	"lockinMonths",
	"womenBenefit",
	"seniorCitizen",
	"imageUrl",
	"isMostbought",
	"description",
	"insuredAmount"
FROM
	RankedPlans
WHERE
	"row_num" = 1;
`
	FetchFDMinIR = `WITH RankedPlans AS (
	SELECT
		p.fsi AS "fsi",
		b.name AS "name",
		p.plan_type AS "type",
		p.tenure_years AS "tenureYears",
		p.tenure_months AS "tenureMonths",
		p.tenure_days AS "tenureDays",
		p.interest_rate AS "interestRate",
		p.lockin_months AS "lockinMonths",
		COALESCE(p.women_benefit, 0) AS "womenBenefit",
		COALESCE(p.senior_citizen_benefit, 0) AS "seniorCitizen",
		b.image_url AS "imageUrl",
		p.is_mostbought as "isMostbought",
		'' AS "description",
		CASE
			WHEN p.is_insured = true THEN COALESCE(b.insured_amount, 0)
			ELSE 0
		END AS "insuredAmount",
		ROW_NUMBER() OVER (PARTITION BY p.fsi ORDER BY p.interest_rate ASC) AS "row_num"
	FROM
		plans p
	LEFT JOIN
		banks b ON p.fsi = b.fsi
	WHERE p.is_active = true
)
SELECT
	"fsi",
	"name",
	"type",
	"tenureYears",
	"tenureMonths",
	"tenureDays",
	"interestRate",
	"lockinMonths",
	"womenBenefit",
	"seniorCitizen",
	"imageUrl",
	"isMostbought",
	"description",
	"insuredAmount"
FROM
	RankedPlans
WHERE
	"row_num" = 1;

`
	BaseFetchPlanQuery = `SELECT
	p.fsi AS "fsi",
	b.name AS "name",
	p.plan_type AS "type",
	p.tenure_years AS "tenureYears",
	p.tenure_months AS "tenureMonths",
	p.tenure_days AS "tenureDays",
	p.interest_rate AS "interestRate",
	p.lockin_months AS "lockinMonths",
	COALESCE(p.women_benefit, 0) AS "womenBenefit",
	COALESCE(p.senior_citizen_benefit, 0) AS "seniorCitizen",
	b.image_url AS "imageUrl",
	'' AS "description",
	CASE
		WHEN p.is_insured = true THEN COALESCE(b.insured_amount, 0)
		ELSE 0
	END AS "insuredAmount"
	FROM
	plans p
	LEFT JOIN
	banks b ON p.fsi = b.fsi
	WHERE
	p.is_active = true`

	FetchFsiPlansDetails = `WITH active_plans AS (
		SELECT
			p.fsi,
			b.name,
			p.plan_type,
			p.tenure_years,
			p.tenure_months,
			p.tenure_days,
			p.interest_rate,
			p.lockin_months,
			COALESCE(p.women_benefit, 0) AS women_benefit,
			COALESCE(p.senior_citizen_benefit, 0) AS senior_citizen_benefit,
			b.image_url,
			b.about,
			b.calculator,
			CASE
				WHEN p.is_insured THEN COALESCE(b.insured_amount, 0)
				ELSE 0
			END AS insured_amount,
			b.min_investment_amount
		FROM
			plans p
		LEFT JOIN
			banks b ON p.fsi = b.fsi
		WHERE
			p.is_active = true
	),
	max_interest_rate_fsi AS (
		SELECT fsi, name, interest_rate, image_url
		FROM active_plans
		where fsi <> $1 
		ORDER BY interest_rate DESC
		LIMIT 1
	)
	SELECT
		ap.fsi AS "fsi",
		ap.name AS "name",
		ap.plan_type AS "type",
		ap.tenure_years AS "tenureYears",
		ap.tenure_months AS "tenureMonths",
		ap.tenure_days AS "tenureDays",
		ap.interest_rate AS "interestRate",
		ap.lockin_months AS "lockinMonths",
		ap.women_benefit AS "womenBenefit",
		ap.senior_citizen_benefit AS "seniorCitizen",
		ap.image_url AS "imageUrl",
		ap.about AS "about",
		ap.calculator AS "calculator",
		'' AS "description",
		ap.insured_amount AS "insuredAmount",
		ap.min_investment_amount AS "minInvestment",
		(SELECT fsi FROM max_interest_rate_fsi) AS "compareFsi",
		(SELECT name FROM max_interest_rate_fsi) AS "compareFsiName",
		(SELECT interest_rate FROM max_interest_rate_fsi) AS "compareFsiInterestRate",
		(SELECT image_url FROM max_interest_rate_fsi) AS "compareFsiImageUrl"
	FROM
		active_plans ap
	WHERE
		ap.fsi = $1;`

	FetchAllPlansDetails = BaseFetchPlanQuery + ";"

	FetchMostBoughtPlanDetails = BaseFetchPlanQuery + " AND p.is_mostbought = true"

	FetchPendingJourneyDetails = `select pending, payment_pending, kyc_pending from pending_journey`

	FetchPendingForClient = FetchPendingJourneyDetails + ` where client_code = $1 and provider = $2`

	GetFAQsByTag = `select faq from faqs where tag=$1 and is_active=true`

	CompareLandingPageQuery = `WITH RankedPlans AS (
		SELECT
			p.fsi AS "fsi",
			b.name AS "name",
			b.image_url as "imageUrl",
			p.interest_rate AS "interestRate",
			ROW_NUMBER() OVER (PARTITION BY p.fsi ORDER BY p.interest_rate DESC) AS "row_num"
		FROM
			plans p
		LEFT JOIN
			banks b ON p.fsi = b.fsi
		WHERE p.is_active = true
	)
	SELECT
		"fsi",
		"name",
		"imageUrl",
		"interestRate"
	FROM
		RankedPlans
	WHERE
		"row_num" = 1;`

	CompareFsiQuery = `SELECT 
    b.fsi AS "fsi",
    b.name AS "name",
    p.tenure_years AS "tenureYears",
	p.tenure_months AS "tenureMonths",
	p.tenure_days AS "tenureDays",
    MAX(p.interest_rate) AS "interestRate",
    b.min_investment_amount AS "minDeposit",
    CASE
        WHEN COUNT(p.senior_citizen_benefit) = COUNT(*) THEN TRUE
        ELSE FALSE
    END AS "seniorCitizenBenefit",
    'Not Required' AS "bankAccount",
    CASE
		WHEN p.is_insured = true THEN COALESCE(b.insured_amount, 0)
		ELSE 0
	END AS "insuredAmount",
	b.image_url AS "imageUrl"
FROM 
    plans AS p
JOIN
    banks AS b ON p.fsi = b.fsi
WHERE 
    p.is_active = TRUE AND p.fsi IN (%s)
GROUP BY 
    b.fsi, b.name, p.tenure_years, p.tenure_months, p.tenure_days, b.min_investment_amount, p.is_insured, b.insurance_description
ORDER BY 
    p.tenure_years;`

	FsiDetailsQuery = `SELECT 
    b.fsi AS "fsi",
    b.name AS "name",
	p.plan_type AS "type",
    MAX(p.interest_rate) AS "interestRate",
	p.lockin_months AS "lockinMonths",
	COALESCE(p.women_benefit, 0) AS "womenBenefit",
	COALESCE(p.senior_citizen_benefit, 0) AS "seniorCitizen",
	b.image_url AS "imageUrl",
	b.insurance_description AS "description",
	CASE
		WHEN p.is_insured = true THEN COALESCE(b.insured_amount, 0)
		ELSE 0
	END AS "insuredAmount",
	b.about,
	b.calculator
	FROM 
    	plans AS p
	JOIN
    	banks AS b ON p.fsi = b.fsi
	WHERE 
    	p.is_active = TRUE AND p.fsi IN (%s)
	GROUP BY 
		p.fsi, b.fsi, p.plan_type, p.lockin_months, p.women_benefit, p.senior_citizen_benefit, p.is_insured`
)

// portfolio
const (
	SelectPortfolio = " select client_code, total_active_deposits, provider, invested_value, current_value, interest_earned, returns_value, returns_percentage from portfolio"

	PortfolioByClientCode = SelectPortfolio + " where client_code = $1 and provider = $2"

	FetchPortfolioClientListByProvider = "select client_code from portfolio where provider = $1 and invalid_client = $2"

	FetchRefreshPortfolioClientListByProvider = "select client_code from portfolio where provider = $1 and invalid_client = $2 and to_be_refreshed = $3"

	UpdateRefreshPortfolioClientList = "UPDATE portfolio SET to_be_refreshed = false, updated_by = 'portfolio_refresher_api', updated_at = current_timestamp WHERE client_code IN (%s);"

	CleanStalePortfolioRecords = "delete from portfolio where total_active_deposits = 0"

	InsertClientPortfolio = `INSERT INTO portfolio (client_code, provider, total_active_deposits, invested_value, current_value, interest_earned, returns_value, returns_percentage, created_by, updated_by, invalid_client, api_error)
	VALUES `

	UpdateClientPortfolio = `ON CONFLICT (client_code, provider) DO UPDATE SET
    total_active_deposits = EXCLUDED.total_active_deposits,
    invested_value = EXCLUDED.invested_value,
    current_value = EXCLUDED.current_value,
    interest_earned = EXCLUDED.interest_earned,
    returns_value = EXCLUDED.returns_value,
    returns_percentage = EXCLUDED.returns_percentage,
	updated_by = EXCLUDED.updated_by,
	updated_at =current_timestamp,
	invalid_client = EXCLUDED.invalid_client,
	api_error = EXCLUDED.api_error;`
)

// pending journey
const (
	FetchPendingJourneyClientListByProvider = "select client_code from pending_journey where provider = $1  and invalid_client = $2"

	FetchRefreshPendingJourneyClientListByProvider = "select client_code from pending_journey where provider = $1 and invalid_client = $2 and to_be_refreshed = $3"

	CleanStalePendingJourneyRecords = "delete from pending_journey where pending = false"

	UpdateRefreshPendingJourneyClientList = "UPDATE pending_journey SET to_be_refreshed = false, updated_by = 'pending_journey_refresher_api', updated_at = current_timestamp WHERE client_code IN (%s);"

	InsertPendingJourneyDetails = `INSERT INTO pending_journey (
		client_code, provider, pending, payment_pending, kyc_pending, created_by, updated_by, invalid_client, api_error
	) VALUES `

	UpdatePendingJourneyDetails = `ON CONFLICT (client_code, provider) DO UPDATE SET
    pending = EXCLUDED.pending,
    payment_pending = EXCLUDED.payment_pending,
    kyc_pending = EXCLUDED.kyc_pending,
	updated_by = EXCLUDED.updated_by,
    updated_at =current_timestamp,
	invalid_client = EXCLUDED.invalid_client,
	api_error = EXCLUDED.api_error;`
)
