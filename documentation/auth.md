# auth

- register
- login
- forgot password
- authenticating API calls and authorizing the actions


User comes on the home page.
sees two buttons -> Login and Sign Up

First let's talk about Sign Up page.

Sign Up page will show user a form:

- Take username
- Take Email
- Take Password
- Confirm Password

- On submitting the form, the details will be verified on frontend and sent to the backend.
- After one more round of verification, details will be stored in the database and 201 created response will be sent.
- if details are not proper, 400 bad request will be sent.
