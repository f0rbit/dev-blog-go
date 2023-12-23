import { useContext, useState } from "react";
import { PostContext } from "../App";

type CategorySetter = (category: string) => void;

function CategoryInput({ categories, setValue, value }: { categories: { name: string, parent: string }[], setValue: CategorySetter, value: string }) {
    const [input, setInput] = useState<string>(value);
    const [open, setOpen] = useState<boolean>(false);
    const [hovering, setHovering] = useState<number>(0);

    const cat_list = Object.values(categories).map((c) => c.name).filter((c) => c.includes(input));

    function handleKeyPress(e: React.KeyboardEvent<HTMLInputElement>) {
        if (e.key == 'Enter' || e.key == 'Tab') {
            if (cat_list.length > 0) {
                const index = hovering >= 0 ? hovering : 0;
                setValue(cat_list[index]);
                setInput(cat_list[index]);
                setHovering(-1);
            } else {
                // cry
            }
            // typescript hack to trigger blur
            requestAnimationFrame(() => (e.target as any).blur());
            e.preventDefault();
        } else if (e.key == 'ArrowUp') {
            setHovering(hovering - 1);
            e.preventDefault();
            return false;
        } else if (e.key == 'ArrowDown') {
            setHovering(hovering + 1);
            e.preventDefault();
            return false;
        }
    }

    if (hovering < -1) setHovering(-1);
    if (hovering > cat_list.length) setHovering(cat_list.length - 1);


    return <div style={{ position: "relative" }}>
        <input 
            style={{ width: "100%" }} 
            type="text" 
            value={input} 
            onInput={() => { setHovering(-1); if (!open) setOpen(true); }} 
            onChange={(e) => setInput(e.target.value)} 
            onBlur={() => { setValue(input); setOpen(false); setHovering(-1); }} 
            onKeyDown={handleKeyPress} 
        />
        {(open && cat_list.length > 0) && <div className="category-input-list flex-col">
            {cat_list.map((cat, index) => <CategoryOption key={index} hovered={hovering == index} category={cat} setValue={(v) => { setInput(v); setValue(v); setOpen(false); setHovering(-1) }} />)}
        </div>}
    </div>
}

function CategoryOption({ category, setValue, hovered }: { category: string, setValue: CategorySetter, hovered: boolean }) {
    const { posts } = useContext(PostContext);
    const count = posts.posts.filter((p) => p.category == category).length;
    return <button onClick={() => setValue(category)} className={hovered ? "hovered" : ""}>
        <span>{category}</span>
        <span style={{ marginLeft: "auto" }} className="post-count">{count} posts</span>
    </button>;
}

export default CategoryInput;
