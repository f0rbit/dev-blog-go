import { expect, test, describe, afterAll } from "bun:test";
import { exit } from "process";

type Category = {
    name: string,
    parent: string
}

type Post = {
    id: number,
    slug: string,
    title: string,
    content: string,
    category: string,
    created_at: string,
    updated_at: string
}

type PostsResponse = {
    posts: Post[],
    total_posts: number,
    total_pages: number,
    per_page: number,
    current_page: number
}

const test_post = {
    slug: "bun-test-post",
    title: "Bun Test Post",
    content: "this post should be removed after tests completed.",
    category: "root"
}
let test_post_id: number | null = null;

describe("posts", () => {
    describe("simple operations", () => {
        test("create", async () => {
            const response = await fetch("localhost:8080/post/new", { method: "POST", body: JSON.stringify(test_post) });
            expect(!!response);
            const result = (await response.json()) as Post;
            expect(result.id > 0);
            expect(result.slug == "bun-test-post");
            test_post_id = result.id;
        })
        test("duplicate slug", async () => {
            const response = await fetch("localhost:8080/post/new", { method: "POST", body: JSON.stringify({ ...test_post, title: "Duplicated Bun Post" }) })
            expect(!!response);
            expect(response.ok == false);
            expect(response.status >= 400 && response.status < 500); // some 4xxx status code
        })
        test("update", async () => {
            // we are going to update the content to be "Updated Post Content"
            const response = await fetch("localhost:8080/post/edit", { method: "PUT", body: JSON.stringify({ ...test_post, content:" Updated Post Content" }) });
            expect(!!response);

            // then re-fetch the post via id
            const fetch_response = await fetch(`localhost:8080/post/${test_post_id}`, { method: "GET" });
            expect(!!fetch_response);
            const result = (await fetch_response.json()) as Post;
            expect(result.id == test_post_id);
            expect(result.content == "Updated Post Content");            
        })
        test("delete", async () => {
            const response = await fetch(`localhost:8080/post/delete/${test_post_id}`, { method: "DELETE" });
            expect(!!response);
            expect(response.ok);
        })
    }),
        describe("pagination", () => {
            /** @todo this entire test block */
            test("mass-creation", () => {

            })
            test("get", () => {

            })
            test("get via category", () => {

            })
            afterAll(() => {

            })
        })
});
