import { expect, test } from "bun:test";

type Category = {
    name: string,
    parent: string
}

type CategoryNode = {
    name: string,
    children: CategoryNode[]
}

interface CategoryResponse {
    graph: CategoryNode,
    categories: Category[]
}

test("categories", async () => {
    const response = await fetch("localhost:8080/categories", { method: "GET" });
    expect(response).toBeTruthy();
    const result = (await response.json()) as CategoryResponse;
    expect(result).toBeTruthy();
    expect(result.graph).toBeTruthy();
    expect(result.categories).toBeTruthy();
    expect(result.categories.length).toBeGreaterThan(0);
    expect(result.categories.find((cat) => cat.name == "coding")).toBeTruthy();
    expect(result.categories.find((cat) => cat.name == "webdev" && cat.parent == "devlog")).toBeTruthy();
});
