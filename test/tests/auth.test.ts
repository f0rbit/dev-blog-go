import { describe, test, expect } from "bun:test";

describe("auth tests", () => {
    test("invalid-token fails", async () => {
        const response = await fetch(`localhost:8080/auth/test?token=123123`, { method: "GET" }); 
        expect(response).toBeTruthy();
        expect(response.ok).toBeFalse();
    });
    test("valid-token succeeds", async () => {
        const response = await fetch(`localhost:8080/auth/test?token=${process.env.AUTH_TOKEN}`, { method: "GET" });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
    });
});
