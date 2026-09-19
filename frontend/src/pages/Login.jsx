import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../services/api";

function Login() {
    const navigate = useNavigate();

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const handleSubmit = async (event) => {
        event.preventDefault();

        setError("");

        if (!email.trim() || !password) {
            setError("Please enter your email and password.");
            return;
        }

        try {
            setLoading(true);

            const data = await login(
                email.trim(),
                password
            );

            localStorage.setItem("token", data.token);

            if (data.user) {
                localStorage.setItem(
                    "user",
                    JSON.stringify(data.user)
                );
            }

            navigate("/create");
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page">
            <div className="card auth-card">
                <div className="card-header">
                    <h1>Welcome back</h1>
                    <p>
                        Sign in to create and manage your polls.
                    </p>
                </div>

                <form onSubmit={handleSubmit}>
                    <label htmlFor="email">Email</label>

                    <input
                        id="email"
                        type="email"
                        placeholder="you@example.com"
                        value={email}
                        onChange={(event) =>
                            setEmail(event.target.value)
                        }
                        autoComplete="email"
                    />

                    <label htmlFor="password">Password</label>

                    <input
                        id="password"
                        type="password"
                        placeholder="Enter your password"
                        value={password}
                        onChange={(event) =>
                            setPassword(event.target.value)
                        }
                        autoComplete="current-password"
                    />

                    {error && (
                        <div className="error-message">
                            {error}
                        </div>
                    )}

                    <button
                        type="submit"
                        className="primary-button"
                        disabled={loading}
                    >
                        {loading ? "Signing in..." : "Sign in"}
                    </button>
                </form>

                <p className="auth-footer">
                    Don't have an account?{" "}
                    <Link to="/signup">
                        Create one
                    </Link>
                </p>
            </div>
        </div>
    );
}

export default Login;
