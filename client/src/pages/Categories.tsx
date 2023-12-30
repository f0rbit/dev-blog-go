import { Dispatch, SetStateAction, useContext, useState } from "react";
import { AuthContext, PostContext } from "../App";
import { CategoryNode } from "../../schema";
import { Check, Plus, Save, Trash, X } from "lucide-react";
import { Oval } from "react-loader-spinner";

export function CategoriesPage() {
    const { categories } = useContext(PostContext);
    const [open, setOpen] = useState(false);

    if (!categories) return <main>No Categories!</main>;

    const elements = categories.graph.children.flatMap((c) => getCategoryElements(c, 0, "root"));

    return <main id='category-list'>
        {elements}
        <div className="category-card">
            <CategoryCreator right={false} open={open} setOpen={setOpen} parent={"root"} />
        </div>
    </main>
}

function getCategoryElements(root: CategoryNode, depth: number, parent: string) {
    const list: JSX.Element[] = [];

    list.push(<CategoryCard cat={root.name} depth={depth} parent={parent} />);
    if (root.children.length > 0) {
        for (const node of root.children) {
            list.push(...getCategoryElements(node, depth + 1, root.name));
        }
    }
    return list;
}

function CategoryCard({ cat, depth, parent }: { cat: string, depth: number, parent: string }) {
    const [open, setOpen] = useState(false);

    return <div style={{ marginLeft: (depth * 40) + "px" }} className="category-card">
        <span>{cat}</span>
        <CategoryCreator parent={parent} right={true} open={open} setOpen={setOpen} />
    </div>
}

function CategoryCreator({ parent, right, open, setOpen }: { right: boolean, open: boolean, setOpen: Dispatch<SetStateAction<boolean>>, parent: string }) {
    const [input, setInput] = useState("");
    const [saving, setSaving] = useState(false);
    const { categories, setCategories } = useContext(PostContext);
    const { user } = useContext(AuthContext)

    function save() {
        const cat = { name: input, parent, owner_id: user.user_id };
        setSaving(true);
        setOpen(false);
        // call new category
        // response should return the usual GET /categories

        setSaving(false);

    }
    if (saving) return <div className={"flex-row center" + (right ? "right" : "")}>
        <Oval width={16} height={16} strokeWidth={8} />
    </div>

    if (!open) return <div className={"flex-row center " + (right ? "right" : "")}>
        {right && <button title="Delete category"><Trash /></button>}
        <button onClick={() => setOpen(true)} title="Create child category"><Plus /></button>
    </div>

    const valid = input.length > 3 && !Object.values(categories.categories).map((c) => c.name).includes(input);

    return <div className={`flex-row ${right ? "right" : ""}`}>
        <input type="text" autoFocus={true} onChange={(e) => setInput(e.target.value)} />
        {valid ? <Check color="green" /> : <X className="error-message" />}
        <button onClick={() => save()}><Save /></button>
    </div>
}
