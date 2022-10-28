# Backend Hexagonal Architecture

The expense system backend is very flat today.
* cmd - separate executables for local and Lambda-based invocation
* api.go - HTTP port, request/response bodies, logging concerns, request validation
* auth.go - adapter for Google Auth OAuth flow and token validation
* config.go - domain object (?) for a common config struct passed from executables
* orgs.go - adapter for DynamoDB and org creation/retrieval

This is quickly becoming unwieldy, and we've barely implemented anything!  I've heard good things about hexagonal architectures, so let's prototype a refactored architecture.

At the core of our system are domains.  Things like:
* Organizations
* Users
* Expense reimbursement requests
* Approvals
* Comments
* Invoices

Then comes the application.  This is the core of the business logic, and ties together various ports, adapters, and domains.  It owns validation and decides when to call which repository for a given request.  Notably, everything external relies on an interface:
* Repositories - OrgRepo, UserRepo, ReceiptRepository, etc.
* Services - AuthService, EmailService, etc.

These interfaces define the ports of our architecture to the outside world.  We then define adapters in a separate package which meet the requirements of each port.  A DynamoDB client which fetches organization details is an adapter that adheres to the OrgRepo port.
