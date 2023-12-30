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
});
