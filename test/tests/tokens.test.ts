import { expect, test, describe } from "bun:test";
import { AUTH_HEADERS } from "user";

const test_token = {
    user_id: 1,
    name: "New Test Token",
    note: "Generated during tests",
    enabled: true
}

let token_id: number | null = null;
describe("tokens", () => {
    test("create", async () => {
        const response = await fetch("localhost:8080/token/new", { method: "POST", headers: AUTH_HEADERS, body: JSON.stringify(test_token) });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
        const result = await response.json();
        expect(result).toBeTruthy();
        expect(result.id).toBeTruthy();
        token_id = result.id;
    })
    test("update", async () => {
        const response = await fetch("localhost:8080/token/edit", { method: "PUT", headers: AUTH_HEADERS, body: JSON.stringify({ ...test_token, id: token_id, name: "Updated Token" }) });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
    })
    test("get", async () => {
        const response = await fetch("localhost:8080/tokens", { method: "GET", headers: AUTH_HEADERS });
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
        const result = await response.json();
        expect(result).toBeTruthy();
        expect(Object.values(result).map((r: any) => r.name)).toContain("Updated Token");
    })
    test("delete", async () => {
        const response = await fetch(`localhost:8080/token/delete/${token_id}`, { method: "DELETE", headers: AUTH_HEADERS })
        expect(response).toBeTruthy();
        expect(response.ok).toBeTrue();
        const get_response = await fetch("localhost:8080/tokens", { method: "GET", headers: AUTH_HEADERS });
        expect(get_response).toBeTruthy();
        expect(get_response.ok).toBeTrue();
        const result = await get_response.json();
        expect(Object.values(result).find((r: any) => r.name == "Updated Token")).toBeFalsy();
    })
    test("delete not-existant", async () => {
        const response = await fetch(`localhost:8080/token/delete/${token_id}`, { method: "DELETE", headers: AUTH_HEADERS })
        expect(response).toBeTruthy();
        expect(response.ok).toBeFalse();
    })
});
