import { Routes, Route, Navigate } from "react-router-dom";
import Login from "./pages/Login.jsx";
import Signup from "./pages/Signup.jsx";
import CreatePoll from "./pages/CreatePoll.jsx";
import Poll from "./pages/Poll.jsx";
import Navbar from "./components/Navbar.jsx";

function App() {
    return (
        <>
            <Navbar />

            <Routes>
                <Route
                    path="/"
                    element={<Navigate to="/login" replace />}
                />

                <Route
                    path="/login"
                    element={<Login />}
                />

                <Route
                    path="/signup"
                    element={<Signup />}
                />

                <Route
                    path="/create"
                    element={<CreatePoll />}
                />

                <Route
                    path="/poll/:id"
                    element={<Poll />}
                />
            </Routes>
        </>
    );
}

export default App;
