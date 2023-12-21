import { Dispatch, SetStateAction, useContext, useEffect, useState } from 'react'
import './App.css'
import React from 'react';
import { ArrowDownNarrowWide, Filter, FolderTree, Home, LibraryBig, Plus, Save, Search, Settings, Tags, Trash } from 'lucide-react';
import Modal from "./components/Modal";
import CategoryInput from './components/CategoryInput';

type Post = {
    id: number,
    slug: string,
    title: string,
    content: string,
    category: string,
    created_at: string,
    updated_at: string
}

type PostResponse = {
    posts: Post[],
    total_pages: number,
    current_page: number,
    total_posts: number,
    per_page: number,
}

interface PostContext {
    posts: PostResponse,
    setPosts: Dispatch<SetStateAction<PostResponse>>,
    categories: any,
    setCategories: Dispatch<SetStateAction<any>>,
    tags: string[],
    setTags: Dispatch<SetStateAction<string[]>>
}

const PAGES = ["home", "posts", "categories", "tags", "settings"] as const;
type Page = (typeof PAGES)[keyof typeof PAGES]

const API_URL: string = import.meta.env.VITE_API_URL;
const VERSION = "v0.5";
console.log("Version: " + VERSION);

const PostContext = React.createContext<PostContext>({ posts: {} as PostResponse, setPosts: () => { }, categories: [], setCategories: () => { }, tags: [], setTags: () => { } });

function App() {
    const [posts, setPosts] = useState<PostResponse>({} as PostResponse);
    const [categories, setCategories] = useState([]);
    const [page, setPage] = useState<Page>("home");
    const [tags, setTags] = useState<string[]>([]);

    useEffect(() => {
        (async () => {
            const response = await fetch(`${API_URL}/posts?limit=-1`);
            const result = await response.json();
            setPosts(result);

            const cat_res = await fetch(`${API_URL}/categories`);
            const cat_result = await cat_res.json();
            setCategories(cat_result);

            const tag_res = await fetch(`${API_URL}/tags`);
            const tag_result = await tag_res.json();
            setTags(tag_result);
        })();
    }, []);

    return (
        <PostContext.Provider value={{ posts, setPosts, categories, setCategories, tags, setTags }}>
            <nav>
                <Sidebar page={page} setPage={setPage} />
            </nav>
            <section>
                <TitleBar page={page} />
                <Content page={page} />
            </section>
        </PostContext.Provider>
    )
}

function TitleBar({ page }: { page: Page }) {
    const title = () => {
        switch (page) {
            case "home": return "Home";
            case "categories": return "Categories";
            case "posts": return "Posts";
            case "tags": return "Tags";
            case "settings": return "Settings";
        }
    }
    const description = () => {
        switch (page) {
            case "home": return "Overview & analytics";
            case "categories": return "Manage available categories";
            case "posts": return "Create, edit, delete posts";
            case "tags": return "View tag information";
            case "settings": return "Modify settings";
        }
    }
    return <header>
        <h1>{title()}</h1>
        <p>{description()}</p>
    </header>

}

function Content({ page }: { page: Page }) {
    switch (page) {
        case "home": return <HomePage />
        case "tags": return <TagsPage />
        case "posts": return <PostsPage />
        case "categories": return <CategoriesPage />
        case "settings": return <SettingsPage />
    }
}

function LinkContent({ page }: { page: Page }) {
    switch (page) {
        case "home":
            return <>
                <Home />
                <span>Home</span>
            </>
        case "tags":
            return <><Tags /><span>Tags</span></>;
        case "posts":
            return <><LibraryBig /><span>Posts</span></>;
        case "categories":
            return <><FolderTree /><span>Categories</span></>;
        case "settings":
            return <><Settings /><span>Settings</span></>;
    }
}

function Sidebar({ page, setPage }: { page: Page, setPage: Dispatch<SetStateAction<Page>> }) {
    // render the icons for each page
    return <>
        {PAGES.map((p) => (<button key={p} data-page={p} className={page == p ? "selected" : ""} onClick={() => setPage(p)}><LinkContent page={p} /></button>))}
    </>;
}

function HomePage() {
    return <>
        <h1>Home Page</h1>
        <p>we so up</p>
    </>;
}

type PostSort = "created" | "edited" | "published" | "oldest";

function PostsPage() {
    const { posts, categories } = useContext(PostContext);
    const [selected, setSelected] = useState<PostSort>("created");
    const [openCreatePost, setOpenCreatePost] = useState(false);

    if (!posts || !posts.posts) return <p>No Posts Found!</p>;


    async function createPost({ title, slug, content, category }: { title: string, slug: string, content: string, category: string }): Promise<{ success: true, error: null } | { success: false, error: string }> {
        if (slug.includes(" ")) return { error: "Invalid Slug!", success: false } 
        const cat_list = Object.values(categories).map((cat: any) => cat.name);
        console.log({ category, cat_list, categories });
        if (!(cat_list.includes(category))) return { error: "Invalid Category!", success: false }
        const data = { title, slug, content, category };
        const response = await fetch("http://localhost:8080/post/new", { method: "POST", body: JSON.stringify(data) });
        const result = await response.json();
        posts.posts.push(result);
        return { error: null, success: true }
    }

    // sort posts
    let sorted_posts = structuredClone(posts.posts);
    switch (selected) {
        case "edited": 
            sorted_posts = sorted_posts.sort((a,b) => (new Date(b.updated_at).getTime()) - (new Date(a.updated_at).getTime()));
            break;
        case "created":
        case "oldest":
            sorted_posts = sorted_posts.sort((a,b) => (new Date(b.updated_at).getTime()) - (new Date(a.updated_at).getTime()));
            if (selected == "oldest") sorted_posts.reverse();
            break;
        case "published":
            break;
    }

    return <main>
        <section id="post-filters">
            <label htmlFor='post-search'><Search /></label>
            <input type="text" id='post-search' />
            <SortControl selected={selected} setSelected={setSelected} />
            <button><Filter /><span>Filter</span></button>
            <button style={{ marginLeft: "auto" }} onClick={() => setOpenCreatePost(true)}><Plus /><span>Create</span></button>
        </section>
        <section id='post-grid'>
            {sorted_posts.map((p: any) => <PostCard key={p.id} post={p} />)}
        </section>
        <Modal openModal={openCreatePost} closeModal={() => setOpenCreatePost(false)}>
            <CreatePost create={createPost} />
        </Modal>
    </main>
}

function CreatePost({ create }: { create: (data: { title: string, slug: string, content: string, category: string }) => Promise<{ success: true, error: null } | { success: false, error: string }> }) {
    const [title, setTitle] = useState<string>("");
    const [slug, setSlug] = useState<string>("");
    const [manualSlug, setManualSlug] = useState<boolean>(false);
    const [content, setContent] = useState<string>("");
    const [category, setCategory] = useState<string>("");
    const [error, setError] = useState<string|null>(null);

    const { categories } = useContext(PostContext);

    function updateTitle(value: string) {
        if (!manualSlug) {
            setSlug(value.replaceAll(" ", "-").toLowerCase());
        }
        setTitle(value);
    }

    function updateSlug(value: string) {
        if (manualSlug == false) setManualSlug(true);
        setSlug(value);
    }


    return <div className="flex-col">
        <h3>Create Post</h3>
        <div className="input-grid">
            <label>Title</label><input type="text" value={title} onChange={(e) => updateTitle(e.target.value)} />
            <label>Slug</label><input type="text" value={slug} onChange={(e) => updateSlug(e.target.value)}/>
            <label>Category</label><CategoryInput categories={categories} setValue={(c) => setCategory(c)} />
            <label>Publish</label><input type="date" />
            <label style={{ placeSelf: "stretch" }}>Content</label><textarea style={{ gridColumn: "span 3", fontFamily: "monospace" }} rows={10} value={content} onChange={(e) => setContent(e.target.value)} />
        </div>
        {error && <p className="error-message">{error}</p>} 
        <div className="flex-row center">
            <button onClick={() => create({ title, slug, content, category }).then((res) => setError(res.error))}><Save />Save</button><button><Trash />Cancel</button>
        </div>
    </div>
}

function PostCard({ post }: { post: Post }) {

    return <div className='post-card'>
        <h2>{post.title}</h2>
        <pre>{post.slug}</pre>
        <p>{post.content}</p>
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

function CategoriesPage() {
    const { categories } = useContext(PostContext);

    // construct a graph of categories
    // let's start at the root node
    const graph = getChildrenCategories(categories, 'root');

    const elements = getCategoryElements(graph, 0);

    return <main id='category-list'>
        {elements}
    </main>
}

function getChildrenCategories(categories: { name: string, parent: string }[], root: string): any {
    const graph: any = {};
    for (const cat of categories) {
        if (cat.parent == root) graph[cat.name] = getChildrenCategories(categories, cat.name);
    }
    return graph;
}
function getCategoryElements(values: any, depth: number) {
    const list: JSX.Element[] = [];

    const CategoryCard = ({ cat, depth }: { cat: string, depth: number }) => {
        return <div style={{ marginLeft: (depth * 40) + "px" }} className="category-card">{cat}</div>
    }

    for (const [key, value] of Object.entries(values)) {
        list.push(<CategoryCard cat={key} depth={depth} />);
        list.push(...getCategoryElements(value, depth + 1));
    }
    return list;
}

function TagsPage() {
    const { tags } = useContext(PostContext);
    return <main>
        <pre>{JSON.stringify(tags, null, 2)}</pre>
    </main>
}
function SettingsPage() {
    return <>
        <h1>Settings</h1>
        <p>idk what kinda settings we want</p>
    </>
}

export default App
