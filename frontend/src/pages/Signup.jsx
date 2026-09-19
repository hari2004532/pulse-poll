import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { signup } from "../services/api";

function Signup() {
    const navigate = useNavigate();

    const [name, setName] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    async function handleSubmit(event) {
        event.preventDefault();

        setError("");
        setLoading(true);

        try {
            await signup(name, email, password);
            navigate("/login");
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }

    return (
        <div className="auth-page">
            <div className="auth-card">
                <div className="brand">
                    <div className="brand-icon">P</div>
                    <h1>PulsePoll</h1>
                </div>

                <h2>Create your account</h2>
                <p className="subtitle">
                    Start creating interactive live polls.
                </p>

                {error && (
                    <div className="error-message">
                        {error}
                    </div>
                )}

                <form onSubmit={handleSubmit}>
                    <label htmlFor="name">Name</label>

                    <input
                        id="name"
                        type="text"
                        placeholder="Your name"
                        value={name}
                        onChange={(event) => setName(event.target.value)}
                        required
                    />

                    <label htmlFor="email">Email</label>

                    <input
                        id="email"
                        type="email"
                        placeholder="you@example.com"
                        value={email}
                        onChange={(event) => setEmail(event.target.value)}
                        required
                    />

                    <label htmlFor="password">Password</label>

                    <input
                        id="password"
                        type="password"
                        placeholder="At least 6 characters"
                        value={password}
                        onChange={(event) => setPassword(event.target.value)}
                        minLength={6}
                        required
                    />

                    <button
                        className="primary-button"
                        type="submit"
                        disabled={loading}
                    >
                        {loading ? "Creating..." : "Create account"}
                    </button>
                </form>

                <p className="auth-footer">
                    Already have an account?{" "}
                    <Link to="/login">Sign in</Link>
                </p>
            </div>
        </div>
    );
}

export default Signup;
