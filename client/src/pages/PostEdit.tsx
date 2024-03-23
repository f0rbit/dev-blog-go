import { useContext, useState } from "react";
import { PostUpdate, toIsoString } from "../../schema"
import { PostContext, FunctionResponse } from "../App"
import CategoryInput from "../components/CategoryInput";
import { TagEditor } from "./Posts";
import { Save, X } from "lucide-react";


export function PostEdit({ initial, save, cancel }: { initial: PostUpdate, save: (post: PostUpdate) => Promise<FunctionResponse>, cancel: () => void }) {
    const [post, setPost] = useState<PostUpdate>(initial);
    const [manual_slug, setManualSlug] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const { categories } = useContext(PostContext)

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

    return <main className="flex-col input-grid">
        <h3 style={{ gridColumn: "span 4" }}>{post.id ? `Editing Post ${post.title}` : "Creating New Post"}</h3>

        <label>Title</label><input type="text" value={post.title} onChange={(e) => updateTitle(e.target.value)} />
        <label>Slug</label><input type="text" value={post.slug} onChange={(e) => updateSlug(e.target.value)} />
        <label>Category</label><CategoryInput value={post.category} categories={categories.categories} setValue={(c) => setPost({ ...post, category: c })} />
        <label>Publish</label><input type="datetime-local" value={edit_time} onChange={(e) => setPublishDate(new Date(e.target.value).toISOString())} />

        <label id="content-label">Content</label>
        <textarea style={{ gridColumn: "span 4", fontFamily: "monospace", height: "50vh" }} value={post.content} onChange={(e) => setPost({ ...post, content: e.target.value })} />
        <label>Tags</label><TagEditor tags={post.tags} setTags={(tags) => setPost({ ...post, tags })} />

        {error && <p className="error-message">{error}</p>}
        <div className="flex-row center" style={{ gridColumn: "span 4"}}>
            <button onClick={() => save(post).then((res) => setError(res.error))}><SaveContent /></button><button onClick={cancel}><X />Cancel</button>
        </div>

    </main>
}

