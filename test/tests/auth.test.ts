import { describe, test, expect } from "bun:test";
import { AUTH_HEADERS } from "user";


describe("auth tests", () => {
    test("invalid-token fails", async () => {
        const response = await fetch(`localhost:8080/auth/user`, { method: "GET", headers: { 'Auth-Token': "rubbish-token" } }); 
        expect(response).toBeTruthy();
        expect(response.ok).toBeFalse();
    });
    test("valid-token succeeds", async () => {
        const response = await fetch(`localhost:8080/auth/user`, { method: "GET", headers: AUTH_HEADERS });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
    });
    test("unauthorized-acces", async () => {
        const unauth_routes = [
            ["GET", "/posts"],
            ["GET", "/posts/coding"],
            ["GET", "/post/test-post"],
            ["POST", "/post/new"],
            ["PUT", "/post/edit"],
            ["DELETE", "/post/delete/1"],
            ["GET", "/categories"], 
            ["POST", "/category/new"], 
            ["DELETE", "/category/delete/coding"], 
            ["PUT", "/post/tag"],
            ["DELETE", "/post/tag"],
            ["GET", "/tags"], 
            ["GET", "/tokens"],
            ["POST", "/token/new"],
            ["PUT", "/token/edit"],
            ["DELETE", "/token/delete/1"]
        ];
        for (const [method, path] of unauth_routes) {
            const response = await fetch(`localhost:8080${path}`, { method });
            expect(response).toBeTruthy();
            expect(response.ok).toBeFalse();
            expect(response.status).toBe(401);
            expect(response.statusText).toBe("Unauthorized");
        }
    })
});
