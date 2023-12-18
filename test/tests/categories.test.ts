import { expect, test } from "bun:test";

type Category = {
    name: string,
    parent: string
}

test("categories", async () => {
    const response = await fetch("localhost:8080/categories", { method: "GET" });
    expect(response).toBeTruthy();
    const result = (await response.json()) as Category[];
    expect(result).toBeTruthy();
    expect(result.length).toBeGreaterThan(0);
    expect(result.find((cat) => cat.name == "coding")).toBeTruthy();
    expect(result.find((cat) => cat.name == "webdev" && cat.parent == "devlog")).toBeTruthy();
});
