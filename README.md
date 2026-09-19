# PulsePoll

A real-time live polling application where users can create polls, share them with others, vote, and see results update instantly without refreshing the page.

## Live Demo

Frontend:
https://pulse-poll-frontend.onrender.com

Backend:
https://pulse-poll-bcep.onrender.com

GitHub:
https://github.com/hari2004532/pulse-poll

---

## Features

- User signup and login
- JWT-based authentication
- Create polls with 2–10 options
- Shareable poll links
- Vote on active polls
- One vote per authenticated user per poll
- Real-time result updates without page refresh
- Redis-based live vote counters
- Redis Pub/Sub for real-time events
- WebSocket communication between backend and clients
- MongoDB persistence
- Responsive React UI
- Backend-side validation
- React Router client-side routing

---

## Tech Stack

### Frontend

- React
- Vite
- JavaScript
- React Router
- CSS

### Backend

- Go
- Gin Web Framework
- JWT Authentication
- bcrypt
- Gorilla WebSocket

### Database

- MongoDB Atlas

### Real-Time Layer

- Redis
- Redis Pub/Sub
- WebSockets

### Deployment

- Render

---

## Architecture

```text
                    ┌──────────────────────┐
                    │     React Frontend   │
                    │       (Vite)         │
                    └──────────┬───────────┘
                               │
                         REST API / HTTP
                               │
                               ▼
                    ┌──────────────────────┐
                    │    Go + Gin Backend  │
                    └───────┬───────┬──────┘
                            │       │
                     MongoDB       Redis
                            │       │
                            │       ├── Vote Counters
                            │       │
                            │       └── Redis Pub/Sub
                            │                │
                            │                ▼
                            │        WebSocket Hub
                            │                │
                            └────────────────┘
                                     │
                                     ▼
                              Connected Clients

How Real-Time Voting Works

The application uses Redis Pub/Sub and WebSockets to deliver live poll results.

When a user votes:

User selects an option
        ↓
React sends vote to Go backend
        ↓
Backend validates the poll and option
        ↓
Vote is stored in MongoDB
        ↓
Redis increments the option's vote count
        ↓
Redis Pub/Sub publishes the updated results
        ↓
WebSocket subscriber receives the update
        ↓
WebSocket Hub broadcasts the update
        ↓
Connected React clients receive the update
        ↓
Poll results update without page refresh

Redis is actively used for the real-time functionality through:

Vote counters
Redis Pub/Sub
Real-time event distribution

MongoDB is used for persistent application data.

Authentication

PulsePoll uses JWT-based authentication.

Users can:

Create an account
Login
Receive an authentication token
Use the token for protected operations such as poll creation and voting

Passwords are securely hashed using bcrypt before being stored.

Poll Creation

Authenticated users can create polls.

Poll validation includes:

Question cannot be empty
Question length validation
Minimum 2 options
Maximum 10 options
Options cannot be empty
Duplicate options are rejected

Each poll option receives a unique identifier.

Voting

Users must be authenticated before voting.

The backend validates:

Poll ID
Poll existence
Poll active status
Option ID
Existing vote

The application prevents a user from voting more than once on the same poll.

A unique MongoDB index on:

pollId + voterId

provides database-level protection against duplicate votes.

API Endpoints
Authentication
POST /api/auth/signup
POST /api/auth/login
Polls
GET  /api/polls/:id
POST /api/polls
Voting
POST /api/polls/:id/vote
WebSocket
/ws/polls/:id
Health Check
GET /health
Project Structure
pulse-poll/
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   └── Navbar.jsx
│   │   ├── pages/
│   │   │   ├── Login.jsx
│   │   │   ├── Signup.jsx
│   │   │   ├── CreatePoll.jsx
│   │   │   └── Poll.jsx
│   │   ├── services/
│   │   │   └── api.js
│   │   ├── App.jsx
│   │   ├── main.jsx
│   │   └── index.css
│   ├── .env.example
│   ├── package.json
│   └── vite.config.js
│
├── backend/
│   ├── config/
│   │   ├── config.go
│   │   └── database.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── poll.go
│   │   └── vote.go
│   ├── middleware/
│   │   └── auth.go
│   ├── models/
│   │   ├── user.go
│   │   ├── poll.go
│   │   └── vote.go
│   ├── routes/
│   │   └── routes.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── poll_service.go
│   │   ├── redis_service.go
│   │   └── vote_service.go
│   ├── websocket/
│   │   ├── client.go
│   │   └── hub.go
│   ├── main.go
│   ├── go.mod
│   └── .env.example
│
├── README.md
└── .gitignore
Environment Variables
Backend

Create a .env file inside the backend directory.

PORT=8080
MONGODB_URI=your_mongodb_connection_string
MONGODB_DATABASE=pulsepoll
REDIS_ADDR=your_redis_host:6379
REDIS_PASSWORD=your_redis_password
JWT_SECRET=your_jwt_secret
FRONTEND_URL=http://localhost:5173
Frontend

Create a .env file inside the frontend directory.

VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080

Do not commit real credentials, passwords, API keys, or secrets to GitHub.

Running Locally
1. Clone the repository
git clone https://github.com/hari2004532/pulse-poll.git
cd pulse-poll
2. Start the Backend
cd backend
go mod download
go run .

The backend runs on:

http://localhost:8080
3. Start the Frontend

Open another terminal:

cd frontend
npm install
npm run dev

The frontend runs on:

http://localhost:5173
Testing

The application can be tested using two browser windows.

Authentication
Create a new account.
Login with the account.
Open the Create Poll page.
Create a Poll
Enter a poll question.
Add at least two options.
Create the poll.
Copy the generated poll URL.
Real-Time Voting
Open the poll in Browser 1.
Open the same poll in Browser 2.
Login as another user in Browser 2.
Vote from Browser 2.
Observe Browser 1.
The results should update automatically without refreshing.
Duplicate Vote Protection

Attempt to vote again using the same account.

The backend rejects the duplicate vote.

Deployment

The application is deployed using Render.

Frontend

https://pulse-poll-frontend.onrender.com

Backend

https://pulse-poll-bcep.onrender.com

The frontend is deployed as a Render Static Site and the backend is deployed as a Render Web Service.

MongoDB Atlas provides persistent storage, while Redis provides real-time vote counters and Pub/Sub functionality.

Key Design Decisions
MongoDB

MongoDB is used for persistent application data:

Users
Polls
Votes
Redis

Redis is used for:

Real-time vote counters
Redis Pub/Sub
Distributing poll updates
WebSockets

WebSockets provide a persistent connection between the browser and backend so that connected users receive poll updates immediately.

JWT

JWT is used to authenticate protected API requests.

AI Usage

AI tools were used during development as a development assistant for:

Project scaffolding
Code suggestions
Debugging assistance
Reviewing implementation approaches
UI improvements
Deployment troubleshooting

The implemented application was tested locally and in production, including authentication, poll creation, voting, duplicate-vote protection, Redis-based real-time updates, WebSocket communication, and deployment.

The developer understands the implemented architecture and functionality and can explain the design and implementation decisions.

Author

Hari

GitHub:

https://github.com/hari2004532