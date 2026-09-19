import { Link, useNavigate } from "react-router-dom";

function Navbar() {
    const navigate = useNavigate();
    const token = localStorage.getItem("token");

    const handleLogout = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        navigate("/login");
    };

    return (
        <nav className="navbar">
            <Link to="/" className="brand">
                PulsePoll
            </Link>

            <div className="nav-actions">
                {token ? (
                    <>
                        <Link to="/create" className="nav-link">
                            Create Poll
                        </Link>

                        <button
                            type="button"
                            className="logout-button"
                            onClick={handleLogout}
                        >
                            Logout
                        </button>
                    </>
                ) : (
                    <>
                        <Link to="/login" className="nav-link">
                            Login
                        </Link>

                        <Link to="/signup" className="nav-signup">
                            Sign up
                        </Link>
                    </>
                )}
            </div>
        </nav>
    );
}

export default Navbar;
