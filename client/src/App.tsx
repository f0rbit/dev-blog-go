import { Dispatch, SetStateAction, useContext, useEffect, useState } from 'react'
import './App.css'
import React from 'react';
import { FolderTree, Home, LibraryBig, Settings, Tags } from 'lucide-react';

interface PostContext {
    posts: any,
    setPosts: Dispatch<SetStateAction<any>>,
    categories: any,
    setCategories: Dispatch<SetStateAction<any>>
}

const PAGES = ["home", "posts", "categories", "tags", "settings"] as const;
type Page = (typeof PAGES)[keyof typeof PAGES]

const PostContext = React.createContext<PostContext>({ posts: [], setPosts: () => { }, categories: [], setCategories: () => { } });

function App() {
    const [posts, setPosts] = useState({});
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
            <main>
                <TitleBar page={page} />
                <Content page={page} />
            </main>
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
    return <div id="title-bar">
        <h1>{title()}</h1>
        <p>{description()}</p>
    </div>

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

function PostsPage() {
    const { posts } = useContext(PostContext);
    console.log(posts);
    if (!posts || !posts.posts) return <p>No Posts Found!</p>
    return <>
        <h1>Posts</h1>
        <p>{posts.posts.map((p: any) => (<pre key={p.id}>{JSON.stringify(p, null, 2)}</pre>))}</p>
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
