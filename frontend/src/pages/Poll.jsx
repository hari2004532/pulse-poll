import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import {
    getPoll,
    vote,
    connectToPoll,
} from "../services/api";

function Poll() {
    const { id } = useParams();

    const [poll, setPoll] = useState(null);
    const [results, setResults] = useState({});
    const [selectedOption, setSelectedOption] = useState("");
    const [loading, setLoading] = useState(true);
    const [voting, setVoting] = useState(false);
    const [voted, setVoted] = useState(false);
    const [error, setError] = useState("");
    const [message, setMessage] = useState("");
    const [copied, setCopied] = useState(false);

    const token = localStorage.getItem("token");

    useEffect(() => {
        let socket;

        const loadPoll = async () => {
            try {
                setLoading(true);
                setError("");

                const data = await getPoll(id);

                setPoll(data.poll);
                setResults(data.results || {});

                socket = connectToPoll(
                    id,
                    (update) => {
                        if (update.pollId === id) {
                            setResults(update.results || {});
                        }
                    },
                    () => {
                        console.log("Live connection unavailable");
                    }
                );
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        loadPoll();

        return () => {
            if (socket) {
                socket.close();
            }
        };
    }, [id]);

    const handleVote = async () => {
        if (!selectedOption) {
            setError("Please select an option.");
            return;
        }

        if (!token) {
            setError("Please login before voting.");
            return;
        }

        try {
            setVoting(true);
            setError("");
            setMessage("");

            const data = await vote(token, id, selectedOption);

            setResults(data.results || {});
            setVoted(true);
            setMessage("Your vote has been recorded.");
        } catch (err) {
            if (err.message === "you have already voted") {
                setVoted(true);
            }

            setError(err.message);
        } finally {
            setVoting(false);
        }
    };

    const copyShareLink = async () => {
        const link = window.location.href;

        try {
            await navigator.clipboard.writeText(link);
            setCopied(true);

            setTimeout(() => {
                setCopied(false);
            }, 2000);
        } catch {
            setError("Unable to copy the link.");
        }
    };

    if (loading) {
        return (
            <div className="page">
                <div className="card">
                    <p>Loading poll...</p>
                </div>
            </div>
        );
    }

    if (error && !poll) {
        return (
            <div className="page">
                <div className="card">
                    <h1>Poll unavailable</h1>
                    <p className="error-message">{error}</p>
                    <Link to="/create">Create a new poll</Link>
                </div>
            </div>
        );
    }

    if (!poll) {
        return null;
    }

    const totalVotes = poll.options.reduce(
        (total, option) => total + (results[option.id] || 0),
        0
    );

    return (
        <div className="page">
            <div className="card poll-card">
                <div className="poll-top">
                    <div>
                        <span className="live-badge">
                            LIVE
                        </span>

                        <h1>{poll.question}</h1>

                        <p className="vote-count">
                            {totalVotes}{" "}
                            {totalVotes === 1 ? "vote" : "votes"}
                        </p>
                    </div>
                </div>

                <div className="share-section">
                    <input
                        type="text"
                        value={window.location.href}
                        readOnly
                    />

                    <button
                        type="button"
                        className="secondary-button"
                        onClick={copyShareLink}
                    >
                        {copied ? "Copied!" : "Copy Link"}
                    </button>
                </div>

                <div className="poll-options">
                    {poll.options.map((option) => {
                        const count = results[option.id] || 0;

                        const percentage =
                            totalVotes > 0
                                ? Math.round(
                                      (count / totalVotes) * 100
                                  )
                                : 0;

                        return (
                            <button
                                type="button"
                                key={option.id}
                                className={`poll-option ${
                                    selectedOption === option.id
                                        ? "selected"
                                        : ""
                                }`}
                                onClick={() =>
                                    !voted &&
                                    setSelectedOption(option.id)
                                }
                                disabled={voted}
                            >
                                <div className="option-content">
                                    <span>{option.text}</span>

                                    <strong>
                                        {percentage}%
                                    </strong>
                                </div>

                                <div className="progress-track">
                                    <div
                                        className="progress-bar"
                                        style={{
                                            width: `${percentage}%`,
                                        }}
                                    />
                                </div>

                                <small>
                                    {count}{" "}
                                    {count === 1
                                        ? "vote"
                                        : "votes"}
                                </small>
                            </button>
                        );
                    })}
                </div>

                {!voted && (
                    <button
                        type="button"
                        className="primary-button"
                        onClick={handleVote}
                        disabled={voting}
                    >
                        {voting ? "Submitting..." : "Submit Vote"}
                    </button>
                )}

                {voted && (
                    <div className="success-message">
                        You have already voted in this poll.
                    </div>
                )}

                {message && (
                    <div className="success-message">
                        {message}
                    </div>
                )}

                {error && (
                    <div className="error-message">
                        {error}
                    </div>
                )}
            </div>
        </div>
    );
}

export default Poll;
