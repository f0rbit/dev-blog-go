import { useState } from "react";

export function LoginPage({ attemptLogin, error }: { attemptLogin: (i: string) => void, error: string }) {
    const [input, setInput] = useState("");

    return <section className="flex-col center">
        <label>Password</label>
        <input type="text" value={input} onChange={(e) => setInput(e.target.value)} />
        <button type="button" onClick={() => attemptLogin(input)}>Login</button>
        <p className="error-message" style={{ height: "16px"}}>{error}</p>
    </section>
}
