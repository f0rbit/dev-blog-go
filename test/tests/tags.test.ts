import { beforeAll, expect, test, describe } from "bun:test";

const test_post = {
    slug: 'tag-test-post',
    title: "Tag Test Post",
    content: "this post should be removed after tests complete.",
    category: "root"
}
let test_id = null as number | null;

const tags = ["tag-1", "tag-2", "tag-3"] as const;
const headers = { 'Auth-Token': process.env.AUTH_TOKEN } as any as Headers;

beforeAll(async () => {
    // insert a fresh post so we have something to test against
    const response = await fetch("localhost:8080/post/new", { method: "POST", body: JSON.stringify(test_post), headers });
    expect(response).toBeTruthy();
    const result = (await response.json()) as any;
    expect(result.id).toBeTruthy();
    expect(result.slug).toBe(test_post.slug);
    test_id = result.id;
    console.log(test_id);
});

describe("tags", () => {
    test("add", async () => {
        const response = await fetch(`localhost:8080/post/tag?id=${test_id}&tag=test-tag`, { method: "PUT", headers });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
    })
    test("add duplicate", async () => {
        const response = await fetch(`localhost:8080/post/tag?id=${test_id}&tag=test-tag`, { method: "PUT", headers });
        expect(response).toBeTruthy();
        expect(response.ok).toBeFalse();
    })
    test("add multiple", async () => {
        for (const tag of tags) {
            const response = await fetch(`localhost:8080/post/tag?id=${test_id}&tag=${tag}`, { method: "PUT", headers });
            expect(response).toBeTruthy();
            expect(response.ok).toBeTrue();
        }
    })
    test("get tags", async () => {
        const response = await fetch("localhost:8080/tags", { method: "GET" });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
        const result = await response.json() as string[];
        for (const tag of tags) {
            expect(result).toContain(tag)
        }
        expect(result).toContain("test-tag");
    })
    test("get tagged posts", async () => {
        const response = await fetch("localhost:8080/posts?tag=test-tag", { method: "GET" });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
        const result = await response.json();
        expect(result).toBeTruthy();
        expect(result.total_posts).toBe(1);
        expect(result.posts[0].id).toBe(test_id);
    })
    test("delete tag", async () => {
        const response = await fetch(`localhost:8080/post/tag?id=${test_id}&tag=test-tag`, { method: "DELETE", headers });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
        // there should no longer be a "test-tag" in the tags
        const check = await fetch("localhost:8080/tags", { method: "GET" });
        const check_result = await check.json();
        expect(check_result).not.toContain("test-tag");
    })
})
