import { expect, test, describe, afterAll } from "bun:test";
import type { Post, PostsResponse } from "@client/schema";
import { AUTH_HEADERS } from "user";

const test_post = {
    author_id: 1,
    slug: "bun-test-post",
    title: "Bun Test Post",
    content: "this post should be removed after tests completed.",
    category: "root"
}
let test_post_id: number | null = null;

const pagination_posts = [
    { slug: "pagepost-1", title: "Page Post 1", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-2", title: "Page Post 2", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-3", title: "Page Post 3", content: "Lorem Ipsum", category: "learning" },
    { slug: "pagepost-4", title: "Page Post 4", content: "Lorem Ipsum", category: "devlog" },
    { slug: "pagepost-5", title: "Page Post 5", content: "Lorem Ipsum", category: "devlog" },
    { slug: "pagepost-6", title: "Page Post 6", content: "Lorem Ipsum", category: "gamedev" },
    { slug: "pagepost-7", title: "Page Post 7", content: "Lorem Ipsum", category: "gamedev" },
    { slug: "pagepost-8", title: "Page Post 8", content: "Lorem Ipsum", category: "hobbies" },
    { slug: "pagepost-9", title: "Page Post 9", content: "Lorem Ipsum", category: "hobbies" },
    { slug: "pagepost-10", title: "Page Post 10", content: "Lorem Ipsum", category: "story" },
    { slug: "pagepost-11", title: "Page Post 11", content: "Lorem Ipsum", category: "learning" },
    { slug: "pagepost-12", title: "Page Post 12", content: "Lorem Ipsum", category: "gamedev" },
    { slug: "pagepost-13", title: "Page Post 13", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-14", title: "Page Post 14", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-15", title: "Page Post 15", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-16", title: "Page Post 16", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-17", title: "Page Post 17", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-18", title: "Page Post 18", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-19", title: "Page Post 19", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-20", title: "Page Post 20", content: "Lorem Ipsum", category: "coding" },
    { slug: "pagepost-21", title: "Page Post 21", content: "Lorem Ipsum", category: "coding" },
];
const pagination_ids = new Map<string, number>();

const headers = AUTH_HEADERS;

describe("posts", () => {
    describe("simple operations", () => {
        test("create", async () => {
            const response = await fetch("localhost:8080/post/new", { method: "POST", body: JSON.stringify(test_post), headers });
            expect(!!response).toBeTrue();
            const result = (await response.json()) as Post;
            expect(result.id).toBeGreaterThan(0);
            expect(result.slug).toBe("bun-test-post");
            test_post_id = result.id;
        })
        test("duplicate slug", async () => {
            const response = await fetch("localhost:8080/post/new", { method: "POST", body: JSON.stringify({ ...test_post, title: "Duplicated Bun Post" }), headers })
            expect(response).toBeTruthy();
            expect(response.ok).toBeFalse();
            expect(response.status).toBeGreaterThanOrEqual(500);
        })
        test("update", async () => {
            // we are going to update the content to be "Updated Post Content"
            const response = await fetch("localhost:8080/post/edit", { method: "PUT", body: JSON.stringify({ ...test_post, content: " Updated Post Content" }), headers });
            expect(!!response).toBeTrue();

            // then re-fetch the post via id
            const fetch_response = await fetch(`localhost:8080/post/${test_post.slug}`, { method: "GET", headers });
            expect(!!fetch_response).toBeTrue();
            const result = (await fetch_response.json()) as Post;
            expect(result.id).toBe(test_post_id as number);
            expect(result.content == "Updated Post Content");
        })
        test("delete", async () => {
            const response = await fetch(`localhost:8080/post/delete/${test_post_id}`, { method: "DELETE", headers });
            expect(!!response).toBeTrue();
            expect(response.ok).toBeTrue();
        })
    })
    describe("errors", () => {
        test("not found post", async () => {
            const response = await fetch("localhost:8080/post/invalid-post", { method: "GET", headers });
            expect(response).toBeTruthy();
            expect(response.ok).toBeFalse();
            expect(response.status).toBe(404);
            expect(response.statusText).toBe("Not Found");
        });
    })
    describe("pagination", () => {
        test("mass-creation", async () => {
            for (const post of pagination_posts) {
                const response = await fetch("localhost:8080/post/new", { method: "POST", body: JSON.stringify({ ...post, author_id: 1 }), headers });
                expect(!!response).toBeTrue();
                expect(response.ok).toBeTrue();
                const result = (await response.json()) as Post;
                expect(result.id).toBeGreaterThan(0);
                pagination_ids.set(post.slug, result.id);
            }
            expect(pagination_ids.size == pagination_posts.length);
        })
        test("get", async () => {
            const response = await fetch("localhost:8080/posts", { method: "GET", headers });
            expect(!!response).toBeTrue();
            expect(response.ok).toBeTrue();
            const result = (await response.json()) as PostsResponse;
            expect(result.posts.length).toBe(10);
            expect(result.per_page).toBe(10);
            expect(result.current_page).toBe(1);
            expect(result.total_pages).toBeGreaterThan(1);

            const page2_response = await fetch("localhost:8080/posts?offset=10", { method: "GET", headers });
            expect(!!page2_response).toBeTrue();
            expect(page2_response.ok).toBeTrue();
            const page2 = (await page2_response.json()) as PostsResponse;
            expect(page2.current_page).toBe(2);
            expect(page2.posts.length).toBe(10);
        })
        test("get via category", async () => {
            // there should only be 2 'learning' posts
            const learning_response = await fetch("localhost:8080/posts/learning", { method: "GET", headers });
            expect(!!learning_response).toBeTrue();
            expect(learning_response.ok).toBeTrue();
            const learning = (await learning_response.json()) as PostsResponse;
            expect(learning.total_posts).toBe(2);
            expect(learning.total_pages).toBe(1);
            expect(learning.current_page).toBe(1);

            // now we get all 'coding' posts
            // these should include the child categories of coding
            const coding_response = await fetch("localhost:8080/posts/coding", { method: "GET", headers });
            expect(!!coding_response).toBeTrue();
            expect(coding_response.ok).toBeTrue();
            const coding = (await coding_response.json()) as PostsResponse;
            expect(coding.current_page).toBe(1);
            expect(coding.total_pages).toBeGreaterThan(1);
            // check for child categoriers
            expect(coding.posts.find((p) => p.category == "devlog")).toBeTruthy();
            expect(coding.posts.find((p) => p.category == "gamedev")).toBeTruthy();
        })
        afterAll(async () => {
            for (const [_, id] of pagination_ids) {
                const response = await fetch(`localhost:8080/post/delete/${id}`, { method: "DELETE", headers });
                expect(!!response).toBeTrue();
                expect(response.ok).toBeTrue();
            }
        })
    })
});
