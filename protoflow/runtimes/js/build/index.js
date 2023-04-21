// This file is written in TypeScript
// Import necessary libraries
import express from "express";
// Create express app and set port
var app = express();
var port = 3000;
// Define the chatbot response given user input
function chatbotResponse(message) {
    // Add logic to determine response based on message 
    return "Hello, how can I assist you?";
}
// Handle incoming HTTP requests
app.get("/", function(req, res) {
    // Extract user input from request parameters
    var message = req.query.message;
    // Generate corresponding chatbot response
    var response = chatbotResponse(message);
    // Send the response back to the client
    res.send(response);
});
// Start the HTTP server
app.listen(port, function() {
    console.log("Server running at http://localhost:".concat(port));
});

