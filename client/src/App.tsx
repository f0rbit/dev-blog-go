import { Dispatch, SetStateAction, useEffect, useState } from 'react'
import './App.css'
import React from 'react';
import { FolderTree, Home, LibraryBig, LogOut, Settings, Tags } from 'lucide-react';
import { PostsResponse, SCHEMA, CategoryResponse, ProjectsResponse } from "../schema";
import { HomePage } from './pages/Home';
import { PostsPage } from './pages/Posts';
import { CategoriesPage } from './pages/Categories';
import { TagsPage } from './pages/Tags';
import { SettingsPage } from './pages/Settings';
import { LoginPage } from './pages/Login';


export interface PostContext {
  posts: PostsResponse,
  setPosts: Dispatch<SetStateAction<PostsResponse>>,
  categories: CategoryResponse
  setCategories: Dispatch<SetStateAction<CategoryResponse>>,
  tags: string[],
  setTags: Dispatch<SetStateAction<string[]>>,
  projects: ProjectsResponse,
}
export interface AuthContext {
  user: any | null,
}

const PAGES = ["home", "posts", "categories", "tags", "settings"] as const;
type Page = (typeof PAGES)[keyof typeof PAGES]

export type FunctionResponse = { success: true, error: null } | { success: false, error: string }

export const API_URL: string = import.meta.env.VITE_API_URL;
const VERSION = "v0.5.0";
console.log("Version: " + VERSION);

export const PostContext = React.createContext<PostContext>({ posts: {} as PostsResponse, setPosts: () => { }, categories: {} as CategoryResponse, setCategories: () => { }, tags: [], setTags: () => { }, projects: []});
export const AuthContext = React.createContext<AuthContext>({ user: null });

function App() {
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  // fetch user
  useEffect(() => {
    (async () => {
      const response = await fetch(`${API_URL}/auth/user`, { credentials: "include" });
      if (response.ok) {
        setUser(await response.json());
        setLoading(false);
      } else {
        setUser(null);
        setLoading(false);
      }
    })();
  }, []);

  if (loading) {
    return <LoadingPage />
  }

  if (user == null) {
    return <LoginPage />
  }

  return <AuthContext.Provider value={{ user }}>
    <MainContent />
  </AuthContext.Provider>

}

function MainContent() {
  const [posts, setPosts] = useState({} as PostsResponse);
  const [categories, setCategories] = useState({} as CategoryResponse);
  const [page, setPage] = useState<Page>("home");
  const [tags, setTags] = useState<string[]>([]);
  const [projects, setProjects] = useState({} as ProjectsResponse);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      const response = await fetch(`${API_URL}/posts?limit=-1`, { credentials: "include" });
      if (!response.ok) throw new Error("Couldn't fetch posts");
      const result = await response.json();
      setPosts(SCHEMA.POSTS_RESPONSE.parse(result));

      const cat_res = await fetch(`${API_URL}/categories`, { credentials: "include" });
      if (!cat_res.ok) throw new Error("Couldn't fetch categories");
      setCategories(SCHEMA.CATEGORY_RESPONSE.parse(await cat_res.json()));

      const tag_res = await fetch(`${API_URL}/tags`, { credentials: "include" });
      const tag_result = await tag_res.json();
      setTags(tag_result);

      const project_res = await fetch(`${API_URL}/projects`, { credentials: "include" });
      if (!project_res.ok) throw new Error("Couldn't fetch projects");
      setProjects(SCHEMA.PROJECTS_RESPONSE.parse(await project_res.json()));

      setLoading(false);
    })();
  }, []);

  if (loading) return <LoadingPage />

  return (
    <PostContext.Provider value={{ posts, setPosts, categories, setCategories, tags, setTags, projects }}>
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

function LoadingPage() {
  return <section className='flex-col center'>
    <h4>Loading...</h4>
  </section>
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
    {<button onClick={() => window.location.href = `${API_URL}/auth/logout`}><LogOut /><span>Logout</span></button>}
  </>;
}

export default App
