import { z } from "zod";


const post_schema = z.object({
    id: z.number(),
    slug: z.string(),
    title: z.string(),
    content: z.string(),
    category: z.string(),
    tags: z.array(z.string()),
    archived: z.union([z.literal(0), z.literal(1)]),
    publish_at: z.string(),
    created_at: z.string(),
    updated_at: z.string()
});

export type Post = z.infer<typeof post_schema>;

const posts_response_schema = z.object({
    posts: z.array(post_schema),
    total_posts: z.number(),
    total_pages: z.number(),
    per_page: z.number(),
    current_page: z.number(),
})

export type PostsResponse = z.infer<typeof posts_response_schema>;

const category_schema = z.object({
    name: z.string(),
    parent: z.string()
});

const category_node_schema = z.object({
    name: z.string(),
    children: z.array(z.any())
})

export type Category = z.infer<typeof category_schema>;

// i don't think there's a way for zod to do self-recursive types yet.
export type CategoryNode = { name: string, children: CategoryNode[] };

const category_response = z.object({
    graph: category_node_schema,
    categories: z.array(category_schema)
});

export type CategoryResponse = {
    categories: Category[]
    graph: CategoryNode
};

export const SCHEMA = {
    POST: post_schema,
    POSTS_RESPONSE: posts_response_schema,
    CATEGORY: category_schema,
    CATEGORY_NODE: category_node_schema,
    CATEGORY_RESPONSE: category_response
}


