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


Assumptions:
- For the sake of simplicity, typical British Visa and Mastercard payment card characteristics have been assumed to be the only accepted payment method my solution for this exercise. The characteristics are: 16-digit long card number, 3-digit long CVV, mandatory expiry date. In reality, these fields vary across countries and card payment services card, e.g. American Express cards having 15-digit card numbers and CVV of 4 digit length.
- I have assumed that the Aquiring Bank handles expiry date validation, so my solution does not cover this otherwise an edge case can arise when using localised timezone where that localised time is in a different month (and possibly different year) to UTC time, leading to incorrect validation. E.g. UTC time is Jan 2025, the card expiry date is set to Dec 2024 and the local time of the card issuer is also Dec 2024. For productionisation, it would be worth ensuring checking if the Aquiring Bank does indeed handle validating card expiry. If it is agreed that validation should be added to this program, then to avoid that edge case, then a mandatory timezone field should be added to the ProcessPaymentRequest that corresponds with the card issuer timezone.
