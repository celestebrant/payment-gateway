# Building a Payment Gateway

Steps:
1. Generate a `PaymentRequest`
1. Validate `PaymentRequest`

Run server:
`go run ./cmd/server`

Testing:
`go test ./...`

- endpoint happy
- endpoint any error
- concurrency?
- mocked happy
- validation:
    - valid card number
    - contains alphabet
    - too short
    - too long
    - expiry year = 999


## Assumptions
- For simplicity, typical British Visa and Mastercard payment card characteristics are the only accepted payment method for my solution. The characteristics are:
    - Card nyumber must be 16 numeric digits long
    - CVV must be 3 numeric digits long
    - In reality, these fields vary across countries and card payment services card, e.g. American Express cards having 15-digit card numbers and CVV of 4 digit length.
- The supported currency is either GBP or EUR, and the currency may be set to 2 decimal places, and must be positive.
- The "Aquiring Bank" handles expiry date validation, so my solution does not cover this. An edge case could otherwise arise when using localised time. If the local time of the card issuer is different to the local time of this server's deployment, there is a risk of incorrect validation. It could be possible for the card issuer date to be in a different month, or even year, to the time the server operates with. E.g. UTC time is Jan 2025, the card expiry date is set to Dec 2024 and the local time of the card issuer is also Dec 2024. For productionisation, it would be worth ensuring checking if the Aquiring Bank does indeed handle validating card expiry. If it is agreed that validation should be added to this program, then a mandatory timezone field should be added to the ProcessPaymentRequest that corresponds with the card issuer timezone.

## Areas for improvement
- Persistent storage for payments e.g. relational (SQL) database, and cache utilisation for frequently fetched payment IDs or other frequently fetched data.
- Supported currencies not hard-coded but either set via a config map, or alternatively, a separate store plus caching.
- Stronger security for stored payment data, e.g. encryption.
- Deployment in a containerised manner (e.g. Docker) and containter orchestration (e.g. Kubernetes) to handle high load.
- Deployment on a cloud instance for reduced overhead on hardware maintenance, although requiring platform engineering experience.
- Observability, performance monitoring and logging. Useful for understanding load requirements when performance analysis is done for production deployments - might find answers to questions like "High read? Or high writing? Or both?", "Why do my deployments fail sometimes?", etc.
- Load testing: stress tests, peak load, soak testing for perfomance degradation and race condition identification.