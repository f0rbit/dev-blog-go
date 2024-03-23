import { Dispatch, SetStateAction, useContext, useState } from "react";
import { API_URL, FunctionResponse, PostContext } from "../App";
import { ArrowDownNarrowWide, Edit, Filter, FolderTree, Plus, RefreshCw, Search, Tags, Trash, X } from "lucide-react";
import CategoryInput from "../components/CategoryInput";
import { Post, PostCreation, PostUpdate } from "../../schema";
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
    const { posts } = useContext(PostContext);
    const [sort, setSort] = useState<PostSort>("created");
    const [loading, _] = useState(false);
    const [filters, setFilters] = useState<PostFilters>({ category: null, tag: null });
    const [editingPost, setEditingPost] = useState<PostUpdate | null>(null);


    function newPost() {
        setEditingPost(EMPTY_POST_CREATION)
    }

    if (editingPost != null) return <PostEdit initial={editingPost} save={async () => ({ success: false, error: 'test' })} cancel={() => setEditingPost(null) } />

    if (!posts || !posts.posts) return <div>
        <button style={{ marginLeft: "auto" }} onClick={newPost}><Plus /><span>Create</span></button>
        <p>No Posts Found!</p>
    </div>;

    /*
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
    */

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
                <button style={{ marginLeft: "auto" }} onClick={newPost}><Plus /><span>Create</span></button>
            </div>
        </section>
        <section id='post-grid'>
            {filtered_posts.map((p: any) => <PostCard key={p.id} post={p} edit={() => setEditingPost(p) }/>)}
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

function PostCard({ post, edit }: { post: Post, edit: () => void }) {
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


    return <div className='post-card hidden-parent'>
        <div className='flex-row center top-right hidden-child'>
            {post.archived == true ? <button onClick={() => revivePost()}><RefreshCw /></button> : <button onClick={() => deletePost()}><Trash /></button>}
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

