import { Dispatch, SetStateAction, useContext, useEffect, useState } from 'react'
import './App.css'
import React from 'react';
import { FolderTree, Home, LibraryBig, Settings, Tags } from 'lucide-react';
import { Post, PostsResponse, SCHEMA, CategoryResponse } from "../schema";
import { HomePage } from './pages/Home';
import { PostsPage } from './pages/Posts';
import { CategoriesPage } from './pages/Categories';
import { TagsPage } from './pages/Tags';
import { SettingsPage } from './pages/Settings';
import { LoginPage } from './pages/Login';

export type PostCreation = Omit<Post, "id" | "created_at" | "updated_at">

export interface PostContext {
    posts: PostsResponse,
    setPosts: Dispatch<SetStateAction<PostsResponse>>,
    categories: CategoryResponse
    setCategories: Dispatch<SetStateAction<CategoryResponse>>,
    tags: string[],
    setTags: Dispatch<SetStateAction<string[]>>
}
export interface AuthContext {
    token: string | null,
}

const PAGES = ["home", "posts", "categories", "tags", "settings"] as const;
type Page = (typeof PAGES)[keyof typeof PAGES]

export type FunctionResponse = { success: true, error: null } | { success: false, error: string }

export const API_URL: string = import.meta.env.VITE_API_URL;
const VERSION = "v0.5.0";
console.log("Version: " + VERSION);

export const PostContext = React.createContext<PostContext>({ posts: {} as PostsResponse, setPosts: () => { }, categories: {} as CategoryResponse, setCategories: () => { }, tags: [], setTags: () => { } });
export const AuthContext = React.createContext<AuthContext>({ token: null });

function App() {
    const [token, setToken] = useState<string | null>(null);

    async function attemptLogin(input: string) {
        console.log("attempting login with input: " + input);
        setToken(input);
    }

    if (token == null) {
        return <LoginPage attemptLogin={attemptLogin}/>
    }
    return <AuthContext.Provider value={{token }}>
        <MainContent />
    </AuthContext.Provider>
    
}

function MainContent() {
    const [posts, setPosts] = useState({} as PostsResponse);
    const [categories, setCategories] = useState({} as CategoryResponse);
    const [page, setPage] = useState<Page>("home");
    const [tags, setTags] = useState<string[]>([]);

    useEffect(() => {
        (async () => {
            const response = await fetch(`${API_URL}/posts?limit=-1`);
            if (!response.ok) throw new Error("Couldn't fetch posts");
            const result = await response.json();
            setPosts(SCHEMA.POSTS_RESPONSE.parse(result));

            const cat_res = await fetch(`${API_URL}/categories`);
            if (!cat_res.ok) throw new Error("Couldn't fetch categories");
            setCategories(SCHEMA.CATEGORY_RESPONSE.parse(await cat_res.json()));

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
            return <><Home /><span>Home</span></>;
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

export default App
