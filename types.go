package main

type GreetInput struct {
	Name string `json:"name" jsonschema:"il nome della persona da salutare"`
}

type GreetOutput struct {
	Message string `json:"message" jsonschema:"il messaggio di saluto"`
}

