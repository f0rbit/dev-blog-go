import { Dispatch, SetStateAction, useContext, useEffect, useState } from 'react'
import './App.css'
import React from 'react';
import { ArrowDownNarrowWide, Filter, FolderTree, Home, LibraryBig, Plus, Search, Settings, Tags } from 'lucide-react';

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
    setCategories: Dispatch<SetStateAction<any>>
}

const PAGES = ["home", "posts", "categories", "tags", "settings"] as const;
type Page = (typeof PAGES)[keyof typeof PAGES]

const PostContext = React.createContext<PostContext>({ posts: {} as PostResponse, setPosts: () => { }, categories: [], setCategories: () => { } });

function App() {
    const [posts, setPosts] = useState<PostResponse>({} as PostResponse);
    const [categories, setCategories] = useState([]);
    const [page, setPage] = useState<Page>("home");

    useEffect(() => {
        (async () => {
            const response = await fetch("http://localhost:8080/posts");
            const result = await response.json();
            setPosts(result);

            const cat_res = await fetch("http://localhost:8080/categories");
            const cat_result = await cat_res.json();
            setCategories(cat_result);
        })();
    }, []);

    return (
        <PostContext.Provider value={{ posts, setPosts, categories, setCategories }}>
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
    const { posts } = useContext(PostContext);
    const [selected, setSelected] = useState<PostSort>("created");

    if (!posts || !posts.posts) return <p>No Posts Found!</p>;

    return <main>
        <section id="post-filters">
            <label htmlFor='post-search'><Search /></label>
            <input type="text" id='post-search' />
            <SortControl selected={selected} setSelected={setSelected} />
            <button><Filter /><span>Filter</span></button>
            <button style={{ marginLeft: "auto" }}><Plus /><span>Create</span></button>
        </section>
        <section id='post-grid'>
            {posts.posts.map((p: any) => <PostCard key={p.id} post={p} />)}
        </section>
    </main>
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
    return <>
        <h1>Categories Page</h1>
        {categories.map((c: any) => (<pre key={c.name}>{JSON.stringify(c, null, 2)}</pre>))}
    </>
}

function TagsPage() {
    return <>
        <h1>Tags Page</h1>
        <p>list of tags</p>
    </>
}
function SettingsPage() {
    return <>
        <h1>Settings</h1>
        <p>idk what kinda settings we want</p>
    </>
}

export default App
