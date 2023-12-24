import { describe, test, expect } from "bun:test";

const routes = [
    ["POST","/post/new"],
    ["PUT","/post/edit"],
    ["PUT", "/post/tag"],
    ["DELETE","/post/tag"]
] as const;
const headers = { 'Auth-Token': process.env.AUTH_TOKEN } as any as Headers;

describe("invalid json", () => {
    for (const route of routes) {
        test(`${route[0]} ${route[1]}`, async () => {
            const response = await fetch(`localhost:8080${route[1]}`, { method: route[0], body: "invalid json", headers });
            expect(response).toBeTruthy();
            expect(response.ok).toBeFalse();
            expect(response.status).toBe(400);
            expect(response.statusText).toBe("Bad Request");
        });
    }
});
