import { Check, Edit, Palette, Plus, Save, Search, Shield, User } from "lucide-react";
import { useContext, useEffect, useState } from "react";
import { API_URL, AuthContext } from "../App";
import { AccessKey } from "../../schema";
import { Oval } from "react-loader-spinner";
import { BuildingPage } from "../components/Building";

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
    return <BuildingPage />
}

function Profile() {
    return <BuildingPage />
}

function Tokens() {
    const [tokens, setTokens] = useState<AccessKey[] | null>(null);
    const [creating, setCreating] = useState(false);
    const { user } = useContext(AuthContext);

    useEffect(() => {
        (async () => {
            const response = await fetch(`${API_URL}/tokens`, { method: "GET", credentials: "include" });
            if (!response.ok) throw new Error("Couldn't fetch tokens");
            const result = await response.json();
            setTokens(result);
        })();
    }, []);

    async function save(token: TokenCreation) {
        const mode = token.id < 0 ? "new" : "edit";
        const response = await fetch(`${API_URL}/token/${mode}`, { method: mode == "new" ? "POST" : "PUT", credentials: "include", body: JSON.stringify(token) });
        if (!response.ok) return;
        const result = await response.json() as AccessKey;

        if (tokens == null) {
            setTokens([ result ]);
            return;
        }
        if (mode == "new") {
            setTokens([ ...tokens, result ]);
            setCreating(false);
        } else {
            setTokens(tokens.map((t) => {
                if (t.id != result.id) return t;
                return result;
            }));
        }
    }

    if (tokens == null) return <section className="page-center"><Oval height={20} width={20} strokeWidth={8} /></section>

    return <div className="flex-col">
        <div className="flex-row">
            <Search />
            <input type="text" />
            <button style={{ marginLeft: "auto" }} onClick={() => setCreating(true)}><Plus /><span>Create</span></button>
        </div>
        <div className="token-grid">
            <div>Token</div>
            <div>Name</div>
            <div>Note</div>
            <div>Enabled</div>
            <div></div>
            {tokens.map((t) => <TokenRow token={t} save={save} />)}
            {creating && <TokenRow token={{ ...EMPTY_TOKEN, user_id: user.user_id }} save={save} />}
        </div>
    </div>;
}

type TokenCreation = Omit<AccessKey, "created_at" | "updated_at">

const EMPTY_TOKEN: Omit<TokenCreation, "user_id"> = {
    id: -1,
    value: "<generated>",
    name: "",
    note: "",
    enabled: true,
}

function TokenRow({ token, save }: { token: TokenCreation, save: (token: TokenCreation) => void }) {
    const [editing, setEditing] = useState(token);
    const mode = token.id < 0 ? "create" : "edit";
    const [enabled, setEnabled] = useState(mode == "create");
    const [saving, setSaving] = useState(false);
    const icon = saving ? <Oval width={18} height={18} strokeWidth={8} /> : (mode == "create" ? <Check /> : <Save />);

    return <>
        <div className="api-token" style={{ fontFamily: "monospace" }}>{editing.value}</div>
        <input type="text" value={editing.name} onChange={(e) => setEditing({ ...editing, name: e.target.value })} disabled={!enabled} />
        <input type="text" value={editing.note} onChange={(e) => setEditing({ ...editing, note: e.target.value })} disabled={!enabled} />
        <input type="checkbox" checked={editing.enabled} onChange={(e) => setEditing({ ...editing, enabled: e.target.checked })} disabled={!enabled} />
        {saving ? <button>{icon}</button> : enabled ? <button onClick={() => { setEnabled(false); setSaving(true); save(editing) }}>{icon}</button> : <button onClick={() => setEnabled(true)}><Edit /></button>}
    </>;
}
