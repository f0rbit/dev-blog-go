import { useContext } from "react";
import { PostContext } from "../App";

export function TagsPage() {
    const { tags } = useContext(PostContext);
    return <main>
        <pre>{JSON.stringify(tags, null, 2)}</pre>
    </main>
}
