import { Dispatch, SetStateAction, useContext, useState } from "react";
import { API_URL, AuthContext, FunctionResponse, PostContext } from "../App";
import { ArrowDownNarrowWide, Edit, Filter, FolderTree, Plus, RefreshCw, Search, Tags, Trash, X } from "lucide-react";
import CategoryInput from "../components/CategoryInput";
import { Post, PostUpdate } from "../../schema";
import { Oval } from "react-loader-spinner";
import TagInput from "../components/TagInput";
import { PostEdit } from "./PostEdit";

type PostSort = "created" | "edited" | "published" | "oldest";
const EMPTY_POST_CREATION: PostUpdate = {
    id: null,
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
    const [sort, setSort] = useState<PostSort>("created");
    const [loading, _] = useState(false);
    const [filters, setFilters] = useState<PostFilters>({ category: null, tag: null });
    const [editingPost, setEditingPost] = useState<PostUpdate | null>(null);
    const { posts, categories, setPosts } = useContext(PostContext)
    const { user } = useContext(AuthContext)

    if (editingPost != null) return <PostEdit initial={editingPost} save={save} cancel={() => setEditingPost(null)} />

    if (!posts || !posts.posts) return <div className="flex-col">
        <p>No Posts Found!</p>
        <button onClick={() => setEditingPost(EMPTY_POST_CREATION)}><Plus /><span>Create</span></button>
    </div>

    async function save(post: PostUpdate): Promise<FunctionResponse> {
        if (!post) throw new Error("Called save() without having an updated post");
        // input validation
        const { slug, category } = post;
        if (slug.includes(" ")) return { error: "Invalid Slug!", success: false }
        if (categories == null) return { error: "Not fetched categories", success: false };
        const cat_list = Object.values(categories.categories).map((cat: any) => cat.name);
        if (!(cat_list.includes(category))) return { error: "Invalid Category!", success: false }

        if (post.id == null || post.id < 0) post.id = -1;
        if (post.id < 0) post.author_id = user.user_id;

        const url = API_URL + (post.id < 0 ? "/post/new" : "/post/edit");
        const method = post.id < 0 ? "POST" : "PUT";
        const response = await fetch(url, { method, body: JSON.stringify(post), credentials: "include" });

        try {
            if (!post.id || post.id < 0) {
                const result = await response.json();
                setPosts({ ...posts, posts: [...posts.posts, { ...result } as Post] });
            } else {
                const update_post = post as (Omit<PostUpdate, "id"> & { id: number });
                if (!response || !response.ok) return { error: "Invalid Response", success: false };

                setPosts({
                    ...posts, posts: posts.posts.map((p) => {
                        if (p.id != post.id) return p;
                        return { ...p, ...update_post };
                    })
                });
            }
            return { error: null, success: true }
        } catch (err) {
            console.error(err);
            return { error: "Failed to decode response", success: false };
        }
    }

    // sort posts
    let sorted_posts = structuredClone(posts.posts);
    switch (sort) {
        case "edited":
            sorted_posts = sorted_posts.sort((a, b) => (new Date(b.updated_at).getTime()) - (new Date(a.updated_at).getTime()));
            break;
        case "created":
        case "oldest":
            sorted_posts = sorted_posts.sort((a, b) => (new Date(b.updated_at).getTime()) - (new Date(a.updated_at).getTime()));
            if (sort == "oldest") sorted_posts.reverse();
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


    return <main>
        <section id="post-filters">
            <label htmlFor='post-search'><Search /></label>
            <input type="text" id='post-search' />
            <SortControl selected={sort} setSelected={setSort} />
            <FilterControl filters={filters} setFilters={setFilters} />
            <div style={{ marginLeft: "auto" }} className="flex-row" >
                {loading && <Oval height={20} width={20} strokeWidth={8} />}
                <button style={{ marginLeft: "auto" }} onClick={() => setEditingPost(null)}><Plus /><span>Create</span></button>
            </div>
        </section>
        <section id='post-grid'>
            {filtered_posts.map((p: any) => <PostCard key={p.id} post={p} edit={() => setEditingPost(p)} save={save} />)}
        </section>
    </main>
}


export function TagEditor({ tags, setTags }: { tags: Post['tags'], setTags: (tags: Post['tags']) => void }) {
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

function PostCard({ post: initial, edit, save }: { post: Post, edit: () => void, save: (post: PostUpdate) => Promise<FunctionResponse> }) {
    const [post, setPost] = useState<PostUpdate>({ ...initial });

    const trash = async (): Promise<FunctionResponse> => {
        post.archived = true;
        setPost({ ...post });
        return await save(post);
    }

    const revive = async (): Promise<FunctionResponse> => {
        post.archived = false;
        setPost({ ...post });
        return await save(post);
    }

    return <div className='post-card hidden-parent'>
        <div className='flex-row center top-right hidden-child'>
            {post.archived == true ? <button onClick={revive}><RefreshCw /></button> : <button onClick={trash}><Trash /></button>}
            <button onClick={edit}><Edit /></button>
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

