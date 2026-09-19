import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPoll } from "../services/api";

function CreatePoll() {
    const navigate = useNavigate();

    const [question, setQuestion] = useState("");
    const [options, setOptions] = useState(["", ""]);
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const token = localStorage.getItem("token");

    const handleOptionChange = (index, value) => {
        const updatedOptions = [...options];
        updatedOptions[index] = value;
        setOptions(updatedOptions);
    };

    const addOption = () => {
        if (options.length < 10) {
            setOptions([...options, ""]);
        }
    };

    const removeOption = (index) => {
        if (options.length <= 2) {
            return;
        }

        setOptions(options.filter((_, i) => i !== index));
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        setError("");

        if (!token) {
            navigate("/login");
            return;
        }

        const cleanedQuestion = question.trim();
        const cleanedOptions = options
            .map((option) => option.trim())
            .filter((option) => option.length > 0);

        if (!cleanedQuestion) {
            setError("Please enter a question.");
            return;
        }

        if (cleanedOptions.length < 2) {
            setError("Please provide at least 2 options.");
            return;
        }

        if (cleanedOptions.length > 10) {
            setError("You can have a maximum of 10 options.");
            return;
        }

        setLoading(true);

        try {
            const data = await createPoll(
                token,
                cleanedQuestion,
                cleanedOptions
            );

            navigate(`/poll/${data.poll.id}`);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page">
            <div className="card create-card">
                <div className="card-header">
                    <h1>Create a Poll</h1>
                    <p>Ask a question and let your audience vote live.</p>
                </div>

                <form onSubmit={handleSubmit}>
                    <label htmlFor="question">Question</label>

                    <input
                        id="question"
                        type="text"
                        placeholder="What should we build next?"
                        value={question}
                        onChange={(event) =>
                            setQuestion(event.target.value)
                        }
                        maxLength={300}
                    />

                    <div className="options-header">
                        <label>Options</label>
                        <span>{options.length}/10</span>
                    </div>

                    <div className="options-list">
                        {options.map((option, index) => (
                            <div className="option-row" key={index}>
                                <input
                                    type="text"
                                    placeholder={`Option ${index + 1}`}
                                    value={option}
                                    onChange={(event) =>
                                        handleOptionChange(
                                            index,
                                            event.target.value
                                        )
                                    }
                                    maxLength={100}
                                />

                                {options.length > 2 && (
                                    <button
                                        type="button"
                                        className="remove-button"
                                        onClick={() =>
                                            removeOption(index)
                                        }
                                    >
                                        Remove
                                    </button>
                                )}
                            </div>
                        ))}
                    </div>

                    {options.length < 10 && (
                        <button
                            type="button"
                            className="secondary-button"
                            onClick={addOption}
                        >
                            + Add Option
                        </button>
                    )}

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
                        {loading ? "Creating..." : "Create Poll"}
                    </button>
                </form>
            </div>
        </div>
    );
}

export default CreatePoll;
