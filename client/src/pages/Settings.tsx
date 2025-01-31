import { Check, Edit, Link, Palette, Plus, RefreshCw, Save, Search, Shield, Trash, Unlink, User, X } from "lucide-react";
import { useContext, useEffect, useState } from "react";
import { API_URL, AuthContext, PostContext } from "../App";
import { AccessKey } from "../../schema";
import { Oval } from "react-loader-spinner";
import { BuildingPage } from "../components/Building";
import React from "react";
import { type IntegrationLink } from "../../schema";

const subpages = ["profile", "theme", "integrations", "tokens"] as const;

type SubPage = (typeof subpages)[keyof typeof subpages];

function SubPageIcon({ page }: { page: SubPage }) {
  switch (page) {
    case "theme": return <><Palette /><span>Theme</span></>;
    case "profile": return <><User /><span>Profile</span></>;
    case "tokens": return <><Shield /><span>Tokens</span></>;
    case "integrations": return <><Link /><span>Integrations</span></>;
  }
}

function SubPageContent({ page }: { page: SubPage }) {
  switch (page) {
    case "theme": return <Theme />;
    case "profile": return <Profile />;
    case "tokens": return <Tokens />
    case "integrations": return <Integrations />;
  }
}

export function SettingsPage() {
  const [page, setPage] = useState<SubPage>(subpages[0]);
  return <section className="settings-container">
    <nav className="flex-col">
      {subpages.map((p) => <button key={p} onClick={() => setPage(p)}><SubPageIcon page={p} /></button>)}
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

const INTEGRATION_STATES = {
  devto: { name: "devto", enabled: true },
  medium: { name: "medium", enabled: false },
  devpad: { name: "devpad", enabled: true }
} as const;

const DEFAULT_INTEGRATIONS_CONTEXT = {
  links: [] as IntegrationLink[],
  refetch: (() => { }),
  devpad_key: "",
  last_cache: {} as any
}

type LinksResponse = {
  integrations: IntegrationLink[];
  devpad_key: string;
  last_cache: any; // TODO: zod validate the links response
}

const IntegrationsContext = React.createContext<typeof DEFAULT_INTEGRATIONS_CONTEXT>(DEFAULT_INTEGRATIONS_CONTEXT);

function Integrations() {
  const [links, setLinks] = useState<IntegrationLink[]>([]);
  const [devpad_key, setDevpadKey] = useState<string>("");
  const [last_cache, setLastCache] = useState<any>({});

  const refetch = async () => {
    const response = await fetch(`${API_URL}/links`, { method: "GET", credentials: "include" });
    if (!response.ok) throw new Error("Couldn't fetch integrations");
    const result = ((await response.json()) ?? {}) as LinksResponse;
    setLinks(result.integrations ?? []);
    setDevpadKey(result.devpad_key ?? "");
    setLastCache(result.last_cache ?? {});
  }

  useEffect(() => {
    refetch();
  }, []);

  return <IntegrationsContext.Provider value={{ links, refetch, devpad_key, last_cache }}>
    <div id="integration-container" className="flex-col">
      <div id="integration-grid">
        {Object.entries(INTEGRATION_STATES).map(([key, data]) => (<IntegrationCard key={key} name={data.name} enabled={data.enabled} link={links.find((l) => l.source == data.name)} refetch={refetch} />))}
      </div>
      <div className="divider" />
      <div style={{ height: "100%" }}>
        <BuildingPage />
      </div>
    </div>
  </IntegrationsContext.Provider>
}

function IntegrationCard({ name, enabled, link, refetch }: { name: string, enabled: boolean, link: IntegrationLink | undefined, refetch: () => Promise<void> }) {
  const Header = () => <div className="flex-row">
    <span className={"status-indicator " + (enabled ? "ok" : "bad")}></span>
    <h2>{name}</h2>
  </div>;

  const unlink = async () => {
    if (!link) return;
    // make delete request to server
    const response = await fetch(`${API_URL}/links/delete/${link.id}`, { method: "DELETE", credentials: "include" });
    if (!response.ok) return;
    await refetch();
  }

  const sync = async () => {
    // make get request to "/links/fetch/{source}"
    const response = await fetch(`${API_URL}/links/fetch/${name}`, { method: "GET", credentials: "include" });
    if (!response || !response.ok) return;
    await refetch();

  }

  const Content = () => {
    if (!enabled) return <BuildingPage />;
    if (name == "devpad") return <DevpadCard refetch={refetch} />
    if (!link) return <LinkingInterface name={name as "devto" | "medium" | "substack"} />;

    const last_fetch = link.last_fetch ? new Date(link.last_fetch).toLocaleString() : "Never";
    const fetch_failed = link.data && JSON.parse(link.data)?.status === "failed";

    return <div className="flex-col" style={{ height: "100%" }}>
      <div style={{ height: "100%" }} className="flex-col">
        <div className="flex-row">
          <span>Last Fetch:</span>
          <span>{last_fetch}</span>
        </div>
        <div className="flex-row">
          <span>URL:</span>
          <span>{link.location}</span>
        </div>
        <div className="flex-row">
          <span>Posts:</span>
          {fetch_failed ? (
            <span style={{ color: "#ff7676", fontStyle: "italic" }}>Fetch failed</span>
          ) : (
            <span>{link.fetch_links?.length ?? "0"}</span>
          )}
        </div>
      </div>
      <div className="flex-row center">
        <button onClick={sync}><RefreshCw />Fetch</button>
        <button onClick={unlink}><Unlink /><span>Unlink</span></button>
      </div>
    </div>
  }


  return <div className="integration-card">
    <Header />
    <Content />
  </div>
}

function LinkingInterface({ name }: { name: "devto" | "medium" | "substack" }) {
  const [open, setOpen] = useState(false);
  const [pending, setPending] = useState(false);
  const { refetch } = useContext(IntegrationsContext);

  const devto = () => {
    const [token, setToken] = useState("");

    async function upload() {
      setPending(true);
      // push token to server
      // on success, update the context with body from response
      const input = { data: JSON.stringify({ token }), source: name };
      const response = await fetch(`${API_URL}/links/upsert`, { method: "PUT", body: JSON.stringify(input), credentials: "include" });
      if (!response.ok) {
        setPending(false);
        return false;
      }
      refetch();
      return true;
    }


    return open ?
      <div className="flex-col center" style={{ justifyContent: "space-between", height: "100%" }}>
        <input type="text" placeholder="API Key" value={token} onChange={(e) => setToken(e.target.value)} style={{ width: "100%" }} />
        <div className="flex-row center">
          {pending ? <button disabled><Oval width={18} height={18} strokeWidth={8} />Confirm</button> : <button onClick={upload}><Check />Confirm</button>}
          <button onClick={() => setOpen(false)}><X />Cancel</button>
        </div>
      </div>
      :
      <div className="flex-row center" style={{ height: "100%" }}>
        <button onClick={() => setOpen(true)}>
          <Link />
          <span>Link</span>
        </button>
      </div>;
  }



  switch (name) {
    case "devto": return devto();
    case "medium": return <BuildingPage />;
    case "substack": return <BuildingPage />;
  }

  return <p>Internal Error! Couldn't find integration for {name}</p>;
}

function Tokens() {
  const [tokens, setTokens] = useState<TokenCreation[] | null>(null);
  const [creating, setCreating] = useState(false);
  const { user } = useContext(AuthContext);

  useEffect(() => {
    (async () => {
      const response = await fetch(`${API_URL}/tokens`, { method: "GET", credentials: "include" });
      if (!response.ok) throw new Error("Couldn't fetch tokens");
      const result = (await response.json()) as AccessKey[];
      setTokens(result.map((r) => ({ ...r, saving: false })));
    })();
  }, []);


  async function save(token: TokenCreation) {
    const mode = token.id < 0 ? "new" : "edit";
    const response = await fetch(`${API_URL}/token/${mode}`, { method: mode == "new" ? "POST" : "PUT", credentials: "include", body: JSON.stringify(token) });
    if (!response.ok) return;

    if (tokens == null) {
      const new_token = { ...(await response.json()), saving: false };
      setTokens([new_token]);
      return;
    }
    if (mode == "new") {
      const new_token = { ...(await response.json()), saving: false };
      setTokens([...tokens, new_token]);
      setCreating(false);
    } else {
      setTokens(tokens.map((t) => {
        if (t.id != token.id) return t;
        return { ...token, saving: false };
      }));
    }
  }

  async function remove(token: TokenCreation) {
    const response = await fetch(`${API_URL}/token/delete/${token.id}`, { method: "DELETE", credentials: "include" });
    if (!response.ok) return;
    if (tokens == null) return;
    setTokens(tokens.filter((t) => t.id != token.id));
  }

  if (tokens == null) return <section className="page-center"><Oval height={20} width={20} strokeWidth={8} /></section>

  return <div className="flex-col" style={{ padding: "10px" }}>
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
      <div style={{ gridColumn: "span 2" }}></div>
      {tokens.map((t) => <TokenRow key={t.id} token={t} save={save} remove={remove} />)}
      {creating && <TokenRow token={{ ...EMPTY_TOKEN, user_id: user.user_id }} save={save} remove={remove} />}
    </div>
  </div>;
}

type TokenCreation = Omit<AccessKey, "created_at" | "updated_at"> & { saving: boolean };

const EMPTY_TOKEN: Omit<TokenCreation, "user_id"> = {
  id: -1,
  value: "<generated>",
  name: "",
  note: "",
  enabled: true,
  saving: false,
}

function TokenRow({ token, save, remove }: { token: TokenCreation, save: (token: TokenCreation) => void, remove: (token: TokenCreation) => void }) {
  const [editing, setEditing] = useState(token);
  const mode = token.id < 0 ? "create" : "edit";
  const [enabled, setEnabled] = useState(mode == "create");
  const icon = token.saving ? <Oval width={18} height={18} strokeWidth={8} /> : (mode == "create" ? <Check /> : <Save />);

  return <>
    <div className="api-token" style={{ fontFamily: "monospace" }}>{editing.value}</div>
    <input type="text" value={editing.name} onChange={(e) => setEditing({ ...editing, name: e.target.value })} disabled={!enabled} />
    <input type="text" value={editing.note} onChange={(e) => setEditing({ ...editing, note: e.target.value })} disabled={!enabled} />
    <input type="checkbox" checked={editing.enabled} onChange={(e) => setEditing({ ...editing, enabled: e.target.checked })} disabled={!enabled} />
    {token.saving ? <button style={{ gridColumn: "span 2" }}>{icon}</button> : enabled ? <button onClick={() => { setEnabled(false); token.saving = true; save(editing) }} style={{ gridColumn: "span 2" }}>{icon}</button> : <><button onClick={() => setEnabled(true)}><Edit /></button><button onClick={() => { setEnabled(false); token.saving = true; remove(editing) }}><Trash /></button></>}
  </>;
}

function DevpadCard({ refetch }: { refetch: () => Promise<void> }) {
  const { refetchProjects } = useContext(PostContext);
  const { devpad_key, last_cache } = useContext(IntegrationsContext);
  const [open, setOpen] = useState(false);

  const [loading, setLoading] = useState(false);

  const saveApiKey = async (key = devpad_key) => {
    setLoading(true);
    const response = await fetch(`${API_URL}/project/key`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ api_key: key }),
    });
    if (!response.ok) {
      setLoading(false);
      return;
    }
    await refetch();
    await refetchProjects();
    setLoading(false);
  };

  const sync = async () => {
    // TODO: make a request to force refetch
    await refetch();
  }

  const unlink = async () => {
    await saveApiKey("");
  }

  const Linker = () => {
    const [token, setToken] = useState("");

    async function upload() {
      await saveApiKey(token);
    }


    return open ?
      <div className="flex-col center" style={{ justifyContent: "space-between", height: "100%" }}>
        <input type="text" placeholder="API Key" value={token} onChange={(e) => setToken(e.target.value)} style={{ width: "100%" }} />
        <div className="flex-row center">
          {loading ? <button disabled><Oval width={18} height={18} strokeWidth={8} />Confirm</button> : <button onClick={upload}><Check />Confirm</button>}
          <button onClick={() => setOpen(false)}><X />Cancel</button>
        </div>
      </div>
      :
      <div className="flex-row center" style={{ height: "100%" }}>
        <button onClick={() => setOpen(true)}>
          <Link />
          <span>Link</span>
        </button>
      </div>;
  }

  const projects = last_cache.data ? JSON.parse(last_cache.data ?? "[]") : [];
  const last_fetch = last_cache.fetched_at ? new Date(last_cache.fetched_at).toLocaleString() : "Never";
  const fetch_failed = last_cache.status === "failed";

  if (!devpad_key) return <Linker />;

  return (
    <div className="flex-col" style={{ height: "100%" }}>
      <div style={{ height: "100%" }} className="flex-col">
        <div className="flex-row">
          <span>Last Fetch:</span>
          <span>{last_fetch}</span>
        </div>
        <div className="flex-row">
          <span>URL:</span>
          <span>{last_cache.url}</span>
        </div>
        <div className="flex-row">
          <span>Projects:</span>
          {fetch_failed ? (
            <span style={{ color: "#ff7676", fontStyle: "italic" }}>Fetch failed</span>
          ) : (
            <span>{projects.length}</span>
          )}
        </div>
      </div>
      <div className="flex-row center">
        <button onClick={sync}><RefreshCw />Fetch</button>
        <button onClick={unlink}><Unlink /><span>Unlink</span></button>
      </div>
    </div>
  );
}
