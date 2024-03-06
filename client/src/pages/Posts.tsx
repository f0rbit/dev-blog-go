import { Dispatch, SetStateAction, useContext, useState } from "react";
import { API_URL, AuthContext, FunctionResponse, PostContext, PostCreation } from "../App";
import { ArrowDownNarrowWide, Edit, Filter, FolderTree, Plus, RefreshCw, Save, Search, Tags, Trash, X } from "lucide-react";
import Modal from "../components/Modal";
import CategoryInput from "../components/CategoryInput";
import { Post } from "../../schema";
import { Oval } from "react-loader-spinner";
import TagInput from "../components/TagInput";

type PostSort = "created" | "edited" | "published" | "oldest";
const EMPTY_POST_CREATION: PostCreation = {
    author_id: -1,
    slug: "",
    title: "",
    content: "",
    category: "",
    tags: [],
    publish_at: "",
    archived: false
}

export function PostsPage() {
    const { posts, setPosts, categories } = useContext(PostContext);
    const [selected, setSelected] = useState<PostSort>("created");
    const [openCreatePost, setOpenCreatePost] = useState(false);
    const [loading, setLoading] = useState(false);
    const [filters, setFilters] = useState<PostFilters>({ category: null, tag: null });
    const { user } = useContext(AuthContext);
    const [creatingPost, setCreatingPost] = useState<PostCreation>(EMPTY_POST_CREATION);

    if (!posts || !posts.posts) return <p>No Posts Found!</p>;

    async function createPost(): Promise<FunctionResponse> {
        // input validation
        const { slug, category } = creatingPost;
        if (slug.includes(" ")) return { error: "Invalid Slug!", success: false }
        if (categories == null) return { error: "Not fetched categories", success: false };
        const cat_list = Object.values(categories.categories).map((cat: any) => cat.name);
        if (!(cat_list.includes(category))) return { error: "Invalid Category!", success: false }
        const new_post: Post & { loading: boolean } = { ...creatingPost, id: -1, created_at: toIsoString(new Date()), updated_at: toIsoString(new Date()), loading: true };

        setPosts({ ...posts, posts: [...posts.posts, new_post] });
        // send request
        const response = await fetch(`${API_URL}/post/new`, { method: "POST", body: JSON.stringify({ ...creatingPost, author_id: user.user_id }), credentials: "include" });
        const result = await response.json();
        // update state?
        setLoading(false);
        setPosts({ ...posts, posts: [...posts.posts, { ...result, loading: false } as Post] });
        return { error: null, success: true }
    }

    // sort posts
    let sorted_posts = structuredClone(posts.posts);
    switch (selected) {
        case "edited":
            sorted_posts = sorted_posts.sort((a, b) => (new Date(b.updated_at).getTime()) - (new Date(a.updated_at).getTime()));
            break;
        case "created":
        case "oldest":
            sorted_posts = sorted_posts.sort((a, b) => (new Date(b.updated_at).getTime()) - (new Date(a.updated_at).getTime()));
            if (selected == "oldest") sorted_posts.reverse();
            break;
        case "published":
            break;
    }

    let filtered_posts = sorted_posts.filter((p) => {
        // tag filtering
        if (filters.tag) {
            const contains_tag = p.tags.includes(filters.tag);
            if (!contains_tag) return false;
        }
        // category filtering
        if (filters.category) {
            if (p.category != filters.category) return false;
        }
        return true;
    });

    function closeEditor() {
        setOpenCreatePost(false);
        setCreatingPost(EMPTY_POST_CREATION);
    }

    return <main>
        <section id="post-filters">
            <label htmlFor='post-search'><Search /></label>
            <input type="text" id='post-search' />
            <SortControl selected={selected} setSelected={setSelected} />
            <FilterControl filters={filters} setFilters={setFilters} />
            <div style={{ marginLeft: "auto" }} className="flex-row" >
                {loading && <Oval height={20} width={20} strokeWidth={8} />}
                <button style={{ marginLeft: "auto" }} onClick={() => setOpenCreatePost(true)}><Plus /><span>Create</span></button>
            </div>
        </section>
        <section id='post-grid'>
            {filtered_posts.map((p: any) => <PostCard key={p.id} post={p} />)}
        </section>
        <Modal openModal={openCreatePost} closeModal={closeEditor}>
            <PostEditor post={creatingPost} setPost={setCreatingPost} save={createPost} type={"create"} cancel={closeEditor} />
        </Modal>
    </main>
}


interface PostEditorProps {
    post: PostCreation,
    setPost: Dispatch<SetStateAction<PostCreation>>,
    save: () => Promise<FunctionResponse>,
    type: "create" | "edit",
    cancel: () => void
}

function PostEditor({ post, setPost, save, type, cancel }: PostEditorProps) {
    const [manualSlug, setManualSlug] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const { categories } = useContext(PostContext);

    function updateTitle(value: string) {
        const update_post = { ...post, title: value };
        if (!manualSlug) update_post['slug'] = value.replaceAll(" ", "-").toLowerCase();
        setPost(update_post);
    }

    function updateSlug(value: string) {
        if (manualSlug == false) setManualSlug(true);
        setPost({ ...post, slug: value });
    }

    function SaveContent() {
        switch (type) {
            case "create": return <><Save />Create</>
            case "edit": return <><Save />Save</>
        }
    }

    function setPublishDate(value: any) {
        setPost({ ...post, publish_at: value });
    }

    const edit_time = post.publish_at?.length > 1 ? toIsoString(new Date(post.publish_at)) : "";

    if (categories == null) {
        return <p>No Categories!</p>
    }

    return <div className="flex-col">
        <h3 style={{ textTransform: "capitalize" }}>{type} Post</h3>
        <div className="input-grid">
            <label>Title</label><input type="text" value={post.title} onChange={(e) => updateTitle(e.target.value)} />
            <label>Slug</label><input type="text" value={post.slug} onChange={(e) => updateSlug(e.target.value)} />
            <label>Category</label><CategoryInput value={post.category} categories={categories.categories} setValue={(c) => setPost({ ...post, category: c })} />
            <label>Publish</label><input type="datetime-local" value={edit_time} onChange={(e) => setPublishDate(new Date(e.target.value).toISOString())} />
            <label style={{ placeSelf: "stretch" }}>Content</label><textarea style={{ gridColumn: "span 3", fontFamily: "monospace" }} rows={10} value={post.content} onChange={(e) => setPost({ ...post, content: e.target.value })} />
            <label>Tags</label><TagEditor tags={post.tags} setTags={(tags) => setPost({ ...post, tags })} />
        </div>
        {error && <p className="error-message">{error}</p>}
        <div className="flex-row center">
            <button onClick={() => save().then((res) => setError(res.error))}><SaveContent /></button><button onClick={cancel}><X />Cancel</button>
        </div>
    </div>
}

function toIsoString(date: Date) {
    const pad = function(num: number) {
        return (num < 10 ? '0' : '') + num;
    };

    return date.getFullYear() +
        '-' + pad(date.getMonth() + 1) +
        '-' + pad(date.getDate()) +
        'T' + pad(date.getHours()) +
        ':' + pad(date.getMinutes()) +
        ':' + pad(date.getSeconds())
}


function TagEditor({ tags, setTags }: { tags: Post['tags'], setTags: (tags: Post['tags']) => void }) {
    const [input, setInput] = useState("");
    const { tags: all_tags } = useContext(PostContext);

    return <div className="flex-row tag-editor">
        <TagInput tags={all_tags} value={input} setValue={(t) => { setTags([...tags, t]); setInput(""); }} />
        <button onClick={() => { setTags([...tags, input]); setInput("") }}><Plus /></button>
        {tags.map((tag, index) => <Tag key={index} tag={tag} remove={() => setTags(tags.toSpliced(index, 1))} />)}
    </div>
}

function Tag({ tag, remove }: { tag: string, remove: (() => void) | null }) {
    return <div className="flex-row center tag">
        <span>{tag}</span>
        {remove && <button onClick={remove}><X /></button>}
    </div>;
}

function PostCard({ post }: { post: Post }) {
    const [editorOpen, setEditorOpen] = useState(false);
    const [editingPost, setEditingPost] = useState<PostCreation>({ ...post });

    const { posts, setPosts } = useContext(PostContext);

    const deletePost = async (): Promise<FunctionResponse> => {
        editingPost.archived = true;
        setEditingPost({ ...editingPost });
        return await savePost();
    }

    const revivePost = async (): Promise<FunctionResponse> => {
        editingPost.archived = false;
        setEditingPost({ ...editingPost });
        return await savePost();
    }

    const savePost = async (): Promise<FunctionResponse> => {
        const id = post.id;
        const upload_post: PostCreation & { id: number, loading: boolean } = {
            ...editingPost,
            id,
            loading: true,
        }
        setPosts({
            ...posts, posts: posts.posts.map((p) => {
                if (p.id != id) return p;
                return { ...p, ...upload_post, loading: true };
            })
        });
        // send upload post to server
        const response = await fetch(`${API_URL}/post/edit`, { method: "PUT", body: JSON.stringify(upload_post), credentials: "include" });
        if (!response || !response.ok) return { error: "Invalid Response", success: false };
        // update state
        setPosts({
            ...posts, posts: posts.posts.map((p) => {
                if (p.id != id) return p;
                return { ...p, ...upload_post, loading: false };
            })
        });
        return { error: null, success: true };
    }

    function close() {
        setEditorOpen(false);
        setEditingPost({ ...post });
    }

    return <div className='post-card hidden-parent'>
        <Modal openModal={editorOpen} closeModal={close}>
            <PostEditor post={editingPost} setPost={setEditingPost} save={savePost} type={"edit"} cancel={close} />
        </Modal>
        <div className='flex-row center top-right hidden-child'>
            {post.archived == true ? <button onClick={() => revivePost()}><RefreshCw /></button> : <button onClick={() => deletePost()}><Trash /></button>}
            <button onClick={() => setEditorOpen(true)}><Edit /></button>
        </div>
        <div className="flex-row center bottom-right">
            {(post as any).loading && <Oval width={16} height={16} strokeWidth={8} />}
        </div>
        <h2>{post.title}</h2>
        <pre>{post.slug}</pre>
        <p className="clamp-text-3">{post.content}</p>
    </div>
}

function SortControl({ selected, setSelected }: { selected: PostSort, setSelected: Dispatch<SetStateAction<PostSort>> }) {
    const [open, setOpen] = useState(false);

    const SortButton = ({ sort }: { sort: PostSort }) => {
        return <button className={(selected == sort ? "selected " : "") + " sort-button"} onClick={() => { setSelected(sort); setOpen(false) }}>{sort}</button>
    }

    return <>
        <button onClick={() => setOpen(!open)}><ArrowDownNarrowWide /><span>Sort</span></button>
        {open && <>
            <SortButton sort="created" />
            <SortButton sort="edited" />
            <SortButton sort="published" />
            <SortButton sort="oldest" />
        </>}
    </>
}

type PostFilters = {
    category: Post['category'] | null,
    tag: string | null
}

function FilterControl({ filters, setFilters }: { filters: PostFilters, setFilters: Dispatch<SetStateAction<PostFilters>> }) {
    const [open, setOpen] = useState(false);
    const { categories, tags } = useContext(PostContext)

    const applied = filters.category || filters.tag;

    return <>
        <button onClick={() => setOpen(!open)}><Filter /><span>Filter</span></button>
        {open && <>
            <FolderTree />
            {categories && <CategoryInput value={filters.category ?? ""} categories={categories.categories} key={0} setValue={(c) => setFilters({ ...filters, category: c })} />}
            <Tags />
            <TagInput value={filters.tag ?? ""} tags={tags} setValue={(t) => setFilters({ ...filters, tag: t })} />
        </>}
        {(applied || open) && <button title="Clear" onClick={() => { setFilters({ category: null, tag: null }); setOpen(false) }}><X /></button>}
    </>
}

