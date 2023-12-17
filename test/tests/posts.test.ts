import { expect, test, describe, afterAll } from "bun:test";

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
            /** @todo update test */
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
test("categories", async () => {
    const response = await fetch("localhost:8080/categories", { method: "GET" });
    expect(!!response);
    const result = (await response.json()) as Category[];
    expect(!!result);
    expect(result.length > 0);
    expect(result.find((cat) => cat.name == "coding"));
    expect(result.find((cat) => cat.name == "webdev" && cat.parent == "devlog"));
});
