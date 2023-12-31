import { Github } from "lucide-react";
import { API_URL } from "../App";

export function LoginPage() {
    return <section className="flex-col center">
        <a href={`${API_URL}/auth/github/login`} style={{textDecoration: "none"}}>
            <button><Github /><span>Login</span></button>
        </a>
    </section>
}
