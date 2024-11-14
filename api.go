package main

// ApiResponse represents a list of API messages containing temperature and last modification date.
type ApiResponse []ApiMessage

// ApiMessage represents a structured message for the API containing temperature and last modification date.
type ApiMessage struct {
	Temperature  Temperature `json:"temperature"`
	LastModified MessageDate `json:"datetime"`
}
