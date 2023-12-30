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
            const response = await fetch(`${API_URL}/auth/tokens`, { method: "GET", credentials: "include" });
            if (!response.ok) throw new Error("Couldn't fetch tokens");
            const result = await response.json();
            setTokens(result);
        })();
    }, []);

    if (tokens == null) return <section className="page-center"><Oval height={20} width={20} strokeWidth={8} /></section>

    return <div className="flex-col">
        <div className="flex-row">
            <Search />
            <input type="text" />
            <button style={{ marginLeft: "auto"}} onClick={() => setCreating(true)}><Plus /><span>Create</span></button>
        </div>
        <div className="token-grid">
            <div>Token</div>
            <div>Name</div>
            <div>Note</div>
            <div>Enabled</div>
            <div></div>
            {tokens.map((t) => <TokenRow token={t} />)}
            {creating && <TokenRow token={{ ...EMPTY_TOKEN, user_id: user.id }} />}
        </div>
    </div>;
}

type TokenCreation = Omit<AccessKey, "created_at" | "updated_at">

const EMPTY_TOKEN: Omit<TokenCreation, "user_id"> = {
    id: -1,
    value: "",
    name: "",
    note: "",
    enabled: true,
}

function TokenRow({ token }: { token: TokenCreation }) {
    const [editing, setEditing] = useState(token);
    const mode = token.id < 0 ? "create" : "edit";
    const [enabled, setEnabled] = useState(mode == "create");
    const icon = mode == "create" ? <Check /> : <Save />

    return <>
        <input type="text" value={editing.value} onChange={(e) => setEditing({ ...editing, value: e.target.value})} disabled={!enabled}/>
        <input type="text" value={editing.name} onChange={(e) => setEditing({ ...editing, name: e.target.value })} disabled={!enabled}/>
        <input type="text" value={editing.note} onChange={(e) => setEditing({ ...editing, note: e.target.value })} disabled={!enabled}/>
        <input type="checkbox" checked={editing.enabled} onChange={(e) => setEditing({...editing, enabled: e.target.checked})} disabled={!enabled}/>
        {enabled ? <button onClick={() => {setEnabled(false)}}>{icon}</button> : <button onClick={() => setEnabled(true)}><Edit /></button>}
    </>;
}
