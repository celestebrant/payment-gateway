# Payment Gateway

## API Documentation
The Payment Gateway API allows merchants to process new payments and retrieve details of previously made payments. This section provides an overview of the endpoints, request and response formats, and working examples.

### Base URL
```
http://localhost:8000
```

### Endpoints

There are 2 endpoints:
1. Process payment
2. Get payment

#### Process payment

- `POST /process-payment`
- Processes a new payment through the payment gateway.
- Headers: `Content-Type: application/json`
- Example request body
  ```json
  {
    "card_number": "1234123412341234",
    "expiry_year": 2028,
    "expiry_month": 12,
    "cvv": "123",
    "amount": 12.05,
    "currency": "GBP"
  }
  ```

*Definitions:*
- `card_number` - (mandatory) String with exactly 16 digits of numbers only.
- `expiry_year` - (mandatory) Integer value that is 4 digits long and greater than 0.
- `expiry_month` - (mandatory) Integer with value of 1 to 12, inclusive.
- `cvv` - (mandatory) String with exactly 3 digits of numbers only.
- `amount` - (mandatory) Floating-point number with a positive value and up to 2 decimal places.
- `currency` - (mandatory) String (3 characters long) with value `"GBP"` or `"EUR"`.

**Response**

Status Code
- `200 OK`, success
- `400 Bad Request`, validation error
- `500 Internal Server Error`, server error

Example body
  ```json
  {
    "id": "c08a3e62-ab97-43fc-a633-5b49f929e235",
    "status": "SUCCESS",
    "masked_card_number": "************1234",
    "expiry_year": 2028,
    "expiry_month": 12,
    "amount": 12.05,
    "currency": "GBP"
  }
  ```

*Definitions:*
- `id` - A generated ID for the payment set by the bank.
- `status` - Denotes the success of the payment. Has value `"SUCCESS"` or `"FAILED`.
- `masked_card_number` - The card number as requested, with the first 12 digits masked with `*`.
- `expiry_year` - The expiry year of the card as requested.
- `expiry_month` - The expiry month of the card as requested.
- `cvv` - (mandatory) Must be 3 digits long of numbers only.
- `amount` - (mandatory) Must be a positive number with up to 2 decimal places.
- `currency` - (mandatory) Must be either `"GBP"` or `"EUR"`.

*Example cURL request*

```sh
curl -X POST http://localhost:8000/payments \
    -H "Content-Type: application/json" \
    -d '{"card_number":"1234123412341234", "expiry_year":2028, "expiry_month":12, "cvv":"123", "amount":12.05, "currency":"GBP"}'
```

#### Get payment

- `GET /process-payment/{id}`
- Retrieves details of a previously made payment using its ID.
- Headers: `Content-Type: application/json`
- Path parameters
    - `id` - the ID of the payment to retrieve.

**Response**

Status Code
- `200 OK`, success
- `404 Not Found`, validation error

Example body
  ```json
  {
    "id": "c08a3e62-ab97-43fc-a633-5b49f929e235",
    "status": "SUCCESS",
    "masked_card_number": "************1234",
    "expiry_year": 2028,
    "expiry_month": 12,
    "amount": 12.05,
    "currency": "GBP"
  }
  ```

*Definitions:*
- `id` - A generated ID for the payment set by the bank.
- `status` - Denotes the success of the payment. Has value `"SUCCESS"` or `"FAILED`.
- `masked_card_number` - The card number as requested, with the first 12 digits masked with `*`.
- `expiry_year` - The expiry year of the card as requested.
- `expiry_month` - The expiry month of the card as requested.
- `cvv` - (mandatory) Must be 3 digits long of numbers only.
- `amount` - (mandatory) Must be a positive number with up to 2 decimal places.
- `currency` - (mandatory) Must be either `"GBP"` or `"EUR"`.

*Example cURL request*

```sh
curl -X GET http://localhost:8000/payments/c08a3e62-ab97-43fc-a633-5b49f929e235 \
    -H "Content-Type: application/json"
```

## How to run the server
Run the server locally with `go run ./cmd/server`. You should see the output "hang" like this
```
$ go run ./cmd/server
2024/07/11 22:03:55 server listening on port 8000...
```

This runs the server locally on port `8000`. You are now able to make requests.

## How to interact with the server
You can call the server by opening a separate terminal window and running a CURL command. The response will be printed:
```
$ curl -X POST http://localhost:8000/payments -H "Content-Type: application/json" -d '{"card_number":"1234567812345678", "expiry_year":2028, "expiry_month":12, "cvv":"987", "amount":12.05, "currency":"GBP"}'
{"id":"9fdbd34c-3082-4ce7-9718-369f541fa317","status":"FAILED","masked_card_number":"************5678","expiry_year":2028,"expiry_month":12,"amount":12.05,"currency":"GBP"}
```

If you look at the terminal window running the server, you will see a logged output:
```
celeste@Celestes-MacBook-Pro processout-payment-gateway % go run ./cmd/server
2024/07/11 22:03:55 server listening on port 8000...
2024/07/11 22:04:40 {9fdbd34c-3082-4ce7-9718-369f541fa317 FAILED ************5678 2028 12 12.05 GBP}
```

You can now fetch the existing payment by ID, which will also output to the server console:
```
celeste@Celestes-MacBook-Pro processout-payment-gateway % curl -X GET http://localhost:8000/payments/9fdbd34c-3082-4ce7-9718-369f541fa317 -H "Content-Type: application/json"
{"id":"9fdbd34c-3082-4ce7-9718-369f541fa317","status":"FAILED","masked_card_number":"************5678","expiry_year":2028,"expiry_month":12,"amount":12.05,"currency":"GBP"}
```

## How does the application work?
The application has two endpoints, with a handler for each:
1. A request to `POST /process-payment` calls `ProcessPaymentHandler`, for processing a new payment.
1. A request to `GET /process-payment/{id}` calls `GetPaymentHandler`, for fetching individual payments by payment ID.

### How processing payments works
`ProcessPaymentHandler` is the handler for processing payments. It works by:
1. Decoding the request into a new `ProcessPaymentRequest`, or returns a http 400 response if an error is encountered along with an error message in the body.
1. Once successfully decoded, `ProcessPaymentRequest` is validated - each field undergoes specific validation. Details on this is covered in the API documentation. If a validation error is encountered, http 400 response is returned with the validation error in the response body.
1. If validation succeeds, the "Aquiring Bank" is called to request a new payment is made with the request data provided. This is currently mocked with `MockCallBank`. The assumed response body is structured like:
    ```json
    {
        "payment_id": "c08a3e62-ab97-43fc-a633-5b49f929e235", 
        "status": "SUCCESS"
    }
    ```
    `MockCallBank` also handles decoding this response into `CallBankResponse`. If an error is encountered when making a call to the Bank, a http 500 response is returned with an error message in the reponse body.
1. If the previous step succeeds and the bank response has been successfully decoded, a new `MaskedPayment` is created and populated with data from the original http request, and the bank response (the payment ID and status).
1. This data is stored locally in memory (via `PaymentStore.AddPayment`), and also logged in the server (which you can see in the terminal window that runs the server). `MaskedPayment` contains a the masked card number and does not contain the CVV number. 
1. Finally the `MaskedPayment` is written to the response body with a http status code of 200. This is to confirm the payment has been handled successfully while providing data that could be useful for merchant accounting purposes. Reaching this point does not necessarily mean that the payment was successful on the bank's side as the payment status can either be `"SUCCESS"` OR `"FAILED"`.

### How payment retrieval works
`GetPaymentHandler` is the handler for fetching individual payments by payment ID. It works by:
1. Accessing the payment ID which is in the path of the endpoint call.
1. Fetching the payment from the payment store via `PaymentStore.GetPayment`, and if not found then returns a http 400 response with an error message expressing that the payment is not found. Similar to `PaymentStore.AddPayment`, the operation for obtaining the data in the map is surrounded by a mutex lock and unlock.
1. If payment in form `MaskedPayment` is found, it is then written to the response and the http response code is 200.

### Design choices
- `MaskedPayment` as a data structure: Payment data generally is very sensitive and the payment data that is stored in this application is masked to reduce the risk in the event of a data breach, such as masking the card number and omitting CVV.
- The in-memory payment data store, `PaymentStore.AddPayment`: The in-memory payment store, `PaymentStore`, holds a map containing masked payment data, and a mutex. Any moment the entire set of payment data is changed, the mutex is locked, the operation is performed, and then the mutex is unlocked. This is to prevent race conditions where a payment has not yet completed processing and an attempted fetch is performed concurrently (although in this current design, the payment ID is only returned upon process completion so this situation would not be possible in reality). This approach would be especially handy if the application became more complex and support for amending individual payments was added, as it would prevent fetching stale payment data.
- Logging is implemented in each handler which outputs to the server console every time a payment is processed and fetched. This would aid debugging.

## Testing
Run all tests with `go test ./...`. This runs all test files (ending with `_test.go`).
- Unit tests reside in each package.
- End to end (e2e) tests reside in `/tests`. Note that payment retrieval is best tested e2e due to usage of the store.
- There are 2 levels of testing currently: unit and integration tests.

To clean the test cache, run `go clean -testcache`.
Run run a specific test with `go test -run TestName ./path/to/test`, e.g. `go test -run TestEndToEndPaymentFlow ./tests`, for example.

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
- Concurrency tests to ensure race conditions are prevented, and suitable usage of mutex locks is in order.
- Deployment in a containerised manner (e.g. Docker) and containter orchestration (e.g. Kubernetes) to handle high load.
- Deployment on a cloud instance for reduced overhead on hardware maintenance, although requiring platform engineering experience.
- Observability, performance monitoring and logging. Useful for understanding load requirements when performance analysis is done for production deployments - might find answers to questions like "High read? Or high writing? Or both?", "Why do my deployments fail sometimes?", etc.
- Load testing: stress tests, peak load, soak testing for perfomance degradation and race condition identification.
- Test utils package and helper functions.

## Cloud technologies
It could be possible to run the computation proceses using AWS Elastic Load Balancing as this would automatically distribute traffic across EC2 instances, and may be suitable for generally high and variable load. Alternatively AWS Lambda could be an option, although would not suitable for containerised deployment. The payment details could be stored using AWS RDS with an SQL database like MySQL or PostgreSQL. The API could be explosed with AWS API Gateway and AWS CloudWatch and Performance Insights could be used for monitoring and logging including monitoring database usage.

For old payment data, it might become unnecessary and expensive to retain so an archival strategy for old (and very infrequently fetched) data could be useful for reducing the cloud bill.