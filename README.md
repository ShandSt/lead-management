<title>README - Lead Management API</title>
<h1>Overview</h1>
<p>This document describes the endpoints available in the Lead Management API, which allows clients to manage leads, client data, and lead assignment processes efficiently.</p>

<h2>API Endpoints</h2>

1. Create Client
   POST /clients

Creates a new client in the system.
<code>
{
	"ID": "string",
	"Name": "string",
	"WorkingHours": {
		"Start": "datetime",
		"End": "datetime"
	},
	"Priority": "integer",
	"LeadCount": "integer",
	"Capacity": "integer"
}
</code>

Response
201 Created: Successfully created client.
400 Bad Request: Invalid request data.

<h2>2. Get All Clients</h2>
   GET /clients

Retrieves a list of all clients.

Response
200 OK: Successfully retrieved all clients. Returns an array of clients.
404 Not Found: No clients found.
Sample Response
<code>
[
	{
		"ID": "string",
		"Name": "string",
		"WorkingHours": {
			"Start": "datetime",
			"End": "datetime"
		},
		"Priority": "integer",
		"LeadCount": "integer",
		"Capacity": "integer"
	}
]
</code>

<h2>3. Get Client</h2>
   GET /clients/:id

Retrieves details of a specific client.

Path Parameters
id (string): Unique identifier of the client.
Response
200 OK: Successfully retrieved client.
404 Not Found: Client not found.
Sample Response
<code>
{
	"ID": "string",
	"Name": "string",
	"WorkingHours": {
		"Start": "datetime",
		"End": "datetime"
	},
	"Priority": "integer",
	"LeadCount": "integer",
	"Capacity": "integer"
}
</code>

<h2>4. Assign Lead</h2>
   GET /clients/assignLead

Assigns a lead to the specified client.

Query Parameters
clientID (string): Unique identifier of the client to assign a lead to.
Response
200 OK: Successfully assigned a lead.
404 Not Found: Client not found.
Sample Response
<code>
{
	"message": "Lead assigned",
	"client": {
		"ID": "string",
		"Name": "string",
		"WorkingHours": {
		"Start": "datetime",
		"End": "datetime"
	},
		"Priority": "integer",
		"LeadCount": "integer",
		"Capacity": "integer"
	}
}
</code>

Development and Testing
This API is developed using Go with the Gin framework. For testing, use tools such as Postman or cURL to make requests to the local server usually running at http://localhost:8001.

Running Tests
Execute the following command in the terminal to run the automated tests:

<code>
go test ./api_test
</code>
