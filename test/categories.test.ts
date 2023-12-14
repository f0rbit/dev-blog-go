import { describe, expect, test, beforeAll } from "bun:test";

type Category = {
    name: string,
    parent: string
}

test("categories", async () => {
    const response = await fetch("localhost:8080/categories", { method: "GET" });
    expect(!!response);
    const result = (await response.json()) as Category[];
    expect(!!result);
    expect(result.length > 0);
    expect(result.find((cat) => cat.name == "coding"));
    expect(result.find((cat) => cat.name == "webdev" && cat.parent == "devlog"));
});
