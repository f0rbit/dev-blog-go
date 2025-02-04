import { useContext, useEffect, useState } from "react";
import { PostUpdate, ProjectsResponse, toIsoString } from "../../schema"
import { PostContext, FunctionResponse } from "../App"
import CategoryInput from "../components/CategoryInput";
import { TagEditor } from "./Posts";
import { Save, X } from "lucide-react";
import { remark } from "remark";
import remarkHtml from "remark-html";
import asciidoctor from "asciidoctor";

type Views = "metadata" | "content" | "preview";

export function PostEdit({ initial, save, cancel }: { initial: PostUpdate, save: (post: PostUpdate) => Promise<FunctionResponse>, cancel: () => void }) {
  const [post, setPost] = useState<PostUpdate>(initial);
  const [view, setView] = useState<Views>("metadata");
  const [manual_slug, setManualSlug] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const { categories, projects } = useContext(PostContext);

  function updateTitle(value: string) {
    const update_post = { ...post, title: value };
    if (!manual_slug) update_post['slug'] = value.replaceAll(" ", "-").toLowerCase();
    setPost(update_post);
  }

  function updateSlug(value: string) {
    if (manual_slug == false) setManualSlug(true);
    setPost({ ...post, slug: value });
  }

  function setPublishDate(value: any) {
    setPost({ ...post, publish_at: value });
  }
  function SaveContent() {
    return <><Save />{post.id == null ? "Create" : "Save"}</>
  }

  const edit_time = post.publish_at?.length > 1 ? toIsoString(new Date(post.publish_at)) : "";

  if (categories == null) {
    return <p>No categories found, please create some before posting.</p>
  }

  function Preview({ content, format }: { content: string, format: "md" | "adoc" }) {
    const [parsed, setParsed] = useState<string | null>(null);

    useEffect(() => {
      if (format == "md") {
        remark().use(remarkHtml).process(content).then((p: any) => setParsed(p));
      } else if (format == "adoc") {
        const result = asciidoctor().convert(content);
        setParsed(result.toString());
      } else {
        throw new Error(`unknown content format: ${format}`);
      }
    }, []);

    if (!parsed) return <p>Loading...</p>;
    return <article dangerouslySetInnerHTML={{ __html: parsed }} ></article>
  }

  console.log(post.project_id);

  return <main className="flex-col">
    <nav className="vertical-tabs">
      <button onClick={() => setView("metadata")} className={view == "metadata" ? 'selected' : ''}>Metadata</button>
      <button onClick={() => setView("content")} className={view == "content" ? "selected" : ""}>Content</button>
      <button onClick={() => setView("preview")} className={view == "preview" ? "selected" : ""}>Preview</button>
    </nav >
    {view == "metadata" && <div className="flex-col input-grid">
      <h3 style={{ gridColumn: "span 4" }}>{post.id ? `Editing Post ${post.title}` : "Creating New Post"}</h3>

      <label>Title</label><input type="text" value={post.title} onChange={(e) => updateTitle(e.target.value)} />
      <label>Slug</label><input type="text" value={post.slug} onChange={(e) => updateSlug(e.target.value)} />
      <label>Category</label><CategoryInput value={post.category} categories={categories.categories} setValue={(c) => setPost({ ...post, category: c })} />
      <label>Publish</label><input type="datetime-local" value={edit_time} onChange={(e) => setPublishDate(new Date(e.target.value).toISOString())} />
      <label>Format</label><select value={post.format} onChange={(e) => setPost({ ...post, format: e.target.value as PostUpdate['format'] })}><option value="md">Markdown</option><option value="adoc">ASCII Doc</option></select>
      <label>Project</label><ProjectSelect value={post.project_id ?? null} projects={projects} setValue={(p) => setPost({ ...post, project_id: p == "" ? null : p })} />
      <label>Description</label><input style={{ gridColumn: "span 4" }} value={post.description} onChange={(e) => setPost({ ...post, description: e.target.value })} placeholder={post.content.substring(0, 20)} />
      <label>Tags</label><TagEditor tags={post.tags} setTags={(tags) => setPost({ ...post, tags })} />

      {error && <p className="error-message">{error}</p>}
      <div className="flex-row center" style={{ gridColumn: "span 4" }}>
        <button onClick={() => save(post).then((res) => setError(res.error))}><SaveContent /></button><button onClick={cancel}><X />Cancel</button>
      </div>
    </div>}
    {view == "content" && <div className="flex-col input-grid">
      <textarea style={{ gridColumn: "span 4", fontFamily: "monospace", height: "50vh" }} value={post.content} onChange={(e) => setPost({ ...post, content: e.target.value })} />
    </div>}
    {view == "preview" && <Preview content={post.content} format={post.format} />}
  </main >
}


function ProjectSelect({ value, projects, setValue }: { value: string | null, projects: ProjectsResponse, setValue: (value: string | null) => void }) {

  return <>
    <select value={value ?? ""} onChange={(e) => setValue(e.target.value)}>
      <option value={""}>-</option>
      {projects.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}
    </select>
  </>
}



