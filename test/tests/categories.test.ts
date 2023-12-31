import { CategoryResponse, SCHEMA } from "@client/schema";
import { expect, test,describe} from "bun:test";
import { AUTH_HEADERS } from "user";

test("get categories", async () => {
    const response = await fetch("localhost:8080/categories", { method: "GET", headers: AUTH_HEADERS });
    expect(response).toBeTruthy();
    const result = SCHEMA.CATEGORY_RESPONSE.parse(await response.json()) as CategoryResponse;
    expect(result).toBeTruthy();
    expect(result.graph).toBeTruthy();
    expect(result.categories).toBeTruthy();
    expect(result.categories.length).toBeGreaterThan(0);
    expect(result.categories.find((cat) => cat.name == "coding")).toBeTruthy();
    expect(result.categories.find((cat) => cat.name == "webdev" && cat.parent == "devlog")).toBeTruthy();
});


const test_category = {
    name: "Test",
    parent: "coding",
    owner_id: 1,
}

describe("categories", () => {
    test("create", async () => {
        const response = await fetch(`localhost:8080/category/new`, { method: "POST", headers: AUTH_HEADERS, body: JSON.stringify(test_category) });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
        const result = await response.json();
        const categories = SCHEMA.CATEGORY_RESPONSE.parse(result) as CategoryResponse;
        expect(categories.categories.map((c) => c.name)).toContain("Test");
        expect(categories.graph.children.find((c) => c.children.map((r) => r.name).includes("Test"))).toBeTruthy();
    })
    test("delete", async () => {
        const response = await fetch(`localhost:8080/category/delete/${test_category.name}`, { method: "DELETE", headers: AUTH_HEADERS });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
    })
});
