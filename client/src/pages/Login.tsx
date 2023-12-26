import { useState } from "react";

export function LoginPage({ attemptLogin }: { attemptLogin: (i: string) => void }) {
    const [input, setInput] = useState("");

    return <section className="flex-col center">
        <label>Password</label>
        <input type="text" value={input} onChange={(e) => setInput(e.target.value)} />
        <button type="button" onClick={() => attemptLogin(input)}>Login</button>
    </section>
}
