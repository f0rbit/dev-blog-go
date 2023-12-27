import { useContext, useState } from "react";
import { PostContext } from "../App";

/** @todo make this structure generic so we can use the same component for categories & tags */
/** @todo also rename all the category stuff */

type TagSetter = (tag: string) => void;

function TagInput({ tags, setValue, value }: { tags: string[], setValue: TagSetter, value: string }) {
    const [input, setInput] = useState<string>(value);
    const [open, setOpen] = useState<boolean>(false);
    const [hovering, setHovering] = useState<number>(0);

    const cat_list = tags.filter((c) => c.includes(input));

    function scrollView() {
        const elements = document.querySelector(".category-input-list");
        if (!elements) return;
        const selected = elements?.children?.[hovering];
        if (!selected) return;
        selected.scrollIntoView({
            behavior: "auto",
            block: "center",
            inline: "center"
        });
    }

    function handleKeyPress(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key == 'Enter' || e.key == 'Tab') {
            if (cat_list.length > 0) {
                const index = hovering >= 0 ? hovering : 0;
                setInput("");
                setValue(cat_list[index]);
                setOpen(false);
                setHovering(-1);
            } else {
                // cry
                const v = input;
                setInput("");
                setValue(v);
                setOpen(false);
                setHovering(-1);
            }
            e.preventDefault();
        } else if (e.key == 'ArrowUp') {
            setHovering(hovering - 1);
            e.preventDefault();
            scrollView();
            return false;
        } else if (e.key == 'ArrowDown') {
            setHovering(hovering + 1);
            e.preventDefault();
            scrollView();
            return false;
        }
    }

    if (hovering < -1) setHovering(-1);
    if (hovering >= cat_list.length) setHovering(cat_list.length - 1);


    return <div style={{ position: "relative" }}>
        <input 
            style={{ width: "100%" }} 
            type="text" 
            value={input} 
            onInput={() => { setHovering(-1); if (!open) setOpen(true); }} 
            onChange={(e) => { setInput(e.target.value); }} 
            onKeyDown={handleKeyPress} 
        />
        {(open && cat_list.length > 0) && <div className="category-input-list flex-col">
            {cat_list.map((cat, index) => <CategoryOption key={index} hovered={hovering == index} category={cat} setValue={(v) => { setInput(v); setValue(v); setOpen(false); setHovering(-1) }} />)}
        </div>}
    </div>
}

function CategoryOption({ category, setValue, hovered }: { category: string, setValue: TagSetter, hovered: boolean }) {
    const { posts } = useContext(PostContext);
    const count = posts.posts.filter((p) => p.tags.includes(category)).length;
    return <button onClick={(e) => { e.preventDefault(); e.stopPropagation(); setValue(category) }} className={hovered ? "hovered" : ""}>
        <span>{category}</span>
        <span style={{ marginLeft: "auto" }} className="post-count">{count} posts</span>
    </button>;
}

export default TagInput;
