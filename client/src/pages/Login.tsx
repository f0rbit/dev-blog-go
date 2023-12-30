import { Github } from "lucide-react";

export function LoginPage() {
    return <section className="flex-col center">
        <a href="http://localhost:8080/auth/github/login" style={{textDecoration: "none"}}>
            <button><Github /><span>Login</span></button>
        </a>
    </section>
}
