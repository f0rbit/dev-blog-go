import { Palette, Plus, Search, Shield, User } from "lucide-react";
import { useEffect, useState } from "react";
import { API_URL } from "../App";

const subpages = ["profile", "theme", "tokens"] as const;

type SubPage = (typeof subpages)[keyof typeof subpages];

function SubPageIcon({ page }: { page: SubPage }) {
    switch (page) {
        case "theme": return <><Palette /><span>Theme</span></>;
        case "profile": return <><User /><span>Profile</span></>;
        case "tokens": return <><Shield /><span>Tokens</span></>;
    }
}

function SubPageContent({ page }: { page: SubPage }) {
    switch (page) {
        case "theme": return <Theme />;
        case "profile": return <Profile />;
        case "tokens": return <Tokens />
    }
}

export function SettingsPage() {
    const [page, setPage] = useState<SubPage>(subpages[0]);
    return <section className="flex-row settings-container">
        <nav className="flex-col">
            {subpages.map((p) => <button onClick={() => setPage(p)}><SubPageIcon page={p} /></button>)}
        </nav>
        <main>
            <SubPageContent page={page} />
        </main>
    </section>
}

function Theme() {
    return <div>Theme</div>;
}

function Profile() {
    return <div>Profile</div>
}

function Tokens() {
    const [tokens, setTokens] = useState(null);

    useEffect(() => {
        (async () => {
            const response = await fetch(`${API_URL}/auth/tokens`, { method: "GET", credentials: "include" });
            if (!response.ok) throw new Error("Couldn't fetch tokens");
            const result = await response.json();
            setTokens(result);
        })();
    }, []);
    return <div className="flex-col">
        <div className="flex-row">
            <Search />
            <input type="text" />
            <button style={{ marginLeft: "auto"}}><Plus /><span>Create</span></button>
        </div>
        <div>
            <pre>{JSON.stringify(tokens, null, 2)}</pre>
        </div>
    </div>;
}
