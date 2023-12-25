import { CategoryResponse, SCHEMA } from "@client/schema";
import { expect, test } from "bun:test";

test("categories", async () => {
    const response = await fetch("localhost:8080/categories", { method: "GET" });
    expect(response).toBeTruthy();
    const result = SCHEMA.CATEGORY_RESPONSE.parse(await response.json()) as CategoryResponse;
    expect(result).toBeTruthy();
    expect(result.graph).toBeTruthy();
    expect(result.categories).toBeTruthy();
    expect(result.categories.length).toBeGreaterThan(0);
    expect(result.categories.find((cat) => cat.name == "coding")).toBeTruthy();
    expect(result.categories.find((cat) => cat.name == "webdev" && cat.parent == "devlog")).toBeTruthy();
});
