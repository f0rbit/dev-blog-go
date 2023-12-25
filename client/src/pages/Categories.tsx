import { useContext } from "react";
import { PostContext } from "../App";
import { CategoryNode } from "../../schema";

export function CategoriesPage() {
    const { categories } = useContext(PostContext);

    if (!categories) return <main>No Categories!</main>;

    const elements = categories.graph.children.flatMap((c) => getCategoryElements(c, 0));

    return <main id='category-list'>
        {elements}
    </main>
}

function getCategoryElements(root: CategoryNode, depth: number) {
    const list: JSX.Element[] = [];

    const CategoryCard = ({ cat, depth }: { cat: string, depth: number }) => {
        return <div style={{ marginLeft: (depth * 40) + "px" }} className="category-card">{cat}</div>
    }

    list.push(<CategoryCard cat={root.name} depth={depth} />);
    if (root.children.length > 0) {
        for (const node of root.children) {
            list.push(...getCategoryElements(node, depth + 1));
        }
    }
    return list;
}
