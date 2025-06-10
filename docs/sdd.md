# System Design Document

## Context
### Problem
Building a service which allows users to subscribe on periodical email notification about the weather. Users can subscribe on their city and select preferred notification frequency. Service should send weather data in user's region in letter on email.

### Requirements
 - Functional requirements
	 - User should be able to register their email, city and preferred frequency.
	 - Each user can select only one city to get notification for.
	 - User should be able to select daily or hourly frequency.
	 - User should be able to unsubscribe from mailing list.
	 - Service should authorize user's email after subscription.
 - Non-functional requirements
	 - Accessibility: Service should be available for >99% of time.
	 - Extensibility: Service should be easily scalable and modifiable.
	 - Reliability: \>95% of users should receive email successfully.
	 - Scalability: Up to 10000 users, 100000 emails per day.
 - Resources:
	 - Software Developer, QA support, Designer, Marketing specialist, external mailing API, external weather API.
	 - Mailing service has limit of 500 transactions per hour. Weather service has 1000 transactions per day.
	 - A little budget.

### Motivation
Users can easily get recent weather updates to be informed.

## Solutions

### User Stories
- As user, I want to subscribe on my city weather updates to periodically receive current weather conditions on my email.
- As user, I want to be able to unsubscribe in any time from notificator in order to stop receiving update letters on my email or change city.
- As user, I want to be able to select a period of receiving letters to be informed only in periods I want.

### Integration

![image]("./images/Editor _ Mermaid Chart-2025-06-09-191322.png")

- Raw HTML/CSS Templates. In order to keep simplicity and speed, raw HTML/CSS templates are used. It's totally enough to use them for the task.
- Main Backend Server. It is responsible for providing API for users. The server is built in 2-layer monolithic way with thick controllers. It's lightweight, easily extendable and simple.
- External Mailing Service. It's used to send email letters to customers for authorization and weather updates. FE, Twilio Sendgrid can be used.
- External Weather Service. Service should be able to get current weather conditions (like temperature and humidity) for given location.
- Database. Relational database is used to store users profile data and available tokens.
- Notifier. Just a separate thread that runs CRON jobs for notifying user periodically.

### Database Schema
**Subscriptions**
| Column | Type | Description |
|----------|------|--------|
|ID|Serial|Unique user identifier|
|Email|String(255)|User email. Should be unique, as user is allowed to subscribe on only one city update|
|City|String(255)|City subscribed on|
|Frequency|String(255)| Notification period `daily` \| `hourly`|
|Confirmed|Boolean(False)|If user is validated |

**Tokens**
| Column | Type | Description |
|----------|------|--------|
|ID|UUIDv4|Unique token
|Expires|Timestamp|Token expiration datetime
|Subscription ID|Serial|User which owns token
|Created At|Timestamp|Whan token created

### RestFul API Description

GET _/weather_
---
Returns the current weather forecast for the specified city using WeatherAPI.com.
#### Parameters
**Query**:
- city* (_string_)
	City name for weather forecast

**Responses**

200 Successful operation - current weather forecast returned. Example:
```json
{
  "temperature": 0,
  "humidity": 0,
  "description": "string"
}
```
400 Invalid Request
404 City not found
POST _/subscribe_
---
Subscribe an email to receive weather updates for a specific city with chosen frequency.
#### Parameters
**Form Data**:
|Name|Description|
|-------|--|
|email *(string)|Email address to subscribe
|city *(string)|City for weather updates
|frequency *(string)|Frequency of updates (hourly or daily)
#### Responses
200 Subscription successful. Confirmation email sent.
400 Invalid input
409 Email already subscribed

GET _/confirm/{token}_
---
Confirms a subscription using the token sent in the confirmation email.
#### Parameters
**Path**
token * (string) -- Confirmation token
#### Responses
200 Subscription confirmed successfully
400 Invalid token
404 Token not found

GET _/unsubscribe/{token}_
---
Unsubscribes an email from weather updates using the token sent in emails.
#### Parameters
**Path**
token * (string) -- Confirmation token
#### Responses
200 Subscription confirmed successfully
400 Invalid token
404 Token not found

### Stack
Programming language: Golang. Golang is simple and straightforward programming language with easy access to build concurrent systems (goroutines).
Database: Postgres, GORM (In future should be removed, as raw SQL queries are preferred).
Web-Framework: Gin (In future should be replaced with native `net/http` implementation).
Caching: Internal (Hash-Map with read-write-mutex).

## Reviewers
 - [ ] DzOlha
 - [ ] Nikita Dmytriyenko

##  C4 Levels

### Level 1 -- System context diagram

![image]("./images/C4-1.png")

### Level 2 -- System context diagram

![image]("./images/C4-2.drawio.png")

### Level 3 -- System context diagram

![image]("./images/C4-3.drawio.png")

### Level 4 -- System context diagram

![image]("./images/C4-4.drawio.png")
